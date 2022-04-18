package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/src/finding/model"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
)

func (f *findingService) PutFindingBatch(ctx context.Context, req *finding.PutFindingBatchRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Base entity
	var findings []*model.Finding
	var resources []*model.Resource
	var recommends []*model.Recommend
	findingIDCache := map[int]uint64{}
	resourceIDCache := map[int]uint64{}
	recommendIDCache := map[int]uint32{}
	for idx, d := range req.Finding {
		fi, err := f.getFindingDataForUpsert(ctx, d.Finding)
		if err != nil {
			return nil, err
		}
		if !zero.IsZeroVal(fi.FindingID) {
			findingIDCache[idx] = fi.FindingID
		}
		findings = append(findings, fi)

		r, err := f.getResourceForUpsert(ctx, d.Finding.ProjectId, d.Finding.ResourceName)
		if err != nil {
			return nil, err
		}
		if !zero.IsZeroVal(r.ResourceID) {
			resourceIDCache[idx] = r.ResourceID
		}
		resources = append(resources, r)

		if d.Recommend != nil {
			storedRecommendID, err := f.getStoredRecommendID(ctx, d.Finding.DataSource, d.Recommend.Type)
			if err != nil {
				return nil, err
			}
			var recommendID uint32
			if storedRecommendID != nil {
				recommendID = *storedRecommendID
				recommendIDCache[idx] = *storedRecommendID
			}
			recommends = append(recommends, &model.Recommend{
				RecommendID:    recommendID,
				DataSource:     d.Finding.DataSource,
				Type:           d.Recommend.Type,
				Risk:           d.Recommend.Risk,
				Recommendation: d.Recommend.Recommendation,
			})
		}
	}
	if err := f.repository.BulkUpsertFinding(ctx, findings); err != nil {
		return nil, fmt.Errorf("Failed to BulkUpsertFinding, err=%+w", err)
	}
	if err := f.repository.BulkUpsertResource(ctx, resources); err != nil {
		return nil, fmt.Errorf("Failed to BulkUpsertResource, err=%+w", err)
	}
	if err := f.repository.BulkUpsertRecommend(ctx, recommends); err != nil {
		return nil, fmt.Errorf("Failed to BulkUpsertRecommend, err=%+w", err)
	}

	// Sub entity
	var recommendFindings []*model.RecommendFinding
	var findingTags []*model.FindingTag
	var resourceTags []*model.ResourceTag
	for idx, d := range req.Finding {
		var storedFindingID uint64
		if id, ok := findingIDCache[idx]; ok {
			storedFindingID = id
		} else {
			newFinding, err := f.repository.GetFindingByDataSource(ctx, d.Finding.ProjectId, d.Finding.DataSource, d.Finding.DataSourceId)
			if err != nil {
				return nil, err
			}
			storedFindingID = newFinding.FindingID
		}

		var storedResourceID uint64
		if id, ok := resourceIDCache[idx]; ok {
			storedResourceID = id
		} else {
			newResource, err := f.repository.GetResourceByName(ctx, d.Finding.ProjectId, d.Finding.ResourceName)
			if err != nil {
				return nil, err
			}
			storedResourceID = newResource.ResourceID
		}

		if d.Recommend != nil {
			var storedRecommendID uint32
			if id, ok := recommendIDCache[idx]; ok {
				storedRecommendID = id
			} else {
				newRecommend, err := f.repository.GetRecommendByDataSourceType(ctx, d.Finding.DataSource, d.Recommend.Type)
				if err != nil {
					return nil, err
				}
				storedRecommendID = newRecommend.RecommendID
			}
			recommendFindings = append(recommendFindings, &model.RecommendFinding{
				FindingID:   storedFindingID,
				RecommendID: storedRecommendID,
				ProjectID:   d.Finding.ProjectId,
			})
		}
		storedFindingTags, err := f.repository.ListFindingTagByFindingID(ctx, d.Finding.ProjectId, storedFindingID)
		if err != nil {
			return nil, err
		}
		storedResourceTags, err := f.repository.ListResourceTagByResourceID(ctx, d.Finding.ProjectId, storedResourceID)
		if err != nil {
			return nil, err
		}
		for _, t := range d.Tag {
			findingTag := &model.FindingTag{
				FindingID: storedFindingID,
				ProjectID: d.Finding.ProjectId,
				Tag:       t.Tag,
			}
			for _, storedTag := range *storedFindingTags {
				if t.Tag == storedTag.Tag {
					findingTag.FindingTagID = storedTag.FindingTagID
					break
				}
			}
			findingTags = append(findingTags, findingTag)

			resourceTag := &model.ResourceTag{
				ResourceID: storedResourceID,
				ProjectID:  d.Finding.ProjectId,
				Tag:        t.Tag,
			}
			for _, storedTag := range *storedResourceTags {
				if t.Tag == storedTag.Tag {
					resourceTag.ResourceTagID = storedTag.ResourceTagID
					break
				}
			}
			resourceTags = append(resourceTags, resourceTag)
		}
	}
	if err := f.repository.BulkUpsertRecommendFinding(ctx, recommendFindings); err != nil {
		return nil, fmt.Errorf("Failed to BulkUpsertRecommendFinding, err=%+w", err)
	}
	if err := f.repository.BulkUpsertFindingTag(ctx, findingTags); err != nil {
		return nil, fmt.Errorf("Failed to BulkUpsertFindingTag, err=%+w", err)
	}
	if err := f.repository.BulkUpsertResourceTag(ctx, resourceTags); err != nil {
		return nil, fmt.Errorf("Failed to BulkUpsertResourceTag, err=%+w", err)
	}
	appLogger.Infof("Succeded PutFindingBatch, project_id=%d, findings=%d", req.ProjectId, len(req.Finding))
	return &empty.Empty{}, nil
}
