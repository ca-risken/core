package finding

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (f *FindingService) CleanOldResource(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	if err := f.repository.DeleteOldResource(ctx, f.excludeDeleteDataSource); err != nil {
		return nil, err
	}
	if err := f.repository.DeleteNoResourceIdTag(ctx); err != nil {
		return nil, err
	}
	if err := f.repository.DeleteOldFinding(ctx, f.excludeDeleteDataSource); err != nil {
		return nil, err
	}
	if err := f.repository.DeleteNoFindingIdTag(ctx); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
