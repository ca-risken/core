package org_alert

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
	"github.com/ca-risken/core/proto/org_alert"
	"gorm.io/gorm"
)

func TestListOrgNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *org_alert.ListOrgNotificationRequest
		want     *org_alert.ListOrgNotificationResponse
		wantErr  bool
		mockResp []*model.OrganizationNotification
		mockErr  error
	}{
		{
			name:  "OK",
			input: &org_alert.ListOrgNotificationRequest{OrganizationId: 1},
			want: &org_alert.ListOrgNotificationResponse{
				OrganizationNotification: []*org_alert.OrganizationNotification{
					{NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack", NotifySetting: `{"channel_id":"ch1"}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResp: []*model.OrganizationNotification{
				{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"channel_id":"ch1"}`, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:     "OK - empty",
			input:    &org_alert.ListOrgNotificationRequest{OrganizationId: 1},
			want:     &org_alert.ListOrgNotificationResponse{},
			mockResp: []*model.OrganizationNotification{},
		},
		{
			name:    "OK - record not found",
			input:   &org_alert.ListOrgNotificationRequest{OrganizationId: 1},
			want:    &org_alert.ListOrgNotificationResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG - DB error",
			input:   &org_alert.ListOrgNotificationRequest{OrganizationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &org_alert.ListOrgNotificationRequest{OrganizationId: 0},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrgAlertRepository(t)
			svc := OrgAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListOrgNotification", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListOrgNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrgNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *org_alert.GetOrgNotificationRequest
		want     *org_alert.GetOrgNotificationResponse
		wantErr  bool
		mockResp *model.OrganizationNotification
		mockErr  error
	}{
		{
			name:  "OK",
			input: &org_alert.GetOrgNotificationRequest{OrganizationId: 1, NotificationId: 1},
			want: &org_alert.GetOrgNotificationResponse{
				OrganizationNotification: &org_alert.OrganizationNotification{
					NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack",
					NotifySetting: `{"channel_id":"ch1"}`,
					CreatedAt:     now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockResp: &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"channel_id":"ch1"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "OK - record not found",
			input:   &org_alert.GetOrgNotificationRequest{OrganizationId: 1, NotificationId: 999},
			want:    &org_alert.GetOrgNotificationResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG - DB error",
			input:   &org_alert.GetOrgNotificationRequest{OrganizationId: 1, NotificationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &org_alert.GetOrgNotificationRequest{OrganizationId: 0, NotificationId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrgAlertRepository(t)
			svc := OrgAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("GetOrgNotification", test.RepeatMockAnything(3)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.GetOrgNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrgNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *org_alert.PutOrgNotificationRequest
		want        *org_alert.PutOrgNotificationResponse
		wantErr     bool
		mockGetResp *model.OrganizationNotification
		mockGetErr  error
		mockUpResp  *model.OrganizationNotification
		mockUpErr   error
	}{
		{
			name: "OK Insert",
			input: &org_alert.PutOrgNotificationRequest{
				OrganizationId: 1,
				Name:           "notif1",
				Type:           "slack",
				NotifySetting:  `{"webhook_url":"https://example.com"}`,
			},
			want: &org_alert.PutOrgNotificationResponse{
				OrganizationNotification: &org_alert.OrganizationNotification{
					NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack",
					NotifySetting: `{"webhook_url":"https://e**********","channel_id":"","data":{},"locale":""}`,
					CreatedAt:     now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockUpResp: &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "OK Update",
			input: &org_alert.PutOrgNotificationRequest{
				OrganizationId: 1,
				NotificationId: 1,
				Name:           "notif1",
				Type:           "slack",
				NotifySetting:  `{"webhook_url":"https://example.com"}`,
			},
			want: &org_alert.PutOrgNotificationResponse{
				OrganizationNotification: &org_alert.OrganizationNotification{
					NotificationId: 1, Name: "notif1", OrganizationId: 1, Type: "slack",
					NotifySetting: `{"webhook_url":"https://e**********","channel_id":"","data":{},"locale":""}`,
					CreatedAt:     now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockGetResp: &model.OrganizationNotification{NotificationID: 1, OrganizationID: 1, Name: "notif1", Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`, CreatedAt: now, UpdatedAt: now},
			mockUpResp:  &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: `{"webhook_url":"https://example.com"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "NG - record not found on update",
			input: &org_alert.PutOrgNotificationRequest{
				OrganizationId: 1,
				NotificationId: 1,
				Name:           "notif1",
				Type:           "slack",
				NotifySetting:  `{"webhook_url":"https://example.com"}`,
			},
			want:       &org_alert.PutOrgNotificationResponse{},
			mockGetErr: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrgAlertRepository(t)
			svc := OrgAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetOrgNotification", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertOrgNotification", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutOrgNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrgNotification(t *testing.T) {
	cases := []struct {
		name    string
		input   *org_alert.DeleteOrgNotificationRequest
		wantErr bool
		mockErr error
	}{
		{
			name:  "OK",
			input: &org_alert.DeleteOrgNotificationRequest{OrganizationId: 1, NotificationId: 1},
		},
		{
			name:    "NG - DB error",
			input:   &org_alert.DeleteOrgNotificationRequest{OrganizationId: 1, NotificationId: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:    "NG - validation error",
			input:   &org_alert.DeleteOrgNotificationRequest{OrganizationId: 0, NotificationId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrgAlertRepository(t)
			svc := OrgAlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockErr != nil || c.input.OrganizationId != 0 && !c.wantErr {
				mockDB.On("DeleteOrgNotification", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrgNotification(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
