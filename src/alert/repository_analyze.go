package main

import (
	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (f *alertDB) ListAlertRuleByAlertConditionID(projectID, alertConditionID uint32) (*[]model.AlertRule, error) {
	query := `select * from alert_rule where alert_rule_id = any (select alert_rule_id from alert_cond_rule where project_id = ? and alert_condition_id = ?);`
	var params []interface{}
	params = append(params, projectID, alertConditionID)
	var data []model.AlertRule
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) ListNotificationByAlertConditionID(projectID, alertConditionID uint32) (*[]model.Notification, error) {
	query := `select * from notification where notification_id = any (select notification_id from alert_cond_notification where project_id = ? and alert_condition_id = ?);`
	var params []interface{}
	params = append(params, projectID, alertConditionID)
	var data []model.Notification
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) DeactivateAlert(data *model.Alert) error {
	if err := f.Master.Model(&model.Alert{}).Where("project_id = ? AND alert_id = ?", data.ProjectID, data.AlertID).Update("status", "DEACTIVE").Error; err != nil {
		return err
	}
	return nil
}

func (f *alertDB) GetAlertByAlertConditionIDStatus(projectID uint32, AlertConditionID uint32, status []string) (*model.Alert, error) {
	var data model.Alert
	if err := f.Slave.Where("project_id = ? AND alert_condition_id = ? AND status in (?)", projectID, AlertConditionID, status).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListFinding string = `
select
  * 
from
  finding f 
where
  f.project_id = ?
  and not exists(select * from pend_finding pf where pf.finding_id=f.finding_id)
`

func (f *alertDB) ListFinding(projectID uint32) (*[]model.Finding, error) {
	var data []model.Finding
	if err := f.Slave.Raw(selectListFinding, projectID).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) ListFindingTag(projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := f.Slave.Where("project_id = ? AND finding_id = ?", projectID, findingID).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) ListEnabledAlertCondition(projectID uint32, alertConditionID []uint32) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and enabled = ?`
	var params []interface{}
	params = append(params, projectID, true)
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id in (?)"
		params = append(params, alertConditionID)
	}
	var data []model.AlertCondition
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) ListDisabledAlertCondition(projectID uint32, alertConditionID []uint32) (*[]model.AlertCondition, error) {
	query := `select * from alert_condition where project_id = ? and enabled = ?`
	var params []interface{}
	params = append(params, projectID, false)
	if !zero.IsZeroVal(alertConditionID) {
		query += " and alert_condition_id in (?)"
		params = append(params, alertConditionID)
	}
	var data []model.AlertCondition
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) GetProject(projectID uint32) (*model.Project, error) {
	var data model.Project
	if err := f.Slave.Where("project_id = ?", projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
