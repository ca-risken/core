package org_alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/org_alert"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *OrgAlertService) ListOrgNotification(ctx context.Context, req *org_alert.ListOrgNotificationRequest) (*org_alert.ListOrgNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOrgNotification(ctx, req.OrganizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_alert.ListOrgNotificationResponse{}, nil
		}
		return nil, err
	}
	data := org_alert.ListOrgNotificationResponse{}
	for _, d := range list {
		converted, err := convertOrgNotification(d, true)
		if err != nil {
			s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
			return nil, err
		}
		data.OrganizationNotification = append(data.OrganizationNotification, converted)
	}
	return &data, nil
}

func (s *OrgAlertService) GetOrgNotification(ctx context.Context, req *org_alert.GetOrgNotificationRequest) (*org_alert.GetOrgNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.GetOrgNotification(ctx, req.OrganizationId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_alert.GetOrgNotificationResponse{}, nil
		}
		return nil, err
	}
	converted, err := convertOrgNotification(data, true)
	if err != nil {
		s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
		return nil, err
	}
	return &org_alert.GetOrgNotificationResponse{OrganizationNotification: converted}, nil
}

func (s *OrgAlertService) PutOrgNotification(ctx context.Context, req *org_alert.PutOrgNotificationRequest) (*org_alert.PutOrgNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.NotificationId == 0 {
		if err := validateNewNotifySetting(req.NotifySetting); err != nil {
			return nil, err
		}
	} else {
		if err := validateExistingNotifySetting(req.NotifySetting); err != nil {
			return nil, err
		}
	}

	var existData *model.OrganizationNotification
	var err error
	if req.NotificationId != 0 {
		existData, err = s.repository.GetOrgNotification(ctx, req.OrganizationId, req.NotificationId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &org_alert.PutOrgNotificationResponse{}, nil
			}
			return nil, err
		}
	}

	data := &model.OrganizationNotification{
		NotificationID: req.NotificationId,
		Name:           req.Name,
		OrganizationID: req.OrganizationId,
		Type:           req.Type,
		NotifySetting:  req.NotifySetting,
	}

	if existData != nil {
		switch existData.Type {
		case "slack":
			convertedNotifySetting, err := s.replaceSlackNotifySetting(ctx, existData.NotifySetting, data.NotifySetting)
			if err != nil {
				return nil, err
			}
			newNotifySetting, err := json.Marshal(convertedNotifySetting)
			if err != nil {
				s.logger.Errorf(ctx, "Error occured when marshal update.NotifySetting. err: %v", err)
				return nil, err
			}
			data.NotifySetting = string(newNotifySetting)
		default:
			s.logger.Warnf(ctx, "This notification_type is unimprement. type: %v", existData.Type)
		}
	}

	registeredData, err := s.repository.UpsertOrgNotification(ctx, data)
	if err != nil {
		return nil, err
	}
	converted, err := convertOrgNotification(registeredData, true)
	if err != nil {
		s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
		return nil, err
	}
	return &org_alert.PutOrgNotificationResponse{OrganizationNotification: converted}, nil
}

func (s *OrgAlertService) DeleteOrgNotification(ctx context.Context, req *org_alert.DeleteOrgNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOrgNotification(ctx, req.OrganizationId, req.NotificationId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func validateNewNotifySetting(value string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(value), &setting); err != nil {
		return fmt.Errorf("invalid json, %w", err)
	}
	if strings.TrimSpace(setting.WebhookURL) == "" && strings.TrimSpace(setting.ChannelID) == "" {
		return errors.New("required webhook_url or channel_id in json")
	}
	if err := validation.Validate(strings.TrimSpace(setting.WebhookURL), is.URL); err != nil {
		return err
	}
	return nil
}

func validateExistingNotifySetting(value string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(value), &setting); err != nil {
		return fmt.Errorf("invalid json, %w", err)
	}
	if strings.TrimSpace(setting.WebhookURL) != "" {
		if err := validation.Validate(strings.TrimSpace(setting.WebhookURL), validation.Required, is.URL); err != nil {
			return err
		}
	}
	return nil
}

func convertOrgNotification(n *model.OrganizationNotification, mask bool) (*org_alert.OrganizationNotification, error) {
	if n == nil {
		return &org_alert.OrganizationNotification{}, nil
	}
	setting := n.NotifySetting
	if mask {
		var err error
		setting, err = maskingNotifySetting(n.Type, setting)
		if err != nil {
			return nil, err
		}
	}
	return &org_alert.OrganizationNotification{
		NotificationId: n.NotificationID,
		Name:           n.Name,
		OrganizationId: n.OrganizationID,
		Type:           n.Type,
		NotifySetting:  setting,
		CreatedAt:      n.CreatedAt.Unix(),
		UpdatedAt:      n.UpdatedAt.Unix(),
	}, nil
}
