// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	project "github.com/ca-risken/core/proto/project"
)

// ProjectServiceClient is an autogenerated mock type for the ProjectServiceClient type
type ProjectServiceClient struct {
	mock.Mock
}

// CleanProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) CleanProject(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) CreateProject(ctx context.Context, in *project.CreateProjectRequest, opts ...grpc.CallOption) (*project.CreateProjectResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *project.CreateProjectResponse
	if rf, ok := ret.Get(0).(func(context.Context, *project.CreateProjectRequest, ...grpc.CallOption) *project.CreateProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.CreateProjectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.CreateProjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) DeleteProject(ctx context.Context, in *project.DeleteProjectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *project.DeleteProjectRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.DeleteProjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsActive provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) IsActive(ctx context.Context, in *project.IsActiveRequest, opts ...grpc.CallOption) (*project.IsActiveResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *project.IsActiveResponse
	if rf, ok := ret.Get(0).(func(context.Context, *project.IsActiveRequest, ...grpc.CallOption) *project.IsActiveResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.IsActiveResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.IsActiveRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) ListProject(ctx context.Context, in *project.ListProjectRequest, opts ...grpc.CallOption) (*project.ListProjectResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *project.ListProjectResponse
	if rf, ok := ret.Get(0).(func(context.Context, *project.ListProjectRequest, ...grpc.CallOption) *project.ListProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.ListProjectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.ListProjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TagProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) TagProject(ctx context.Context, in *project.TagProjectRequest, opts ...grpc.CallOption) (*project.TagProjectResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *project.TagProjectResponse
	if rf, ok := ret.Get(0).(func(context.Context, *project.TagProjectRequest, ...grpc.CallOption) *project.TagProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.TagProjectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.TagProjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UntagProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) UntagProject(ctx context.Context, in *project.UntagProjectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *project.UntagProjectRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.UntagProjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProject provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) UpdateProject(ctx context.Context, in *project.UpdateProjectRequest, opts ...grpc.CallOption) (*project.UpdateProjectResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *project.UpdateProjectResponse
	if rf, ok := ret.Get(0).(func(context.Context, *project.UpdateProjectRequest, ...grpc.CallOption) *project.UpdateProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*project.UpdateProjectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *project.UpdateProjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProjectServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewProjectServiceClient creates a new instance of ProjectServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProjectServiceClient(t mockConstructorTestingTNewProjectServiceClient) *ProjectServiceClient {
	mock := &ProjectServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
