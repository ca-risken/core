package main

import (
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
)

func (f *alertDB) ListAlertCondition(projectID uint32, severity []string, enabled bool, fromAt, toAt int64) (*[]model.AlertCondition, error) {
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
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) GetAlertCondition(projectID uint32, alertConditionID uint32) (*model.AlertCondition, error) {
	var data model.AlertCondition
	if err := f.Slave.Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) UpsertAlertCondition(data *model.AlertCondition) (*model.AlertCondition, error) {
	var retData model.AlertCondition
	update := alertConditionToMap(data)
	if err := f.Master.Where("project_id = ? AND alert_condition_id = ?", data.ProjectID, data.AlertConditionID).Assign(update).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (f *alertDB) DeleteAlertCondition(projectID uint32, alertConditionID uint32) error {
	if err := f.Master.Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).Delete(model.AlertCondition{}).Error; err != nil {
		return err
	}
	return nil
}

func (f *alertDB) ListAlertRule(projectID uint32, fromScore, toScore float32, fromAt, toAt int64) (*[]model.AlertRule, error) {
	query := `select * from alert_rule where project_id = ? and score between ? and ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, fromScore, toScore, time.Unix(fromAt, 0), time.Unix(toAt, 0))

	var data []model.AlertRule
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) GetAlertRule(projectID uint32, alertRuleID uint32) (*model.AlertRule, error) {
	var data model.AlertRule
	if err := f.Slave.Where("project_id = ? AND alert_rule_id = ?", projectID, alertRuleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) UpsertAlertRule(data *model.AlertRule) (*model.AlertRule, error) {
	var retData model.AlertRule
	update := alertRuleToMap(data)
	if err := f.Master.Where("project_id = ? AND alert_rule_id = ?", data.ProjectID, data.AlertRuleID).Assign(update).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (f *alertDB) DeleteAlertRule(projectID uint32, alertRuleID uint32) error {
	if err := f.Master.Where("project_id = ? AND alert_rule_id = ?", projectID, alertRuleID).Delete(model.AlertRule{}).Error; err != nil {
		return err
	}
	return nil
}

func (f *alertDB) ListAlertCondRule(projectID, alertConditionID, alertRuleID uint32, fromAt, toAt int64) (*[]model.AlertCondRule, error) {
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
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) GetAlertCondRule(projectID, alertConditionID, alertRuleID uint32) (*model.AlertCondRule, error) {
	var data model.AlertCondRule
	if err := f.Slave.Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", projectID, alertConditionID, alertRuleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) UpsertAlertCondRule(data *model.AlertCondRule) (*model.AlertCondRule, error) {
	var retData model.AlertCondRule
	if err := f.Master.Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", data.ProjectID, data.AlertConditionID, data.AlertRuleID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (f *alertDB) DeleteAlertCondRule(projectID, alertConditionID, alertRuleID uint32) error {
	if err := f.Master.Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", projectID, alertConditionID, alertRuleID).Delete(model.AlertCondRule{}).Error; err != nil {
		return err
	}
	return nil
}

func (f *alertDB) ListNotification(projectID uint32, notifyType string, fromAt, toAt int64) (*[]model.Notification, error) {
	query := `select * from notification where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(notifyType) {
		query += " and type = ?"
		params = append(params, notifyType)
	}
	var data []model.Notification
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) GetNotification(projectID uint32, NotificationID uint32) (*model.Notification, error) {
	var data model.Notification
	if err := f.Slave.Where("project_id = ? AND notification_id = ?", projectID, NotificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) UpsertNotification(data *model.Notification) (*model.Notification, error) {
	var retData model.Notification
	if err := f.Master.Where("project_id = ? AND notification_id = ?", data.ProjectID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (f *alertDB) DeleteNotification(projectID uint32, NotificationID uint32) error {
	if err := f.Master.Where("project_id = ? AND notification_id = ?", projectID, NotificationID).Delete(model.Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func (f *alertDB) ListAlertCondNotification(projectID, alertConditionID, notificationID uint32, fromAt, toAt int64) (*[]model.AlertCondNotification, error) {
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
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) GetAlertCondNotification(projectID, alertConditionID, notificationID uint32) (*model.AlertCondNotification, error) {
	var data model.AlertCondNotification
	if err := f.Slave.Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", projectID, alertConditionID, notificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) UpsertAlertCondNotification(data *model.AlertCondNotification) (*model.AlertCondNotification, error) {
	var retData model.AlertCondNotification
	if err := f.Master.Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", data.ProjectID, data.AlertConditionID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (f *alertDB) DeleteAlertCondNotification(projectID, alertConditionID, notificationID uint32) error {
	if err := f.Master.Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", projectID, alertConditionID, notificationID).Delete(model.AlertCondNotification{}).Error; err != nil {
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
