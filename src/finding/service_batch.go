package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/src/finding/model"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
)

const (
	findingIDCacheKeyFmt   = "finding/%d/%s/%s"
	resourceIDCacheKeyFmt  = "resource/%d/%s"
	recommendIDCacheKeyFmt = "recommend/%s/%s"
)

func (f *findingService) PutFindingBatch(ctx context.Context, req *finding.PutFindingBatchRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Base entity
	var findings []*model.Finding
	var resources []*model.Resource
	var recommends []*model.Recommend
	findingIDCache := map[string]uint64{}
	resourceIDCache := map[string]uint64{}
	recommendIDCache := map[string]uint32{}
	for _, d := range req.Finding {
		fi, err := f.getFindingDataForUpsert(ctx, d.Finding)
		if err != nil {
			return nil, err
		}
		if !zero.IsZeroVal(fi.FindingID) {
			findingIDCache[fmt.Sprintf(findingIDCacheKeyFmt, d.Finding.ProjectId, d.Finding.DataSource, d.Finding.DataSourceId)] = fi.FindingID
		}
		findings = append(findings, fi)

		r, err := f.getResourceForUpsert(ctx, d.Finding.ProjectId, d.Finding.ResourceName)
		if err != nil {
			return nil, err
		}
		if !zero.IsZeroVal(r.ResourceID) {
			resourceIDCache[fmt.Sprintf(resourceIDCacheKeyFmt, d.Finding.ProjectId, d.Finding.ResourceName)] = r.ResourceID
		}
		resources = append(resources, r)

		storedRecommendID, err := f.getStoredRecommendID(ctx, d.Finding.DataSource, d.Recommend.Type)
		if err != nil {
			return nil, err
		}
		var recommendID uint32
		if storedRecommendID != nil {
			recommendID = *storedRecommendID
			recommendIDCache[fmt.Sprintf(recommendIDCacheKeyFmt, d.Finding.DataSource, d.Recommend.Type)] = *storedRecommendID
		}
		recommends = append(recommends, &model.Recommend{
			RecommendID:    recommendID,
			DataSource:     d.Finding.DataSource,
			Type:           d.Recommend.Type,
			Risk:           d.Recommend.Risk,
			Recommendation: d.Recommend.Recommendation,
		})
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
	for _, d := range req.Finding {
		var storedFindingID uint64
		if id, ok := findingIDCache[fmt.Sprintf(findingIDCacheKeyFmt, d.Finding.ProjectId, d.Finding.DataSource, d.Finding.DataSourceId)]; ok {
			storedFindingID = id
		} else {
			newFinding, err := f.repository.GetFindingByDataSource(ctx, d.Finding.ProjectId, d.Finding.DataSource, d.Finding.DataSourceId)
			if err != nil {
				return nil, err
			}
			storedFindingID = newFinding.FindingID
		}

		var storedResourceID uint64
		if id, ok := resourceIDCache[fmt.Sprintf(resourceIDCacheKeyFmt, d.Finding.ProjectId, d.Finding.ResourceName)]; ok {
			storedResourceID = id
		} else {
			newResource, err := f.repository.GetResourceByName(ctx, d.Finding.ProjectId, d.Finding.ResourceName)
			if err != nil {
				return nil, err
			}
			storedResourceID = newResource.ResourceID
		}

		var storedRecommendID uint32
		if id, ok := recommendIDCache[fmt.Sprintf(recommendIDCacheKeyFmt, d.Finding.DataSource, d.Recommend.Type)]; ok {
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
		for _, t := range d.Tag {
			storedFindingTag, err := f.getFindingTagForUpsert(ctx, d.Finding.ProjectId, storedFindingID, t.Tag)
			if err != nil {
				return nil, err
			}
			findingTags = append(findingTags, &model.FindingTag{
				FindingTagID: storedFindingTag.FindingTagID,
				FindingID:    storedFindingID,
				ProjectID:    d.Finding.ProjectId,
				Tag:          t.Tag,
			})
			storedResourceTag, err := f.getResourceTagForUpsert(ctx, d.Finding.ProjectId, storedResourceID, t.Tag)
			if err != nil {
				return nil, err
			}
			resourceTags = append(resourceTags, &model.ResourceTag{
				ResourceTagID: storedResourceTag.ResourceTagID,
				ResourceID:    storedResourceID,
				ProjectID:     d.Finding.ProjectId,
				Tag:           t.Tag,
			})
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
