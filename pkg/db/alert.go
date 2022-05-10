package db

import (
	"context"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
)

type AlertRepository interface {
	// Alert
	ListAlert(context.Context, uint32, []string, []string, string, int64, int64) (*[]model.Alert, error)
	GetAlert(context.Context, uint32, uint32) (*model.Alert, error)
	UpsertAlert(context.Context, *model.Alert) (*model.Alert, error)
	DeleteAlert(context.Context, uint32, uint32) error
	ListAlertHistory(context.Context, uint32, uint32, []string, []string, int64, int64) (*[]model.AlertHistory, error)
	GetAlertHistory(context.Context, uint32, uint32) (*model.AlertHistory, error)
	UpsertAlertHistory(context.Context, *model.AlertHistory) (*model.AlertHistory, error)
	DeleteAlertHistory(context.Context, uint32, uint32) error
	ListRelAlertFinding(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.RelAlertFinding, error)
	GetRelAlertFinding(context.Context, uint32, uint32, uint32) (*model.RelAlertFinding, error)
	UpsertRelAlertFinding(context.Context, *model.RelAlertFinding) (*model.RelAlertFinding, error)
	DeleteRelAlertFinding(context.Context, uint32, uint32, uint32) error
	ListAlertCondition(context.Context, uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error)
	GetAlertCondition(context.Context, uint32, uint32) (*model.AlertCondition, error)
	UpsertAlertCondition(context.Context, *model.AlertCondition) (*model.AlertCondition, error)
	DeleteAlertCondition(context.Context, uint32, uint32) error
	ListAlertRule(context.Context, uint32, float32, float32, int64, int64) (*[]model.AlertRule, error)
	GetAlertRule(context.Context, uint32, uint32) (*model.AlertRule, error)
	UpsertAlertRule(context.Context, *model.AlertRule) (*model.AlertRule, error)
	DeleteAlertRule(context.Context, uint32, uint32) error
	ListAlertCondRule(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error)
	GetAlertCondRule(context.Context, uint32, uint32, uint32) (*model.AlertCondRule, error)
	UpsertAlertCondRule(context.Context, *model.AlertCondRule) (*model.AlertCondRule, error)
	DeleteAlertCondRule(context.Context, uint32, uint32, uint32) error
	ListNotification(context.Context, uint32, string, int64, int64) (*[]model.Notification, error)
	GetNotification(context.Context, uint32, uint32) (*model.Notification, error)
	UpsertNotification(context.Context, *model.Notification) (*model.Notification, error)
	DeleteNotification(context.Context, uint32, uint32) error
	ListAlertCondNotification(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error)
	GetAlertCondNotification(context.Context, uint32, uint32, uint32) (*model.AlertCondNotification, error)
	UpsertAlertCondNotification(context.Context, *model.AlertCondNotification) (*model.AlertCondNotification, error)
	DeleteAlertCondNotification(context.Context, uint32, uint32, uint32) error

	// forAnalyze
	ListAlertRuleByAlertConditionID(context.Context, uint32, uint32) (*[]model.AlertRule, error)
	DeactivateAlert(context.Context, *model.Alert) error
	GetAlertByAlertConditionIDStatus(context.Context, uint32, uint32, []string) (*model.Alert, error)
	ListEnabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error)
	ListDisabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error)
}

var _ AlertRepository = (*Client)(nil)

