package ai

import (
	"context"
	"errors"
	"strings"
	"time"

	aiservice "github.com/ca-risken/core/pkg/ai"
	"github.com/ca-risken/core/proto/ai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const sensitiveContentTimeout = 10 * time.Second

type sensitiveContentClassifier interface {
	Classify(context.Context, *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error)
}

type SensitiveContentService struct {
	classifier sensitiveContentClassifier
}

var _ ai.SensitiveContentServiceServer = (*SensitiveContentService)(nil)

func NewSensitiveContentService(openAIToken string) *SensitiveContentService {
	if openAIToken == "" {
		return &SensitiveContentService{}
	}
	return &SensitiveContentService{
		classifier: aiservice.NewSensitiveContentClassifier(openAIToken),
	}
}

func (s *SensitiveContentService) ClassifySensitiveContent(ctx context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
	if s.classifier == nil {
		return nil, status.Error(codes.FailedPrecondition, "sensitive content classification is disabled")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request required")
	}
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if len(req.GetImage()) == 0 && strings.TrimSpace(req.GetText()) == "" ||
		len(req.GetImage()) > 0 && req.GetText() != "" {
		return nil, status.Error(codes.InvalidArgument, "provide either non-blank text or one image")
	}
	ctx, cancel := context.WithTimeout(ctx, sensitiveContentTimeout)
	defer cancel()
	result, err := s.classifier.Classify(ctx, req)
	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}
	if errors.Is(err, aiservice.ErrInvalidSensitiveContentImage) {
		return nil, status.Error(codes.InvalidArgument, aiservice.ErrInvalidSensitiveContentImage.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Unavailable, "sensitive content classification failed")
	}
	return result, nil
}
