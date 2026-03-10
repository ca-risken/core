package organization_alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization_alert"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *OrgAlertService) ListOrganizationNotification(ctx context.Context, req *organization_alert.ListOrganizationNotificationRequest) (*organization_alert.ListOrganizationNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOrgNotification(ctx, req.OrganizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_alert.ListOrganizationNotificationResponse{}, nil
		}
		return nil, err
	}
	data := organization_alert.ListOrganizationNotificationResponse{}
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

func (s *OrgAlertService) GetOrganizationNotification(ctx context.Context, req *organization_alert.GetOrganizationNotificationRequest) (*organization_alert.GetOrganizationNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.GetOrgNotification(ctx, req.OrganizationId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_alert.GetOrganizationNotificationResponse{}, nil
		}
		return nil, err
	}
	converted, err := convertOrgNotification(data, true)
	if err != nil {
		s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
		return nil, err
	}
	return &organization_alert.GetOrganizationNotificationResponse{OrganizationNotification: converted}, nil
}

func (s *OrgAlertService) PutOrganizationNotification(ctx context.Context, req *organization_alert.PutOrganizationNotificationRequest) (*organization_alert.PutOrganizationNotificationResponse, error) {
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
				return &organization_alert.PutOrganizationNotificationResponse{}, nil
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
	return &organization_alert.PutOrganizationNotificationResponse{OrganizationNotification: converted}, nil
}

func (s *OrgAlertService) DeleteOrganizationNotification(ctx context.Context, req *organization_alert.DeleteOrganizationNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOrgNotification(ctx, req.OrganizationId, req.NotificationId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OrgAlertService) ListOrganizationNotificationByProject(ctx context.Context, req *organization_alert.ListOrganizationNotificationByProjectRequest) (*organization_alert.ListOrganizationNotificationByProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOrgNotificationByProjectID(ctx, req.ProjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_alert.ListOrganizationNotificationByProjectResponse{}, nil
		}
		return nil, err
	}
	data := organization_alert.ListOrganizationNotificationByProjectResponse{}
	for _, d := range list {
		converted, err := convertOrgNotification(d, false)
		if err != nil {
			s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
			return nil, err
		}
		data.OrganizationNotification = append(data.OrganizationNotification, converted)
	}
	return &data, nil
}

func (s *OrgAlertService) TestOrganizationNotification(ctx context.Context, req *organization_alert.TestOrganizationNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.GetOrgNotification(ctx, req.OrganizationId, req.NotificationId)
	if err != nil {
		return nil, err
	}
	switch data.Type {
	case "slack":
		if err := s.sendSlackTestNotification(ctx, data.NotifySetting); err != nil {
			return nil, fmt.Errorf("failed to send test notification: %w", err)
		}
	default:
		s.logger.Warnf(ctx, "Unsupported notification type for test: %s", data.Type)
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

func convertOrgNotification(n *model.OrganizationNotification, mask bool) (*organization_alert.OrganizationNotification, error) {
	if n == nil {
		return &organization_alert.OrganizationNotification{}, nil
	}
	setting := n.NotifySetting
	if mask {
		var err error
		setting, err = maskingNotifySetting(n.Type, setting)
		if err != nil {
			return nil, err
		}
	}
	return &organization_alert.OrganizationNotification{
		NotificationId: n.NotificationID,
		Name:           n.Name,
		OrganizationId: n.OrganizationID,
		Type:           n.Type,
		NotifySetting:  setting,
		CreatedAt:      n.CreatedAt.Unix(),
		UpdatedAt:      n.UpdatedAt.Unix(),
	}, nil
}
