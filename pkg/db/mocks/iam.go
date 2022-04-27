package mocks

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
	"github.com/stretchr/testify/mock"
)

/*
 * Mock Repository
 */
type MockIAMRepository struct {
	mock.Mock
}

func (m *MockIAMRepository) ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32, admin bool) (*[]model.User, error) {
	args := m.Called()
	return args.Get(0).(*[]model.User), args.Error(1)
}
func (m *MockIAMRepository) GetUser(context.Context, uint32, string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockIAMRepository) GetUserBySub(context.Context, string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockIAMRepository) GetUserPolicy(context.Context, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *MockIAMRepository) GetTokenPolicy(context.Context, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *MockIAMRepository) PutUser(context.Context, *model.User) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockIAMRepository) GetActiveUserCount(ctx context.Context) (*int, error) {
	args := m.Called()
	return args.Get(0).(*int), args.Error(1)
}
func (m *MockIAMRepository) ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Role), args.Error(1)
}
func (m *MockIAMRepository) GetRole(context.Context, uint32, uint32) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *MockIAMRepository) GetRoleByName(context.Context, uint32, string) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *MockIAMRepository) PutRole(ctx context.Context, r *model.Role) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *MockIAMRepository) DeleteRole(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) AttachRole(context.Context, uint32, uint32, uint32) (*model.UserRole, error) {
	args := m.Called()
	return args.Get(0).(*model.UserRole), args.Error(1)
}
func (m *MockIAMRepository) AttachAllAdminRole(context.Context, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) DetachRole(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) ListPolicy(context.Context, uint32, string, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *MockIAMRepository) GetPolicy(context.Context, uint32, uint32) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *MockIAMRepository) GetPolicyByName(context.Context, uint32, string) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *MockIAMRepository) PutPolicy(context.Context, *model.Policy) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *MockIAMRepository) DeletePolicy(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) AttachPolicy(context.Context, uint32, uint32, uint32) (*model.RolePolicy, error) {
	args := m.Called()
	return args.Get(0).(*model.RolePolicy), args.Error(1)
}
func (m *MockIAMRepository) DetachPolicy(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) GetAdminPolicy(context.Context, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *MockIAMRepository) ListAccessToken(ctx context.Context, projectID uint32, name string, accessTokenID uint32) (*[]model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AccessToken), args.Error(1)
}
func (m *MockIAMRepository) GetAccessTokenByID(ctx context.Context, projectID, accessTokenID uint32) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *MockIAMRepository) GetAccessTokenByUniqueKey(ctx context.Context, projectID uint32, name string) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *MockIAMRepository) GetActiveAccessTokenHash(ctx context.Context, projectID, accessTokenID uint32, tokenHash string) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *MockIAMRepository) PutAccessToken(ctx context.Context, r *model.AccessToken) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *MockIAMRepository) DeleteAccessToken(ctx context.Context, projectID, accessTokenID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) AttachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) (*model.AccessTokenRole, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessTokenRole), args.Error(1)
}
func (m *MockIAMRepository) GetAccessTokenRole(ctx context.Context, accessTokenID, roleID uint32) (*model.AccessTokenRole, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessTokenRole), args.Error(1)
}
func (m *MockIAMRepository) DetachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockIAMRepository) ExistsAccessTokenMaintainer(ctx context.Context, projectID, accessTokenID uint32) (bool, error) {
	args := m.Called()
	return args.Get(0).(bool), args.Error(1)
}
func (m *MockIAMRepository) ListExpiredAccessToken(ctx context.Context) (*[]model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AccessToken), args.Error(1)
}
