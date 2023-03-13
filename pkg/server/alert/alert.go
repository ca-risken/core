package alert

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

const (
	alertHistoryTypeCreated = "created"
	alertHistoryTypeUpdated = "updated"
	alertHistoryTypeDeleted = "deleted"
)

func (a *AlertService) ListAlert(ctx context.Context, req *alert.ListAlertRequest) (*alert.ListAlertResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	converted := convertListAlertRequest(req)
	list, err := a.repository.ListAlert(ctx, converted.ProjectId, getStrings(converted.Status), converted.Severity, converted.Description, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListAlertResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListAlertResponse{}
	for _, d := range *list {
		data.Alert = append(data.Alert, convertAlert(&d))
	}
	return &data, nil
}

func convertListAlertRequest(req *alert.ListAlertRequest) *alert.ListAlertRequest {
	converted := alert.ListAlertRequest{
		ProjectId:   req.ProjectId,
		Status:      req.Status,
		Severity:    req.Severity,
		Description: req.Description,
		FromAt:      req.FromAt,
		ToAt:        req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetAlert(ctx context.Context, req *alert.GetAlertRequest) (*alert.GetAlertResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetAlert(ctx, req.ProjectId, req.AlertId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetAlertResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertResponse{Alert: convertAlert(data)}, nil
}

func (a *AlertService) PutAlert(ctx context.Context, req *alert.PutAlertRequest) (*alert.PutAlertResponse, error) {
	if err := req.Alert.Validate(); err != nil {
		return nil, err
	}
	var alertID uint32
	alertHistoryStatus := alertHistoryTypeCreated
	// AlertIdのパラメータがリクエストに存在する場合、レコードの存在チェック
	// 存在しなければエラー終了
	if !zero.IsZeroVal(req.Alert.AlertId) {
		savedData, err := a.repository.GetAlert(ctx, req.Alert.ProjectId, req.Alert.AlertId)
		if err != nil {
			return nil, err
		}
		alertID = savedData.AlertID
		alertHistoryStatus = alertHistoryTypeUpdated
	}

	data := &model.Alert{
		AlertID:          alertID,
		AlertConditionID: req.Alert.AlertConditionId,
		Description:      req.Alert.Description,
		Severity:         req.Alert.Severity,
		ProjectID:        req.Alert.ProjectId,
		Status:           req.Alert.Status.String(),
	}

	// Fiding upsert
	registeredData, err := a.repository.UpsertAlert(ctx, data)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	// list RelAlertFinding
	relAlertFindings, err := a.repository.ListRelAlertFinding(ctx, registeredData.ProjectID, registeredData.AlertID, 0, 0, now)
	if err != nil {
		a.logger.Errorf(ctx, "Failed listRelAlertFinding when PutAlert. err: %v", err)
		return &alert.PutAlertResponse{Alert: convertAlert(registeredData)}, err
	}
	findingIDs := []uint64{}
	for _, relAlertFinding := range *relAlertFindings {
		findingIDs = append(findingIDs, uint64(relAlertFinding.FindingID))
	}
	findingHistory, err := a.makeFindingIDs(ctx, findingIDs)
	if err != nil {
		a.logger.Errorf(ctx, "Failed makeFindingIDs when PutAlert. err: %v", err)
		return &alert.PutAlertResponse{Alert: convertAlert(registeredData)}, err
	}
	dataHistory := &model.AlertHistory{
		AlertID:        registeredData.AlertID,
		HistoryType:    alertHistoryStatus,
		Description:    registeredData.Description,
		Severity:       registeredData.Severity,
		FindingHistory: findingHistory,
		ProjectID:      registeredData.ProjectID,
	}

	// Fiding upsert
	_, err = a.repository.UpsertAlertHistory(ctx, dataHistory)
	if err != nil {
		a.logger.Errorf(ctx, "Failed PutAlertHistory when PutAlert. err: %v", err)
		return &alert.PutAlertResponse{Alert: convertAlert(registeredData)}, err
	}

	return &alert.PutAlertResponse{Alert: convertAlert(registeredData)}, nil
}

func (a *AlertService) DeleteAlert(ctx context.Context, req *alert.DeleteAlertRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteAlert(ctx, req.ProjectId, req.AlertId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * AlertHistory
 */

func (a *AlertService) ListAlertHistory(ctx context.Context, req *alert.ListAlertHistoryRequest) (*alert.ListAlertHistoryResponse, error) {
	listLimit := uint32(9)
	list, err := a.repository.ListAlertHistory(ctx, req.ProjectId, req.AlertId, "", listLimit)
	if err != nil {
		return nil, err
	}
	var histories []*alert.AlertHistory
	createdContain := false
	for _, d := range *list {
		if d.HistoryType == alertHistoryTypeCreated {
			createdContain = true
		}
		converted, err := convertAlertHistory(&d, true)
		if err != nil {
			a.logger.Errorf(ctx, "Error occurred in convertAlertHistory. err: %v", err)
			return nil, err
		}
		histories = append(histories, converted)
	}
	if createdContain {
		return &alert.ListAlertHistoryResponse{AlertHistory: histories}, nil
	}
	listCreated, err := a.repository.ListAlertHistory(ctx, req.ProjectId, req.AlertId, alertHistoryTypeCreated, 1)
	if err != nil {
		return nil, err
	}
	if *listCreated != nil && len(*listCreated) > 0 {
		converted, err := convertAlertHistory(&(*listCreated)[0], true)
		if err != nil {
			a.logger.Errorf(ctx, "Error occurred in convertAlertHistory. err: %v", err)
			return nil, err
		}
		histories = append(histories, converted)

	}
	return &alert.ListAlertHistoryResponse{AlertHistory: histories}, nil
}

func (a *AlertService) GetAlertHistory(ctx context.Context, req *alert.GetAlertHistoryRequest) (*alert.GetAlertHistoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := a.repository.GetAlertHistory(ctx, req.ProjectId, req.AlertHistoryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetAlertHistoryResponse{}, nil
		}
		return nil, err
	}
	convertedAlertHistory, err := convertAlertHistory(data, false)
	if err != nil {
		a.logger.Errorf(ctx, "Error occurred in convertAlertHistory. err: %v", err)
		return nil, err
	}
	return &alert.GetAlertHistoryResponse{AlertHistory: convertedAlertHistory}, nil
}

func (a *AlertService) PutAlertHistory(ctx context.Context, req *alert.PutAlertHistoryRequest) (*alert.PutAlertHistoryResponse, error) {
	if err := req.AlertHistory.Validate(); err != nil {
		return nil, err
	}

	data := &model.AlertHistory{
		AlertID:        req.AlertHistory.AlertId,
		HistoryType:    req.AlertHistory.HistoryType,
		Description:    req.AlertHistory.Description,
		Severity:       req.AlertHistory.Severity,
		FindingHistory: req.AlertHistory.FindingHistory,
		ProjectID:      req.AlertHistory.ProjectId,
	}

	// Fiding upsert
	registeredData, err := a.repository.UpsertAlertHistory(ctx, data)
	if err != nil {
		return nil, err
	}
	convertedAlertHistory, err := convertAlertHistory(registeredData, false)
	if err != nil {
		a.logger.Errorf(ctx, "Error occurred in convertAlertHistory. err: %v", err)
		return nil, err
	}
	return &alert.PutAlertHistoryResponse{AlertHistory: convertedAlertHistory}, nil
}

func (a *AlertService) DeleteAlertHistory(ctx context.Context, req *alert.DeleteAlertHistoryRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteAlertHistory(ctx, req.ProjectId, req.AlertHistoryId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * RelAlertFinding
 */

func (a *AlertService) ListRelAlertFinding(ctx context.Context, req *alert.ListRelAlertFindingRequest) (*alert.ListRelAlertFindingResponse, error) {
	converted := convertListRelAlertFindingRequest(req)
	list, err := a.repository.ListRelAlertFinding(ctx, converted.ProjectId, converted.AlertId, converted.FindingId, converted.FromAt, converted.ToAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.ListRelAlertFindingResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListRelAlertFindingResponse{}
	for _, d := range *list {
		data.RelAlertFinding = append(data.RelAlertFinding, convertRelAlertFinding(&d))
	}
	return &data, nil
}

func convertListRelAlertFindingRequest(req *alert.ListRelAlertFindingRequest) *alert.ListRelAlertFindingRequest {
	converted := alert.ListRelAlertFindingRequest{
		ProjectId: req.ProjectId,
		AlertId:   req.AlertId,
		FindingId: req.FindingId,
		FromAt:    req.FromAt,
		ToAt:      req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (a *AlertService) GetRelAlertFinding(ctx context.Context, req *alert.GetRelAlertFindingRequest) (*alert.GetRelAlertFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := a.repository.GetRelAlertFinding(ctx, req.ProjectId, req.AlertId, req.FindingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &alert.GetRelAlertFindingResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetRelAlertFindingResponse{RelAlertFinding: convertRelAlertFinding(data)}, nil
}

func (a *AlertService) PutRelAlertFinding(ctx context.Context, req *alert.PutRelAlertFindingRequest) (*alert.PutRelAlertFindingResponse, error) {
	if err := req.RelAlertFinding.Validate(); err != nil {
		return nil, err
	}
	data := &model.RelAlertFinding{
		AlertID:   req.RelAlertFinding.AlertId,
		FindingID: req.RelAlertFinding.FindingId,
		ProjectID: req.RelAlertFinding.ProjectId,
	}

	// Fiding upsert
	registeredData, err := a.repository.UpsertRelAlertFinding(ctx, data)
	if err != nil {
		return nil, err
	}

	return &alert.PutRelAlertFindingResponse{RelAlertFinding: convertRelAlertFinding(registeredData)}, nil
}

func (a *AlertService) DeleteRelAlertFinding(ctx context.Context, req *alert.DeleteRelAlertFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteRelAlertFinding(ctx, req.ProjectId, req.AlertId, req.FindingId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * Converter
 */

func convertAlert(a *model.Alert) *alert.Alert {
	if a == nil {
		return &alert.Alert{}
	}
	return &alert.Alert{
		AlertId:          a.AlertID,
		AlertConditionId: a.AlertConditionID,
		Description:      a.Description,
		Severity:         a.Severity,
		ProjectId:        a.ProjectID,
		Status:           getStatus(a.Status),
		CreatedAt:        a.CreatedAt.Unix(),
		UpdatedAt:        a.UpdatedAt.Unix(),
	}
}

func convertAlertHistory(a *model.AlertHistory, getCount bool) (*alert.AlertHistory, error) {
	if a == nil {
		return &alert.AlertHistory{}, nil
	}
	findingHistory := a.FindingHistory
	var err error
	if getCount {
		findingHistory, err = convertIDsToCount(findingHistory)
		if err != nil {
			return nil, err
		}
	}
	return &alert.AlertHistory{
		AlertHistoryId: a.AlertHistoryID,
		AlertId:        a.AlertID,
		HistoryType:    a.HistoryType,
		Description:    a.Description,
		Severity:       a.Severity,
		FindingHistory: findingHistory,
		ProjectId:      a.ProjectID,
		CreatedAt:      a.CreatedAt.Unix(),
		UpdatedAt:      a.UpdatedAt.Unix(),
	}, nil
}

type FindingHistory struct {
	FindingIDs []uint32 `json:"finding_id"`
}

func convertIDsToCount(history string) (string, error) {
	var findingHistory FindingHistory
	err := json.Unmarshal([]byte(history), &findingHistory)
	if err != nil {
		return "", err
	}
	converted, err := json.Marshal(struct {
		Count int `json:"count"`
	}{Count: len(findingHistory.FindingIDs)})
	if err != nil {
		return "", err
	}
	return string(converted), nil
}

func convertRelAlertFinding(f *model.RelAlertFinding) *alert.RelAlertFinding {
	if f == nil {
		return &alert.RelAlertFinding{}
	}
	return &alert.RelAlertFinding{
		AlertId:   f.AlertID,
		FindingId: f.FindingID,
		ProjectId: f.ProjectID,
		CreatedAt: f.CreatedAt.Unix(),
		UpdatedAt: f.UpdatedAt.Unix(),
	}
}

func getStatus(s string) alert.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := alert.Status_value[statusKey]; !ok {
		return alert.Status_UNKNOWN
	}
	switch statusKey {
	case alert.Status_ACTIVE.String():
		return alert.Status_ACTIVE
	case alert.Status_PENDING.String():
		return alert.Status_PENDING
	case alert.Status_DEACTIVE.String():
		return alert.Status_DEACTIVE
	default:
		return alert.Status_UNKNOWN
	}
}

func getStrings(statusSlice []alert.Status) []string {
	if len(statusSlice) == 0 {
		return nil
	}
	ret := []string{}
	for _, status := range statusSlice {
		ret = append(ret, status.String())
	}
	return ret
}
