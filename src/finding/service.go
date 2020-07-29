package main

import (
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

type findingService struct {
	repository findingRepository
}

func newFindingService() finding.FindingServiceServer {
	return &findingService{
		repository: newFindingRepository(),
	}
}
