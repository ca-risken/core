package finding

import (
	"context"
	"errors"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
)

func TestPutFindingBatch(t *testing.T) {
	var ctx context.Context
	type mockResp struct {
		GetFindingByDataSourceCall bool
		GetFindingByDataSourceResp *model.Finding
		GetFindingByDataSourceErr  error

		GetListFindingSettingCall bool
		GetListFindingSettingResp *[]model.FindingSetting
		GetListFindingSettingErr  error

		GetResourceByNameCall bool
		GetResourceByNameResp *model.Resource
		GetResourceByNameErr  error

		GetRecommendByDataSourceTypeCall bool
		GetRecommendByDataSourceTypeResp *model.Recommend
		GetRecommendByDataSourceTypeErr  error

		BulkUpsertFindingCall bool
		BulkUpsertFindingErr  error

		BulkUpsertResourceCall bool
		BulkUpsertResourceErr  error

		BulkUpsertRecommendCall bool
		BulkUpsertRecommendErr  error

		ListFindingTagByFindingIDCall bool
		ListFindingTagByFindingIDResp *[]model.FindingTag
		ListFindingTagByFindingIDErr  error

		ListResourceTagByResourceIDCall bool
		ListResourceTagByResourceIDResp *[]model.ResourceTag
		ListResourceTagByResourceIDErr  error

		BulkUpsertFindingTagCall bool
		BulkUpsertFindingTagErr  error

		BulkUpsertResourceTagCall bool
		BulkUpsertResourceTagErr  error

		BulkUpsertRecommendFindingCall bool
		BulkUpsertRecommendFindingErr  error
	}
	cases := []struct {
		name    string
		input   *finding.PutFindingBatchRequest
		wantErr bool
		mock    *mockResp
	}{
		{
			name: "OK",
			input: &finding.PutFindingBatchRequest{ProjectId: 1, Finding: []*finding.FindingBatchForUpsert{
				{
					Finding:   &finding.FindingForUpsert{ProjectId: 1, DataSource: "ds", DataSourceId: "1", ResourceName: "r", OriginalScore: 1.0, OriginalMaxScore: 1.0},
					Recommend: &finding.RecommendForBatch{Type: "type", Risk: "risk", Recommendation: "recommend"},
					Tag:       []*finding.FindingTagForBatch{{Tag: "tag1"}, {Tag: "tag2"}},
				},
			}},
			wantErr: false,
			mock: &mockResp{
				GetFindingByDataSourceCall:       true,
				GetFindingByDataSourceResp:       &model.Finding{FindingID: 1},
				GetListFindingSettingCall:        true,
				GetListFindingSettingResp:        &[]model.FindingSetting{},
				GetResourceByNameCall:            true,
				GetResourceByNameResp:            &model.Resource{ResourceID: 1},
				GetRecommendByDataSourceTypeCall: true,
				GetRecommendByDataSourceTypeResp: &model.Recommend{RecommendID: 1},

				BulkUpsertFindingCall:   true,
				BulkUpsertResourceCall:  true,
				BulkUpsertRecommendCall: true,

				ListFindingTagByFindingIDCall:   true,
				ListFindingTagByFindingIDResp:   &[]model.FindingTag{{FindingTagID: 1, Tag: "tag1"}, {FindingTagID: 2, Tag: "tag2"}},
				ListResourceTagByResourceIDCall: true,
				ListResourceTagByResourceIDResp: &[]model.ResourceTag{{ResourceTagID: 1, Tag: "tag1"}, {ResourceTagID: 2, Tag: "tag2"}},

				BulkUpsertRecommendFindingCall: true,
				BulkUpsertFindingTagCall:       true,
				BulkUpsertResourceTagCall:      true,
			},
		},
		{
			name: "OK/No recommend & tags",
			input: &finding.PutFindingBatchRequest{ProjectId: 1, Finding: []*finding.FindingBatchForUpsert{
				{
					Finding: &finding.FindingForUpsert{ProjectId: 1, DataSource: "ds", DataSourceId: "1", ResourceName: "r", OriginalScore: 1.0, OriginalMaxScore: 1.0},
				},
			}},
			wantErr: false,
			mock: &mockResp{
				GetFindingByDataSourceCall: true,
				GetFindingByDataSourceResp: &model.Finding{FindingID: 1},
				GetListFindingSettingCall:  true,
				GetListFindingSettingResp:  &[]model.FindingSetting{},
				GetResourceByNameCall:      true,
				GetResourceByNameResp:      &model.Resource{ResourceID: 1},

				BulkUpsertFindingCall:          true,
				BulkUpsertResourceCall:         true,
				BulkUpsertRecommendCall:        true,
				BulkUpsertRecommendFindingCall: true,
				BulkUpsertFindingTagCall:       true,
				BulkUpsertResourceTagCall:      true,

				ListFindingTagByFindingIDCall:   true,
				ListFindingTagByFindingIDResp:   &[]model.FindingTag{{FindingTagID: 1, Tag: "tag1"}, {FindingTagID: 2, Tag: "tag2"}},
				ListResourceTagByResourceIDCall: true,
				ListResourceTagByResourceIDResp: &[]model.ResourceTag{{ResourceTagID: 1, Tag: "tag1"}, {ResourceTagID: 2, Tag: "tag2"}},
			},
		},
		{
			name: "NG/Invalid request",
			input: &finding.PutFindingBatchRequest{ProjectId: 999, Finding: []*finding.FindingBatchForUpsert{
				{
					Finding: &finding.FindingForUpsert{ProjectId: 1, DataSource: "ds", DataSourceId: "1", ResourceName: "r", OriginalScore: 1.0, OriginalMaxScore: 1.0},
				},
			}},
			wantErr: true,
			mock:    &mockResp{},
		},
		{
			name: "NG/DB error",
			input: &finding.PutFindingBatchRequest{ProjectId: 1, Finding: []*finding.FindingBatchForUpsert{
				{
					Finding: &finding.FindingForUpsert{ProjectId: 1, DataSource: "ds", DataSourceId: "1", ResourceName: "r", OriginalScore: 1.0, OriginalMaxScore: 1.0},
				},
			}},
			wantErr: true,
			mock: &mockResp{
				GetFindingByDataSourceCall: true,
				GetFindingByDataSourceErr:  errors.New("DB error"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB, logger: logging.NewLogger()}

			if c.mock.GetFindingByDataSourceCall {
				mockDB.On("GetFindingByDataSource", test.RepeatMockAnything(4)...).Return(c.mock.GetFindingByDataSourceResp, c.mock.GetFindingByDataSourceErr).Once()
			}
			if c.mock.GetListFindingSettingCall {
				mockDB.On("ListFindingSetting", test.RepeatMockAnything(3)...).Return(c.mock.GetListFindingSettingResp, c.mock.GetListFindingSettingErr).Once()
			}
			if c.mock.GetResourceByNameCall {
				mockDB.On("GetResourceByName", test.RepeatMockAnything(3)...).Return(c.mock.GetResourceByNameResp, c.mock.GetResourceByNameErr).Once()
			}
			if c.mock.GetRecommendByDataSourceTypeCall {
				mockDB.On("GetRecommendByDataSourceType", test.RepeatMockAnything(3)...).Return(c.mock.GetRecommendByDataSourceTypeResp, c.mock.GetRecommendByDataSourceTypeErr).Once()
			}
			if c.mock.BulkUpsertFindingCall {
				mockDB.On("BulkUpsertFinding", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertFindingErr).Once()
			}
			if c.mock.BulkUpsertResourceCall {
				mockDB.On("BulkUpsertResource", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertResourceErr).Once()
			}
			if c.mock.BulkUpsertRecommendCall {
				mockDB.On("BulkUpsertRecommend", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertRecommendErr).Once()
			}
			if c.mock.ListFindingTagByFindingIDCall {
				mockDB.On("ListFindingTagByFindingID", test.RepeatMockAnything(3)...).Return(c.mock.ListFindingTagByFindingIDResp, c.mock.ListFindingTagByFindingIDErr)
			}
			if c.mock.ListResourceTagByResourceIDCall {
				mockDB.On("ListResourceTagByResourceID", test.RepeatMockAnything(3)...).Return(c.mock.ListResourceTagByResourceIDResp, c.mock.ListResourceTagByResourceIDErr)
			}
			if c.mock.BulkUpsertFindingTagCall {
				mockDB.On("BulkUpsertFindingTag", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertFindingTagErr).Once()
			}
			if c.mock.BulkUpsertResourceTagCall {
				mockDB.On("BulkUpsertResourceTag", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertResourceTagErr).Once()
			}
			if c.mock.BulkUpsertRecommendFindingCall {
				mockDB.On("BulkUpsertRecommendFinding", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertRecommendFindingErr).Once()
			}

			_, err := svc.PutFindingBatch(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestPutResourceBatch(t *testing.T) {
	var ctx context.Context
	type mockResp struct {
		GetResourceByNameCall bool
		GetResourceByNameResp *model.Resource
		GetResourceByNameErr  error

		BulkUpsertResourceCall bool
		BulkUpsertResourceErr  error

		ListResourceTagByResourceIDCall bool
		ListResourceTagByResourceIDResp *[]model.ResourceTag
		ListResourceTagByResourceIDErr  error

		BulkUpsertResourceTagCall bool
		BulkUpsertResourceTagErr  error
	}
	cases := []struct {
		name    string
		input   *finding.PutResourceBatchRequest
		wantErr bool
		mock    *mockResp
	}{
		{
			name: "OK",
			input: &finding.PutResourceBatchRequest{ProjectId: 1, Resource: []*finding.ResourceBatchForUpsert{
				{
					Resource: &finding.ResourceForUpsert{ProjectId: 1, ResourceName: "r"},
					Tag:      []*finding.ResourceTagForBatch{{Tag: "tag1"}, {Tag: "tag2"}},
				},
			}},
			wantErr: false,
			mock: &mockResp{
				GetResourceByNameCall:  true,
				GetResourceByNameResp:  &model.Resource{ResourceID: 1},
				BulkUpsertResourceCall: true,

				ListResourceTagByResourceIDCall: true,
				ListResourceTagByResourceIDResp: &[]model.ResourceTag{{ResourceTagID: 1, Tag: "tag1"}, {ResourceTagID: 2, Tag: "tag2"}},
				BulkUpsertResourceTagCall:       true,
			},
		},
		{
			name: "OK/No tags",
			input: &finding.PutResourceBatchRequest{ProjectId: 1, Resource: []*finding.ResourceBatchForUpsert{
				{
					Resource: &finding.ResourceForUpsert{ProjectId: 1, ResourceName: "r"},
				},
			}},
			wantErr: false,
			mock: &mockResp{
				GetResourceByNameCall: true,
				GetResourceByNameResp: &model.Resource{ResourceID: 1},

				BulkUpsertResourceCall:    true,
				BulkUpsertResourceTagCall: true,

				ListResourceTagByResourceIDCall: true,
				ListResourceTagByResourceIDResp: &[]model.ResourceTag{{ResourceTagID: 1, Tag: "tag1"}, {ResourceTagID: 2, Tag: "tag2"}},
			},
		},
		{
			name: "NG/Invalid request",
			input: &finding.PutResourceBatchRequest{ProjectId: 999, Resource: []*finding.ResourceBatchForUpsert{
				{
					Resource: &finding.ResourceForUpsert{ProjectId: 1, ResourceName: "r"},
				},
			}},
			wantErr: true,
			mock:    &mockResp{},
		},
		{
			name: "NG/DB error",
			input: &finding.PutResourceBatchRequest{ProjectId: 1, Resource: []*finding.ResourceBatchForUpsert{
				{
					Resource: &finding.ResourceForUpsert{ProjectId: 1, ResourceName: "r"},
				},
			}},
			wantErr: true,
			mock: &mockResp{
				GetResourceByNameCall: true,
				GetResourceByNameErr:  errors.New("DB error"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB, logger: logging.NewLogger()}

			if c.mock.GetResourceByNameCall {
				mockDB.On("GetResourceByName", test.RepeatMockAnything(3)...).Return(c.mock.GetResourceByNameResp, c.mock.GetResourceByNameErr).Once()
			}
			if c.mock.BulkUpsertResourceCall {
				mockDB.On("BulkUpsertResource", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertResourceErr).Once()
			}
			if c.mock.ListResourceTagByResourceIDCall {
				mockDB.On("ListResourceTagByResourceID", test.RepeatMockAnything(3)...).Return(c.mock.ListResourceTagByResourceIDResp, c.mock.ListResourceTagByResourceIDErr)
			}
			if c.mock.BulkUpsertResourceTagCall {
				mockDB.On("BulkUpsertResourceTag", test.RepeatMockAnything(2)...).Return(c.mock.BulkUpsertResourceTagErr).Once()
			}

			_, err := svc.PutResourceBatch(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
