package finding

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (f *FindingService) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListFindingRequest(req)
	total, err := f.repository.ListFindingCount(
		ctx,
		param.ProjectId, param.AlertId,
		param.FromScore, param.ToScore,
		param.FindingId,
		param.DataSource, param.ResourceName, param.Tag,
		param.Status,
	)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListFindingResponse{FindingId: []uint64{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	list, err := f.repository.ListFinding(ctx, param)
	if err != nil {
		return nil, err
	}
	var ids []uint64
	for _, data := range *list {
		ids = append(ids, uint64(data.FindingID))
	}
	return &finding.ListFindingResponse{FindingId: ids, Count: uint32(len(ids)), Total: convertToUint32(total)}, nil
}

func convertListFindingRequest(req *finding.ListFindingRequest) *finding.ListFindingRequest {
	converted := req
	if zero.IsZeroVal(converted.ToScore) {
		converted.ToScore = 1.0
	}
	if zero.IsZeroVal(converted.Sort) {
		converted.Sort = "finding_id"
	}
	if zero.IsZeroVal(converted.Direction) {
		converted.Direction = defaultSortDirection
	}
	if zero.IsZeroVal(converted.Limit) {
		converted.Limit = defaultLimit
	}
	return converted
}

func (f *FindingService) BatchListFinding(ctx context.Context, req *finding.BatchListFindingRequest) (*finding.BatchListFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertBatchListFindingRequest(req)
	total, err := f.repository.ListFindingCount(
		ctx,
		param.ProjectId, param.AlertId,
		param.FromScore, param.ToScore,
		param.FindingId,
		param.DataSource, param.ResourceName, param.Tag,
		param.Status)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.BatchListFindingResponse{FindingId: []uint64{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	list, err := f.repository.BatchListFinding(ctx, param)
	if err != nil {
		return nil, err
	}
	var ids []uint64
	for _, data := range *list {
		ids = append(ids, uint64(data.FindingID))
	}
	return &finding.BatchListFindingResponse{FindingId: ids, Count: uint32(len(ids)), Total: convertToUint32(total)}, nil
}

func convertBatchListFindingRequest(req *finding.BatchListFindingRequest) *finding.BatchListFindingRequest {
	converted := req
	if zero.IsZeroVal(converted.ToScore) {
		converted.ToScore = 1.0
	}
	return converted
}

func (f *FindingService) GetFinding(ctx context.Context, req *finding.GetFindingRequest) (*finding.GetFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetFinding(ctx, req.ProjectId, req.FindingId, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &finding.GetFindingResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetFindingResponse{Finding: convertFinding(data)}, nil
}

func (f *FindingService) PutFinding(ctx context.Context, req *finding.PutFindingRequest) (*finding.PutFindingResponse, error) {
	if err := req.Finding.Validate(); err != nil {
		return nil, err
	}
	data, err := f.getFindingDataForUpsert(ctx, req.Finding)
	if err != nil {
		return nil, err
	}

	// Fiding upsert
	registerdData, err := f.repository.UpsertFinding(ctx, data)
	if err != nil {
		return nil, err
	}

	// Resource
	_, err = f.PutResource(ctx, &finding.PutResourceRequest{
		Resource: &finding.ResourceForUpsert{
			ResourceName: req.Finding.ResourceName,
			ProjectId:    req.Finding.ProjectId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &finding.PutFindingResponse{Finding: convertFinding(registerdData)}, nil
}

func (f *FindingService) DeleteFinding(ctx context.Context, req *finding.DeleteFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteFinding(ctx, req.ProjectId, req.FindingId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f *FindingService) ListFindingTag(ctx context.Context, req *finding.ListFindingTagRequest) (*finding.ListFindingTagResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListFindingTagRequest(req)
	total, err := f.repository.ListFindingTagCount(ctx, param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListFindingTagResponse{Tag: []*finding.FindingTag{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	list, err := f.repository.ListFindingTag(ctx, param)
	if err != nil {
		return nil, err
	}
	var tags []*finding.FindingTag
	for _, tag := range *list {
		tags = append(tags, convertFindingTag(&tag))
	}
	return &finding.ListFindingTagResponse{Tag: tags, Count: uint32(len(tags)), Total: convertToUint32(total)}, nil
}

func convertListFindingTagRequest(req *finding.ListFindingTagRequest) *finding.ListFindingTagRequest {
	converted := req
	if zero.IsZeroVal(converted.Sort) {
		converted.Sort = "tag"
	}
	if zero.IsZeroVal(converted.Direction) {
		converted.Direction = defaultSortDirection
	}
	if zero.IsZeroVal(converted.Limit) {
		converted.Limit = defaultLimit
	}
	return converted
}

func (f *FindingService) ListFindingTagName(ctx context.Context, req *finding.ListFindingTagNameRequest) (*finding.ListFindingTagNameResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListFindingTagNameRequest(req)
	total, err := f.repository.ListFindingTagNameCount(ctx, param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListFindingTagNameResponse{Tag: []string{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	tags, err := f.repository.ListFindingTagName(ctx, param)
	if err != nil {
		return nil, err
	}
	var tagNames []string
	for _, tag := range *tags {
		tagNames = append(tagNames, tag.Tag)
	}
	return &finding.ListFindingTagNameResponse{Tag: tagNames, Count: uint32(len(tagNames)), Total: convertToUint32(total)}, nil
}

func convertListFindingTagNameRequest(req *finding.ListFindingTagNameRequest) *finding.ListFindingTagNameRequest {
	converted := req
	if zero.IsZeroVal(converted.Sort) {
		converted.Sort = "tag"
	}
	if zero.IsZeroVal(converted.Direction) {
		converted.Direction = defaultSortDirection
	}
	if zero.IsZeroVal(converted.Limit) {
		converted.Limit = defaultLimit
	}
	return converted
}

func (f *FindingService) TagFinding(ctx context.Context, req *finding.TagFindingRequest) (*finding.TagFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	fData, err := f.repository.GetFinding(ctx, req.ProjectId, req.Tag.FindingId, true)
	if err != nil {
		return nil, fmt.Errorf("Failed to GetFinding, finding_id=%d, err=%+v", req.Tag.FindingId, err)
	}

	findingTag, err := f.getFindingTagForUpsert(ctx, req.ProjectId, req.Tag.FindingId, req.Tag.Tag)
	if err != nil {
		return nil, err
	}

	registerd, err := f.repository.TagFinding(ctx, findingTag)
	if err != nil {
		return nil, err
	}
	// Put resource tag
	r, err := f.repository.GetResourceByName(ctx, req.ProjectId, fData.ResourceName)
	if err != nil {
		return nil, fmt.Errorf("Failed to GetResourceByName, resource_name=%s, err=%+v", fData.ResourceName, err)
	}
	_, err = f.repository.TagResource(ctx, &model.ResourceTag{
		ResourceID: r.ResourceID,
		ProjectID:  r.ProjectID,
		Tag:        req.Tag.Tag,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to TagResource, resource_id=%d, err=%+v", r.ResourceID, err)
	}

	return &finding.TagFindingResponse{Tag: convertFindingTag(registerd)}, nil
}

func (f *FindingService) UntagFinding(ctx context.Context, req *finding.UntagFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.UntagFinding(ctx, req.ProjectId, req.FindingTagId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f *FindingService) ClearScore(ctx context.Context, req *finding.ClearScoreRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.ClearScoreFinding(ctx, req)
	if err != nil {
		return nil, err
	}
	f.logger.Infof(ctx, "Finding score cleared, param=%+v", req)

	return &empty.Empty{}, nil
}

func convertFinding(f *model.Finding) *finding.Finding {
	if f == nil {
		return &finding.Finding{}
	}
	return &finding.Finding{
		FindingId:     f.FindingID,
		Description:   f.Description,
		DataSource:    f.DataSource,
		DataSourceId:  f.DataSourceID,
		ResourceName:  f.ResourceName,
		ProjectId:     f.ProjectID,
		OriginalScore: f.OriginalScore,
		Score:         f.Score,
		Data:          f.Data,
		CreatedAt:     f.CreatedAt.Unix(),
		UpdatedAt:     f.UpdatedAt.Unix(),
	}
}

func convertFindingTag(f *model.FindingTag) *finding.FindingTag {
	if f == nil {
		return &finding.FindingTag{}
	}
	return &finding.FindingTag{
		FindingTagId: f.FindingTagID,
		FindingId:    f.FindingID,
		ProjectId:    f.ProjectID,
		Tag:          f.Tag,
		CreatedAt:    f.CreatedAt.Unix(),
		UpdatedAt:    f.UpdatedAt.Unix(),
	}
}

func calculateScore(score, maxScore float32, setting *findingSetting) float32 {
	baseScore := float32(math.Round(float64(score/maxScore*100)) / 100)
	if setting == nil || zero.IsZeroVal(setting.ScoreCoefficient) {
		return baseScore
	}
	calculated := baseScore * setting.ScoreCoefficient
	if calculated > 1.0 {
		return 1.0
	}
	if calculated < 0 {
		return 0.0
	}
	return calculated
}

func (f *FindingService) getFindingDataForUpsert(ctx context.Context, req *finding.FindingForUpsert) (*model.Finding, error) {
	storedData, err := f.repository.GetFindingByDataSource(
		ctx, req.ProjectId, req.DataSource, req.DataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}

	// Specify the ID in PK as much as possible to avoid unnecessary AUTO_INCREMENT.
	// https://dev.mysql.com/doc/refman/5.6/ja/insert-on-duplicate.html
	var findingID uint64
	if !noRecord {
		findingID = storedData.FindingID
	}
	fs, err := f.getFindingSettingByResource(ctx, req.ProjectId, req.ResourceName)
	if err != nil {
		return nil, err
	}
	return &model.Finding{
		FindingID:     findingID,
		Description:   req.Description,
		DataSource:    req.DataSource,
		DataSourceID:  req.DataSourceId,
		ResourceName:  req.ResourceName,
		ProjectID:     req.ProjectId,
		OriginalScore: req.OriginalScore,
		Score:         calculateScore(req.OriginalScore, req.OriginalMaxScore, fs),
		Data:          req.Data,
	}, nil
}

func (f *FindingService) getFindingTagForUpsert(ctx context.Context, projectID uint32, findingID uint64, tag string) (*model.FindingTag, error) {
	storedData, err := f.repository.GetFindingTagByKey(ctx, projectID, findingID, tag)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, fmt.Errorf("Failed to GetFindingTagByKey, finding_id=%d, tag=%s, err=%+v", findingID, tag, err)
	}

	// Specify the ID in PK as much as possible to avoid unnecessary AUTO_INCREMENT.
	// https://dev.mysql.com/doc/refman/5.6/ja/insert-on-duplicate.html
	var findingTagID uint64
	if !noRecord {
		findingTagID = storedData.FindingTagID
	}
	return &model.FindingTag{
		FindingTagID: findingTagID,
		FindingID:    findingID,
		ProjectID:    projectID,
		Tag:          tag,
	}, nil
}
