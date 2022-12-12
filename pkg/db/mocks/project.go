package mocks

import (
	"context"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/model"
	"github.com/stretchr/testify/mock"
)

/**
 * Mock Repository
**/
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) ListProject(ctx context.Context, userID, projectID uint32, name string) (*[]db.ProjectWithTag, error) {
	args := m.Called()
	return args.Get(0).(*[]db.ProjectWithTag), args.Error(1)
}
func (m *MockProjectRepository) CreateProject(ctx context.Context, name string) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}
func (m *MockProjectRepository) UpdateProject(ctx context.Context, projectID uint32, name string) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}
func (m *MockProjectRepository) ListProjectTag(ctx context.Context, projectID uint32) (*[]model.ProjectTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ProjectTag), args.Error(1)
}
func (m *MockProjectRepository) TagProject(ctx context.Context, projectID uint32, tag, color string) (*model.ProjectTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ProjectTag), args.Error(1)
}
func (m *MockProjectRepository) UntagProject(ctx context.Context, projectID uint32, tag string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockProjectRepository) DeleteProject(ctx context.Context, projectID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockProjectRepository) CleanWithNoProject(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}
