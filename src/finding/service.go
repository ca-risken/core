package main

import (
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-core/proto/iam"
)

type findingService struct {
	iamClient  iam.IAMServiceClient
	repository findingRepoInterface
}

func newFindingService(repo findingRepoInterface, client iam.IAMServiceClient) finding.FindingServiceServer {
	return &findingService{
		iamClient:  client,
		repository: repo,
	}
}
