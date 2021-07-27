package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func (f *findingDB) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	var param []interface{}
	query := "select * from finding_setting where project_id=?"
	param = append(param, req.ProjectId)
	if req.Status != finding.FindingSettingStatus_SETTING_UNKNOWN {
		query += " and status=?"
		param = append(param, getStatusString(req.Status))
	}
	var data []model.FindingSetting
	if err := f.Slave.WithContext(ctx).Raw(query, param...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingSetting = `select * from finding_setting where project_id=? and finding_setting_id=?`

func (f *findingDB) GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	var data model.FindingSetting
	if err := f.Slave.WithContext(ctx).Raw(selectGetFindingSetting, projectID, findingSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingSettingByResource = `select * from finding_setting where project_id=? and resource_name=?`

func (f *findingDB) GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error) {
	var data model.FindingSetting
	if err := f.Master.WithContext(ctx).Raw(selectGetFindingSettingByResource, projectID, resourceName).First(&data).Error; err != nil {
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

func (f *findingDB) UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error) {
	var retData model.FindingSetting
	if err := f.Master.WithContext(ctx).Where("project_id=? AND resource_name=?", data.ProjectID, data.ResourceName).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	appLogger.Info(retData)
	return &retData, nil
}

const deleteDeleteFindingSetting = `delete from finding_setting where project_id = ? and finding_setting_id = ?`

func (f *findingDB) DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error {
	return f.Master.WithContext(ctx).Exec(deleteDeleteFindingSetting, projectID, findingSettingID).Error
}
