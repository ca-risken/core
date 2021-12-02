package main

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
)

func (a *alertDB) ListAlertRuleByAlertConditionID(ctx context.Context, projectID, alertConditionID uint32) (*[]model.AlertRule, error) {
	query := `select * from alert_rule where alert_rule_id = any (select alert_rule_id from alert_cond_rule where project_id = ? and alert_condition_id = ?);`
	var params []interface{}
	params = append(params, projectID, alertConditionID)
	var data []model.AlertRule
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) DeactivateAlert(ctx context.Context, data *model.Alert) error {
	if err := a.Master.WithContext(ctx).Model(&model.Alert{}).Where("project_id = ? AND alert_id = ?", data.ProjectID, data.AlertID).Update("status", "DEACTIVE").Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) GetAlertByAlertConditionIDStatus(ctx context.Context, projectID uint32, AlertConditionID uint32, status []string) (*model.Alert, error) {
	var data model.Alert
	if err := a.Slave.WithContext(ctx).Where("project_id = ? AND alert_condition_id = ? AND status in (?)", projectID, AlertConditionID, status).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) ListEnabledAlertCondition(ctx context.Context, projectID uint32, alertConditionID []uint32) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and enabled = ?`
	var params []interface{}
	params = append(params, projectID, true)
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id in (?)"
		params = append(params, alertConditionID)
	}
	var data []model.AlertCondition
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) ListDisabledAlertCondition(ctx context.Context, projectID uint32, alertConditionID []uint32) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and enabled = ?`
	var params []interface{}
	params = append(params, projectID, false)
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id in (?)"
		params = append(params, alertConditionID)
	}
	var data []model.AlertCondition
	if err := a.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetProject(ctx context.Context, projectID uint32) (*model.Project, error) {
	var data model.Project
	if err := a.Slave.WithContext(ctx).Where("project_id = ?", projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
