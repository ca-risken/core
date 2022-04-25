package finding_client

import (
	"context"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/finding/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	sampleProjectID uint32 = 1

	sampleFinding1 = &finding.FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "1", ResourceName: "1", ProjectId: sampleProjectID, OriginalScore: 0.1, OriginalMaxScore: 1.0}
	sampleFinding2 = &finding.FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "2", ResourceName: "2", ProjectId: sampleProjectID, OriginalScore: 0.2, OriginalMaxScore: 1.0}
	sampleFinding3 = &finding.FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "3", ResourceName: "3", ProjectId: sampleProjectID, OriginalScore: 0.3, OriginalMaxScore: 1.0}

	sampleResource1 = &finding.ResourceForUpsert{ResourceName: "1", ProjectId: sampleProjectID}
	sampleResource2 = &finding.ResourceForUpsert{ResourceName: "2", ProjectId: sampleProjectID}
	sampleResource3 = &finding.ResourceForUpsert{ResourceName: "3", ProjectId: sampleProjectID}
)

func TestPutFindingBatch(t *testing.T) {
	type mockResponse struct {
		PutFindingBatchResp *emptypb.Empty
		PutFindingBatchErr  error
	}

	cases := []struct {
		name     string
		input    []*finding.FindingBatchForUpsert
		mockResp mockResponse
		wantErr  bool
	}{
		{
			name: "OK",
			input: []*finding.FindingBatchForUpsert{
				{Finding: sampleFinding1},
				{Finding: sampleFinding2},
				{Finding: sampleFinding3},
			},
			mockResp: mockResponse{
				PutFindingBatchResp: &emptypb.Empty{},
				PutFindingBatchErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "OK(limit over)",
			input: []*finding.FindingBatchForUpsert{
				{Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, // 10
				{Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, // 20
				{Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, // 30
				{Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, // 40
				{Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, {Finding: sampleFinding1}, // 50
				{Finding: sampleFinding1}, // 51
			},
			mockResp: mockResponse{
				PutFindingBatchResp: &emptypb.Empty{},
				PutFindingBatchErr:  nil,
			},
			wantErr: false,
		},
		{
			name:  "OK(no data)",
			input: []*finding.FindingBatchForUpsert{},
			mockResp: mockResponse{
				PutFindingBatchResp: &emptypb.Empty{},
				PutFindingBatchErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "NG(something error)",
			input: []*finding.FindingBatchForUpsert{
				{Finding: sampleFinding1},
			},
			mockResp: mockResponse{
				PutFindingBatchResp: nil,
				PutFindingBatchErr:  errors.New("something wrong"),
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := &mocks.FindingServiceClient{}
			if c.mockResp.PutFindingBatchResp != nil || c.mockResp.PutFindingBatchErr != nil {
				mockClient.On("PutFindingBatch", mock.Anything, mock.Anything).Return(
					c.mockResp.PutFindingBatchResp, c.mockResp.PutFindingBatchErr)
			}

			// exec
			ctx := context.Background()
			err := PutFindingBatch(ctx, mockClient, sampleProjectID, c.input)
			if err == nil && c.wantErr {
				t.Fatalf("Unexpected no error: wantErr=%t", c.wantErr)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestPutResourceBatch(t *testing.T) {
	type mockResponse struct {
		PutResourceBatchResp *emptypb.Empty
		PutResourceBatchErr  error
	}

	cases := []struct {
		name     string
		input    []*finding.ResourceBatchForUpsert
		mockResp mockResponse
		wantErr  bool
	}{
		{
			name: "OK",
			input: []*finding.ResourceBatchForUpsert{
				{Resource: sampleResource1},
				{Resource: sampleResource2},
				{Resource: sampleResource3},
			},
			mockResp: mockResponse{
				PutResourceBatchResp: &emptypb.Empty{},
				PutResourceBatchErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "OK(limit over)",
			input: []*finding.ResourceBatchForUpsert{
				{Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, // 10
				{Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, // 20
				{Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, // 30
				{Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, // 40
				{Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, {Resource: sampleResource1}, // 50
				{Resource: sampleResource1}, // 51
			},
			mockResp: mockResponse{
				PutResourceBatchResp: &emptypb.Empty{},
				PutResourceBatchErr:  nil,
			},
			wantErr: false,
		},
		{
			name:  "OK(no data)",
			input: []*finding.ResourceBatchForUpsert{},
			mockResp: mockResponse{
				PutResourceBatchResp: &emptypb.Empty{},
				PutResourceBatchErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "NG(something error)",
			input: []*finding.ResourceBatchForUpsert{
				{Resource: sampleResource1},
			},
			mockResp: mockResponse{
				PutResourceBatchResp: nil,
				PutResourceBatchErr:  errors.New("something wrong"),
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := &mocks.FindingServiceClient{}
			if c.mockResp.PutResourceBatchResp != nil || c.mockResp.PutResourceBatchErr != nil {
				mockClient.On("PutResourceBatch", mock.Anything, mock.Anything).Return(
					c.mockResp.PutResourceBatchResp, c.mockResp.PutResourceBatchErr)
			}

			// exec
			ctx := context.Background()
			err := PutResourceBatch(ctx, mockClient, sampleProjectID, c.input)
			if err == nil && c.wantErr {
				t.Fatalf("Unexpected no error: wantErr=%t", c.wantErr)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
