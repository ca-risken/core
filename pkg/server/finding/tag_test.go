package finding

import (
	"context"
	"errors"
	"testing"

	dbmocks "github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
)

func TestUntagByResourceName(t *testing.T) {
	type MockListFinding struct {
		Resp *[]model.Finding
		Err  error
	}
	type GetFindingTagByKey struct {
		Resp *model.FindingTag
		Err  error
	}
	type MockUntagFinding struct {
		Err error
	}
	type MockGetResourceByName struct {
		Resp *model.Resource
		Err  error
	}
	type MockGetResourceTagByKey struct {
		Resp *model.ResourceTag
		Err  error
	}
	type MockUntagResource struct {
		Err error
	}

	cases := []struct {
		name    string
		input   *finding.UntagByResourceNameRequest
		wantErr bool

		mockListFinding         *MockListFinding
		mockGetFindingTagByKey  *GetFindingTagByKey
		mockUntagFinding        *MockUntagFinding
		mockGetResourceByName   *MockGetResourceByName
		mockGetResourceTagByKey *MockGetResourceTagByKey
		mockUntagResource       *MockUntagResource
	}{
		{
			name:  "OK",
			input: &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			mockListFinding: &MockListFinding{
				Resp: &[]model.Finding{{FindingID: 1}},
			},
			mockGetFindingTagByKey: &GetFindingTagByKey{
				Resp: &model.FindingTag{FindingTagID: 1},
			},
			mockUntagFinding: &MockUntagFinding{
				Err: nil,
			},
			mockGetResourceByName: &MockGetResourceByName{
				Resp: &model.Resource{ResourceID: 1},
			},
			mockGetResourceTagByKey: &MockGetResourceTagByKey{
				Resp: &model.ResourceTag{ResourceTagID: 1},
			},
			mockUntagResource: &MockUntagResource{
				Err: nil,
			},
		},
		{
			name:    "NG Invalid param",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name"},
			wantErr: true,
		},
		{
			name:    "NG DB error(GetFindingByName)",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			wantErr: true,
			mockListFinding: &MockListFinding{
				Err: errors.New("something error"),
			},
		},
		{
			name:    "NG DB error(GetFindingTagByKey)",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			wantErr: true,
			mockListFinding: &MockListFinding{
				Resp: &[]model.Finding{{FindingID: 1}},
			},
			mockGetFindingTagByKey: &GetFindingTagByKey{
				Err: errors.New("something error"),
			},
		},
		{
			name:    "NG DB error(UntagFinding)",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			wantErr: true,
			mockListFinding: &MockListFinding{
				Resp: &[]model.Finding{{FindingID: 1}},
			},
			mockGetFindingTagByKey: &GetFindingTagByKey{
				Resp: &model.FindingTag{FindingTagID: 1},
			},
			mockUntagFinding: &MockUntagFinding{
				Err: errors.New("something error"),
			},
		},
		{
			name:    "NG DB error(GetResourceByName)",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			wantErr: true,
			mockListFinding: &MockListFinding{
				Resp: &[]model.Finding{{FindingID: 1}},
			},
			mockGetFindingTagByKey: &GetFindingTagByKey{
				Resp: &model.FindingTag{FindingTagID: 1},
			},
			mockUntagFinding: &MockUntagFinding{
				Err: nil,
			},
			mockGetResourceByName: &MockGetResourceByName{
				Err: errors.New("something error"),
			},
		},
		{
			name:    "NG DB error(GetResourceTagByKey)",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			wantErr: true,
			mockListFinding: &MockListFinding{
				Resp: &[]model.Finding{{FindingID: 1}},
			},
			mockGetFindingTagByKey: &GetFindingTagByKey{
				Resp: &model.FindingTag{FindingTagID: 1},
			},
			mockUntagFinding: &MockUntagFinding{
				Err: nil,
			},
			mockGetResourceByName: &MockGetResourceByName{
				Resp: &model.Resource{ResourceID: 1},
			},
			mockGetResourceTagByKey: &MockGetResourceTagByKey{
				Err: errors.New("something error"),
			},
		},
		{
			name:    "NG DB error(UntagResource)",
			input:   &finding.UntagByResourceNameRequest{ProjectId: 1, ResourceName: "name", Tag: "tag"},
			wantErr: true,
			mockListFinding: &MockListFinding{
				Resp: &[]model.Finding{{FindingID: 1}},
			},
			mockGetFindingTagByKey: &GetFindingTagByKey{
				Resp: &model.FindingTag{FindingTagID: 1},
			},
			mockUntagFinding: &MockUntagFinding{
				Err: nil,
			},
			mockGetResourceByName: &MockGetResourceByName{
				Resp: &model.Resource{ResourceID: 1},
			},
			mockGetResourceTagByKey: &MockGetResourceTagByKey{
				Resp: &model.ResourceTag{ResourceTagID: 1},
			},
			mockUntagResource: &MockUntagResource{
				Err: errors.New("something error"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := dbmocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockListFinding != nil {
				mockDB.
					On("ListFinding", test.RepeatMockAnything(15)...).
					Return(c.mockListFinding.Resp, c.mockListFinding.Err).
					Once()
			}
			if c.mockGetFindingTagByKey != nil {
				mockDB.
					On("GetFindingTagByKey", test.RepeatMockAnything(4)...).
					Return(c.mockGetFindingTagByKey.Resp, c.mockGetFindingTagByKey.Err).
					Once()
			}
			if c.mockUntagFinding != nil {
				mockDB.
					On("UntagFinding", test.RepeatMockAnything(3)...).
					Return(c.mockUntagFinding.Err).
					Once()
			}
			if c.mockGetResourceByName != nil {
				mockDB.
					On("GetResourceByName", test.RepeatMockAnything(3)...).
					Return(c.mockGetResourceByName.Resp, c.mockGetResourceByName.Err).
					Once()
			}
			if c.mockGetResourceTagByKey != nil {
				mockDB.
					On("GetResourceTagByKey", test.RepeatMockAnything(4)...).
					Return(c.mockGetResourceTagByKey.Resp, c.mockGetResourceTagByKey.Err).
					Once()
			}
			if c.mockUntagResource != nil {
				mockDB.
					On("UntagResource", test.RepeatMockAnything(3)...).
					Return(c.mockUntagResource.Err).
					Once()
			}
			_, err := svc.UntagByResourceName(context.TODO(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatalf("Expected error but got nil")
			}
		})
	}
}
