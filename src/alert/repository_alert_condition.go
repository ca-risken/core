package main

import (
	"context"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
)

func (a *alertDB) ListAlertCondition(ctx context.Context, projectID uint32, severity []string, enabled bool, fromAt, toAt int64) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(severity) {
		query += " and severity in (?)"
		params = append(params, severity)
	}
	if !zero.IsZeroVal(enabled) {
		query += " and enabled = ?"
		params = append(params, enabled)
	}
	var data []model.AlertCondition
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlertCondition(ctx context.Context, projectID uint32, alertConditionID uint32) (*model.AlertCondition, error) {
	var data model.AlertCondition
	if err := a.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertAlertCondition(ctx context.Context, data *model.AlertCondition) (*model.AlertCondition, error) {
	var retData model.AlertCondition
	update := alertConditionToMap(data)
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ?", data.ProjectID, data.AlertConditionID).Assign(update).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteAlertCondition(ctx context.Context, projectID uint32, alertConditionID uint32) error {
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).Delete(model.AlertCondition{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) ListAlertRule(ctx context.Context, projectID uint32, fromScore, toScore float32, fromAt, toAt int64) (*[]model.AlertRule, error) {
	query := `select * from alert_rule where project_id = ? and score between ? and ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, fromScore, toScore, time.Unix(fromAt, 0), time.Unix(toAt, 0))

	var data []model.AlertRule
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlertRule(ctx context.Context, projectID uint32, alertRuleID uint32) (*model.AlertRule, error) {
	var data model.AlertRule
	if err := a.Slave.WithContext(ctx).Where("project_id = ? AND alert_rule_id = ?", projectID, alertRuleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertAlertRule(ctx context.Context, data *model.AlertRule) (*model.AlertRule, error) {
	var retData model.AlertRule
	update := alertRuleToMap(data)
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_rule_id = ?", data.ProjectID, data.AlertRuleID).Assign(update).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteAlertRule(ctx context.Context, projectID uint32, alertRuleID uint32) error {
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_rule_id = ?", projectID, alertRuleID).Delete(model.AlertRule{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) ListAlertCondRule(ctx context.Context, projectID, alertConditionID, alertRuleID uint32, fromAt, toAt int64) (*[]model.AlertCondRule, error) {
	query := `select * from alert_cond_rule where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id = ?"
		params = append(params, alertConditionID)
	}
	if !zero.IsZeroVal(alertRuleID) {
		query += " and alert_rule_id in (?)"
		params = append(params, alertRuleID)
	}
	var data []model.AlertCondRule
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlertCondRule(ctx context.Context, projectID, alertConditionID, alertRuleID uint32) (*model.AlertCondRule, error) {
	var data model.AlertCondRule
	if err := a.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", projectID, alertConditionID, alertRuleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertAlertCondRule(ctx context.Context, data *model.AlertCondRule) (*model.AlertCondRule, error) {
	var retData model.AlertCondRule
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", data.ProjectID, data.AlertConditionID, data.AlertRuleID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteAlertCondRule(ctx context.Context, projectID, alertConditionID, alertRuleID uint32) error {
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", projectID, alertConditionID, alertRuleID).Delete(model.AlertCondRule{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) ListNotification(ctx context.Context, projectID uint32, notifyType string, fromAt, toAt int64) (*[]model.Notification, error) {
	query := `select * from notification where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(notifyType) {
		query += " and type = ?"
		params = append(params, notifyType)
	}
	var data []model.Notification
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetNotification(ctx context.Context, projectID uint32, NotificationID uint32) (*model.Notification, error) {
	var data model.Notification
	if err := a.Slave.WithContext(ctx).Where("project_id = ? AND notification_id = ?", projectID, NotificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertNotification(ctx context.Context, data *model.Notification) (*model.Notification, error) {
	var retData model.Notification
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND notification_id = ?", data.ProjectID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteNotification(ctx context.Context, projectID uint32, NotificationID uint32) error {
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND notification_id = ?", projectID, NotificationID).Delete(model.Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) ListAlertCondNotification(ctx context.Context, projectID, alertConditionID, notificationID uint32, fromAt, toAt int64) (*[]model.AlertCondNotification, error) {
	query := `select * from alert_cond_notification where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id = ?"
		params = append(params, alertConditionID)
	}
	if !zero.IsZeroVal(notificationID) {
		query += " and notification_id in (?)"
		params = append(params, notificationID)
	}
	var data []model.AlertCondNotification
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlertCondNotification(ctx context.Context, projectID, alertConditionID, notificationID uint32) (*model.AlertCondNotification, error) {
	var data model.AlertCondNotification
	if err := a.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", projectID, alertConditionID, notificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertAlertCondNotification(ctx context.Context, data *model.AlertCondNotification) (*model.AlertCondNotification, error) {
	var retData model.AlertCondNotification
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", data.ProjectID, data.AlertConditionID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteAlertCondNotification(ctx context.Context, projectID, alertConditionID, notificationID uint32) error {
	if err := a.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", projectID, alertConditionID, notificationID).Delete(model.AlertCondNotification{}).Error; err != nil {
		return err
	}
	return nil
}

func alertConditionToMap(alertCondition *model.AlertCondition) map[string]interface{} {
	return map[string]interface{}{
		"alert_condition_id": alertCondition.AlertConditionID,
		"description":        alertCondition.Description,
		"severity":           alertCondition.Severity,
		"project_id":         alertCondition.ProjectID,
		"and_or":             alertCondition.AndOr,
		"enabled":            alertCondition.Enabled,
	}
}

func alertRuleToMap(alertRule *model.AlertRule) map[string]interface{} {
	return map[string]interface{}{
		"alert_rule_id": alertRule.AlertRuleID,
		"name":          alertRule.Name,
		"project_id":    alertRule.ProjectID,
		"score":         alertRule.Score,
		"resource_name": alertRule.ResourceName,
		"tag":           alertRule.Tag,
		"finding_cnt":   alertRule.FindingCnt,
	}
}
