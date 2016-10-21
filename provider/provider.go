package provider

import "io"

type Provider interface {
	AppCreate(name string, opts AppCreateOptions) (*App, error)
	AppDelete(name string) error

	BlobFetch(app, key string) (io.ReadCloser, error)
	BlobStore(app, key string, r io.Reader, opts BlobStoreOptions) (string, error)

	BuildCreate(app, url string, opts BuildCreateOptions) (*Build, error)
	BuildLoad(app, id string) (*Build, error)
	BuildLogs(app, id string) (io.ReadCloser, error)
	BuildSave(build *Build) error

	ProcessStart(app, service string, opts ProcessRunOptions) (*Process, error)
	ProcessWait(app, pid string) (int, error)

	TableDelete(app, table, id string) error
	TableLoad(app, table, id string) (map[string]string, error)
	TableSave(app, table, id string, attrs map[string]string) error
}
