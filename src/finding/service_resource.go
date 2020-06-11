package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

const (
	maxSumScore = 999999.9
)

func (f *findingService) ListResource(ctx context.Context, req *finding.ListResourceRequest) (*finding.ListResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := f.repository.ListResource(convertListResourceRequest(req))
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.ListResourceResponse{}, nil
		}
		return nil, err
	}

	// TODO authz
	var ids []uint64
	for _, id := range *list {
		ids = append(ids, uint64(id.ResourceID))
	}
	return &finding.ListResourceResponse{ResourceId: ids}, nil
}

func convertListResourceRequest(req *finding.ListResourceRequest) *finding.ListResourceRequest {
	converted := finding.ListResourceRequest{
		ProjectId:    req.ProjectId,
		ResourceName: req.ResourceName,
		FromSumScore: req.FromSumScore,
		ToSumScore:   req.ToSumScore,
		FromAt:       req.FromAt,
		ToAt:         req.ToAt,
	}
	if converted.ToSumScore == 0 {
		converted.ToSumScore = maxSumScore
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (f *findingService) GetResource(ctx context.Context, req *finding.GetResourceRequest) (*finding.GetResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *findingService) PutResource(ctx context.Context, req *finding.PutResourceRequest) (*finding.PutResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	savedData, err := f.repository.GetResourceByName(req.Resource.ProjectId, req.Resource.ResourceName)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var resourceID uint64
	if !noRecord {
		if req.Resource.ResourceId != 0 && req.Resource.ResourceId != savedData.ResourceID {
			return nil, fmt.Errorf("Invalid resource_id, want=%d, got=%d", savedData.ResourceID, req.Resource.ResourceId)
		}
		resourceID = savedData.ResourceID
	}

	// TODO: Authz
	// upsert
	registerdData, err := f.repository.UpsertResource(
		&model.Resource{
			ResourceID:   resourceID,
			ResourceName: req.Resource.ResourceName,
			ProjectID:    req.Resource.ProjectId,
		})
	if err != nil {
		return nil, err
	}
	return &finding.PutResourceResponse{Resource: convertResource(registerdData)}, nil
}

func (f *findingService) DeleteResource(ctx context.Context, req *finding.DeleteResourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (f *findingService) ListResourceTag(ctx context.Context, req *finding.ListResourceTagRequest) (*finding.ListResourceTagResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *findingService) TagResource(ctx context.Context, req *finding.TagResourceRequest) (*finding.TagResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *findingService) UntagResource(ctx context.Context, req *finding.UntagResourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func convertResource(r *model.Resource) *finding.Resource {
	if r == nil {
		return &finding.Resource{}
	}
	return &finding.Resource{
		ResourceId:   r.ResourceID,
		ResourceName: r.ResourceName,
		ProjectId:    r.ProjectID,
		CreatedAt:    r.CreatedAt.Unix(),
		UpdatedAt:    r.UpdatedAt.Unix(),
	}
}
