package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/convox/praxis/cmd/build/source"
	"github.com/convox/praxis/manifest"
	"github.com/convox/praxis/provider"
	"github.com/convox/praxis/provider/local"
	"github.com/convox/praxis/provider/models"
)

var (
	flagApp      string
	flagAuth     string
	flagCache    string
	flagId       string
	flagManifest string
	flagMethod   string
	flagPush     string
	flagRelease  string
	flagUrl      string

	currentBuild    *models.Build
	currentLogs     string
	currentManifest string
	currentProvider provider.Provider
)

func init() {
	currentProvider = providerFromEnv()
}

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.StringVar(&flagApp, "app", "example", "app name")
	fs.StringVar(&flagAuth, "auth", "", "docker auth data (base64 encoded)")
	fs.StringVar(&flagCache, "cache", "true", "use docker cache")
	fs.StringVar(&flagId, "id", "", "build id")
	fs.StringVar(&flagManifest, "manifest", "convox.yml", "path to manifest")
	fs.StringVar(&flagMethod, "method", "", "source method")
	fs.StringVar(&flagPush, "push", "", "push to registry")
	fs.StringVar(&flagUrl, "url", "", "source url")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fail(err)
	}

	if v := os.Getenv("BUILD_APP"); v != "" {
		flagApp = v
	}

	if v := os.Getenv("BUILD_AUTH"); v != "" {
		flagAuth = v
	}

	if v := os.Getenv("BUILD_ID"); v != "" {
		flagId = v
	}

	if v := os.Getenv("BUILD_MANIFEST"); v != "" {
		flagManifest = v
	}

	if v := os.Getenv("BUILD_PUSH"); v != "" {
		flagPush = v
	}

	if v := os.Getenv("BUILD_URL"); v != "" {
		flagUrl = v
	}

	if flagId == "" {
		fail(fmt.Errorf("no build id"))
	}

	if err := execute(); err != nil {
		fail(err)
	}

	if err := success(); err != nil {
		fail(err)
	}
}

func execute() error {
	b, err := currentProvider.BuildLoad(flagApp, flagId)
	if err != nil {
		return err
	}

	currentBuild = b

	if err := login(); err != nil {
		return err
	}

	dir, err := fetch()
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	data, err := ioutil.ReadFile(filepath.Join(dir, flagManifest))
	if err != nil {
		return err
	}

	currentBuild.Manifest = string(data)

	if err := build(dir); err != nil {
		return err
	}

	return nil
}

func fetch() (string, error) {
	var s source.Source

	switch flagMethod {
	case "git":
		s = &source.SourceGit{flagUrl}
	// case "index":
	//   s = &source.SourceIndex{flagUrl}
	case "tgz":
		s = &source.SourceTgz{flagUrl}
	case "zip":
		s = &source.SourceZip{flagUrl}
	default:
		return "", fmt.Errorf("unknown method: %s", flagMethod)
	}

	var buf bytes.Buffer

	dir, err := s.Fetch(&buf)
	log(strings.TrimSpace(buf.String()))
	if err != nil {
		return "", err
	}

	return dir, nil
}

func login() error {
	var auth map[string]struct {
		Username string
		Password string
	}

	if err := json.Unmarshal([]byte(flagAuth), &auth); err != nil {
		return err
	}

	for host, entry := range auth {
		out, err := exec.Command("docker", "login", "-u", entry.Username, "-p", entry.Password, host).CombinedOutput()
		log(fmt.Sprintf("Authenticating %s: %s", host, strings.TrimSpace(string(out))))
		if err != nil {
			return err
		}
	}

	return nil
}

func build(dir string) error {
	dcy := filepath.Join(dir, flagManifest)

	if _, err := os.Stat(dcy); os.IsNotExist(err) {
		return fmt.Errorf("no such file: %s", flagManifest)
	}

	data, err := ioutil.ReadFile(dcy)
	if err != nil {
		return err
	}

	m, err := manifest.Load(data)
	if err != nil {
		return err
	}

	s := make(chan string)

	go func() {
		for l := range s {
			log(l)
		}
	}()

	defer close(s)

	if err := m.Build(manifest.BuildOptions{Cache: true}); err != nil {
		return err
	}

	if err := m.Push(flagPush, flagId); err != nil {
		return err
	}

	return nil
}

func success() error {
	// release := &models.Release{
	//   App: flagApp,
	// }

	// // TODO use provider.ReleaseFork()

	// if flagRelease != "" {
	//   r, err := currentProvider.ReleaseGet(flagApp, flagRelease)
	//   if err != nil {
	//     return err
	//   }
	//   release = r
	// }

	// release.Build = flagId
	// release.Created = time.Now()
	// release.Id = id("R", 10)
	// release.Manifest = currentBuild.Manifest

	// if err := currentProvider.ReleaseSave(release); err != nil {
	//   return err
	// }

	url, err := currentProvider.BlobStore(flagApp, fmt.Sprintf("convox/builds/%s/logs", currentBuild.Id), bytes.NewReader([]byte(currentLogs)), models.BlobStoreOptions{})
	if err != nil {
		return err
	}

	currentBuild.Ended = time.Now()
	currentBuild.Logs = url
	// currentBuild.Release = release.Id
	currentBuild.Status = "complete"

	if err := currentProvider.BuildSave(currentBuild); err != nil {
		return err
	}

	return nil
}

func fail(err error) {
	log(fmt.Sprintf("ERROR: %s", err))

	url, _ := currentProvider.BlobStore(flagApp, fmt.Sprintf("convox/builds/%s/logs", currentBuild.Id), bytes.NewReader([]byte(currentLogs)), models.BlobStoreOptions{})

	currentBuild.Ended = time.Now()
	currentBuild.Error = err.Error()
	currentBuild.Logs = url
	currentBuild.Status = "failed"

	if err := currentProvider.BuildSave(currentBuild); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

	os.Exit(1)
}

func log(line string) {
	currentLogs += fmt.Sprintf("%s\n", line)
	fmt.Println(line)
}

func providerFromEnv() provider.Provider {
	switch os.Getenv("PROVIDER") {
	default:
		return local.FromEnv()
	}
}
