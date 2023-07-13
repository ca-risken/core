// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	db "github.com/ca-risken/core/pkg/db"
	mock "github.com/stretchr/testify/mock"

	model "github.com/ca-risken/core/pkg/model"
)

// IAMRepository is an autogenerated mock type for the IAMRepository type
type IAMRepository struct {
	mock.Mock
}

// AttachAccessTokenRole provides a mock function with given fields: ctx, projectID, roleID, accessTokenID
func (_m *IAMRepository) AttachAccessTokenRole(ctx context.Context, projectID uint32, roleID uint32, accessTokenID uint32) (*model.AccessTokenRole, error) {
	ret := _m.Called(ctx, projectID, roleID, accessTokenID)

	var r0 *model.AccessTokenRole
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) (*model.AccessTokenRole, error)); ok {
		return rf(ctx, projectID, roleID, accessTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) *model.AccessTokenRole); ok {
		r0 = rf(ctx, projectID, roleID, accessTokenID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessTokenRole)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, roleID, accessTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachAllAdminRole provides a mock function with given fields: ctx, userID
func (_m *IAMRepository) AttachAllAdminRole(ctx context.Context, userID uint32) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AttachPolicy provides a mock function with given fields: ctx, projectID, roleID, policyID
func (_m *IAMRepository) AttachPolicy(ctx context.Context, projectID uint32, roleID uint32, policyID uint32) (*model.RolePolicy, error) {
	ret := _m.Called(ctx, projectID, roleID, policyID)

	var r0 *model.RolePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) (*model.RolePolicy, error)); ok {
		return rf(ctx, projectID, roleID, policyID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) *model.RolePolicy); ok {
		r0 = rf(ctx, projectID, roleID, policyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RolePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, roleID, policyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachRole provides a mock function with given fields: ctx, projectID, roleID, userID
func (_m *IAMRepository) AttachRole(ctx context.Context, projectID uint32, roleID uint32, userID uint32) (*model.UserRole, error) {
	ret := _m.Called(ctx, projectID, roleID, userID)

	var r0 *model.UserRole
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) (*model.UserRole, error)); ok {
		return rf(ctx, projectID, roleID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) *model.UserRole); ok {
		r0 = rf(ctx, projectID, roleID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserRole)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, roleID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, u
func (_m *IAMRepository) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	ret := _m.Called(ctx, u)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) (*model.User, error)); ok {
		return rf(ctx, u)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *model.User); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccessToken provides a mock function with given fields: ctx, projectID, accessTokenID