func (c *Client) ListAlert(ctx context.Context, projectID uint32, status []string, severity []string, description string, fromAt, toAt int64) (*[]model.Alert, error) {
	query := `select * from alert where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(severity) {
		query += " and severity in (?)"
		params = append(params, severity)
	}
	if !zero.IsZeroVal(status) {
		query += " and status in (?)"
		params = append(params, status)
	}
	if !zero.IsZeroVal(description) {
		query += " and description = ?"
		params = append(params, description)
	}
	var data []model.Alert
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetAlert(ctx context.Context, projectID uint32, alertID uint32) (*model.Alert, error) {
	var data model.Alert
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_id = ?", projectID, alertID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertAlert(ctx context.Context, data *model.Alert) (*model.Alert, error) {
	var retData model.Alert
	c.logger.Info("upsertAlert:", data)
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_id = ?", data.ProjectID, data.AlertID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	c.logger.Info(retData)
	return &retData, nil
}

func (c *Client) DeleteAlert(ctx context.Context, projectID uint32, alertID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_id = ?", projectID, alertID).Delete(model.Alert{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAlertHistory(ctx context.Context, projectID, alertID uint32, HistoryType, severity []string, fromAt, toAt int64) (*[]model.AlertHistory, error) {
	query := `select * from alert_history where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(alertID) {
		query += " and alert_id = ?"
		params = append(params, alertID)
	}
	if !zero.IsZeroVal(HistoryType) {
		query += " and history_type in (?)"
		params = append(params, HistoryType)
	}
	if !zero.IsZeroVal(severity) {
		query += " and severity in (?)"
		params = append(params, severity)
	}
	query += " order by alert_history_id desc"
	var data []model.AlertHistory
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetAlertHistory(ctx context.Context, projectID uint32, alertHistoryID uint32) (*model.AlertHistory, error) {
	var data model.AlertHistory
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_history_id = ?", projectID, alertHistoryID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertAlertHistory(ctx context.Context, data *model.AlertHistory) (*model.AlertHistory, error) {
	var retData model.AlertHistory
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_history_id = ?", data.ProjectID, data.AlertHistoryID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteAlertHistory(ctx context.Context, projectID uint32, alertHistoryID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_history_id = ?", projectID, alertHistoryID).Delete(model.AlertHistory{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListRelAlertFinding(ctx context.Context, projectID, alertID, findingID uint32, fromAt, toAt int64) (*[]model.RelAlertFinding, error) {
	query := `select * from rel_alert_finding where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(alertID) {
		query += " and alert_id = ?"
		params = append(params, alertID)
	}
	if !zero.IsZeroVal(findingID) {
		query += " and finding_id in (?)"
		params = append(params, findingID)
	}
	var data []model.RelAlertFinding
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetRelAlertFinding(ctx context.Context, projectID, alertID, findingID uint32) (*model.RelAlertFinding, error) {
	var data model.RelAlertFinding
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_id = ? AND finding_id = ?", projectID, alertID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertRelAlertFinding(ctx context.Context, data *model.RelAlertFinding) (*model.RelAlertFinding, error) {
	var retData model.RelAlertFinding
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_id = ? AND finding_id = ?", data.ProjectID, data.AlertID, data.FindingID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteRelAlertFinding(ctx context.Context, projectID, alertID, findingID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_id = ? AND finding_id = ?", projectID, alertID, findingID).Delete(model.RelAlertFinding{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAlertCondition(ctx context.Context, projectID uint32, severity []string, enabled bool, fromAt, toAt int64) (*[]model.AlertCondition, error) {
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetAlertCondition(ctx context.Context, projectID uint32, alertConditionID uint32) (*model.AlertCondition, error) {
	var data model.AlertCondition
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertAlertCondition(ctx context.Context, data *model.AlertCondition) (*model.AlertCondition, error) {
	var retData model.AlertCondition
	update := alertConditionToMap(data)
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ?", data.ProjectID, data.AlertConditionID).Assign(update).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteAlertCondition(ctx context.Context, projectID uint32, alertConditionID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).Delete(model.AlertCondition{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAlertRule(ctx context.Context, projectID uint32, fromScore, toScore float32, fromAt, toAt int64) (*[]model.AlertRule, error) {
	query := `select * from alert_rule where project_id = ? and score between ? and ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, fromScore, toScore, time.Unix(fromAt, 0), time.Unix(toAt, 0))

	var data []model.AlertRule
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetAlertRule(ctx context.Context, projectID uint32, alertRuleID uint32) (*model.AlertRule, error) {
	var data model.AlertRule
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_rule_id = ?", projectID, alertRuleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertAlertRule(ctx context.Context, data *model.AlertRule) (*model.AlertRule, error) {
	var retData model.AlertRule
	update := alertRuleToMap(data)
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_rule_id = ?", data.ProjectID, data.AlertRuleID).Assign(update).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteAlertRule(ctx context.Context, projectID uint32, alertRuleID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_rule_id = ?", projectID, alertRuleID).Delete(model.AlertRule{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAlertCondRule(ctx context.Context, projectID, alertConditionID, alertRuleID uint32, fromAt, toAt int64) (*[]model.AlertCondRule, error) {
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetAlertCondRule(ctx context.Context, projectID, alertConditionID, alertRuleID uint32) (*model.AlertCondRule, error) {
	var data model.AlertCondRule
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", projectID, alertConditionID, alertRuleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertAlertCondRule(ctx context.Context, data *model.AlertCondRule) (*model.AlertCondRule, error) {
	var retData model.AlertCondRule
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", data.ProjectID, data.AlertConditionID, data.AlertRuleID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteAlertCondRule(ctx context.Context, projectID, alertConditionID, alertRuleID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND alert_rule_id = ?", projectID, alertConditionID, alertRuleID).Delete(model.AlertCondRule{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListNotification(ctx context.Context, projectID uint32, notifyType string, fromAt, toAt int64) (*[]model.Notification, error) {
	query := `select * from notification where project_id = ? and updated_at between ? and ?`
	var params []interface{}
	params = append(params, projectID, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(notifyType) {
		query += " and type = ?"
		params = append(params, notifyType)
	}
	var data []model.Notification
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetNotification(ctx context.Context, projectID uint32, NotificationID uint32) (*model.Notification, error) {
	var data model.Notification
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND notification_id = ?", projectID, NotificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertNotification(ctx context.Context, data *model.Notification) (*model.Notification, error) {
	var retData model.Notification
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND notification_id = ?", data.ProjectID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteNotification(ctx context.Context, projectID uint32, NotificationID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND notification_id = ?", projectID, NotificationID).Delete(model.Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAlertCondNotification(ctx context.Context, projectID, alertConditionID, notificationID uint32, fromAt, toAt int64) (*[]model.AlertCondNotification, error) {
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetAlertCondNotification(ctx context.Context, projectID, alertConditionID, notificationID uint32) (*model.AlertCondNotification, error) {
	var data model.AlertCondNotification
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", projectID, alertConditionID, notificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertAlertCondNotification(ctx context.Context, data *model.AlertCondNotification) (*model.AlertCondNotification, error) {
	var retData model.AlertCondNotification
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", data.ProjectID, data.AlertConditionID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (c *Client) DeleteAlertCondNotification(ctx context.Context, projectID, alertConditionID, notificationID uint32) error {
	if err := c.Master.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND notification_id = ?", projectID, alertConditionID, notificationID).Delete(model.AlertCondNotification{}).Error; err != nil {
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

func (c *Client) ListAlertRuleByAlertConditionID(ctx context.Context, projectID, alertConditionID uint32) (*[]model.AlertRule, error) {
	query := `select * from alert_rule where alert_rule_id = any (select alert_rule_id from alert_cond_rule where project_id = ? and alert_condition_id = ?);`
	var params []interface{}
	params = append(params, projectID, alertConditionID)
	var data []model.AlertRule
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) DeactivateAlert(ctx context.Context, data *model.Alert) error {
	if err := c.Master.WithContext(ctx).Model(&model.Alert{}).Where("project_id = ? AND alert_id = ?", data.ProjectID, data.AlertID).Update("status", "DEACTIVE").Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAlertByAlertConditionIDStatus(ctx context.Context, projectID uint32, AlertConditionID uint32, status []string) (*model.Alert, error) {
	var data model.Alert
	if err := c.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND status in (?)", projectID, AlertConditionID, status).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListEnabledAlertCondition(ctx context.Context, projectID uint32, alertConditionID []uint32) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and enabled = ?`
	var params []interface{}
	params = append(params, projectID, true)
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id in (?)"
		params = append(params, alertConditionID)
	}
	var data []model.AlertCondition
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListDisabledAlertCondition(ctx context.Context, projectID uint32, alertConditionID []uint32) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and enabled = ?`
	var params []interface{}
	params = append(params, projectID, false)
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id in (?)"
		params = append(params, alertConditionID)
	}
	var data []model.AlertCondition
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
