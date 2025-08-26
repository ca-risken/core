package finding

import (
	"context"
	"errors"

	"github.com/ca-risken/core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (f *FindingService) UntagByResourceName(ctx context.Context, req *finding.UntagByResourceNameRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := f.untagFinding(ctx, req.ProjectId, req.ResourceName, req.Tag); err != nil {
		return nil, err
	}
	if err := f.untagResource(ctx, req.ProjectId, req.ResourceName, req.Tag); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f *FindingService) untagFinding(ctx context.Context, projectID uint32, resourceName, tag string) error {
	req := &finding.ListFindingRequest{
		ProjectId:    projectID,
		ResourceName: []string{resourceName},
		Tag:          []string{tag},
	}
	param := convertListFindingRequest(req)
	list, err := f.repository.ListFinding(ctx, param)
	if err != nil {
		return err
	}
	if list == nil {
		return nil
	}
	for _, fi := range *list {
		t, err := f.repository.GetFindingTagByKey(ctx, projectID, fi.FindingID, tag)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if t == nil {
			return nil
		}
		if err := f.repository.UntagFinding(ctx, projectID, t.FindingTagID); err != nil {
			return err
		}
	}
	return nil
}

func (f *FindingService) untagResource(ctx context.Context, projectID uint32, resourceName, tag string) error {
	r, err := f.repository.GetResourceByName(ctx, projectID, resourceName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if r == nil {
		return nil
	}
	t, err := f.repository.GetResourceTagByKey(ctx, projectID, r.ResourceID, tag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if t == nil {
		return nil
	}
	if err := f.repository.UntagResource(ctx, projectID, t.ResourceTagID); err != nil {
		return err
	}
	return nil
}
