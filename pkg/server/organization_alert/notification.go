package organization_alert

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization_alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (s *OrganizationAlertService) ListOrganizationNotification(ctx context.Context, req *organization_alert.ListOrganizationNotificationRequest) (*organization_alert.ListOrganizationNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOrganizationNotification(ctx, req.OrganizationId, req.Type)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_alert.ListOrganizationNotificationResponse{}, nil
		}
		return nil, err
	}
	data := organization_alert.ListOrganizationNotificationResponse{}
	for _, d := range list {
		converted, err := convertOrganizationNotification(d, true)
		if err != nil {
			s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
			return nil, err
		}
		data.OrganizationNotification = append(data.OrganizationNotification, converted)
	}
	return &data, nil
}

func (s *OrganizationAlertService) GetOrganizationNotification(ctx context.Context, req *organization_alert.GetOrganizationNotificationRequest) (*organization_alert.GetOrganizationNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := s.repository.GetOrganizationNotification(ctx, req.OrganizationId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_alert.GetOrganizationNotificationResponse{}, nil
		}
		return nil, err
	}
	converted, err := convertOrganizationNotification(data, true)
	if err != nil {
		s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
		return nil, err
	}
	return &organization_alert.GetOrganizationNotificationResponse{OrganizationNotification: converted}, nil
}

func (s *OrganizationAlertService) PutOrganizationNotification(ctx context.Context, req *organization_alert.PutOrganizationNotificationRequest) (*organization_alert.PutOrganizationNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var existData *model.OrganizationNotification
	var err error
	if !zero.IsZeroVal(req.NotificationId) {
		existData, err = s.repository.GetOrganizationNotification(ctx, req.OrganizationId, req.NotificationId)
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

	if !zero.IsZeroVal(existData) {
		switch existData.Type {
		case "slack":
			convertedNotifySetting, err := replaceSlackNotifySetting(ctx, s.logger, existData.NotifySetting, data.NotifySetting)
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

	registeredData, err := s.repository.UpsertOrganizationNotification(ctx, data)
	if err != nil {
		return nil, err
	}
	converted, err := convertOrganizationNotification(registeredData, true)
	if err != nil {
		s.logger.Errorf(ctx, "Failed to convert OrganizationNotification. error: %v", err)
		return nil, err
	}
	return &organization_alert.PutOrganizationNotificationResponse{OrganizationNotification: converted}, nil
}

func (s *OrganizationAlertService) DeleteOrganizationNotification(ctx context.Context, req *organization_alert.DeleteOrganizationNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOrganizationNotification(ctx, req.OrganizationId, req.NotificationId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OrganizationAlertService) TestOrganizationNotification(ctx context.Context, req *organization_alert.TestOrganizationNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	notification, err := s.repository.GetOrganizationNotification(ctx, req.OrganizationId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		return nil, err
	}
	switch notification.Type {
	case "slack":
		err = s.sendSlackTestNotification(ctx, notification.NotifySetting, s.defaultLocale)
		if err != nil {
			s.logger.Errorf(ctx, "Error occured when sending test slack notification. err: %v", err)
			return nil, err
		}
	default:
		s.logger.Warnf(ctx, "This notification_type is unimplemented. type: %v", notification.Type)
	}
	return &empty.Empty{}, nil
}

func convertOrganizationNotification(n *model.OrganizationNotification, mask bool) (*organization_alert.OrganizationNotification, error) {
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
