// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	project "github.com/ca-risken/core/proto/project"
)

// ProjectServiceServer is an autogenerated mock type for the ProjectServiceServer type
type ProjectServiceServer struct {
	mock.Mock
}

// CleanProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) CleanProject(_a0 context.Context, _a1 *emptypb.Empty) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CleanProject")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) CreateProject(_a0 context.Context, _a1 *project.CreateProjectRequest) (*project.CreateProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateProject")
	}

	var r0 *project.CreateProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.CreateProjectRequest) (*project.CreateProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.CreateProjectRequest) *project.CreateProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.CreateProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.CreateProjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) DeleteProject(_a0 context.Context, _a1 *project.DeleteProjectRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProject")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.DeleteProjectRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.DeleteProjectRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.DeleteProjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsActive provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) IsActive(_a0 context.Context, _a1 *project.IsActiveRequest) (*project.IsActiveResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for IsActive")
	}

	var r0 *project.IsActiveResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.IsActiveRequest) (*project.IsActiveResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.IsActiveRequest) *project.IsActiveResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.IsActiveResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.IsActiveRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) ListProject(_a0 context.Context, _a1 *project.ListProjectRequest) (*project.ListProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListProject")
	}

	var r0 *project.ListProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.ListProjectRequest) (*project.ListProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.ListProjectRequest) *project.ListProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.ListProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.ListProjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TagProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) TagProject(_a0 context.Context, _a1 *project.TagProjectRequest) (*project.TagProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for TagProject")
	}

	var r0 *project.TagProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.TagProjectRequest) (*project.TagProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.TagProjectRequest) *project.TagProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.TagProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.TagProjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UntagProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) UntagProject(_a0 context.Context, _a1 *project.UntagProjectRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UntagProject")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.UntagProjectRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.UntagProjectRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.UntagProjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProject provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) UpdateProject(_a0 context.Context, _a1 *project.UpdateProjectRequest) (*project.UpdateProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProject")
	}

	var r0 *project.UpdateProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *project.UpdateProjectRequest) (*project.UpdateProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *project.UpdateProjectRequest) *project.UpdateProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.UpdateProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *project.UpdateProjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectServiceServer creates a new instance of ProjectServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectServiceServer {
	mock := &ProjectServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
