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

/**
 * Resource
 */

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

	var ids []uint64
	for _, data := range *list {
		// Authz
		if f.isAuthorizedWithResource(ctx, req.UserId, "ListResource", &data) {
			ids = append(ids, uint64(data.ResourceID))
		}
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

	data, err := f.repository.GetResource(req.ResourceId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.GetResourceResponse{}, nil
		}
		return nil, err
	}
	// Authz
	if !f.isAuthorizedWithResource(ctx, req.UserId, "GetResource", data) {
		return &finding.GetResourceResponse{}, nil
	}
	return &finding.GetResourceResponse{Resource: convertResource(data)}, nil
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
		resourceID = savedData.ResourceID
	}
	data := &model.Resource{
		ResourceID:   resourceID,
		ResourceName: req.Resource.ResourceName,
		ProjectID:    req.Resource.ProjectId,
	}
	// Authz
	if !f.isAuthorizedWithResource(ctx, req.UserId, "PutResource", data) {
		return nil, fmt.Errorf("Unauthorized PutResource action for resource=%+v", data)
	}

	// upsert
	registerdData, err := f.repository.UpsertResource(data)
	if err != nil {
		return nil, err
	}
	return &finding.PutResourceResponse{Resource: convertResource(registerdData)}, nil
}

func (f *findingService) DeleteResource(ctx context.Context, req *finding.DeleteResourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return &empty.Empty{}, err
	}
	// Authz
	if !f.isAuthorizedWithResourceID(ctx, req.UserId, "DeleteResource", req.ResourceId) {
		return nil, fmt.Errorf("Unauthorized DeleteResource action for resource_id=%d", req.ResourceId)
	}
	err := f.repository.DeleteResource(req.ResourceId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

/**
 * ResourceTag
 */

func (f *findingService) ListResourceTag(ctx context.Context, req *finding.ListResourceTagRequest) (*finding.ListResourceTagResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// TODO authz
	list, err := f.repository.ListResourceTag(req.ResourceId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.ListResourceTagResponse{}, nil
		}
		return nil, err
	}
	var tags []*finding.ResourceTag
	for _, tag := range *list {
		// Authz
		if f.isAuthorizedWithResourceTag(ctx, req.UserId, "ListResourceTag", &tag) {
			tags = append(tags, convertResourceTag(&tag))
		}
	}
	return &finding.ListResourceTagResponse{Tag: tags}, nil
}

func (f *findingService) TagResource(ctx context.Context, req *finding.TagResourceRequest) (*finding.TagResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	savedData, err := f.repository.GetResourceTagByKey(req.Tag.ResourceId, req.Tag.TagKey)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var resourceTagID uint64
	if !noRecord {
		resourceTagID = savedData.ResourceTagID
	}
	tag := &model.ResourceTag{
		ResourceTagID: resourceTagID,
		ResourceID:    req.Tag.ResourceId,
		ProjectID:     req.Tag.ProjectId,
		TagKey:        req.Tag.TagKey,
		TagValue:      req.Tag.TagValue,
	}
	// Authz
	if !f.isAuthorizedWithResourceTag(ctx, req.UserId, "TagResource", tag) {
		return nil, fmt.Errorf("Unauthorized TagResource action for tag=%+v", tag)
	}

	registerd, err := f.repository.TagResource(tag)
	if err != nil {
		return nil, err
	}
	return &finding.TagResourceResponse{Tag: convertResourceTag(registerd)}, nil
}

func (f *findingService) UntagResource(ctx context.Context, req *finding.UntagResourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return &empty.Empty{}, err
	}
	// Authz
	if !f.isAuthorizedWithResourceTagID(ctx, req.UserId, "UntagResource", req.ResourceTagId) {
		return nil, fmt.Errorf("Unauthorized UntagResource action for resource_tag_id=%d", req.ResourceTagId)
	}

	err := f.repository.UntagResource(req.ResourceTagId)
	if err != nil {
		return nil, err
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

func convertResourceTag(r *model.ResourceTag) *finding.ResourceTag {
	if r == nil {
		return &finding.ResourceTag{}
	}
	return &finding.ResourceTag{
		ResourceTagId: r.ResourceTagID,
		ResourceId:    r.ResourceID,
		ProjectId:     r.ProjectID,
		TagKey:        r.TagKey,
		TagValue:      r.TagValue,
		CreatedAt:     r.CreatedAt.Unix(),
		UpdatedAt:     r.UpdatedAt.Unix(),
	}
}
