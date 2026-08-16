package org_alert

import (
	"context"
	"errors"

	orgalert "github.com/ca-risken/core/proto/org_alert"
	"gorm.io/gorm"
)

func (s *OrgAlertService) ListOrgAlertCondNotification(ctx context.Context, req *orgalert.ListOrgAlertCondNotificationRequest) (*orgalert.ListOrgAlertCondNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOrgAlertCondNotification(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId, req.NotificationId)
	if err != nil {
		return nil, err
	}
	return &orgalert.ListOrgAlertCondNotificationResponse{AlertCondNotification: list}, nil
}

func (s *OrgAlertService) GetOrgAlertCondNotification(ctx context.Context, req *orgalert.GetOrgAlertCondNotificationRequest) (*orgalert.GetOrgAlertCondNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.GetOrgAlertCondNotification(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &orgalert.GetOrgAlertCondNotificationResponse{}, nil
		}
		return nil, err
	}
	return &orgalert.GetOrgAlertCondNotificationResponse{AlertCondNotification: data}, nil
}

func (s *OrgAlertService) UpdateOrgAlertCondNotificationCache(ctx context.Context, req *orgalert.UpdateOrgAlertCondNotificationCacheRequest) (*orgalert.UpdateOrgAlertCondNotificationCacheResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	exists, err := s.repository.ExistsOrgAlertConditionMembership(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &orgalert.UpdateOrgAlertCondNotificationCacheResponse{}, nil
	}
	if _, err := s.repository.GetOrgAlertCondNotification(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId, req.NotificationId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &orgalert.UpdateOrgAlertCondNotificationCacheResponse{}, nil
		}
		return nil, err
	}
	data, err := s.repository.UpdateOrgAlertCondNotificationCache(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId, req.NotificationId, req.CacheSecond)
	if err != nil {
		return nil, err
	}
	return &orgalert.UpdateOrgAlertCondNotificationCacheResponse{AlertCondNotification: data}, nil
}
