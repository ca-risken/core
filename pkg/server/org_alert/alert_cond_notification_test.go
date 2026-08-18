package org_alert

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	orgalert "github.com/ca-risken/core/proto/org_alert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestListOrgAlertCondNotification(t *testing.T) {
	relation := &orgalert.OrgAlertCondNotification{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, CacheSecond: 1800}
	cases := []struct {
		name       string
		input      *orgalert.ListOrgAlertCondNotificationRequest
		mockResp   []*orgalert.OrgAlertCondNotification
		mockErr    error
		want       *orgalert.ListOrgAlertCondNotificationResponse
		wantErr    bool
		expectCall bool
	}{
		{
			name:       "OK - all four keys",
			input:      &orgalert.ListOrgAlertCondNotificationRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4},
			mockResp:   []*orgalert.OrgAlertCondNotification{relation},
			want:       &orgalert.ListOrgAlertCondNotificationResponse{AlertCondNotification: []*orgalert.OrgAlertCondNotification{relation}},
			expectCall: true,
		},
		{
			name:       "OK - organization scope",
			input:      &orgalert.ListOrgAlertCondNotificationRequest{OrganizationId: 1},
			mockResp:   []*orgalert.OrgAlertCondNotification{relation},
			want:       &orgalert.ListOrgAlertCondNotificationResponse{AlertCondNotification: []*orgalert.OrgAlertCondNotification{relation}},
			expectCall: true,
		},
		{
			name:       "NG - DB error",
			input:      &orgalert.ListOrgAlertCondNotificationRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4},
			mockErr:    errors.New("DB error"),
			wantErr:    true,
			expectCall: true,
		},
		{
			name:    "NG - validation error",
			input:   &orgalert.ListOrgAlertCondNotificationRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repository := mocks.NewOrgAlertRepository(t)
			service := OrgAlertService{repository: repository, logger: logging.NewLogger()}
			if c.expectCall {
				repository.On("ListOrgAlertCondNotification", mock.Anything, c.input.OrganizationId, c.input.ProjectId, c.input.AlertConditionId, c.input.NotificationId).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := service.ListOrgAlertCondNotification(context.Background(), c.input)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrgAlertCondNotification(t *testing.T) {
	relation := &orgalert.OrgAlertCondNotification{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, CacheSecond: 1800}
	cases := []struct {
		name     string
		input    *orgalert.GetOrgAlertCondNotificationRequest
		mockResp *orgalert.OrgAlertCondNotification
		mockErr  error
		want     *orgalert.GetOrgAlertCondNotificationResponse
		wantErr  bool
	}{
		{
			name:     "OK",
			input:    &orgalert.GetOrgAlertCondNotificationRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4},
			mockResp: relation,
			want:     &orgalert.GetOrgAlertCondNotificationResponse{AlertCondNotification: relation},
		},
		{
			name:    "OK - missing relation",
			input:   &orgalert.GetOrgAlertCondNotificationRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4},
			mockErr: gorm.ErrRecordNotFound,
			want:    &orgalert.GetOrgAlertCondNotificationResponse{},
		},
		{
			name:    "NG - DB error",
			input:   &orgalert.GetOrgAlertCondNotificationRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4},
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
		{
			name:    "NG - validation error",
			input:   &orgalert.GetOrgAlertCondNotificationRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repository := mocks.NewOrgAlertRepository(t)
			service := OrgAlertService{repository: repository, logger: logging.NewLogger()}
			if c.input.NotificationId != 0 {
				repository.On("GetOrgAlertCondNotification", mock.Anything, uint32(1), uint32(2), uint32(3), uint32(4)).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := service.GetOrgAlertCondNotification(context.Background(), c.input)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestUpdateOrgAlertCondNotificationCache(t *testing.T) {
	relation := &orgalert.OrgAlertCondNotification{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, CacheSecond: 300}
	validInput := &orgalert.UpdateOrgAlertCondNotificationCacheRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, CacheSecond: 300}
	cases := []struct {
		name       string
		input      *orgalert.UpdateOrgAlertCondNotificationCacheRequest
		updateResp *orgalert.OrgAlertCondNotification
		updateErr  error
		want       *orgalert.UpdateOrgAlertCondNotificationCacheResponse
		wantErr    bool
		expectCall bool
	}{
		{
			name: "OK - atomic update", input: validInput, updateResp: relation,
			want: &orgalert.UpdateOrgAlertCondNotificationCacheResponse{AlertCondNotification: relation}, expectCall: true,
		},
		{
			name: "NG - membership or relation is missing", input: validInput,
			updateErr: gorm.ErrRecordNotFound, wantErr: true, expectCall: true,
		},
		{
			name: "NG - update error", input: validInput,
			updateErr: errors.New("DB error"), wantErr: true, expectCall: true,
		},
		{
			name:    "NG - validation error",
			input:   &orgalert.UpdateOrgAlertCondNotificationCacheRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3},
			wantErr: true,
		},
		{
			name:    "NG - cache exceeds one year",
			input:   &orgalert.UpdateOrgAlertCondNotificationCacheRequest{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, CacheSecond: 31536001},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repository := mocks.NewOrgAlertRepository(t)
			service := OrgAlertService{repository: repository, logger: logging.NewLogger()}
			if c.expectCall {
				repository.On("UpdateOrgAlertCondNotificationCache", mock.Anything, uint32(1), uint32(2), uint32(3), uint32(4), uint32(300)).Return(c.updateResp, c.updateErr).Once()
			}
			got, err := service.UpdateOrgAlertCondNotificationCache(context.Background(), c.input)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
