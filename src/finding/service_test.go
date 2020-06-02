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

func (m *mockFindingRepository) List() (*[]string, error) {
	args := m.Called()
	return args.Get(0).(*[]string), args.Error(1)
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
			input: &finding.ListFindingRequest{ProjectId: []string{"hoge"}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:  &finding.ListFindingResponse{FindingId: []string{"aaa", "bbb"}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("List").Return(&[]string{"aaa", "bbb"}, nil)
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
