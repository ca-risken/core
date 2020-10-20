package main

import (
	"context"
	"strings"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

/**
 * Alert
 */

func (f *alertService) ListAlert(ctx context.Context, req *alert.ListAlertRequest) (*alert.ListAlertResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	converted := convertListAlertRequest(req)
	list, err := f.repository.ListAlert(converted.ProjectId, getStrings(converted.Status), converted.Severity, converted.Description, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetAlert(ctx context.Context, req *alert.GetAlertRequest) (*alert.GetAlertResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetAlert(req.ProjectId, req.AlertId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetAlertResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertResponse{Alert: convertAlert(data)}, nil
}

func (f *alertService) PutAlert(ctx context.Context, req *alert.PutAlertRequest) (*alert.PutAlertResponse, error) {
	if err := req.Alert.Validate(); err != nil {
		return nil, err
	}
	savedData, err := f.repository.GetAlertByAlertConditionID(req.Alert.ProjectId, req.Alert.AlertConditionId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var alertID uint32
	if !noRecord {
		alertID = savedData.AlertID
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
	registerdData, err := f.repository.UpsertAlert(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertResponse{Alert: convertAlert(registerdData)}, nil
}

func (f *alertService) DeleteAlert(ctx context.Context, req *alert.DeleteAlertRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteAlert(req.ProjectId, req.AlertId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * AlertHistory
 */

func (f *alertService) ListAlertHistory(ctx context.Context, req *alert.ListAlertHistoryRequest) (*alert.ListAlertHistoryResponse, error) {
	converted := convertListAlertHistoryRequest(req)
	list, err := f.repository.ListAlertHistory(converted.ProjectId, converted.AlertId, converted.HistoryType, converted.Severity, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.ListAlertHistoryResponse{}, nil
		}
		return nil, err
	}
	data := alert.ListAlertHistoryResponse{}
	for _, d := range *list {
		data.AlertHistory = append(data.AlertHistory, convertAlertHistory(&d))
	}
	return &data, nil
}

func convertListAlertHistoryRequest(req *alert.ListAlertHistoryRequest) *alert.ListAlertHistoryRequest {
	converted := alert.ListAlertHistoryRequest{
		ProjectId:   req.ProjectId,
		HistoryType: req.HistoryType,
		AlertId:     req.AlertId,
		Severity:    req.Severity,
		FromAt:      req.FromAt,
		ToAt:        req.ToAt,
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (f *alertService) GetAlertHistory(ctx context.Context, req *alert.GetAlertHistoryRequest) (*alert.GetAlertHistoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := f.repository.GetAlertHistory(req.ProjectId, req.AlertHistoryId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetAlertHistoryResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetAlertHistoryResponse{AlertHistory: convertAlertHistory(data)}, nil
}

func (f *alertService) PutAlertHistory(ctx context.Context, req *alert.PutAlertHistoryRequest) (*alert.PutAlertHistoryResponse, error) {
	if err := req.AlertHistory.Validate(); err != nil {
		return nil, err
	}

	data := &model.AlertHistory{
		AlertID:     req.AlertHistory.AlertId,
		HistoryType: req.AlertHistory.HistoryType,
		Description: req.AlertHistory.Description,
		Severity:    req.AlertHistory.Severity,
		ProjectID:   req.AlertHistory.ProjectId,
	}

	// Fiding upsert
	registerdData, err := f.repository.UpsertAlertHistory(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutAlertHistoryResponse{AlertHistory: convertAlertHistory(registerdData)}, nil
}

func (f *alertService) DeleteAlertHistory(ctx context.Context, req *alert.DeleteAlertHistoryRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteAlertHistory(req.ProjectId, req.AlertHistoryId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * RelAlertFinding
 */

func (f *alertService) ListRelAlertFinding(ctx context.Context, req *alert.ListRelAlertFindingRequest) (*alert.ListRelAlertFindingResponse, error) {
	converted := convertListRelAlertFindingRequest(req)
	list, err := f.repository.ListRelAlertFinding(converted.ProjectId, converted.AlertId, converted.FindingId, converted.FromAt, converted.ToAt)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (f *alertService) GetRelAlertFinding(ctx context.Context, req *alert.GetRelAlertFindingRequest) (*alert.GetRelAlertFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data, err := f.repository.GetRelAlertFinding(req.ProjectId, req.AlertId, req.FindingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &alert.GetRelAlertFindingResponse{}, nil
		}
		return nil, err
	}
	return &alert.GetRelAlertFindingResponse{RelAlertFinding: convertRelAlertFinding(data)}, nil
}

func (f *alertService) PutRelAlertFinding(ctx context.Context, req *alert.PutRelAlertFindingRequest) (*alert.PutRelAlertFindingResponse, error) {
	if err := req.RelAlertFinding.Validate(); err != nil {
		return nil, err
	}
	data := &model.RelAlertFinding{
		AlertID:   req.RelAlertFinding.AlertId,
		FindingID: req.RelAlertFinding.FindingId,
		ProjectID: req.RelAlertFinding.ProjectId,
	}

	// Fiding upsert
	registerdData, err := f.repository.UpsertRelAlertFinding(data)
	if err != nil {
		return nil, err
	}

	return &alert.PutRelAlertFindingResponse{RelAlertFinding: convertRelAlertFinding(registerdData)}, nil
}

func (f *alertService) DeleteRelAlertFinding(ctx context.Context, req *alert.DeleteRelAlertFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteRelAlertFinding(req.ProjectId, req.AlertId, req.FindingId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * Converter
 */

func convertAlert(f *model.Alert) *alert.Alert {
	if f == nil {
		return &alert.Alert{}
	}
	return &alert.Alert{
		AlertId:          f.AlertID,
		AlertConditionId: f.AlertConditionID,
		Description:      f.Description,
		Severity:         f.Severity,
		ProjectId:        f.ProjectID,
		Status:           getStatus(f.Status),
		CreatedAt:        f.CreatedAt.Unix(),
		UpdatedAt:        f.UpdatedAt.Unix(),
	}
}

func convertAlertHistory(f *model.AlertHistory) *alert.AlertHistory {
	if f == nil {
		return &alert.AlertHistory{}
	}
	return &alert.AlertHistory{
		AlertHistoryId: f.AlertHistoryID,
		AlertId:        f.AlertID,
		HistoryType:    f.HistoryType,
		Description:    f.Description,
		Severity:       f.Severity,
		ProjectId:      f.ProjectID,
		CreatedAt:      f.CreatedAt.Unix(),
		UpdatedAt:      f.UpdatedAt.Unix(),
	}
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
