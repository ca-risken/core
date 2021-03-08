package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/stretchr/testify/mock"
)

/**
 * Mock Repository
**/
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) ListProject(userID, projectID uint32, name string) (*[]projectWithTag, error) {
	args := m.Called()
	return args.Get(0).(*[]projectWithTag), args.Error(1)
}
func (m *mockRepository) CreateProject(name string) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}
func (m *mockRepository) UpdateProject(projectID uint32, name string) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}
func (m *mockRepository) ListProjectTag(projectID uint32) (*[]model.ProjectTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ProjectTag), args.Error(1)
}
func (m *mockRepository) TagProject(projectID uint32, tag string) (*model.ProjectTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ProjectTag), args.Error(1)
}
func (m *mockRepository) UntagProject(projectID uint32, tag string) error {
	args := m.Called()
	return args.Error(0)
}

/**
 * Mock GRPC Client
**/
type mockClient struct {
	mock.Mock
}

func (m *mockClient) CreateDefaultRole(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockClient) DeleteAllProjectRole(context.Context, uint32) error {
	args := m.Called()
	return args.Error(0)
}
