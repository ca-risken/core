package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/stretchr/testify/mock"
)

type mockFindingRepository struct {
	mock.Mock
}

func (m *mockFindingRepository) ListFinding(req *finding.ListFindingRequest) (*[]listFindingResult, error) {
	args := m.Called()
	return args.Get(0).(*[]listFindingResult), args.Error(1)
}

func newFindingMockRepository() findingRepoInterface {
	return &mockFindingRepository{}
}
func TestListFinding(t *testing.T) {
	var ctx context.Context
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)

	cases := []struct {
		name  string
		input *finding.ListFindingRequest
		want  *finding.ListFindingResponse
	}{
		{
			name:  "test1",
			input: &finding.ListFindingRequest{ProjectId: []uint32{123}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:  &finding.ListFindingResponse{FindingId: []uint64{111, 222}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("ListFinding").Return(
				&[]listFindingResult{
					{FindingID: 111},
					{FindingID: 222},
				}, nil)
			result, err := svc.ListFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
