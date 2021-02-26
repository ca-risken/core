package main

import (
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

const (
	defaultSortDirection string = "asc"
	defaultLimit         int32  = 2
)

type findingService struct {
	repository findingRepository
}

func newFindingService() finding.FindingServiceServer {
	return &findingService{
		repository: newFindingRepository(),
	}
}