func (_m *IAMRepository) DeleteAccessToken(ctx context.Context, projectID uint32, accessTokenID uint32) error {
	ret := _m.Called(ctx, projectID, accessTokenID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, accessTokenID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePolicy provides a mock function with given fields: ctx, projectID, policyID
func (_m *IAMRepository) DeletePolicy(ctx context.Context, projectID uint32, policyID uint32) error {
	ret := _m.Called(ctx, projectID, policyID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, policyID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRole provides a mock function with given fields: ctx, projectID, roleID
func (_m *IAMRepository) DeleteRole(ctx context.Context, projectID uint32, roleID uint32) error {
	ret := _m.Called(ctx, projectID, roleID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, roleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUserReserved provides a mock function with given fields: ctx, projectID, reservedID
func (_m *IAMRepository) DeleteUserReserved(ctx context.Context, projectID uint32, reservedID uint32) error {
	ret := _m.Called(ctx, projectID, reservedID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, reservedID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DetachAccessTokenRole provides a mock function with given fields: ctx, projectID, roleID, accessTokenID
func (_m *IAMRepository) DetachAccessTokenRole(ctx context.Context, projectID uint32, roleID uint32, accessTokenID uint32) error {
	ret := _m.Called(ctx, projectID, roleID, accessTokenID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, roleID, accessTokenID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DetachPolicy provides a mock function with given fields: ctx, projectID, roleID, policyID
func (_m *IAMRepository) DetachPolicy(ctx context.Context, projectID uint32, roleID uint32, policyID uint32) error {
	ret := _m.Called(ctx, projectID, roleID, policyID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, roleID, policyID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DetachRole provides a mock function with given fields: ctx, projectID, roleID, userID
func (_m *IAMRepository) DetachRole(ctx context.Context, projectID uint32, roleID uint32, userID uint32) error {
	ret := _m.Called(ctx, projectID, roleID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, roleID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExistsAccessTokenMaintainer provides a mock function with given fields: ctx, projectID, accessTokenID
func (_m *IAMRepository) ExistsAccessTokenMaintainer(ctx context.Context, projectID uint32, accessTokenID uint32) (bool, error) {
	ret := _m.Called(ctx, projectID, accessTokenID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (bool, error)); ok {
		return rf(ctx, projectID, accessTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) bool); ok {
		r0 = rf(ctx, projectID, accessTokenID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, accessTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccessTokenByUniqueKey provides a mock function with given fields: ctx, projectID, name
func (_m *IAMRepository) GetAccessTokenByUniqueKey(ctx context.Context, projectID uint32, name string) (*model.AccessToken, error) {
	ret := _m.Called(ctx, projectID, name)

	var r0 *model.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*model.AccessToken, error)); ok {
		return rf(ctx, projectID, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *model.AccessToken); ok {
		r0 = rf(ctx, projectID, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(ctx, projectID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetActiveAccessTokenHash provides a mock function with given fields: ctx, projectID, accessTokenID, tokenHash
func (_m *IAMRepository) GetActiveAccessTokenHash(ctx context.Context, projectID uint32, accessTokenID uint32, tokenHash string) (*model.AccessToken, error) {
	ret := _m.Called(ctx, projectID, accessTokenID, tokenHash)

	var r0 *model.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, string) (*model.AccessToken, error)); ok {
		return rf(ctx, projectID, accessTokenID, tokenHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, string) *model.AccessToken); ok {
		r0 = rf(ctx, projectID, accessTokenID, tokenHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, string) error); ok {
		r1 = rf(ctx, projectID, accessTokenID, tokenHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetActiveUserCount provides a mock function with given fields: ctx
func (_m *IAMRepository) GetActiveUserCount(ctx context.Context) (*int, error) {
	ret := _m.Called(ctx)

	var r0 *int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAdminPolicy provides a mock function with given fields: ctx, userID
func (_m *IAMRepository) GetAdminPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error) {
	ret := _m.Called(ctx, userID)

	var r0 *[]model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (*[]model.Policy, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) *[]model.Policy); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPolicy provides a mock function with given fields: ctx, projectID, policyID
func (_m *IAMRepository) GetPolicy(ctx context.Context, projectID uint32, policyID uint32) (*model.Policy, error) {
	ret := _m.Called(ctx, projectID, policyID)

	var r0 *model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.Policy, error)); ok {
		return rf(ctx, projectID, policyID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.Policy); ok {
		r0 = rf(ctx, projectID, policyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, policyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPolicyByName provides a mock function with given fields: ctx, projectID, name
func (_m *IAMRepository) GetPolicyByName(ctx context.Context, projectID uint32, name string) (*model.Policy, error) {
	ret := _m.Called(ctx, projectID, name)

	var r0 *model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*model.Policy, error)); ok {
		return rf(ctx, projectID, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *model.Policy); ok {
		r0 = rf(ctx, projectID, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(ctx, projectID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRole provides a mock function with given fields: ctx, projectID, roleID
func (_m *IAMRepository) GetRole(ctx context.Context, projectID uint32, roleID uint32) (*model.Role, error) {
	ret := _m.Called(ctx, projectID, roleID)

	var r0 *model.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.Role, error)); ok {
		return rf(ctx, projectID, roleID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.Role); ok {
		r0 = rf(ctx, projectID, roleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, roleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRoleByName provides a mock function with given fields: ctx, projectID, name
func (_m *IAMRepository) GetRoleByName(ctx context.Context, projectID uint32, name string) (*model.Role, error) {
	ret := _m.Called(ctx, projectID, name)

	var r0 *model.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*model.Role, error)); ok {
		return rf(ctx, projectID, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *model.Role); ok {
		r0 = rf(ctx, projectID, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(ctx, projectID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenPolicy provides a mock function with given fields: ctx, accessTokenID
func (_m *IAMRepository) GetTokenPolicy(ctx context.Context, accessTokenID uint32) (*[]model.Policy, error) {
	ret := _m.Called(ctx, accessTokenID)

	var r0 *[]model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (*[]model.Policy, error)); ok {
		return rf(ctx, accessTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) *[]model.Policy); ok {
		r0 = rf(ctx, accessTokenID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, accessTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, userID, sub
func (_m *IAMRepository) GetUser(ctx context.Context, userID uint32, sub string) (*model.User, error) {
	ret := _m.Called(ctx, userID, sub)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*model.User, error)); ok {
		return rf(ctx, userID, sub)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *model.User); ok {
		r0 = rf(ctx, userID, sub)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(ctx, userID, sub)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserBySub provides a mock function with given fields: ctx, sub
func (_m *IAMRepository) GetUserBySub(ctx context.Context, sub string) (*model.User, error) {
	ret := _m.Called(ctx, sub)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, sub)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, sub)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, sub)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUserIdpKey provides a mock function with given fields: ctx, userIdpKey
func (_m *IAMRepository) GetUserByUserIdpKey(ctx context.Context, userIdpKey string) (*model.User, error) {
	ret := _m.Called(ctx, userIdpKey)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, userIdpKey)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, userIdpKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userIdpKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserPolicy provides a mock function with given fields: ctx, userID
func (_m *IAMRepository) GetUserPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error) {
	ret := _m.Called(ctx, userID)

	var r0 *[]model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (*[]model.Policy, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) *[]model.Policy); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAccessToken provides a mock function with given fields: ctx, projectID, name, accessTokenID
func (_m *IAMRepository) ListAccessToken(ctx context.Context, projectID uint32, name string, accessTokenID uint32) (*[]model.AccessToken, error) {
	ret := _m.Called(ctx, projectID, name, accessTokenID)

	var r0 *[]model.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, uint32) (*[]model.AccessToken, error)); ok {
		return rf(ctx, projectID, name, accessTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, uint32) *[]model.AccessToken); ok {
		r0 = rf(ctx, projectID, name, accessTokenID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string, uint32) error); ok {
		r1 = rf(ctx, projectID, name, accessTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListExpiredAccessToken provides a mock function with given fields: ctx
func (_m *IAMRepository) ListExpiredAccessToken(ctx context.Context) (*[]model.AccessToken, error) {
	ret := _m.Called(ctx)

	var r0 *[]model.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]model.AccessToken, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]model.AccessToken); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPolicy provides a mock function with given fields: ctx, projectID, name, roleID
func (_m *IAMRepository) ListPolicy(ctx context.Context, projectID uint32, name string, roleID uint32) (*[]model.Policy, error) {
	ret := _m.Called(ctx, projectID, name, roleID)

	var r0 *[]model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, uint32) (*[]model.Policy, error)); ok {
		return rf(ctx, projectID, name, roleID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, uint32) *[]model.Policy); ok {
		r0 = rf(ctx, projectID, name, roleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string, uint32) error); ok {
		r1 = rf(ctx, projectID, name, roleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRole provides a mock function with given fields: ctx, projectID, name, userID, accessTokenID
func (_m *IAMRepository) ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error) {
	ret := _m.Called(ctx, projectID, name, userID, accessTokenID)

	var r0 *[]model.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, uint32, uint32) (*[]model.Role, error)); ok {
		return rf(ctx, projectID, name, userID, accessTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, uint32, uint32) *[]model.Role); ok {
		r0 = rf(ctx, projectID, name, userID, accessTokenID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, name, userID, accessTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUser provides a mock function with given fields: ctx, activated, projectID, name, userID, admin
func (_m *IAMRepository) ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32, admin bool) (*[]model.User, error) {
	ret := _m.Called(ctx, activated, projectID, name, userID, admin)

	var r0 *[]model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, uint32, string, uint32, bool) (*[]model.User, error)); ok {
		return rf(ctx, activated, projectID, name, userID, admin)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool, uint32, string, uint32, bool) *[]model.User); ok {
		r0 = rf(ctx, activated, projectID, name, userID, admin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool, uint32, string, uint32, bool) error); ok {
		r1 = rf(ctx, activated, projectID, name, userID, admin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserReserved provides a mock function with given fields: ctx, projectID, userIdpKey
func (_m *IAMRepository) ListUserReserved(ctx context.Context, projectID uint32, userIdpKey string) (*[]model.UserReserved, error) {
	ret := _m.Called(ctx, projectID, userIdpKey)

	var r0 *[]model.UserReserved
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*[]model.UserReserved, error)); ok {
		return rf(ctx, projectID, userIdpKey)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *[]model.UserReserved); ok {
		r0 = rf(ctx, projectID, userIdpKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.UserReserved)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(ctx, projectID, userIdpKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserReservedWithProjectID provides a mock function with given fields: ctx, userIdpKey
func (_m *IAMRepository) ListUserReservedWithProjectID(ctx context.Context, userIdpKey string) (*[]db.UserReservedWithProjectID, error) {
	ret := _m.Called(ctx, userIdpKey)

	var r0 *[]db.UserReservedWithProjectID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*[]db.UserReservedWithProjectID, error)); ok {
		return rf(ctx, userIdpKey)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *[]db.UserReservedWithProjectID); ok {
		r0 = rf(ctx, userIdpKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]db.UserReservedWithProjectID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userIdpKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAccessToken provides a mock function with given fields: ctx, r
func (_m *IAMRepository) PutAccessToken(ctx context.Context, r *model.AccessToken) (*model.AccessToken, error) {
	ret := _m.Called(ctx, r)

	var r0 *model.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AccessToken) (*model.AccessToken, error)); ok {
		return rf(ctx, r)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AccessToken) *model.AccessToken); ok {
		r0 = rf(ctx, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AccessToken) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPolicy provides a mock function with given fields: ctx, p
func (_m *IAMRepository) PutPolicy(ctx context.Context, p *model.Policy) (*model.Policy, error) {
	ret := _m.Called(ctx, p)

	var r0 *model.Policy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Policy) (*model.Policy, error)); ok {
		return rf(ctx, p)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Policy) *model.Policy); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Policy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Policy) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRole provides a mock function with given fields: ctx, r
func (_m *IAMRepository) PutRole(ctx context.Context, r *model.Role) (*model.Role, error) {
	ret := _m.Called(ctx, r)

	var r0 *model.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Role) (*model.Role, error)); ok {
		return rf(ctx, r)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Role) *model.Role); ok {
		r0 = rf(ctx, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Role) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUser provides a mock function with given fields: ctx, u
func (_m *IAMRepository) PutUser(ctx context.Context, u *model.User) (*model.User, error) {
	ret := _m.Called(ctx, u)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) (*model.User, error)); ok {
		return rf(ctx, u)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *model.User); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUserReserved provides a mock function with given fields: ctx, u
func (_m *IAMRepository) PutUserReserved(ctx context.Context, u *model.UserReserved) (*model.UserReserved, error) {
	ret := _m.Called(ctx, u)

	var r0 *model.UserReserved
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserReserved) (*model.UserReserved, error)); ok {
		return rf(ctx, u)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserReserved) *model.UserReserved); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserReserved)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserReserved) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIAMRepository creates a new instance of IAMRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAMRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAMRepository {
	mock := &IAMRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
