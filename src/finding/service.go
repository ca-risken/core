package main

import (
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

const (
	defaultSortDirection string = "asc"
	defaultLimit         int32  = 200
)

type findingService struct {
	repository findingRepository
}

func newFindingService() finding.FindingServiceServer {
	return &findingService{
		repository: newFindingRepository(),
	}
}

// TODO remove after fix response interface type change to int64 from uint32
func convertToUint32(v int64) uint32 {
	return uint32(v)
}
