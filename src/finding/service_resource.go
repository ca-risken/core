package main

import (
	"context"
	"errors"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
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
	param := convertListResourceRequest(req)
	total, err := f.repository.ListResourceCount(ctx, param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListResourceResponse{ResourceId: []uint64{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	list, err := f.repository.ListResource(ctx, param)
	if err != nil {
		return nil, err
	}
	var ids []uint64
	for _, data := range *list {
		ids = append(ids, uint64(data.ResourceID))
	}
	return &finding.ListResourceResponse{ResourceId: ids, Count: uint32(len(ids)), Total: convertToUint32(total)}, nil
}

func convertListResourceRequest(req *finding.ListResourceRequest) *finding.ListResourceRequest {
	converted := req
	if converted.ToSumScore == 0 {
		converted.ToSumScore = maxSumScore
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	if zero.IsZeroVal(converted.Sort) {
		converted.Sort = "resource_id"
	}
	if zero.IsZeroVal(converted.Direction) {
		converted.Direction = defaultSortDirection
	}
	if zero.IsZeroVal(converted.Limit) {
		converted.Limit = defaultLimit
	}
	return converted
}

func (f *findingService) GetResource(ctx context.Context, req *finding.GetResourceRequest) (*finding.GetResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetResource(ctx, req.ProjectId, req.ResourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &finding.GetResourceResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetResourceResponse{Resource: convertResource(data)}, nil
}

func (f *findingService) PutResource(ctx context.Context, req *finding.PutResourceRequest) (*finding.PutResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := f.repository.GetResourceByName(ctx, req.Resource.ProjectId, req.Resource.ResourceName)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
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
	// upsert
	registerdData, err := f.repository.UpsertResource(ctx, data)
	if err != nil {
		return nil, err
	}
	return &finding.PutResourceResponse{Resource: convertResource(registerdData)}, nil
}

func (f *findingService) DeleteResource(ctx context.Context, req *finding.DeleteResourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return &empty.Empty{}, err
	}
	err := f.repository.DeleteResource(ctx, req.ProjectId, req.ResourceId)
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
	param := convertListResourceTagRequest(req)
	total, err := f.repository.ListResourceTagCount(ctx, param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListResourceTagResponse{Tag: []*finding.ResourceTag{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	list, err := f.repository.ListResourceTag(ctx, param)
	if err != nil {
		return nil, err
	}
	var tags []*finding.ResourceTag
	for _, tag := range *list {
		tags = append(tags, convertResourceTag(&tag))
	}
	return &finding.ListResourceTagResponse{Tag: tags, Count: uint32(len(tags)), Total: convertToUint32(total)}, nil
}

func convertListResourceTagRequest(req *finding.ListResourceTagRequest) *finding.ListResourceTagRequest {
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

func (f *findingService) ListResourceTagName(ctx context.Context, req *finding.ListResourceTagNameRequest) (*finding.ListResourceTagNameResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	param := convertListResourceTagNameRequest(req)
	total, err := f.repository.ListResourceTagNameCount(ctx, param)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &finding.ListResourceTagNameResponse{Tag: []string{}, Count: 0, Total: convertToUint32(total)}, nil
	}
	tags, err := f.repository.ListResourceTagName(ctx, param)
	if err != nil {
		return nil, err
	}
	var tagNames []string
	for _, tag := range *tags {
		tagNames = append(tagNames, tag.Tag)
	}
	return &finding.ListResourceTagNameResponse{Tag: tagNames, Count: uint32(len(tagNames)), Total: convertToUint32(total)}, nil
}

func convertListResourceTagNameRequest(req *finding.ListResourceTagNameRequest) *finding.ListResourceTagNameRequest {
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
func (f *findingService) TagResource(ctx context.Context, req *finding.TagResourceRequest) (*finding.TagResourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := f.repository.GetResourceTagByKey(ctx, req.ProjectId, req.Tag.ResourceId, req.Tag.Tag)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
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
		Tag:           req.Tag.Tag,
	}
	registerd, err := f.repository.TagResource(ctx, tag)
	if err != nil {
		return nil, err
	}
	return &finding.TagResourceResponse{Tag: convertResourceTag(registerd)}, nil
}

func (f *findingService) UntagResource(ctx context.Context, req *finding.UntagResourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return &empty.Empty{}, err
	}
	err := f.repository.UntagResource(ctx, req.ProjectId, req.ResourceTagId)
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
		Tag:           r.Tag,
		CreatedAt:     r.CreatedAt.Unix(),
		UpdatedAt:     r.UpdatedAt.Unix(),
	}
}
