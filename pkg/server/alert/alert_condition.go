package alert

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (a *AlertService) ListAlertCondition(ctx context.Context, req *alert.ListAlertConditionRequest) (*alert.ListAlertConditionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	converted := convertListAlertConditionRequest(req)
	list, err := a.repository.ListAlertCondition(ctx, converted.ProjectId, converted.Severity, converted.Enabled, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListAlertConditionResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListAlertConditionResponse{}
	for _, d := range *list {
		data.AlertCondition = append(data.AlertCondition, convertAlertCondition(&d))
	}
	return &data, nil
}

func convertListAlertConditionRequest(req *alert.ListAlertConditionRequest) *alert.ListAlertConditionRequest {
	converted := alert.ListAlertConditionRequest{
		ProjectId: req.ProjectId,
		Severity:  req.Severity,
		Enabled:   req.Enabled,
		FromAt:    req.FromAt,
		ToAt:      req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetAlertCondition(ctx context.Context, req *alert.GetAlertConditionRequest) (*alert.GetAlertConditionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetAlertCondition(ctx, req.ProjectId, req.AlertConditionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetAlertConditionResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertConditionResponse{AlertCondition: convertAlertCondition(data)}, nil
}

func (a *AlertService) PutAlertCondition(ctx context.Context, req *alert.PutAlertConditionRequest) (*alert.PutAlertConditionResponse, error) {
	if err := req.AlertCondition.Validate(); err != nil {
		return nil, err
	}

	data := &model.AlertCondition{
		AlertConditionID: req.AlertCondition.AlertConditionId,
		Description:      req.AlertCondition.Description,
		Severity:         req.AlertCondition.Severity,
		ProjectID:        req.AlertCondition.ProjectId,
		Enabled:          req.AlertCondition.Enabled,
		AndOr:            req.AlertCondition.AndOr,
	}

	// Fiding upsert
	registerdData, err := a.repository.UpsertAlertCondition(ctx, data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertConditionResponse{AlertCondition: convertAlertCondition(registerdData)}, nil
}

func (a *AlertService) DeleteAlertCondition(ctx context.Context, req *alert.DeleteAlertConditionRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteAlertCondition(ctx, req.ProjectId, req.AlertConditionId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * AlertRule
 */

func (a *AlertService) ListAlertRule(ctx context.Context, req *alert.ListAlertRuleRequest) (*alert.ListAlertRuleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	converted := convertListAlertRuleRequest(req)
	list, err := a.repository.ListAlertRule(ctx, converted.ProjectId, converted.FromScore, converted.ToScore, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListAlertRuleResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListAlertRuleResponse{}
	for _, d := range *list {
		data.AlertRule = append(data.AlertRule, convertAlertRule(&d))
	}
	return &data, nil
}

func convertListAlertRuleRequest(req *alert.ListAlertRuleRequest) *alert.ListAlertRuleRequest {
	converted := alert.ListAlertRuleRequest{
		ProjectId: req.ProjectId,
		FromScore: req.FromScore,
		ToScore:   req.ToScore,
		FromAt:    req.FromAt,
		ToAt:      req.ToAt,
	}
	if converted.ToScore == 0 {
		converted.ToScore = 1.0
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetAlertRule(ctx context.Context, req *alert.GetAlertRuleRequest) (*alert.GetAlertRuleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetAlertRule(ctx, req.ProjectId, req.AlertRuleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetAlertRuleResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertRuleResponse{AlertRule: convertAlertRule(data)}, nil
}

func (a *AlertService) PutAlertRule(ctx context.Context, req *alert.PutAlertRuleRequest) (*alert.PutAlertRuleResponse, error) {
	if err := req.AlertRule.Validate(); err != nil {
		return nil, err
	}

	data := &model.AlertRule{
		AlertRuleID:  req.AlertRule.AlertRuleId,
		ProjectID:    req.AlertRule.ProjectId,
		Name:         req.AlertRule.Name,
		Score:        req.AlertRule.Score,
		ResourceName: req.AlertRule.ResourceName,
		Tag:          req.AlertRule.Tag,
		FindingCnt:   req.AlertRule.FindingCnt,
	}

	// Fiding upsert
	registerdData, err := a.repository.UpsertAlertRule(ctx, data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertRuleResponse{AlertRule: convertAlertRule(registerdData)}, nil
}

func (a *AlertService) DeleteAlertRule(ctx context.Context, req *alert.DeleteAlertRuleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteAlertRule(ctx, req.ProjectId, req.AlertRuleId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * AlertCondRule
 */

func (a *AlertService) ListAlertCondRule(ctx context.Context, req *alert.ListAlertCondRuleRequest) (*alert.ListAlertCondRuleResponse, error) {
	converted := convertListAlertCondRuleRequest(req)
	list, err := a.repository.ListAlertCondRule(ctx, converted.ProjectId, converted.AlertConditionId, converted.AlertRuleId, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListAlertCondRuleResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListAlertCondRuleResponse{}
	for _, d := range *list {
		data.AlertCondRule = append(data.AlertCondRule, convertAlertCondRule(&d))
	}
	return &data, nil
}

func convertListAlertCondRuleRequest(req *alert.ListAlertCondRuleRequest) *alert.ListAlertCondRuleRequest {
	converted := alert.ListAlertCondRuleRequest{
		ProjectId:        req.ProjectId,
		AlertConditionId: req.AlertConditionId,
		AlertRuleId:      req.AlertRuleId,
		FromAt:           req.FromAt,
		ToAt:             req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetAlertCondRule(ctx context.Context, req *alert.GetAlertCondRuleRequest) (*alert.GetAlertCondRuleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := a.repository.GetAlertCondRule(ctx, req.ProjectId, req.AlertConditionId, req.AlertRuleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetAlertCondRuleResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertCondRuleResponse{AlertCondRule: convertAlertCondRule(data)}, nil
}

func (a *AlertService) PutAlertCondRule(ctx context.Context, req *alert.PutAlertCondRuleRequest) (*alert.PutAlertCondRuleResponse, error) {
	if err := req.AlertCondRule.Validate(); err != nil {
		return nil, err
	}
	data := &model.AlertCondRule{
		AlertConditionID: req.AlertCondRule.AlertConditionId,
		AlertRuleID:      req.AlertCondRule.AlertRuleId,
		ProjectID:        req.AlertCondRule.ProjectId,
	}

	// Fiding upsert
	registerdData, err := a.repository.UpsertAlertCondRule(ctx, data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertCondRuleResponse{AlertCondRule: convertAlertCondRule(registerdData)}, nil
}

func (a *AlertService) DeleteAlertCondRule(ctx context.Context, req *alert.DeleteAlertCondRuleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteAlertCondRule(ctx, req.ProjectId, req.AlertConditionId, req.AlertRuleId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * Notification
 */

func (a *AlertService) ListNotification(ctx context.Context, req *alert.ListNotificationRequest) (*alert.ListNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
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
		data.Notification = append(data.Notification, convertNotification(&d))
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
	return &alert.GetNotificationResponse{Notification: convertNotification(data)}, nil
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
			convertedNotifySetting, err := replaceSlackNotifySetting(existData.NotifySetting, data.NotifySetting)
			if err != nil {
				return nil, err
			}
			newNotifySetting, err := json.Marshal(convertedNotifySetting)
			if err != nil {
				appLogger.Errorf("Error occured when marshal update.NotifySetting. err: %v", err)
				return nil, err
			}
			data.NotifySetting = string(newNotifySetting)
		default:
			appLogger.Warnf("This notification_type is unimprement. type: %v", existData.Type)
		}
	}

	// Fiding upsert
	registerdData, err := a.repository.UpsertNotification(ctx, data)
	if err != nil {
		return nil, err
	}

	return &alert.PutNotificationResponse{Notification: convertNotification(registerdData)}, nil
}

func replaceSlackNotifySetting(jsonNotifySettingExist, jsonNotifySettingUpdate string) (slackNotifySetting, error) {
	var notifySettingUpdate slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingUpdate), &notifySettingUpdate); err != nil {
		if err != nil {
			appLogger.Errorf("Error occured when unmarshal update.NotifySetting. err: %v", err)
			return slackNotifySetting{}, err
		}
	}
	var notifySettingExist slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingExist), &notifySettingExist); err != nil {
		if err != nil {
			appLogger.Errorf("Error occured when unmarshal exist.NotifySetting. err: %v", err)
			return slackNotifySetting{}, err
		}
	}
	if !zero.IsZeroVal(notifySettingUpdate.WebhookURL) {
		return notifySettingUpdate, nil
	}
	notifySettingUpdate.WebhookURL = notifySettingExist.WebhookURL

	return notifySettingUpdate, nil
}

func (a *AlertService) DeleteNotification(ctx context.Context, req *alert.DeleteNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteNotification(ctx, req.ProjectId, req.NotificationId)
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
		err = sendSlackTestNotification(a.notificationAlertURL, notification.NotifySetting)
		if err != nil {
			appLogger.Errorf("Error occured when sending test slack notification. err: %v", err)
			return nil, err
		}
	default:
		appLogger.Warnf("This notification_type is unimprement. type: %v", notification.Type)
	}
	return &empty.Empty{}, nil
}

/**
 * AlertCondNotification
 */

func (a *AlertService) ListAlertCondNotification(ctx context.Context, req *alert.ListAlertCondNotificationRequest) (*alert.ListAlertCondNotificationResponse, error) {
	converted := convertListAlertCondNotificationRequest(req)
	list, err := a.repository.ListAlertCondNotification(ctx, converted.ProjectId, converted.AlertConditionId, converted.NotificationId, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListAlertCondNotificationResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListAlertCondNotificationResponse{}
	for _, d := range *list {
		data.AlertCondNotification = append(data.AlertCondNotification, convertAlertCondNotification(&d))
	}
	return &data, nil
}

func convertListAlertCondNotificationRequest(req *alert.ListAlertCondNotificationRequest) *alert.ListAlertCondNotificationRequest {
	converted := alert.ListAlertCondNotificationRequest{
		ProjectId:        req.ProjectId,
		AlertConditionId: req.AlertConditionId,
		NotificationId:   req.NotificationId,
		FromAt:           req.FromAt,
		ToAt:             req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetAlertCondNotification(ctx context.Context, req *alert.GetAlertCondNotificationRequest) (*alert.GetAlertCondNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := a.repository.GetAlertCondNotification(ctx, req.ProjectId, req.AlertConditionId, req.NotificationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetAlertCondNotificationResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertCondNotificationResponse{AlertCondNotification: convertAlertCondNotification(data)}, nil
}

func (a *AlertService) PutAlertCondNotification(ctx context.Context, req *alert.PutAlertCondNotificationRequest) (*alert.PutAlertCondNotificationResponse, error) {
	if err := req.AlertCondNotification.Validate(); err != nil {
		return nil, err
	}
	data := &model.AlertCondNotification{
		AlertConditionID: req.AlertCondNotification.AlertConditionId,
		NotificationID:   req.AlertCondNotification.NotificationId,
		ProjectID:        req.AlertCondNotification.ProjectId,
		CacheSecond:      req.AlertCondNotification.CacheSecond,
		NotifiedAt:       time.Unix(req.AlertCondNotification.NotifiedAt, 0),
	}

	// Fiding upsert
	registerdData, err := a.repository.UpsertAlertCondNotification(ctx, data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertCondNotificationResponse{AlertCondNotification: convertAlertCondNotification(registerdData)}, nil
}

func (a *AlertService) DeleteAlertCondNotification(ctx context.Context, req *alert.DeleteAlertCondNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteAlertCondNotification(ctx, req.ProjectId, req.AlertConditionId, req.NotificationId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * Converter
 */

func convertAlertCondition(a *model.AlertCondition) *alert.AlertCondition {
	if a == nil {
		return &alert.AlertCondition{}
	}
	return &alert.AlertCondition{
		AlertConditionId: a.AlertConditionID,
		Description:      a.Description,
		Severity:         a.Severity,
		ProjectId:        a.ProjectID,
		AndOr:            a.AndOr,
		Enabled:          a.Enabled,
		CreatedAt:        a.CreatedAt.Unix(),
		UpdatedAt:        a.UpdatedAt.Unix(),
	}
}

func convertAlertRule(a *model.AlertRule) *alert.AlertRule {
	if a == nil {
		return &alert.AlertRule{}
	}
	return &alert.AlertRule{
		AlertRuleId:  a.AlertRuleID,
		Name:         a.Name,
		Score:        a.Score,
		ProjectId:    a.ProjectID,
		ResourceName: a.ResourceName,
		Tag:          a.Tag,
		FindingCnt:   a.FindingCnt,
		CreatedAt:    a.CreatedAt.Unix(),
		UpdatedAt:    a.UpdatedAt.Unix(),
	}
}

func convertAlertCondRule(a *model.AlertCondRule) *alert.AlertCondRule {
	if a == nil {
		return &alert.AlertCondRule{}
	}
	return &alert.AlertCondRule{
		AlertConditionId: a.AlertConditionID,
		AlertRuleId:      a.AlertRuleID,
		ProjectId:        a.ProjectID,
		CreatedAt:        a.CreatedAt.Unix(),
		UpdatedAt:        a.UpdatedAt.Unix(),
	}
}

func convertNotification(n *model.Notification) *alert.Notification {
	if n == nil {
		return &alert.Notification{}
	}
	maskingSetting, err := maskingNotifySetting(n.Type, n.NotifySetting)
	if err != nil {
		appLogger.Errorf("Failed to masking notify setting. %v", err)
		maskingSetting = n.NotifySetting
	}
	return &alert.Notification{
		NotificationId: n.NotificationID,
		Name:           n.Name,
		ProjectId:      n.ProjectID,
		Type:           n.Type,
		NotifySetting:  maskingSetting,
		CreatedAt:      n.CreatedAt.Unix(),
		UpdatedAt:      n.UpdatedAt.Unix(),
	}
}

func convertAlertCondNotification(a *model.AlertCondNotification) *alert.AlertCondNotification {
	if a == nil {
		return &alert.AlertCondNotification{}
	}
	return &alert.AlertCondNotification{
		AlertConditionId: a.AlertConditionID,
		NotificationId:   a.NotificationID,
		ProjectId:        a.ProjectID,
		CacheSecond:      a.CacheSecond,
		NotifiedAt:       a.NotifiedAt.Unix(),
		CreatedAt:        a.CreatedAt.Unix(),
		UpdatedAt:        a.UpdatedAt.Unix(),
	}
}

func maskingNotifySetting(notificationType, notifySetting string) (string, error) {
	switch notificationType {
	case "slack":
		var setting slackNotifySetting
		if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
			return "", err
		}
		//　通常webhook_urlは存在するはずだが、万が一ない場合はそのまま返す
		if zero.IsZeroVal(setting.WebhookURL) {
			return notifySetting, nil
		}
		setting.WebhookURL = maskRight(setting.WebhookURL, len(setting.WebhookURL)/2)
		ret, err := json.Marshal(setting)
		if err != nil {
			return notifySetting, err
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
