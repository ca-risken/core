// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	iam "github.com/ca-risken/core/proto/iam"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"
)

// IAMServiceServer is an autogenerated mock type for the IAMServiceServer type
type IAMServiceServer struct {
	mock.Mock
}

// AnalyzeTokenExpiration provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) AnalyzeTokenExpiration(_a0 context.Context, _a1 *emptypb.Empty) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachAccessTokenRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) AttachAccessTokenRole(_a0 context.Context, _a1 *iam.AttachAccessTokenRoleRequest) (*iam.AttachAccessTokenRoleResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.AttachAccessTokenRoleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachAccessTokenRoleRequest) *iam.AttachAccessTokenRoleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AttachAccessTokenRoleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.AttachAccessTokenRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachPolicy provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) AttachPolicy(_a0 context.Context, _a1 *iam.AttachPolicyRequest) (*iam.AttachPolicyResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.AttachPolicyResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachPolicyRequest) *iam.AttachPolicyResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AttachPolicyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.AttachPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) AttachRole(_a0 context.Context, _a1 *iam.AttachRoleRequest) (*iam.AttachRoleResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.AttachRoleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AttachRoleRequest) *iam.AttachRoleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AttachRoleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.AttachRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticateAccessToken provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) AuthenticateAccessToken(_a0 context.Context, _a1 *iam.AuthenticateAccessTokenRequest) (*iam.AuthenticateAccessTokenResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.AuthenticateAccessTokenResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.AuthenticateAccessTokenRequest) *iam.AuthenticateAccessTokenResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.AuthenticateAccessTokenResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.AuthenticateAccessTokenRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccessToken provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DeleteAccessToken(_a0 context.Context, _a1 *iam.DeleteAccessTokenRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteAccessTokenRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteAccessTokenRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePolicy provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DeletePolicy(_a0 context.Context, _a1 *iam.DeletePolicyRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeletePolicyRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeletePolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DeleteRole(_a0 context.Context, _a1 *iam.DeleteRoleRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteRoleRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUserReserved provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DeleteUserReserved(_a0 context.Context, _a1 *iam.DeleteUserReservedRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteUserReservedRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteUserReservedRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachAccessTokenRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DetachAccessTokenRole(_a0 context.Context, _a1 *iam.DetachAccessTokenRoleRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachAccessTokenRoleRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DetachAccessTokenRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachPolicy provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DetachPolicy(_a0 context.Context, _a1 *iam.DetachPolicyRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachPolicyRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DetachPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) DetachRole(_a0 context.Context, _a1 *iam.DetachRoleRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DetachRoleRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DetachRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPolicy provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) GetPolicy(_a0 context.Context, _a1 *iam.GetPolicyRequest) (*iam.GetPolicyResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.GetPolicyResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetPolicyRequest) *iam.GetPolicyResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.GetPolicyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.GetPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) GetRole(_a0 context.Context, _a1 *iam.GetRoleRequest) (*iam.GetRoleResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.GetRoleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetRoleRequest) *iam.GetRoleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.GetRoleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.GetRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) GetUser(_a0 context.Context, _a1 *iam.GetUserRequest) (*iam.GetUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.GetUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.GetUserRequest) *iam.GetUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.GetUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.GetUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAdmin provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) IsAdmin(_a0 context.Context, _a1 *iam.IsAdminRequest) (*iam.IsAdminResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.IsAdminResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAdminRequest) *iam.IsAdminResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAdminResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAdminRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorized provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) IsAuthorized(_a0 context.Context, _a1 *iam.IsAuthorizedRequest) (*iam.IsAuthorizedResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.IsAuthorizedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedRequest) *iam.IsAuthorizedResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAuthorizedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAuthorizedRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorizedAdmin provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) IsAuthorizedAdmin(_a0 context.Context, _a1 *iam.IsAuthorizedAdminRequest) (*iam.IsAuthorizedAdminResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.IsAuthorizedAdminResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedAdminRequest) *iam.IsAuthorizedAdminResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAuthorizedAdminResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAuthorizedAdminRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorizedToken provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) IsAuthorizedToken(_a0 context.Context, _a1 *iam.IsAuthorizedTokenRequest) (*iam.IsAuthorizedTokenResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.IsAuthorizedTokenResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.IsAuthorizedTokenRequest) *iam.IsAuthorizedTokenResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.IsAuthorizedTokenResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.IsAuthorizedTokenRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAccessToken provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) ListAccessToken(_a0 context.Context, _a1 *iam.ListAccessTokenRequest) (*iam.ListAccessTokenResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.ListAccessTokenResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListAccessTokenRequest) *iam.ListAccessTokenResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListAccessTokenResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListAccessTokenRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPolicy provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) ListPolicy(_a0 context.Context, _a1 *iam.ListPolicyRequest) (*iam.ListPolicyResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.ListPolicyResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListPolicyRequest) *iam.ListPolicyResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListPolicyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) ListRole(_a0 context.Context, _a1 *iam.ListRoleRequest) (*iam.ListRoleResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.ListRoleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListRoleRequest) *iam.ListRoleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListRoleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUser provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) ListUser(_a0 context.Context, _a1 *iam.ListUserRequest) (*iam.ListUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.ListUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListUserRequest) *iam.ListUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserReserved provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) ListUserReserved(_a0 context.Context, _a1 *iam.ListUserReservedRequest) (*iam.ListUserReservedResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.ListUserReservedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListUserReservedRequest) *iam.ListUserReservedResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListUserReservedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListUserReservedRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAccessToken provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) PutAccessToken(_a0 context.Context, _a1 *iam.PutAccessTokenRequest) (*iam.PutAccessTokenResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.PutAccessTokenResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutAccessTokenRequest) *iam.PutAccessTokenResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutAccessTokenResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutAccessTokenRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPolicy provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) PutPolicy(_a0 context.Context, _a1 *iam.PutPolicyRequest) (*iam.PutPolicyResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.PutPolicyResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutPolicyRequest) *iam.PutPolicyResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutPolicyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutPolicyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRole provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) PutRole(_a0 context.Context, _a1 *iam.PutRoleRequest) (*iam.PutRoleResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.PutRoleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutRoleRequest) *iam.PutRoleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutRoleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutRoleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUser provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) PutUser(_a0 context.Context, _a1 *iam.PutUserRequest) (*iam.PutUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.PutUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutUserRequest) *iam.PutUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUserReserved provides a mock function with given fields: _a0, _a1
func (_m *IAMServiceServer) PutUserReserved(_a0 context.Context, _a1 *iam.PutUserReservedRequest) (*iam.PutUserReservedResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *iam.PutUserReservedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *iam.PutUserReservedRequest) *iam.PutUserReservedResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.PutUserReservedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.PutUserReservedRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIAMServiceServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewIAMServiceServer creates a new instance of IAMServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIAMServiceServer(t mockConstructorTestingTNewIAMServiceServer) *IAMServiceServer {
	mock := &IAMServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
