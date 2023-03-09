package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
	"gorm.io/gorm"
)

func TestListFindingSetting(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *finding.ListFindingSettingRequest
		want         *finding.ListFindingSettingResponse
		mockResponce *[]model.FindingSetting
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &finding.ListFindingSettingRequest{ProjectId: 1, Status: finding.FindingSettingStatus_SETTING_ACTIVE},
			want: &finding.ListFindingSettingResponse{FindingSetting: []*finding.FindingSetting{
				{FindingSettingId: 1, ProjectId: 1, Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{FindingSettingId: 2, ProjectId: 1, Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.FindingSetting{
				{FindingSettingID: 1, ProjectID: 1, Setting: "{}", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
				{FindingSettingID: 2, ProjectID: 1, Setting: "{}", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK record not found",
			input:     &finding.ListFindingSettingRequest{ProjectId: 1},
			want:      &finding.ListFindingSettingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &finding.ListFindingSettingRequest{},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &finding.ListFindingSettingRequest{ProjectId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListFindingSetting", test.RepeatMockAnything(2)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListFindingSetting(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetFindingSetting(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *finding.GetFindingSettingRequest
		want         *finding.GetFindingSettingResponse
		mockResponce *model.FindingSetting
		mockError    error
		wantErr      bool
	}{
		{
			name:         "OK",
			input:        &finding.GetFindingSettingRequest{ProjectId: 1, FindingSettingId: 1},
			want:         &finding.GetFindingSettingResponse{FindingSetting: &finding.FindingSetting{FindingSettingId: 1, ProjectId: 1, Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.FindingSetting{FindingSettingID: 1, ProjectID: 1, Setting: "{}", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK record not found",
			input:     &finding.GetFindingSettingRequest{ProjectId: 1, FindingSettingId: 1},
			want:      &finding.GetFindingSettingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &finding.GetFindingSettingRequest{FindingSettingId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &finding.GetFindingSettingRequest{ProjectId: 1, FindingSettingId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetFindingSetting", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetFindingSetting(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutFindingSetting(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name     string
		input    *finding.PutFindingSettingRequest
		want     *finding.PutFindingSettingResponse
		wantErr  bool
		mockResp *model.FindingSetting
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.PutFindingSettingRequest{ProjectId: 1, FindingSetting: &finding.FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE}},
			want:     &finding.PutFindingSettingResponse{FindingSetting: &finding.FindingSetting{FindingSettingId: 1, ProjectId: 1, ResourceName: "rn", Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.FindingSetting{FindingSettingID: 1, ProjectID: 1, ResourceName: "rn", Setting: "{}", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutFindingSettingRequest{ProjectId: 999, FindingSetting: &finding.FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE}},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &finding.PutFindingSettingRequest{ProjectId: 1, FindingSetting: &finding.FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: "{}", Status: finding.FindingSettingStatus_SETTING_ACTIVE}},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("UpsertFindingSetting", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.PutFindingSetting(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteFindingSetting(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name     string
		input    *finding.DeleteFindingSettingRequest
		wantErr  bool
		mockCall bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.DeleteFindingSettingRequest{ProjectId: 1, FindingSettingId: 1},
			mockCall: true,
			wantErr:  false,
		},
		{
			name:     "NG validation error",
			input:    &finding.DeleteFindingSettingRequest{ProjectId: 1},
			mockCall: false,
			wantErr:  true,
		},
		{
			name:     "Invalid DB error",
			input:    &finding.DeleteFindingSettingRequest{ProjectId: 1, FindingSettingId: 1},
			mockCall: true,
			wantErr:  true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockCall {
				mockDB.On("DeleteFindingSetting", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteFindingSetting(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
