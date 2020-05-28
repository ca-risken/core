package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

type findingService struct{}

func (f *findingService) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	ids := []string{"0000000001", "0000000002", "0000000003"}
	return &finding.ListFindingResponse{ProjectIds: ids}, nil
}

func (f *findingService) GetFinding(ctx context.Context, req *finding.GetFindingRequest) (*finding.GetFindingResponse, error) {
	data := finding.Finding{
		FidingId:   "xxx",
		Name:       "xxx",
		DataSource: "aws:guardduty",
		Resource:   "aws:xxx:xxx:::aaa",
		ProjectId:  "1234567890",
		Data:       "aaaaa",
		CreatedAt:  1590598478,
		UpdatedAt:  1590598478,
	}
	return &finding.GetFindingResponse{Data: &data}, nil
}
