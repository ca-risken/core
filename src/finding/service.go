package main

import (
	"context"

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

func (f *findingService) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	ids, err := f.repository.List()
	if err != nil {
		return nil, err
	}
	return &finding.ListFindingResponse{FindingId: *ids}, nil
}

func (f *findingService) GetFinding(ctx context.Context, req *finding.GetFindingRequest) (*finding.GetFindingResponse, error) {
	return &finding.GetFindingResponse{Finding: &finding.Finding{
		FidingId:     1234567890,
		Description:  "xxx",
		DataSource:   "aws:guardduty",
		ResourceName: "aws:xxx:xxx:::aaa",
		Score:        0.6,
		ProjectId:    "1234567890",
		Data:         `{"data": {"key": "value"}}`,
		CreatedAt:    1590598478,
		UpdatedAt:    1590598478,
	}}, nil
}

func (f *findingService) PutFinding(ctx context.Context, req *finding.PutFindingRequest) (*finding.PutFindingResponse, error) {
	return nil, nil
}

func (f *findingService) DeleteFinding(ctx context.Context, req *finding.DeleteFindingRequest) (*finding.Empty, error) {
	return nil, nil
}

func (f *findingService) ListFindingTag(ctx context.Context, req *finding.ListFindingTagRequest) (*finding.ListFindingTagResponse, error) {
	return nil, nil
}

func (f *findingService) TagFinding(ctx context.Context, req *finding.TagFindingRequest) (*finding.TagFindingResponse, error) {
	return nil, nil
}

func (f *findingService) UntagFinding(ctx context.Context, req *finding.UntagFindingRequest) (*finding.Empty, error) {
	return nil, nil
}

func (f *findingService) ListResource(ctx context.Context, req *finding.ListResourceRequest) (*finding.ListResourceResponse, error) {
	return nil, nil
}

func (f *findingService) GetResource(ctx context.Context, req *finding.GetResourceRequest) (*finding.GetResourceResponse, error) {
	return nil, nil
}

func (f *findingService) PutResource(ctx context.Context, req *finding.PutResourceRequest) (*finding.PutResourceResponse, error) {
	return nil, nil
}

func (f *findingService) DeleteResource(ctx context.Context, req *finding.DeleteResourceRequest) (*finding.Empty, error) {
	return nil, nil
}

func (f *findingService) ListResourceTag(ctx context.Context, req *finding.ListResourceTagRequest) (*finding.ListResourceTagResponse, error) {
	return nil, nil
}

func (f *findingService) TagResource(ctx context.Context, req *finding.TagResourceRequest) (*finding.TagResourceResponse, error) {
	return nil, nil
}

func (f *findingService) UntagResource(ctx context.Context, req *finding.UntagResourceRequest) (*finding.Empty, error) {
	return nil, nil
}
