package main

import (
	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
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
	query := `select * from notification where alert_rule_id = any (select alert_rule_id from alert_cond_rule where project_id = ? and alert_condition_id = ?);`
	var params []interface{}
	params = append(params, projectID, alertConditionID)
	var data []model.Notification
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) DeactivateAlert(data *model.Alert) error {
	if err := f.Master.Model(&model.Alert{}).Where("project_id = ? AND alert_id = ?", data.ProjectID, data.AlertID).Update("activated", false).Error; err != nil {
		return err
	}
	return nil
}

func (f *alertDB) GetAlertByAlertConditionIDWithActivated(projectID uint32, AlertConditionID uint32, activated bool) (*model.Alert, error) {
	var data model.Alert
	if err := f.Slave.Where("project_id = ? AND alert_condition_id = ? AND activated = ?", projectID, AlertConditionID, activated).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *alertDB) ListFinding(projectID uint32) (*[]model.Finding, error) {
	var data []model.Finding
	if err := f.Slave.Where("project_id = ?", projectID).Find(&data).Error; err != nil {
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
