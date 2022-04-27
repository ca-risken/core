package iam

import (
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/finding"
)

type IAMService struct {
	repository    db.IAMRepository
	findingClient finding.FindingServiceClient
}

func NewIAMService(repository db.IAMRepository, findingClient finding.FindingServiceClient) *IAMService {
	return &IAMService{
		repository:    repository,
		findingClient: findingClient,
	}
}
