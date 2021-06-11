package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/vikyd/zero"
)

/**
 * AlertCondition
 */

func (f *alertService) ListAlertCondition(ctx context.Context, req *alert.ListAlertConditionRequest) (*alert.ListAlertConditionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	converted := convertListAlertConditionRequest(req)
	list, err := f.repository.ListAlertCondition(converted.ProjectId, converted.Severity, converted.Enabled, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetAlertCondition(ctx context.Context, req *alert.GetAlertConditionRequest) (*alert.GetAlertConditionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetAlertCondition(req.ProjectId, req.AlertConditionId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetAlertConditionResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertConditionResponse{AlertCondition: convertAlertCondition(data)}, nil
}

func (f *alertService) PutAlertCondition(ctx context.Context, req *alert.PutAlertConditionRequest) (*alert.PutAlertConditionResponse, error) {
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
	registerdData, err := f.repository.UpsertAlertCondition(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertConditionResponse{AlertCondition: convertAlertCondition(registerdData)}, nil
}

func (f *alertService) DeleteAlertCondition(ctx context.Context, req *alert.DeleteAlertConditionRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteAlertCondition(req.ProjectId, req.AlertConditionId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * AlertRule
 */

func (f *alertService) ListAlertRule(ctx context.Context, req *alert.ListAlertRuleRequest) (*alert.ListAlertRuleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	converted := convertListAlertRuleRequest(req)
	list, err := f.repository.ListAlertRule(converted.ProjectId, converted.FromScore, converted.ToScore, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetAlertRule(ctx context.Context, req *alert.GetAlertRuleRequest) (*alert.GetAlertRuleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetAlertRule(req.ProjectId, req.AlertRuleId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetAlertRuleResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertRuleResponse{AlertRule: convertAlertRule(data)}, nil
}

func (f *alertService) PutAlertRule(ctx context.Context, req *alert.PutAlertRuleRequest) (*alert.PutAlertRuleResponse, error) {
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
	registerdData, err := f.repository.UpsertAlertRule(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertRuleResponse{AlertRule: convertAlertRule(registerdData)}, nil
}

func (f *alertService) DeleteAlertRule(ctx context.Context, req *alert.DeleteAlertRuleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteAlertRule(req.ProjectId, req.AlertRuleId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * AlertCondRule
 */

func (f *alertService) ListAlertCondRule(ctx context.Context, req *alert.ListAlertCondRuleRequest) (*alert.ListAlertCondRuleResponse, error) {
	converted := convertListAlertCondRuleRequest(req)
	list, err := f.repository.ListAlertCondRule(converted.ProjectId, converted.AlertConditionId, converted.AlertRuleId, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetAlertCondRule(ctx context.Context, req *alert.GetAlertCondRuleRequest) (*alert.GetAlertCondRuleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := f.repository.GetAlertCondRule(req.ProjectId, req.AlertConditionId, req.AlertRuleId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetAlertCondRuleResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertCondRuleResponse{AlertCondRule: convertAlertCondRule(data)}, nil
}

func (f *alertService) PutAlertCondRule(ctx context.Context, req *alert.PutAlertCondRuleRequest) (*alert.PutAlertCondRuleResponse, error) {
	if err := req.AlertCondRule.Validate(); err != nil {
		return nil, err
	}
	data := &model.AlertCondRule{
		AlertConditionID: req.AlertCondRule.AlertConditionId,
		AlertRuleID:      req.AlertCondRule.AlertRuleId,
		ProjectID:        req.AlertCondRule.ProjectId,
	}

	// Fiding upsert
	registerdData, err := f.repository.UpsertAlertCondRule(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertCondRuleResponse{AlertCondRule: convertAlertCondRule(registerdData)}, nil
}

func (f *alertService) DeleteAlertCondRule(ctx context.Context, req *alert.DeleteAlertCondRuleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteAlertCondRule(req.ProjectId, req.AlertConditionId, req.AlertRuleId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * Notification
 */

func (f *alertService) ListNotification(ctx context.Context, req *alert.ListNotificationRequest) (*alert.ListNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	converted := convertListNotificationRequest(req)
	list, err := f.repository.ListNotification(converted.ProjectId, converted.Type, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetNotification(ctx context.Context, req *alert.GetNotificationRequest) (*alert.GetNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetNotification(req.ProjectId, req.NotificationId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetNotificationResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetNotificationResponse{Notification: convertNotification(data)}, nil
}

func (f *alertService) PutNotification(ctx context.Context, req *alert.PutNotificationRequest) (*alert.PutNotificationResponse, error) {
	err := req.Notification.Validate()
	if err != nil {
		return nil, err
	}
	var existData *model.Notification
	if !zero.IsZeroVal(req.Notification.NotificationId) {
		existData, err = f.repository.GetNotification(req.ProjectId, req.Notification.NotificationId)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
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
			break
		default:
			appLogger.Warnf("This notification_type is unimprement. type: %v", existData.Type)
			break
		}
	}

	// Fiding upsert
	registerdData, err := f.repository.UpsertNotification(data)
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

func (f *alertService) DeleteNotification(ctx context.Context, req *alert.DeleteNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteNotification(req.ProjectId, req.NotificationId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f *alertService) TestNotification(ctx context.Context, req *alert.TestNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	notification, err := f.repository.GetNotification(req.ProjectId, req.NotificationId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &empty.Empty{}, nil
		}
		return nil, err
	}
	switch notification.Type {
	case "slack":
		err = sendSlackTestNotification(notification.NotifySetting)
		if err != nil {
			appLogger.Errorf("Error occured when sending test slack notification. err: %v", err)
			return nil, err
		}
		break
	default:
		appLogger.Warnf("This notification_type is unimprement. type: %v", notification.Type)
		break
	}
	return &empty.Empty{}, nil
}

/**
 * AlertCondNotification
 */

func (f *alertService) ListAlertCondNotification(ctx context.Context, req *alert.ListAlertCondNotificationRequest) (*alert.ListAlertCondNotificationResponse, error) {
	converted := convertListAlertCondNotificationRequest(req)
	list, err := f.repository.ListAlertCondNotification(converted.ProjectId, converted.AlertConditionId, converted.NotificationId, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetAlertCondNotification(ctx context.Context, req *alert.GetAlertCondNotificationRequest) (*alert.GetAlertCondNotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := f.repository.GetAlertCondNotification(req.ProjectId, req.AlertConditionId, req.NotificationId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetAlertCondNotificationResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertCondNotificationResponse{AlertCondNotification: convertAlertCondNotification(data)}, nil
}

func (f *alertService) PutAlertCondNotification(ctx context.Context, req *alert.PutAlertCondNotificationRequest) (*alert.PutAlertCondNotificationResponse, error) {
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
	registerdData, err := f.repository.UpsertAlertCondNotification(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertCondNotificationResponse{AlertCondNotification: convertAlertCondNotification(registerdData)}, nil
}

func (f *alertService) DeleteAlertCondNotification(ctx context.Context, req *alert.DeleteAlertCondNotificationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteAlertCondNotification(req.ProjectId, req.AlertConditionId, req.NotificationId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * Converter
 */

func convertAlertCondition(f *model.AlertCondition) *alert.AlertCondition {
	if f == nil {
		return &alert.AlertCondition{}
	}
	return &alert.AlertCondition{
		AlertConditionId: f.AlertConditionID,
		Description:      f.Description,
		Severity:         f.Severity,
		ProjectId:        f.ProjectID,
		AndOr:            f.AndOr,
		Enabled:          f.Enabled,
		CreatedAt:        f.CreatedAt.Unix(),
		UpdatedAt:        f.UpdatedAt.Unix(),
	}
}

func convertAlertRule(f *model.AlertRule) *alert.AlertRule {
	if f == nil {
		return &alert.AlertRule{}
	}
	return &alert.AlertRule{
		AlertRuleId:  f.AlertRuleID,
		Name:         f.Name,
		Score:        f.Score,
		ProjectId:    f.ProjectID,
		ResourceName: f.ResourceName,
		Tag:          f.Tag,
		FindingCnt:   f.FindingCnt,
		CreatedAt:    f.CreatedAt.Unix(),
		UpdatedAt:    f.UpdatedAt.Unix(),
	}
}

func convertAlertCondRule(f *model.AlertCondRule) *alert.AlertCondRule {
	if f == nil {
		return &alert.AlertCondRule{}
	}
	return &alert.AlertCondRule{
		AlertConditionId: f.AlertConditionID,
		AlertRuleId:      f.AlertRuleID,
		ProjectId:        f.ProjectID,
		CreatedAt:        f.CreatedAt.Unix(),
		UpdatedAt:        f.UpdatedAt.Unix(),
	}
}

func convertNotification(f *model.Notification) *alert.Notification {
	if f == nil {
		return &alert.Notification{}
	}
	maskingSetting, err := maskingNotifySetting(f.Type, f.NotifySetting)
	if err != nil {
		appLogger.Errorf("Failed to masking notify setting. %v", err)
		maskingSetting = f.NotifySetting
	}
	return &alert.Notification{
		NotificationId: f.NotificationID,
		Name:           f.Name,
		ProjectId:      f.ProjectID,
		Type:           f.Type,
		NotifySetting:  maskingSetting,
		CreatedAt:      f.CreatedAt.Unix(),
		UpdatedAt:      f.UpdatedAt.Unix(),
	}
}

func convertAlertCondNotification(f *model.AlertCondNotification) *alert.AlertCondNotification {
	if f == nil {
		return &alert.AlertCondNotification{}
	}
	return &alert.AlertCondNotification{
		AlertConditionId: f.AlertConditionID,
		NotificationId:   f.NotificationID,
		ProjectId:        f.ProjectID,
		CacheSecond:      f.CacheSecond,
		NotifiedAt:       f.NotifiedAt.Unix(),
		CreatedAt:        f.CreatedAt.Unix(),
		UpdatedAt:        f.UpdatedAt.Unix(),
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
