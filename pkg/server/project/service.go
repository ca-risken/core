package project

import (
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/iam"
)

type ProjectService struct {
	repository db.ProjectRepository
	iamClient  iam.IAMServiceClient
}

func NewProjectService(repository db.ProjectRepository, iamClient iam.IAMServiceClient) *ProjectService {
	return &ProjectService{
		repository: repository,
		iamClient:  iamClient,
	}
}
