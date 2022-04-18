package main

import (
	"context"
	"errors"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/src/finding/model"
)

func TestPutFindingBatch(t *testing.T) {
	var ctx context.Context
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	type mockResp struct {
		GetFindingByDataSourceCall bool
		GetFindingByDataSourceResp *model.Finding
		GetFindingByDataSourceErr  error

		GetFindingSettingByResourceCall bool
		GetFindingSettingByResourceResp *model.FindingSetting
		GetFindingSettingByResourceErr  error

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

		GetFindingTagByKeyCall bool
		GetFindingTagByKeyResp *model.FindingTag
		GetFindingTagByKeyErr  error

		GetResourceTagByKeyCall bool
		GetResourceTagByKeyResp *model.ResourceTag
		GetResourceTagByKeyErr  error

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
				GetFindingSettingByResourceCall:  true,
				GetFindingSettingByResourceResp:  &model.FindingSetting{},
				GetResourceByNameCall:            true,
				GetResourceByNameResp:            &model.Resource{ResourceID: 1},
				GetRecommendByDataSourceTypeCall: true,
				GetRecommendByDataSourceTypeResp: &model.Recommend{RecommendID: 1},

				BulkUpsertFindingCall:   true,
				BulkUpsertResourceCall:  true,
				BulkUpsertRecommendCall: true,

				GetFindingTagByKeyCall:  true,
				GetFindingTagByKeyResp:  &model.FindingTag{FindingTagID: 1},
				GetResourceTagByKeyCall: true,
				GetResourceTagByKeyResp: &model.ResourceTag{ResourceTagID: 1},

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
				GetFindingByDataSourceCall:      true,
				GetFindingByDataSourceResp:      &model.Finding{FindingID: 1},
				GetFindingSettingByResourceCall: true,
				GetFindingSettingByResourceResp: &model.FindingSetting{},
				GetResourceByNameCall:           true,
				GetResourceByNameResp:           &model.Resource{ResourceID: 1},

				BulkUpsertFindingCall:          true,
				BulkUpsertResourceCall:         true,
				BulkUpsertRecommendCall:        true,
				BulkUpsertRecommendFindingCall: true,
				BulkUpsertFindingTagCall:       true,
				BulkUpsertResourceTagCall:      true,
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
			if c.mock.GetFindingByDataSourceCall {
				mockDB.On("GetFindingByDataSource").Return(c.mock.GetFindingByDataSourceResp, c.mock.GetFindingByDataSourceErr).Once()
			}
			if c.mock.GetFindingSettingByResourceCall {
				mockDB.On("GetFindingSettingByResource").Return(c.mock.GetFindingSettingByResourceResp, c.mock.GetFindingSettingByResourceErr).Once()
			}
			if c.mock.GetResourceByNameCall {
				mockDB.On("GetResourceByName").Return(c.mock.GetResourceByNameResp, c.mock.GetResourceByNameErr).Once()
			}
			if c.mock.GetRecommendByDataSourceTypeCall {
				mockDB.On("GetRecommendByDataSourceType").Return(c.mock.GetRecommendByDataSourceTypeResp, c.mock.GetRecommendByDataSourceTypeErr).Once()
			}
			if c.mock.BulkUpsertFindingCall {
				mockDB.On("BulkUpsertFinding").Return(c.mock.BulkUpsertFindingErr).Once()
			}
			if c.mock.BulkUpsertResourceCall {
				mockDB.On("BulkUpsertResource").Return(c.mock.BulkUpsertResourceErr).Once()
			}
			if c.mock.BulkUpsertRecommendCall {
				mockDB.On("BulkUpsertRecommend").Return(c.mock.BulkUpsertRecommendErr).Once()
			}
			if c.mock.GetFindingTagByKeyCall {
				mockDB.On("GetFindingTagByKey").Return(c.mock.GetFindingTagByKeyResp, c.mock.GetFindingTagByKeyErr)
			}
			if c.mock.GetResourceTagByKeyCall {
				mockDB.On("GetResourceTagByKey").Return(c.mock.GetResourceTagByKeyResp, c.mock.GetResourceTagByKeyErr)
			}
			if c.mock.BulkUpsertFindingTagCall {
				mockDB.On("BulkUpsertFindingTag").Return(c.mock.BulkUpsertFindingTagErr).Once()
			}
			if c.mock.BulkUpsertResourceTagCall {
				mockDB.On("BulkUpsertResourceTag").Return(c.mock.BulkUpsertResourceTagErr).Once()
			}
			if c.mock.BulkUpsertRecommendFindingCall {
				mockDB.On("BulkUpsertRecommendFinding").Return(c.mock.BulkUpsertRecommendFindingErr).Once()
			}

			_, err := svc.PutFindingBatch(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
