package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (f *findingService) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*finding.ListFindingSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := f.repository.ListFindingSetting(req)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.ListFindingSettingResponse{}, nil
		}
		return nil, err
	}
	data := finding.ListFindingSettingResponse{}
	for _, d := range *list {
		data.FindingSetting = append(data.FindingSetting, convertFindingSetting(&d))
	}
	return &data, nil
}

func convertFindingSetting(data *model.FindingSetting) *finding.FindingSetting {
	if data == nil {
		return &finding.FindingSetting{}
	}
	return &finding.FindingSetting{
		FindingSettingId: data.FindingSettingID,
		ProjectId:        data.ProjectID,
		ResourceName:     data.ResourceName,
		Setting:          data.Setting,
		Status:           getStatus(data.Status),
		CreatedAt:        data.CreatedAt.Unix(),
		UpdatedAt:        data.CreatedAt.Unix(),
	}
}

func getStatus(dbStatus string) finding.FindingSettingStatus {
	statusKey := strings.ToUpper(dbStatus)
	switch statusKey {
	case "ACTIVE":
		return finding.FindingSettingStatus_SETTING_ACTIVE
	case "DEACTIVE":
		return finding.FindingSettingStatus_SETTING_DEACTIVE
	default:
		return finding.FindingSettingStatus_SETTING_UNKNOWN
	}
}

func (f *findingService) GetFindingSetting(ctx context.Context, req *finding.GetFindingSettingRequest) (*finding.GetFindingSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetFindingSetting(req.ProjectId, req.FindingSettingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.GetFindingSettingResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetFindingSettingResponse{FindingSetting: convertFindingSetting(data)}, nil
}

func (f *findingService) PutFindingSetting(ctx context.Context, req *finding.PutFindingSettingRequest) (*finding.PutFindingSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := f.repository.UpsertFindingSetting(&model.FindingSetting{
		ProjectID:    req.FindingSetting.ProjectId,
		ResourceName: req.FindingSetting.ResourceName,
		Status:       getStatusString(req.FindingSetting.Status),
		Setting:      req.FindingSetting.Setting,
	})
	if err != nil {
		return nil, err
	}
	return &finding.PutFindingSettingResponse{FindingSetting: convertFindingSetting(registerd)}, nil
}

func getStatusString(status finding.FindingSettingStatus) string {
	switch status {
	case finding.FindingSettingStatus_SETTING_ACTIVE:
		return "ACTIVE"
	case finding.FindingSettingStatus_SETTING_DEACTIVE:
		return "DEACTIVE"
	default:
		return "UNKNOWN"
	}
}

func (f *findingService) DeleteFindingSetting(ctx context.Context, req *finding.DeleteFindingSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := f.repository.DeleteFindingSetting(req.ProjectId, req.FindingSettingId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

type findingSetting struct {
	ScoreCoefficient float32 `json:"score_coefficient,omitempty"`
	// Tag              []string `json:"tag,omitempty"`
}

func (f *findingService) getFindingSettingByResource(projectID uint32, resourceName string) (*findingSetting, error) {
	fs, err := f.repository.GetFindingSettingByResource(projectID, resourceName)
	if gorm.IsRecordNotFoundError(err) {
		return &findingSetting{}, nil
	} else if err != nil {
		return nil, err
	}
	var setting findingSetting
	if err := json.Unmarshal([]byte(fs.Setting), &setting); err != nil {
		appLogger.Warnf("Failed to unmarshal finding setting JSON, projectID=%d, resourceName=%s, err=%+v", projectID, resourceName, err)
		return &findingSetting{}, nil
	}
	return &setting, nil
}