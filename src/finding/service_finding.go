package main

import (
	"context"
	"math"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/vikyd/zero"
)

func (f *findingService) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListFindingRequest(req)
	total, err := f.repository.ListFindingCount(param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListFindingResponse{FindingId: []uint64{}, Count: 0, Total: total}, nil
	}
	list, err := f.repository.ListFinding(param)
	if err != nil {
		return nil, err
	}
	var ids []uint64
	for _, data := range *list {
		ids = append(ids, uint64(data.FindingID))
	}
	return &finding.ListFindingResponse{FindingId: ids, Count: uint32(len(ids)), Total: total}, nil
}

func convertListFindingRequest(req *finding.ListFindingRequest) *finding.ListFindingRequest {
	converted := req
	if zero.IsZeroVal(converted.ToScore) {
		converted.ToScore = 1.0
	}
	if zero.IsZeroVal(converted.ToAt) {
		converted.ToAt = time.Now().Unix()
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

func (f *findingService) GetFinding(ctx context.Context, req *finding.GetFindingRequest) (*finding.GetFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetFinding(req.ProjectId, req.FindingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.GetFindingResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetFindingResponse{Finding: convertFinding(data)}, nil
}

func (f *findingService) PutFinding(ctx context.Context, req *finding.PutFindingRequest) (*finding.PutFindingResponse, error) {
	if err := req.Finding.Validate(); err != nil {
		return nil, err
	}
	savedData, err := f.repository.GetFindingByDataSource(
		req.Finding.ProjectId, req.Finding.DataSource, req.Finding.DataSourceId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var findingID uint64
	if !noRecord {
		findingID = savedData.FindingID
	}
	fs, err := f.getFindingSettingByResource(req.ProjectId, req.Finding.ResourceName)
	if err != nil {
		return nil, err
	}
	data := &model.Finding{
		FindingID:     findingID,
		Description:   req.Finding.Description,
		DataSource:    req.Finding.DataSource,
		DataSourceID:  req.Finding.DataSourceId,
		ResourceName:  req.Finding.ResourceName,
		ProjectID:     req.Finding.ProjectId,
		OriginalScore: req.Finding.OriginalScore,
		Score:         calculateScore(req.Finding.OriginalScore, req.Finding.OriginalMaxScore, fs),
		Data:          req.Finding.Data,
	}

	// Fiding upsert
	registerdData, err := f.repository.UpsertFinding(data)
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

func (f *findingService) DeleteFinding(ctx context.Context, req *finding.DeleteFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.DeleteFinding(req.ProjectId, req.FindingId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f *findingService) ListFindingTag(ctx context.Context, req *finding.ListFindingTagRequest) (*finding.ListFindingTagResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListFindingTagRequest(req)
	total, err := f.repository.ListFindingTagCount(param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListFindingTagResponse{Tag: []*finding.FindingTag{}, Count: 0, Total: total}, nil
	}
	list, err := f.repository.ListFindingTag(param)
	if err != nil {
		return nil, err
	}
	var tags []*finding.FindingTag
	for _, tag := range *list {
		tags = append(tags, convertFindingTag(&tag))
	}
	return &finding.ListFindingTagResponse{Tag: tags, Count: uint32(len(tags)), Total: total}, nil
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

func (f *findingService) ListFindingTagName(ctx context.Context, req *finding.ListFindingTagNameRequest) (*finding.ListFindingTagNameResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListFindingTagNameRequest(req)
	total, err := f.repository.ListFindingTagNameCount(param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListFindingTagNameResponse{Tag: []string{}, Count: 0, Total: total}, nil
	}
	tags, err := f.repository.ListFindingTagName(param)
	if err != nil {
		return nil, err
	}
	var tagNames []string
	for _, tag := range *tags {
		tagNames = append(tagNames, tag.Tag)
	}
	return &finding.ListFindingTagNameResponse{Tag: tagNames, Count: uint32(len(tagNames)), Total: total}, nil
}

func convertListFindingTagNameRequest(req *finding.ListFindingTagNameRequest) *finding.ListFindingTagNameRequest {
	converted := req
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
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

func (f *findingService) TagFinding(ctx context.Context, req *finding.TagFindingRequest) (*finding.TagFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := f.repository.GetFindingTagByKey(req.ProjectId, req.Tag.FindingId, req.Tag.Tag)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var findingTagID uint64
	if !noRecord {
		findingTagID = savedData.FindingTagID
	}

	tag := &model.FindingTag{
		FindingTagID: findingTagID,
		FindingID:    req.Tag.FindingId,
		ProjectID:    req.Tag.ProjectId,
		Tag:          req.Tag.Tag,
	}
	registerd, err := f.repository.TagFinding(tag)
	if err != nil {
		return nil, err
	}
	return &finding.TagFindingResponse{Tag: convertFindingTag(registerd)}, nil
}

func (f *findingService) UntagFinding(ctx context.Context, req *finding.UntagFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := f.repository.UntagFinding(req.ProjectId, req.FindingTagId)
	if err != nil {
		return nil, err
	}
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
