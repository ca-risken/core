package main

import (
	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	_ "github.com/go-sql-driver/mysql"
)

func (f *findingDB) ListFindingSetting(req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	var param []interface{}
	query := "select * from finding_setting where project_id=?"
	param = append(param, req.ProjectId)
	if req.Status != finding.FindingSettingStatus_SETTING_UNKNOWN {
		query += " and status=?"
		param = append(param, req.Status.String())
	}
	var data []model.FindingSetting
	if err := f.Slave.Raw(query, param...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingSetting = `select * from finding_setting where project_id=? and finding_setting_id=?`

func (f *findingDB) GetFindingSetting(projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	var data model.FindingSetting
	if err := f.Slave.Raw(selectGetFindingSetting, projectID, findingSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingSettingByResource = `select * from finding_setting where project_id=? and resource_name=?`

func (f *findingDB) GetFindingSettingByResource(projectID uint32, resourceName string) (*model.FindingSetting, error) {
	var data model.FindingSetting
	if err := f.Master.Raw(selectGetFindingSettingByResource, projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertFindingSetting = `
INSERT INTO finding_setting
  (project_id, resource_name, status, setting)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  status=VALUES(status),
  setting=VALUES(setting)
`

func (f *findingDB) UpsertFindingSetting(data *model.FindingSetting) (*model.FindingSetting, error) {
	var retData model.FindingSetting
	if err := f.Master.Where("project_id=? AND resource_name=?", data.ProjectID, data.ResourceName).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	appLogger.Info(retData)
	return &retData, nil
}

const deleteDeleteFindingSetting = `delete from finding_setting where project_id = ? and finding_setting_id = ?`

func (f *findingDB) DeleteFindingSetting(projectID uint32, findingSettingID uint32) error {
	return f.Master.Exec(deleteDeleteFindingSetting, projectID, findingSettingID).Error
}
