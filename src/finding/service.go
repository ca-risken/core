package main

import (
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

type findingService struct {
	repository findingRepoInterface
}

func newFindingService() finding.FindingServiceServer {
	return &findingService{
		repository: newFindingRepository(),
	}
}
