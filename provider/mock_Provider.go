package provider

import io "io"
import mock "github.com/stretchr/testify/mock"
import types "github.com/convox/praxis/types"

// MockProvider is an autogenerated mock type for the Provider type
type MockProvider struct {
	mock.Mock
}

// AppCreate provides a mock function with given fields: name
func (_m *MockProvider) AppCreate(name string) (*types.App, error) {
	ret := _m.Called(name)

	var r0 *types.App
	if rf, ok := ret.Get(0).(func(string) *types.App); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.App)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AppDelete provides a mock function with given fields: name
func (_m *MockProvider) AppDelete(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AppGet provides a mock function with given fields: name
func (_m *MockProvider) AppGet(name string) (*types.App, error) {
	ret := _m.Called(name)

	var r0 *types.App
	if rf, ok := ret.Get(0).(func(string) *types.App); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.App)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AppList provides a mock function with given fields:
func (_m *MockProvider) AppList() (types.Apps, error) {
	ret := _m.Called()

	var r0 types.Apps
	if rf, ok := ret.Get(0).(func() types.Apps); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Apps)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildCreate provides a mock function with given fields: app, url, opts
func (_m *MockProvider) BuildCreate(app string, url string, opts types.BuildCreateOptions) (*types.Build, error) {
	ret := _m.Called(app, url, opts)

	var r0 *types.Build
	if rf, ok := ret.Get(0).(func(string, string, types.BuildCreateOptions) *types.Build); ok {
		r0 = rf(app, url, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, types.BuildCreateOptions) error); ok {
		r1 = rf(app, url, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildGet provides a mock function with given fields: app, id
func (_m *MockProvider) BuildGet(app string, id string) (*types.Build, error) {
	ret := _m.Called(app, id)

	var r0 *types.Build
	if rf, ok := ret.Get(0).(func(string, string) *types.Build); ok {
		r0 = rf(app, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(app, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildLogs provides a mock function with given fields: app, id
func (_m *MockProvider) BuildLogs(app string, id string) (io.Reader, error) {
	ret := _m.Called(app, id)

	var r0 io.Reader
	if rf, ok := ret.Get(0).(func(string, string) io.Reader); ok {
		r0 = rf(app, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(app, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildUpdate provides a mock function with given fields: app, id, opts
func (_m *MockProvider) BuildUpdate(app string, id string, opts types.BuildUpdateOptions) (*types.Build, error) {
	ret := _m.Called(app, id, opts)

	var r0 *types.Build
	if rf, ok := ret.Get(0).(func(string, string, types.BuildUpdateOptions) *types.Build); ok {
		r0 = rf(app, id, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, types.BuildUpdateOptions) error); ok {
		r1 = rf(app, id, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObjectFetch provides a mock function with given fields: app, key
func (_m *MockProvider) ObjectFetch(app string, key string) (io.ReadCloser, error) {
	ret := _m.Called(app, key)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string, string) io.ReadCloser); ok {
		r0 = rf(app, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(app, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObjectStore provides a mock function with given fields: app, key, r, opts
func (_m *MockProvider) ObjectStore(app string, key string, r io.Reader, opts types.ObjectStoreOptions) (*types.Object, error) {
	ret := _m.Called(app, key, r, opts)

	var r0 *types.Object
	if rf, ok := ret.Get(0).(func(string, string, io.Reader, types.ObjectStoreOptions) *types.Object); ok {
		r0 = rf(app, key, r, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Object)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, io.Reader, types.ObjectStoreOptions) error); ok {
		r1 = rf(app, key, r, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessRun provides a mock function with given fields: app, opts
func (_m *MockProvider) ProcessRun(app string, opts types.ProcessRunOptions) error {
	ret := _m.Called(app, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, types.ProcessRunOptions) error); ok {
		r0 = rf(app, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReleaseCreate provides a mock function with given fields: app, opts
func (_m *MockProvider) ReleaseCreate(app string, opts types.ReleaseCreateOptions) (*types.Release, error) {
	ret := _m.Called(app, opts)

	var r0 *types.Release
	if rf, ok := ret.Get(0).(func(string, types.ReleaseCreateOptions) *types.Release); ok {
		r0 = rf(app, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Release)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, types.ReleaseCreateOptions) error); ok {
		r1 = rf(app, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReleaseGet provides a mock function with given fields: app, id
func (_m *MockProvider) ReleaseGet(app string, id string) (*types.Release, error) {
	ret := _m.Called(app, id)

	var r0 *types.Release
	if rf, ok := ret.Get(0).(func(string, string) *types.Release); ok {
		r0 = rf(app, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Release)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(app, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SystemGet provides a mock function with given fields:
func (_m *MockProvider) SystemGet() (*types.System, error) {
	ret := _m.Called()

	var r0 *types.System
	if rf, ok := ret.Get(0).(func() *types.System); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.System)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

var _ Provider = (*MockProvider)(nil)
