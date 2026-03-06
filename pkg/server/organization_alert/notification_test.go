package organization_alert

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/organization_alert"
	"github.com/jarcoal/httpmock"
	"gorm.io/gorm"
)

func TestListOrganizationNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *organization_alert.ListOrganizationNotificationRequest
		want     *organization_alert.ListOrganizationNotificationResponse
		wantErr  bool
		mockResp []*model.OrganizationNotification
		mockErr  error
	}{
		{
			name:  "OK",
			input: &organization_alert.ListOrganizationNotificationRequest{OrganizationId: 1, Type: ""},
			want: &organization_alert.ListOrganizationNotificationResponse{
				OrganizationNotification: []*organization_alert.OrganizationNotification{
					{NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack", NotifySetting: `{"channel_id":"ch1"}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResp: []*model.OrganizationNotification{
				{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"channel_id":"ch1"}`, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:    "OK - empty",
			input:   &organization_alert.ListOrganizationNotificationRequest{OrganizationId: 1},
			want:    &organization_alert.ListOrganizationNotificationResponse{},
			mockResp: []*model.OrganizationNotification{},
		},
		{
			name:    "OK - record not found",
			input:   &organization_alert.ListOrganizationNotificationRequest{OrganizationId: 1},
			want:    &organization_alert.ListOrganizationNotificationResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG - DB error",
			input:   &organization_alert.ListOrganizationNotificationRequest{OrganizationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &organization_alert.ListOrganizationNotificationRequest{OrganizationId: 0},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrganizationAlertRepository(t)
			svc := OrganizationAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListOrganizationNotification", test.RepeatMockAnything(3)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListOrganizationNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrganizationNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *organization_alert.GetOrganizationNotificationRequest
		want     *organization_alert.GetOrganizationNotificationResponse
		wantErr  bool
		mockResp *model.OrganizationNotification
		mockErr  error
	}{
		{
			name:  "OK",
			input: &organization_alert.GetOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 1},
			want: &organization_alert.GetOrganizationNotificationResponse{
				OrganizationNotification: &organization_alert.OrganizationNotification{
					NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack",
					NotifySetting: `{"channel_id":"ch1"}`,
					CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockResp: &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"channel_id":"ch1"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "OK - record not found",
			input:   &organization_alert.GetOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 999},
			want:    &organization_alert.GetOrganizationNotificationResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG - DB error",
			input:   &organization_alert.GetOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &organization_alert.GetOrganizationNotificationRequest{OrganizationId: 0, NotificationId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrganizationAlertRepository(t)
			svc := OrganizationAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("GetOrganizationNotification", test.RepeatMockAnything(3)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.GetOrganizationNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrganizationNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *organization_alert.PutOrganizationNotificationRequest
		want        *organization_alert.PutOrganizationNotificationResponse
		wantErr     bool
		mockGetResp *model.OrganizationNotification
		mockGetErr  error
		mockUpResp  *model.OrganizationNotification
		mockUpErr   error
	}{
		{
			name: "OK Insert",
			input: &organization_alert.PutOrganizationNotificationRequest{
				OrganizationId: 1,
				OrganizationNotification: &organization_alert.OrganizationNotificationForUpsert{
					OrganizationId: 1, Name: "notif1", Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`,
				},
			},
			want: &organization_alert.PutOrganizationNotificationResponse{
				OrganizationNotification: &organization_alert.OrganizationNotification{
					NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack",
					NotifySetting: `{"webhook_url":"https://e**********","channel_id":"","data":{},"locale":""}`,
					CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockUpResp: &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "OK Update",
			input: &organization_alert.PutOrganizationNotificationRequest{
				OrganizationId: 1,
				OrganizationNotification: &organization_alert.OrganizationNotificationForUpsert{
					NotificationId: 1, OrganizationId: 1, Name: "notif1", Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`,
				},
			},
			want: &organization_alert.PutOrganizationNotificationResponse{
				OrganizationNotification: &organization_alert.OrganizationNotification{
					NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack",
					NotifySetting: `{"webhook_url":"https://e**********","channel_id":"","data":{},"locale":""}`,
					CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockGetResp: &model.OrganizationNotification{NotificationID: 1, OrganizationID: 1, Name: "notif1", Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`, CreatedAt: now, UpdatedAt: now},
			mockUpResp:  &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "NG - record not found on update",
			input: &organization_alert.PutOrganizationNotificationRequest{
				OrganizationId: 1,
				OrganizationNotification: &organization_alert.OrganizationNotificationForUpsert{
					NotificationId: 1, OrganizationId: 1, Name: "notif1", Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`,
				},
			},
			want:       &organization_alert.PutOrganizationNotificationResponse{},
			mockGetErr: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrganizationAlertRepository(t)
			svc := OrganizationAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetOrganizationNotification", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertOrganizationNotification", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutOrganizationNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrganizationNotification(t *testing.T) {
	cases := []struct {
		name    string
		input   *organization_alert.DeleteOrganizationNotificationRequest
		wantErr bool
		mockErr error
	}{
		{
			name:  "OK",
			input: &organization_alert.DeleteOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 1},
		},
		{
			name:    "NG - DB error",
			input:   &organization_alert.DeleteOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &organization_alert.DeleteOrganizationNotificationRequest{OrganizationId: 0, NotificationId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrganizationAlertRepository(t)
			svc := OrganizationAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockErr != nil || c.input.OrganizationId != 0 && !c.wantErr {
				mockDB.On("DeleteOrganizationNotification", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrganizationNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestTestOrganizationNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://example.com", httpmock.NewStringResponder(200, "ok"))

	cases := []struct {
		name     string
		input    *organization_alert.TestOrganizationNotificationRequest
		wantErr  bool
		mockResp *model.OrganizationNotification
		mockErr  error
	}{
		{
			name:     "OK - webhook",
			input:    &organization_alert.TestOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 1},
			mockResp: &model.OrganizationNotification{NotificationID: 1, OrganizationID: 1, Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`},
		},
		{
			name:    "OK - record not found",
			input:   &organization_alert.TestOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 999},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG - DB error",
			input:   &organization_alert.TestOrganizationNotificationRequest{OrganizationId: 1, NotificationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &organization_alert.TestOrganizationNotificationRequest{OrganizationId: 0, NotificationId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrganizationAlertRepository(t)
			svc := OrganizationAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("GetOrganizationNotification", test.RepeatMockAnything(3)...).Return(c.mockResp, c.mockErr).Once()
			}
			_, err := svc.TestOrganizationNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
