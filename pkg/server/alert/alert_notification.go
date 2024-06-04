package alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

/**
 * Notification
 */

func (a *AlertService) ListNotification(ctx context.Context, req *alert.ListNotificationRequest) (*alert.ListNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return a.listNotification(ctx, req, true)
}

func (a *AlertService) ListNotificationForInternal(ctx context.Context, req *alert.ListNotificationForInternalRequest) (*alert.ListNotificationForInternalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	resp, err := a.listNotification(ctx, &alert.ListNotificationRequest{
		ProjectId: req.ProjectId,
		Type:      req.Type,
	}, false)
	if err != nil {
		return nil, err
	}
	return &alert.ListNotificationForInternalResponse{
		Notification: resp.Notification,
	}, nil
}

func (a *AlertService) listNotification(ctx context.Context, req *alert.ListNotificationRequest, mask bool) (*alert.ListNotificationResponse, error) {
	converted := convertListNotificationRequest(req)
	list, err := a.repository.ListNotification(ctx, converted.ProjectId, converted.Type, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListNotificationResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListNotificationResponse{}
	for _, d := range *list {
		convertedNotification, err := a.convertNotification(ctx, &d, mask)
		if err != nil {
			a.logger.Errorf(ctx, "Failed to convert Notification. error: %v", err)
			return nil, err
		}
		data.Notification = append(data.Notification, convertedNotification)
	}
	return &data, nil
}

func convertListNotificationRequest(req *alert.ListNotificationRequest) *alert.ListNotificationRequest {
	converted := alert.ListNotificationRequest{
		ProjectId: req.ProjectId,
		Type:      req.Type,
		FromAt:    req.FromAt,
		ToAt:      req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetNotification(ctx context.Context, req *alert.GetNotificationRequest) (*alert.GetNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetNotification(ctx, req.ProjectId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetNotificationResponse{}, nil
		}
		return nil, err
	}
	convertedNotification, err := a.convertNotification(ctx, data, true)
	if err != nil {
		a.logger.Errorf(ctx, "Failed to convert Notification. error: %v", err)
		return nil, err
	}
	return &alert.GetNotificationResponse{Notification: convertedNotification}, nil
}

func (a *AlertService) PutNotification(ctx context.Context, req *alert.PutNotificationRequest) (*alert.PutNotificationResponse, error) {
	err := req.Notification.Validate()
	if err != nil {
		return nil, err
	}
	var existData *model.Notification
	if !zero.IsZeroVal(req.Notification.NotificationId) {
		existData, err = a.repository.GetNotification(ctx, req.ProjectId, req.Notification.NotificationId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &alert.PutNotificationResponse{}, nil
			}
			return nil, err
		}
	}

	data := &model.Notification{
		NotificationID: req.Notification.NotificationId,
		Name:           req.Notification.Name,
		ProjectID:      req.Notification.ProjectId,
		Type:           req.Notification.Type,
		NotifySetting:  req.Notification.NotifySetting,
	}

	if !zero.IsZeroVal(existData) {
		switch existData.Type {
		case "slack":
			convertedNotifySetting, err := a.replaceSlackNotifySetting(ctx, existData.NotifySetting, data.NotifySetting)
			if err != nil {
				return nil, err
			}
			newNotifySetting, err := json.Marshal(convertedNotifySetting)
			if err != nil {
				a.logger.Errorf(ctx, "Error occured when marshal update.NotifySetting. err: %v", err)
				return nil, err
			}
			data.NotifySetting = string(newNotifySetting)
		default:
			a.logger.Warnf(ctx, "This notification_type is unimprement. type: %v", existData.Type)
		}
	}

	// Fiding upsert
	registerdData, err := a.repository.UpsertNotification(ctx, data)
	if err != nil {
		return nil, err
	}
	convertedNotification, err := a.convertNotification(ctx, registerdData, true)
	if err != nil {
		a.logger.Errorf(ctx, "Failed to convert Notification. error: %v", err)
		return nil, err
	}
	return &alert.PutNotificationResponse{Notification: convertedNotification}, nil
}

func (a *AlertService) replaceSlackNotifySetting(ctx context.Context, jsonNotifySettingExist, jsonNotifySettingUpdate string) (slackNotifySetting, error) {
	var notifySettingUpdate slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingUpdate), &notifySettingUpdate); err != nil {
		if err != nil {
			a.logger.Errorf(ctx, "Error occured when unmarshal update.NotifySetting. err: %v", err)
			return slackNotifySetting{}, err
		}
	}
	var notifySettingExist slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingExist), &notifySettingExist); err != nil {
		if err != nil {
			a.logger.Errorf(ctx, "Error occured when unmarshal exist.NotifySetting. err: %v", err)
			return slackNotifySetting{}, err
		}
	}
	if notifySettingUpdate.WebhookURL != "" {
		notifySettingUpdate.ChannelID = ""
		return notifySettingUpdate, nil
	}
	if notifySettingUpdate.ChannelID != "" {
		notifySettingUpdate.Data = slackNotifyOption{}
		return notifySettingUpdate, nil
	}
	notifySettingUpdate.WebhookURL = notifySettingExist.WebhookURL

	return notifySettingUpdate, nil
}

