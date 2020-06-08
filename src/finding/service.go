package main

import (
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

type findingService struct {
	repository findingRepoInterface
}

func newFindingService(repo findingRepoInterface) finding.FindingServiceServer {
	return &findingService{
		repository: repo,
	}
}
