package finding

import (
	"context"

	"github.com/ca-risken/core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
)

func (f *FindingService) UntagByResourceName(ctx context.Context, req *finding.UntagByResourceNameRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
