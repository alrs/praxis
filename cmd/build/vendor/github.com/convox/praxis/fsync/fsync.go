package fsync

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/convox/praxis/fsync/watch"
	"github.com/fsouza/go-dockerclient"
)

type Stream chan string

type Sync struct {
	Container string
	Local     string
	Remote    string

	docker   *docker.Client
	lock     sync.Mutex
	incoming chan watch.Change
	outgoing chan watch.Change

	incomingBlocks map[string]int
	outgoingBlocks map[string]int
}

func NewSync(container, local, remote string) (*Sync, error) {
	l, err := filepath.Abs(local)

	if err != nil {
		return nil, err
	}

	sync := &Sync{
		Container: container,
		Local:     l,
		Remote:    remote,
	}

	sync.docker, _ = docker.NewClientFromEnv()
	sync.incoming = make(chan watch.Change)
	sync.outgoing = make(chan watch.Change)
	sync.incomingBlocks = make(map[string]int)
	sync.outgoingBlocks = make(map[string]int)

	return sync, nil
}

func (s *Sync) Contains(t Sync) bool {
	if !filepath.HasPrefix(t.Local, s.Local) {
		return false
	}

	lr, err := filepath.Rel(s.Local, t.Local)

	if err != nil {
		return false
	}

	rr, err := filepath.Rel(s.Remote, t.Remote)

	if err != nil {
		return false
	}

	return lr == rr
}

func (s *Sync) Start(st Stream) error {
	s.waitForContainer()

	res, err := s.docker.InspectContainer(s.Container)
	if err != nil {
		return err
	}

	for _, env := range res.Config.Env {
		parts := strings.SplitN(env, "=", 2)

		if len(parts) == 2 {
			s.Remote = strings.Replace(s.Remote, fmt.Sprintf("${%s}", parts[0]), parts[1], -1)
			s.Remote = strings.Replace(s.Remote, fmt.Sprintf("$%s", parts[0]), parts[1], -1)
		}
	}

	if !filepath.IsAbs(s.Remote) {
		fmt.Println("yep")
		s.Remote = filepath.Join(res.Config.WorkingDir, s.Remote)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(wd, s.Local)
	if err != nil {
		return err
	}

	st <- fmt.Sprintf("starting between %s and %s:%s", rel, s.Container, s.Remote)

	go s.watchIncoming(st)
	go s.watchOutgoing(st)

	incoming := []watch.Change{}
	outgoing := []watch.Change{}

	for {
		timeout := time.After(1 * time.Second)

		select {
		case c := <-s.incoming:
			incoming = append(incoming, c)
		case c := <-s.outgoing:
			outgoing = append(outgoing, c)
		case <-timeout:
			if len(incoming) > 0 {
				a, r := watch.Partition(incoming)
				s.syncIncomingAdds(a, st)
				s.syncIncomingRemoves(r, st)
				incoming = []watch.Change{}
			}
			if len(outgoing) > 0 {
				a, r := watch.Partition(outgoing)
				s.syncOutgoingAdds(a, st)
				s.syncOutgoingRemoves(r, st)
				outgoing = []watch.Change{}
			}
		}
	}

	return nil
}

func (s *Sync) syncIncomingAdds(adds []watch.Change, st Stream) {
	if len(adds) == 0 {
		return
	}

	cmd := []string{"tar", "czf", "-"}

	for _, a := range adds {
		cmd = append(cmd, filepath.Join(s.Remote, a.Path))
	}

	res, err := s.docker.CreateExec(docker.CreateExecOptions{
		AttachStdout: true,
		Container:    s.Container,
		Cmd:          cmd,
	})

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
		return
	}

	r, w := io.Pipe()

	go func() {
		gz, err := gzip.NewReader(r)

		if err != nil {
			st <- fmt.Sprintf("error: %s", err)
			return
		}

		tr := tar.NewReader(gz)

		for {
			header, err := tr.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				st <- fmt.Sprintf("error: %s", err)
				return
			}

			switch header.Typeflag {
			case tar.TypeReg:
				rel, err := filepath.Rel(s.Remote, filepath.Join("/", header.Name))

				if err != nil {
					st <- fmt.Sprintf("error: %s", err)
					return
				}

				local := filepath.Join(s.Local, rel)

				s.outgoingBlocks[rel] += 1

				os.MkdirAll(filepath.Dir(local), 0755)

				tmpfile, err := ioutil.TempFile("", filepath.Base(rel))
				if err != nil {
					st <- fmt.Sprintf("error: %s", err)
					return
				}

				_, err = io.Copy(tmpfile, tr)
				if err != nil {
					st <- fmt.Sprintf("error: %s", err)
					return
				}

				err = tmpfile.Close()
				if err != nil {
					st <- fmt.Sprintf("error: %s", err)
					return
				}

				err = os.Chmod(tmpfile.Name(), os.FileMode(header.Mode))
				if err != nil {
					st <- fmt.Sprintf("error: %s", err)
					return
				}

				err = os.Rename(tmpfile.Name(), local)
				if err != nil {
					st <- fmt.Sprintf("error: %s", err)
					return
				}
			}
		}
	}()

	err = s.docker.StartExec(res.ID, docker.StartExecOptions{
		OutputStream: w,
	})

	s.docker.WaitContainer(res.ID)

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
		return
	}

	if len(adds) < 5 {
		for _, a := range adds {
			st <- fmt.Sprintf("syncing %s from %s", a.Path, s.Container)
		}
	} else {
		st <- fmt.Sprintf("syncing %d files from %s", len(adds), s.Remote)
	}
}

