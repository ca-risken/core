package org_alert

import (
	"context"
	"errors"

	orgalert "github.com/ca-risken/core/proto/org_alert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *OrgAlertService) ListOrgAlertCondNotification(ctx context.Context, req *orgalert.ListOrgAlertCondNotificationRequest) (*orgalert.ListOrgAlertCondNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 100
	}
	list, err := s.repository.ListOrgAlertCondNotification(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId, req.NotificationId, pageSize+1, req.PageOffset)
	if err != nil {
		return nil, err
	}
	hasNext := len(list) > int(pageSize)
	if hasNext {
		list = list[:pageSize]
	}
	nextOffset := uint32(0)
	if hasNext {
		nextOffset = req.PageOffset + pageSize
	}
	return &orgalert.ListOrgAlertCondNotificationResponse{AlertCondNotification: list, HasNext: hasNext, NextPageOffset: nextOffset}, nil
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
	data, err := s.repository.UpdateOrgAlertCondNotificationCache(ctx, req.OrganizationId, req.ProjectId, req.AlertConditionId, req.NotificationId, req.CacheSecond)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "organization alert condition notification not found")
		}
		return nil, err
	}
	return &orgalert.UpdateOrgAlertCondNotificationCacheResponse{AlertCondNotification: data}, nil
}

func (s *OrgAlertService) UpdateOrgAlertProjectNotificationEnabled(ctx context.Context, req *orgalert.UpdateOrgAlertProjectNotificationEnabledRequest) (*orgalert.UpdateOrgAlertProjectNotificationEnabledResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.UpdateOrgAlertProjectNotificationEnabled(ctx, req.OrganizationId, req.ProjectId, req.NotificationId, req.Enabled)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "organization alert project notification not found")
		}
		return nil, err
	}
	return &orgalert.UpdateOrgAlertProjectNotificationEnabledResponse{AlertCondNotification: data}, nil
}

func (s *OrgAlertService) UpdateOrgAlertProjectNotificationCache(ctx context.Context, req *orgalert.UpdateOrgAlertProjectNotificationCacheRequest) (*orgalert.UpdateOrgAlertProjectNotificationCacheResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.UpdateOrgAlertProjectNotificationCache(ctx, req.OrganizationId, req.ProjectId, req.NotificationId, req.CacheSecond)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "organization alert project notification not found")
		}
		return nil, err
	}
	return &orgalert.UpdateOrgAlertProjectNotificationCacheResponse{AlertCondNotification: data}, nil
}
