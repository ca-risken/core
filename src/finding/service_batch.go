package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/src/finding/model"
	"github.com/golang/protobuf/ptypes/empty"
)

func (f *findingService) PutFindingBatch(ctx context.Context, req *finding.PutFindingBatchRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Base entity
	var findings []*model.Finding
	var resources []*model.Resource
	var recommends []*model.Recommend
	for _, d := range req.Finding {
		fi, err := f.getFindingDataForUpsert(ctx, d.Finding)
		if err != nil {
			return nil, err
		}
		findings = append(findings, fi)

		r, err := f.getResourceForUpsert(ctx, d.Finding.ProjectId, d.Finding.ResourceName)
		if err != nil {
			return nil, err
		}
		resources = append(resources, r)

		storedRecommendID, err := f.getStoredRecommendID(ctx, d.Finding.DataSource, d.Recommend.Type)
		if err != nil {
			return nil, err
		}
		var recommendID uint32
		if storedRecommendID != nil {
			recommendID = *storedRecommendID
		}
		recommends = append(recommends, &model.Recommend{
			RecommendID:    recommendID,
			DataSource:     d.Finding.DataSource,
			Type:           d.Recommend.Type,
			Risk:           d.Recommend.Risk,
			Recommendation: d.Recommend.Recommendation,
		})
	}
	// Bulk upsert (base entity)
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
		storedFinding, err := f.repository.GetFindingByDataSource(ctx, d.Finding.ProjectId, d.Finding.DataSource, d.Finding.DataSourceId)
		if err != nil {
			return nil, err
		}
		storedResource, err := f.repository.GetResourceByName(ctx, d.Finding.ProjectId, d.Finding.ResourceName)
		if err != nil {
			return nil, err
		}
		storedRecommend, err := f.repository.GetRecommendByDataSourceType(ctx, d.Finding.DataSource, d.Recommend.Type)
		if err != nil {
			return nil, err
		}
		recommendFindings = append(recommendFindings, &model.RecommendFinding{
			FindingID:   storedFinding.FindingID,
			RecommendID: storedRecommend.RecommendID,
			ProjectID:   d.Finding.ProjectId,
		})
		for _, t := range d.Tag {
			storedFindingTag, err := f.getFindingTagForUpsert(ctx, d.Finding.ProjectId, storedFinding.FindingID, t.Tag)
			if err != nil {
				return nil, err
			}
			findingTags = append(findingTags, &model.FindingTag{
				FindingTagID: storedFindingTag.FindingTagID,
				FindingID:    storedFinding.FindingID,
				ProjectID:    d.Finding.ProjectId,
				Tag:          t.Tag,
			})
			storedResourceTag, err := f.getResourceTagForUpsert(ctx, d.Finding.ProjectId, storedResource.ResourceID, t.Tag)
			if err != nil {
				return nil, err
			}
			resourceTags = append(resourceTags, &model.ResourceTag{
				ResourceTagID: storedResourceTag.ResourceTagID,
				ResourceID:    storedResource.ResourceID,
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