func (s *Sync) syncIncomingRemoves(removes []watch.Change, st Stream) {
	// do not sync removes out from the container for safety
}

func (s *Sync) syncOutgoingAdds(adds []watch.Change, st Stream) {
	if len(adds) == 0 {
		return
	}

	var buf bytes.Buffer

	tgz := tar.NewWriter(&buf)

	for _, a := range adds {
		local := filepath.Join(a.Base, a.Path)

		info, err := os.Stat(local)

		if err != nil {
			continue
		}

		remote := filepath.Join(s.Remote, a.Path)

		s.incomingBlocks[a.Path] += 1

		tgz.WriteHeader(&tar.Header{
			Name:    remote,
			Mode:    0644,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})

		fd, err := os.Open(local)

		if err != nil {
			st <- fmt.Sprintf("error: %s", err)
			continue
		}

		io.Copy(tgz, fd)
		fd.Close()
	}

	if len(adds) < 5 {
		for _, a := range adds {
			st <- fmt.Sprintf("syncing %s to %s", a.Path, s.Container)
		}
	} else {
		st <- fmt.Sprintf("syncing %d files to %s", len(adds), s.Container)
	}

	tgz.Close()

	err := s.docker.UploadToContainer(s.Container, docker.UploadToContainerOptions{
		InputStream: &buf,
		Path:        "/",
	})

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
	}
}

func (s *Sync) syncOutgoingRemoves(removes []watch.Change, st Stream) {
	if len(removes) == 0 {
		return
	}

	cmd := []string{"rm", "-f"}

	for _, r := range removes {
		cmd = append(cmd, filepath.Join(s.Remote, r.Path))
	}

	res, err := s.docker.CreateExec(docker.CreateExecOptions{
		Container: s.Container,
		Cmd:       cmd,
	})

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
		return
	}

	err = s.docker.StartExec(res.ID, docker.StartExecOptions{
		Detach: true,
	})

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
		return
	}

	if len(removes) < 5 {
		for _, r := range removes {
			st <- fmt.Sprintf("removing %s from %s", r.Path, s.Container)
		}
	} else {
		st <- fmt.Sprintf("removing %d files from %s", len(removes), s.Container)
	}
}

func (s *Sync) uploadChangesDaemon(st Stream) {
	var buf bytes.Buffer

	tgz := tar.NewWriter(&buf)

	data, err := Asset("watchd")

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
	}

	tgz.WriteHeader(&tar.Header{
		Name: "watchd",
		Mode: 0755,
		Size: int64(len(data)),
	})

	tgz.Write(data)

	tgz.Close()

	err = s.docker.UploadToContainer(s.Container, docker.UploadToContainerOptions{
		InputStream: &buf,
		Path:        "/",
	})

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
	}
}

func (s *Sync) waitForContainer() {
	for {
		if res, err := s.docker.InspectContainer(s.Container); err == nil && res.State.Running {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Sync) watchIncoming(st Stream) {
	s.uploadChangesDaemon(st)

	res, err := s.docker.CreateExec(docker.CreateExecOptions{
		AttachStdout: true,
		Container:    s.Container,
		Cmd:          []string{"/watchd", s.Remote},
	})

	if err != nil {
		st <- fmt.Sprintf("error: %s", err)
		return
	}

	r, w := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			parts := strings.SplitN(scanner.Text(), "|", 3)

			if len(parts) != 3 {
				continue
			}

			// skip incoming removes for now. they make sync hard and not sure we want
			// the container deleting local files anyway
			if parts[0] == "remove" {
				continue
			}

			s.lock.Lock()

			if s.incomingBlocks[parts[2]] > 0 {
				s.incomingBlocks[parts[2]] -= 1
			} else {
				s.incoming <- watch.Change{
					Operation: parts[0],
					Base:      parts[1],
					Path:      parts[2],
				}
			}

			s.lock.Unlock()
		}
	}()

	err = s.docker.StartExec(res.ID, docker.StartExecOptions{
		OutputStream: w,
	})
}

func (s *Sync) watchOutgoing(st Stream) {
	ch := make(chan watch.Change, 1)

	go func() {
		if err := watch.Watch(s.Local, ch); err != nil {
			st <- fmt.Sprintf("error: %s", err)
		}
	}()

	for c := range ch {
		s.lock.Lock()

		if s.outgoingBlocks[c.Path] > 0 {
			s.outgoingBlocks[c.Path] -= 1
		} else {
			s.outgoing <- c
		}

		s.lock.Unlock()
	}
}
