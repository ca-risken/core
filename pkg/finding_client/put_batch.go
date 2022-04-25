package finding_client

import (
	"context"

	"github.com/ca-risken/core/proto/finding"
)

// PutFindingBatch call PutFindingBatch API at each request limit
func PutFindingBatch(ctx context.Context, client finding.FindingServiceClient, projectID uint32, params []*finding.FindingBatchForUpsert) error {
	appLogger.Infof(ctx, "Putting findings(%d)...", len(params))
	for idx := 0; idx < len(params); idx = idx + finding.PutFindingBatchMaxLength {
		lastIdx := idx + finding.PutFindingBatchMaxLength
		if lastIdx > len(params) {
			lastIdx = len(params)
		}
		// request per API limits
		appLogger.Debugf(ctx, "Call PutFindingBatch API, (%d ~ %d / %d)", idx+1, lastIdx, len(params))
		req := &finding.PutFindingBatchRequest{ProjectId: projectID, Finding: params[idx:lastIdx]}
		if _, err := client.PutFindingBatch(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// PutResourceBatch call PutResourceBatch API at each request limit
func PutResourceBatch(ctx context.Context, client finding.FindingServiceClient, projectID uint32, params []*finding.ResourceBatchForUpsert) error {
	appLogger.Infof(ctx, "Putting resources(%d)...", len(params))
	for idx := 0; idx < len(params); idx = idx + finding.PutResourceBatchMaxLength {
		lastIdx := idx + finding.PutResourceBatchMaxLength
		if lastIdx > len(params) {
			lastIdx = len(params)
		}
		// request per API limits
		appLogger.Debugf(ctx, "Call PutResourceBatch API, (%d ~ %d / %d)", idx+1, lastIdx, len(params))
		req := &finding.PutResourceBatchRequest{ProjectId: projectID, Resource: params[idx:lastIdx]}
		if _, err := client.PutResourceBatch(ctx, req); err != nil {
			return err
		}
	}
	return nil
}
