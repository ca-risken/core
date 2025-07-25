// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	iam "github.com/ca-risken/core/proto/iam"

	mock "github.com/stretchr/testify/mock"
)

// IAMServiceClient is an autogenerated mock type for the IAMServiceClient type
type IAMServiceClient struct {
	mock.Mock
}

// AnalyzeTokenExpiration provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) AnalyzeTokenExpiration(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AnalyzeTokenExpiration")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachAccessTokenRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) AttachAccessTokenRole(ctx context.Context, in *iam.AttachAccessTokenRoleRequest, opts ...grpc.CallOption) (*iam.AttachAccessTokenRoleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AttachAccessTokenRole")
	}

	var r0 *iam.AttachAccessTokenRoleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachAccessTokenRoleRequest, ...grpc.CallOption) (*iam.AttachAccessTokenRoleResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachAccessTokenRoleRequest, ...grpc.CallOption) *iam.AttachAccessTokenRoleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AttachAccessTokenRoleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.AttachAccessTokenRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachPolicy provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) AttachPolicy(ctx context.Context, in *iam.AttachPolicyRequest, opts ...grpc.CallOption) (*iam.AttachPolicyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AttachPolicy")
	}

	var r0 *iam.AttachPolicyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachPolicyRequest, ...grpc.CallOption) (*iam.AttachPolicyResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachPolicyRequest, ...grpc.CallOption) *iam.AttachPolicyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AttachPolicyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.AttachPolicyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) AttachRole(ctx context.Context, in *iam.AttachRoleRequest, opts ...grpc.CallOption) (*iam.AttachRoleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AttachRole")
	}

	var r0 *iam.AttachRoleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachRoleRequest, ...grpc.CallOption) (*iam.AttachRoleResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachRoleRequest, ...grpc.CallOption) *iam.AttachRoleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AttachRoleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.AttachRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticateAccessToken provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) AuthenticateAccessToken(ctx context.Context, in *iam.AuthenticateAccessTokenRequest, opts ...grpc.CallOption) (*iam.AuthenticateAccessTokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AuthenticateAccessToken")
	}

	var r0 *iam.AuthenticateAccessTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AuthenticateAccessTokenRequest, ...grpc.CallOption) (*iam.AuthenticateAccessTokenResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AuthenticateAccessTokenRequest, ...grpc.CallOption) *iam.AuthenticateAccessTokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AuthenticateAccessTokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.AuthenticateAccessTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccessToken provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DeleteAccessToken(ctx context.Context, in *iam.DeleteAccessTokenRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAccessToken")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteAccessTokenRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteAccessTokenRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteAccessTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePolicy provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DeletePolicy(ctx context.Context, in *iam.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeletePolicy")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeletePolicyRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeletePolicyRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeletePolicyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DeleteRole(ctx context.Context, in *iam.DeleteRoleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRole")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteRoleRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteRoleRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUserReserved provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DeleteUserReserved(ctx context.Context, in *iam.DeleteUserReservedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUserReserved")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteUserReservedRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteUserReservedRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteUserReservedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachAccessTokenRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DetachAccessTokenRole(ctx context.Context, in *iam.DetachAccessTokenRoleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DetachAccessTokenRole")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachAccessTokenRoleRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachAccessTokenRoleRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DetachAccessTokenRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachPolicy provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DetachPolicy(ctx context.Context, in *iam.DetachPolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DetachPolicy")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachPolicyRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachPolicyRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DetachPolicyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) DetachRole(ctx context.Context, in *iam.DetachRoleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DetachRole")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachRoleRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachRoleRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.DetachRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPolicy provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) GetPolicy(ctx context.Context, in *iam.GetPolicyRequest, opts ...grpc.CallOption) (*iam.GetPolicyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetPolicy")
	}

	var r0 *iam.GetPolicyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetPolicyRequest, ...grpc.CallOption) (*iam.GetPolicyResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetPolicyRequest, ...grpc.CallOption) *iam.GetPolicyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.GetPolicyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.GetPolicyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) GetRole(ctx context.Context, in *iam.GetRoleRequest, opts ...grpc.CallOption) (*iam.GetRoleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetRole")
	}

	var r0 *iam.GetRoleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetRoleRequest, ...grpc.CallOption) (*iam.GetRoleResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetRoleRequest, ...grpc.CallOption) *iam.GetRoleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.GetRoleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.GetRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) GetUser(ctx context.Context, in *iam.GetUserRequest, opts ...grpc.CallOption) (*iam.GetUserResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *iam.GetUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetUserRequest, ...grpc.CallOption) (*iam.GetUserResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetUserRequest, ...grpc.CallOption) *iam.GetUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.GetUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.GetUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAdmin provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) IsAdmin(ctx context.Context, in *iam.IsAdminRequest, opts ...grpc.CallOption) (*iam.IsAdminResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for IsAdmin")
	}

	var r0 *iam.IsAdminResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAdminRequest, ...grpc.CallOption) (*iam.IsAdminResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAdminRequest, ...grpc.CallOption) *iam.IsAdminResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAdminResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAdminRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorized provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) IsAuthorized(ctx context.Context, in *iam.IsAuthorizedRequest, opts ...grpc.CallOption) (*iam.IsAuthorizedResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for IsAuthorized")
	}

	var r0 *iam.IsAuthorizedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedRequest, ...grpc.CallOption) (*iam.IsAuthorizedResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedRequest, ...grpc.CallOption) *iam.IsAuthorizedResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAuthorizedResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAuthorizedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorizedAdmin provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) IsAuthorizedAdmin(ctx context.Context, in *iam.IsAuthorizedAdminRequest, opts ...grpc.CallOption) (*iam.IsAuthorizedAdminResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for IsAuthorizedAdmin")
	}

	var r0 *iam.IsAuthorizedAdminResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedAdminRequest, ...grpc.CallOption) (*iam.IsAuthorizedAdminResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedAdminRequest, ...grpc.CallOption) *iam.IsAuthorizedAdminResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAuthorizedAdminResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAuthorizedAdminRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorizedToken provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) IsAuthorizedToken(ctx context.Context, in *iam.IsAuthorizedTokenRequest, opts ...grpc.CallOption) (*iam.IsAuthorizedTokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for IsAuthorizedToken")
	}

	var r0 *iam.IsAuthorizedTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedTokenRequest, ...grpc.CallOption) (*iam.IsAuthorizedTokenResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedTokenRequest, ...grpc.CallOption) *iam.IsAuthorizedTokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAuthorizedTokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAuthorizedTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAccessToken provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) ListAccessToken(ctx context.Context, in *iam.ListAccessTokenRequest, opts ...grpc.CallOption) (*iam.ListAccessTokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListAccessToken")
	}

	var r0 *iam.ListAccessTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListAccessTokenRequest, ...grpc.CallOption) (*iam.ListAccessTokenResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListAccessTokenRequest, ...grpc.CallOption) *iam.ListAccessTokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListAccessTokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListAccessTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPolicy provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) ListPolicy(ctx context.Context, in *iam.ListPolicyRequest, opts ...grpc.CallOption) (*iam.ListPolicyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListPolicy")
	}

	var r0 *iam.ListPolicyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListPolicyRequest, ...grpc.CallOption) (*iam.ListPolicyResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListPolicyRequest, ...grpc.CallOption) *iam.ListPolicyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListPolicyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListPolicyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) ListRole(ctx context.Context, in *iam.ListRoleRequest, opts ...grpc.CallOption) (*iam.ListRoleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListRole")
	}

	var r0 *iam.ListRoleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListRoleRequest, ...grpc.CallOption) (*iam.ListRoleResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListRoleRequest, ...grpc.CallOption) *iam.ListRoleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListRoleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUser provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) ListUser(ctx context.Context, in *iam.ListUserRequest, opts ...grpc.CallOption) (*iam.ListUserResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListUser")
	}

	var r0 *iam.ListUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListUserRequest, ...grpc.CallOption) (*iam.ListUserResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListUserRequest, ...grpc.CallOption) *iam.ListUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserReserved provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) ListUserReserved(ctx context.Context, in *iam.ListUserReservedRequest, opts ...grpc.CallOption) (*iam.ListUserReservedResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListUserReserved")
	}

	var r0 *iam.ListUserReservedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListUserReservedRequest, ...grpc.CallOption) (*iam.ListUserReservedResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListUserReservedRequest, ...grpc.CallOption) *iam.ListUserReservedResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListUserReservedResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListUserReservedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAccessToken provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) PutAccessToken(ctx context.Context, in *iam.PutAccessTokenRequest, opts ...grpc.CallOption) (*iam.PutAccessTokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PutAccessToken")
	}

	var r0 *iam.PutAccessTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutAccessTokenRequest, ...grpc.CallOption) (*iam.PutAccessTokenResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutAccessTokenRequest, ...grpc.CallOption) *iam.PutAccessTokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutAccessTokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutAccessTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPolicy provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) PutPolicy(ctx context.Context, in *iam.PutPolicyRequest, opts ...grpc.CallOption) (*iam.PutPolicyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PutPolicy")
	}

	var r0 *iam.PutPolicyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutPolicyRequest, ...grpc.CallOption) (*iam.PutPolicyResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutPolicyRequest, ...grpc.CallOption) *iam.PutPolicyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutPolicyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutPolicyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRole provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) PutRole(ctx context.Context, in *iam.PutRoleRequest, opts ...grpc.CallOption) (*iam.PutRoleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PutRole")
	}

	var r0 *iam.PutRoleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutRoleRequest, ...grpc.CallOption) (*iam.PutRoleResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutRoleRequest, ...grpc.CallOption) *iam.PutRoleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutRoleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutRoleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUser provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) PutUser(ctx context.Context, in *iam.PutUserRequest, opts ...grpc.CallOption) (*iam.PutUserResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PutUser")
	}

	var r0 *iam.PutUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutUserRequest, ...grpc.CallOption) (*iam.PutUserResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutUserRequest, ...grpc.CallOption) *iam.PutUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUserReserved provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) PutUserReserved(ctx context.Context, in *iam.PutUserReservedRequest, opts ...grpc.CallOption) (*iam.PutUserReservedResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PutUserReserved")
	}

	var r0 *iam.PutUserReservedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutUserReservedRequest, ...grpc.CallOption) (*iam.PutUserReservedResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutUserReservedRequest, ...grpc.CallOption) *iam.PutUserReservedResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutUserReservedResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutUserReservedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserAdmin provides a mock function with given fields: ctx, in, opts
func (_m *IAMServiceClient) UpdateUserAdmin(ctx context.Context, in *iam.UpdateUserAdminRequest, opts ...grpc.CallOption) (*iam.UpdateUserAdminResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserAdmin")
	}

	var r0 *iam.UpdateUserAdminResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *iam.UpdateUserAdminRequest, ...grpc.CallOption) (*iam.UpdateUserAdminResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *iam.UpdateUserAdminRequest, ...grpc.CallOption) *iam.UpdateUserAdminResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.UpdateUserAdminResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *iam.UpdateUserAdminRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIAMServiceClient creates a new instance of IAMServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAMServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAMServiceClient {
	mock := &IAMServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