func (a *AlertService) DeleteNotification(ctx context.Context, req *alert.DeleteNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	notifications, err := a.repository.ListAlertCondNotification(ctx, req.ProjectId, 0, req.NotificationId, 0, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	for _, n := range *notifications {
		if err = a.repository.DeleteAlertCondNotification(ctx, n.ProjectID, n.AlertConditionID, n.NotificationID); err != nil {
			return nil, err
		}
	}

	err = a.repository.DeleteNotification(ctx, req.ProjectId, req.NotificationId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (a *AlertService) TestNotification(ctx context.Context, req *alert.TestNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	notification, err := a.repository.GetNotification(ctx, req.ProjectId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		return nil, err
	}
	switch notification.Type {
	case "slack":
		err = a.sendSlackTestNotification(ctx, a.baseURL, notification.NotifySetting, a.defaultLocale)
		if err != nil {
			a.logger.Errorf(ctx, "Error occured when sending test slack notification. err: %v", err)
			return nil, err
		}
	default:
		a.logger.Warnf(ctx, "This notification_type is unimplemented. type: %v", notification.Type)
	}
	return &empty.Empty{}, nil
}

func (a *AlertService) RequestProjectRoleNotification(ctx context.Context, req *alert.RequestProjectRoleNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	notifications, err := a.repository.ListNotification(ctx, req.ProjectId, "slack", 0, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	if notifications == nil || len(*notifications) == 0 {
		return nil, fmt.Errorf("notification not found. project_id: %d", req.ProjectId)
	}
	notification := (*notifications)[0]
	projects, err := a.projectClient.ListProject(ctx, &project.ListProjectRequest{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}
	user, err := a.iamClient.GetUser(ctx, &iam.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	switch notification.Type {
	case "slack":
		err = a.sendSlackRequestProjectRoleNotification(ctx, a.baseURL, notification.NotifySetting, a.defaultLocale, user.User.Name, projects.Project[0].Name, req.ProjectId)
		if err != nil {
			a.logger.Errorf(ctx, "Error occured when sending request authz slack notification. err: %v", err)
			return nil, err
		}
	default:
		a.logger.Warnf(ctx, "This notification_type is unimplemented. type: %v", notification.Type)
		return nil, fmt.Errorf("this notification_type is unavailable. type: %v", notification.Type)
	}
	return &empty.Empty{}, nil
}

func (a *AlertService) convertNotification(ctx context.Context, n *model.Notification, mask bool) (*alert.Notification, error) {
	if n == nil {
		return &alert.Notification{}, nil
	}
	var err error
	setting := n.NotifySetting
	if mask {
		setting, err = maskingNotifySetting(n.Type, setting)
		if err != nil {
			a.logger.Errorf(ctx, "Failed to masking notify setting. %v", err)
			return &alert.Notification{}, err
		}
	}
	return &alert.Notification{
		NotificationId: n.NotificationID,
		Name:           n.Name,
		ProjectId:      n.ProjectID,
		Type:           n.Type,
		NotifySetting:  setting,
		CreatedAt:      n.CreatedAt.Unix(),
		UpdatedAt:      n.UpdatedAt.Unix(),
	}, nil
}

func maskingNotifySetting(notificationType, notifySetting string) (string, error) {
	switch notificationType {
	case "slack":
		var setting slackNotifySetting
		if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
			return "", err
		}
		if setting.WebhookURL == "" {
			return notifySetting, nil
		}
		setting.WebhookURL = maskRight(setting.WebhookURL, len(setting.WebhookURL)/2)
		ret, err := json.Marshal(setting)
		if err != nil {
			return "", err
		}
		return string(ret), err
	default:
		return notifySetting, nil
	}
}

func maskRight(s string, num int) string {
	rs := []rune(s)
	for i := num; i < len(rs); i++ {
		rs[i] = '*'
	}
	return string(rs)
}

const (
	MAX_NOTIFY_FINDING_NUM = 3
)

type findingDetail struct {
	FindingCount int
	Exampls      []*findingExample
}

type findingExample struct {
	FindingID    uint64
	Description  string
	ResourceName string
	DataSource   string
	Score        float32
	Tags         []string
}

func (a *AlertService) getFindingDetailsForNotification(ctx context.Context, projectID uint32, findingIDs *[]uint64) (
	*findingDetail, error,
) {
	findings := findingDetail{
		FindingCount: len(*findingIDs),
	}
	for _, id := range *findingIDs {
		if len(findings.Exampls) >= MAX_NOTIFY_FINDING_NUM {
			break
		}

		ex := findingExample{}
		// finding
		resp, err := a.findingClient.GetFinding(ctx, &finding.GetFindingRequest{FindingId: id, ProjectId: projectID})
		if err != nil {
			return nil, fmt.Errorf("get finding error: err=%w", err)
		}
		ex.FindingID = resp.Finding.FindingId
		ex.Description = resp.Finding.Description
		ex.ResourceName = resp.Finding.ResourceName
		ex.DataSource = resp.Finding.DataSource
		ex.Score = resp.Finding.Score

		// finding tag
		tagResp, err := a.findingClient.ListFindingTag(ctx, &finding.ListFindingTagRequest{FindingId: id, ProjectId: projectID})
		if err != nil {
			return nil, fmt.Errorf("get finding tag error: err=%w", err)
		}
		for _, t := range tagResp.Tag {
			ex.Tags = append(ex.Tags, t.Tag)
		}
		findings.Exampls = append(findings.Exampls, &ex)
	}
	return &findings, nil
}
