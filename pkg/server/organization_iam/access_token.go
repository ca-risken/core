package organization_iam

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	maxOrgTokenExpiredAtUnix int64 = 253402268399 // 9999-12-31T23:59:59
)

func (s *OrganizationIAMService) ListOrganizationAccessToken(ctx context.Context, req *organization_iam.ListOrganizationAccessTokenRequest) (*organization_iam.ListOrganizationAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	resp, err := s.repository.ListOrgAccessToken(ctx, req.OrganizationId, req.Name, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.ListOrganizationAccessTokenResponse{}, nil
		}
		return nil, err
	}
	var list []*organization_iam.OrganizationAccessToken
	for _, token := range *resp {
		list = append(list, convertOrgAccessToken(&token))
	}
	return &organization_iam.ListOrganizationAccessTokenResponse{AccessToken: list}, nil
}

func convertOrgAccessToken(token *model.OrgAccessToken) *organization_iam.OrganizationAccessToken {
	return &organization_iam.OrganizationAccessToken{
		AccessTokenId:     token.AccessTokenID,
		Name:              token.Name,
		Description:       token.Description,
		OrganizationId:    token.OrgID,
		ExpiredAt:         token.ExpiredAt.Unix(),
		LastUpdatedUserId: token.LastUpdatedUserID,
		CreatedAt:         token.CreatedAt.Unix(),
		UpdatedAt:         token.UpdatedAt.Unix(),
	}
}

func (s *OrganizationIAMService) PutOrganizationAccessToken(ctx context.Context, req *organization_iam.PutOrganizationAccessTokenRequest) (*organization_iam.PutOrganizationAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	saved, err := s.repository.GetOrgAccessTokenByUniqueKey(ctx, req.OrganizationId, req.Name)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}

	var accessTokenID uint32
	var tokenHash string
	if !noRecord {
		accessTokenID = saved.AccessTokenID
		tokenHash = saved.TokenHash
	} else {
		tokenHash = hashOrgToken(req.PlainTextToken)
	}

	var expiredAt time.Time
	if req.ExpiredAt == 0 {
		expiredAt = time.Unix(maxOrgTokenExpiredAtUnix, 0)
	} else {
		expiredAt = time.Unix(req.ExpiredAt, 0)
	}

	token := &model.OrgAccessToken{
		AccessTokenID:     accessTokenID,
		TokenHash:         tokenHash,
		Name:              req.Name,
		Description:       req.Description,
		OrgID:             req.OrganizationId,
		ExpiredAt:         expiredAt,
		LastUpdatedUserID: req.LastUpdatedUserId,
	}

	registered, err := s.repository.PutOrgAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &organization_iam.PutOrganizationAccessTokenResponse{AccessToken: convertOrgAccessToken(registered)}, nil
}

func (s *OrganizationIAMService) DeleteOrganizationAccessToken(ctx context.Context, req *organization_iam.DeleteOrganizationAccessTokenRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOrgAccessToken(ctx, req.OrganizationId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func hashOrgToken(plainText string) string {
	hash := sha512.Sum512([]byte(plainText))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
