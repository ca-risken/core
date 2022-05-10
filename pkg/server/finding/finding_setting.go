package finding

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (f *FindingService) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*finding.ListFindingSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := f.repository.ListFindingSetting(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (f *FindingService) GetFindingSetting(ctx context.Context, req *finding.GetFindingSettingRequest) (*finding.GetFindingSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetFindingSetting(ctx, req.ProjectId, req.FindingSettingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &finding.GetFindingSettingResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetFindingSettingResponse{FindingSetting: convertFindingSetting(data)}, nil
}

func (f *FindingService) PutFindingSetting(ctx context.Context, req *finding.PutFindingSettingRequest) (*finding.PutFindingSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := f.repository.UpsertFindingSetting(ctx, &model.FindingSetting{
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

func (f *FindingService) DeleteFindingSetting(ctx context.Context, req *finding.DeleteFindingSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := f.repository.DeleteFindingSetting(ctx, req.ProjectId, req.FindingSettingId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

type findingSetting struct {
	ScoreCoefficient float32 `json:"score_coefficient,omitempty"`
	// Tag              []string `json:"tag,omitempty"`
}

func (f *FindingService) getFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*findingSetting, error) {
	fs, err := f.repository.GetFindingSettingByResource(ctx, projectID, resourceName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &findingSetting{}, nil
	} else if err != nil {
		return nil, err
	}
	if fs.Status != "ACTIVE" {
		return &findingSetting{}, nil
	}
	var setting findingSetting
	if err := json.Unmarshal([]byte(fs.Setting), &setting); err != nil {
		appLogger.Warnf("Failed to unmarshal finding setting JSON, projectID=%d, resourceName=%s, err=%+v", projectID, resourceName, err)
		return &findingSetting{}, nil
	}
	return &setting, nil
}
