package main

import (
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
)

func (a *alertDB) ListAlert(projectID uint32, status []string, severity []string, description string, fromAt, toAt int64) (*[]model.Alert, error) {
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
	if err := a.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlert(projectID uint32, alertID uint32) (*model.Alert, error) {
	var data model.Alert
	if err := a.Slave.Where("project_id = ? AND alert_id = ?", projectID, alertID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlertByAlertConditionID(projectID uint32, alertConditionID uint32) (*model.Alert, error) {
	var data model.Alert
	if err := a.Slave.Where("project_id = ? AND alert_condition_id = ?", projectID, alertConditionID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertAlert(data *model.Alert) (*model.Alert, error) {
	var retData model.Alert
	appLogger.Info("upsertAlert:", data)
	if err := a.Master.Where("project_id = ? AND alert_id = ?", data.ProjectID, data.AlertID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	appLogger.Info(retData)
	return &retData, nil
}

func (a *alertDB) DeleteAlert(projectID uint32, alertID uint32) error {
	if err := a.Master.Where("project_id = ? AND alert_id = ?", projectID, alertID).Delete(model.Alert{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) ListAlertHistory(projectID, alertID uint32, HistoryType, severity []string, fromAt, toAt int64) (*[]model.AlertHistory, error) {
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
	if err := a.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetAlertHistory(projectID uint32, alertHistoryID uint32) (*model.AlertHistory, error) {
	var data model.AlertHistory
	if err := a.Slave.Where("project_id = ? AND alert_history_id = ?", projectID, alertHistoryID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertAlertHistory(data *model.AlertHistory) (*model.AlertHistory, error) {
	var retData model.AlertHistory
	if err := a.Master.Where("project_id = ? AND alert_history_id = ?", data.ProjectID, data.AlertHistoryID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteAlertHistory(projectID uint32, alertHistoryID uint32) error {
	if err := a.Master.Where("project_id = ? AND alert_history_id = ?", projectID, alertHistoryID).Delete(model.AlertHistory{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *alertDB) ListRelAlertFinding(projectID, alertID, findingID uint32, fromAt, toAt int64) (*[]model.RelAlertFinding, error) {
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
	if err := a.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) GetRelAlertFinding(projectID, alertID, findingID uint32) (*model.RelAlertFinding, error) {
	var data model.RelAlertFinding
	if err := a.Slave.Where("project_id = ? AND alert_id = ? AND finding_id = ?", projectID, alertID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (a *alertDB) UpsertRelAlertFinding(data *model.RelAlertFinding) (*model.RelAlertFinding, error) {
	var retData model.RelAlertFinding
	if err := a.Master.Where("project_id = ? AND alert_id = ? AND finding_id = ?", data.ProjectID, data.AlertID, data.FindingID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

func (a *alertDB) DeleteRelAlertFinding(projectID, alertID, findingID uint32) error {
	if err := a.Master.Where("project_id = ? AND alert_id = ? AND finding_id = ?", projectID, alertID, findingID).Delete(model.RelAlertFinding{}).Error; err != nil {
		return err
	}
	return nil
}
