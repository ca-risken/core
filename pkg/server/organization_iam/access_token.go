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
	"gorm.io/gorm"
)

const (
	maxOrgTokenExpiredAtUnix int64 = 253402268399 // 9999-12-31T23:59:59
)

func (s *OrganizationIAMService) ListOrgAccessToken(ctx context.Context, req *organization_iam.ListOrgAccessTokenRequest) (*organization_iam.ListOrgAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	resp, err := s.repository.ListOrgAccessToken(ctx, req.OrganizationId, req.Name, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.ListOrgAccessTokenResponse{}, nil
		}
		return nil, err
	}
	var list []*organization_iam.OrganizationAccessToken
	for _, token := range *resp {
		list = append(list, convertOrgAccessToken(&token))
	}
	return &organization_iam.ListOrgAccessTokenResponse{AccessToken: list}, nil
}

func convertOrgAccessToken(token *model.OrganizationAccessToken) *organization_iam.OrganizationAccessToken {
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

func (s *OrganizationIAMService) AuthenticateOrgAccessToken(ctx context.Context, req *organization_iam.AuthenticateOrgAccessTokenRequest) (*organization_iam.AuthenticateOrgAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	token, err := s.repository.GetActiveOrgAccessTokenHash(ctx, req.OrganizationId, req.AccessTokenId, hashOrgToken(req.PlainTextToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.AuthenticateOrgAccessTokenResponse{}, nil
		}
		return nil, err
	}
	return &organization_iam.AuthenticateOrgAccessTokenResponse{AccessToken: convertOrgAccessToken(token)}, nil
}

func (s *OrganizationIAMService) PutOrgAccessToken(ctx context.Context, req *organization_iam.PutOrgAccessTokenRequest) (*organization_iam.PutOrgAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	saved, err := s.repository.GetOrgAccessTokenByUniqueKey(ctx, req.OrganizationId, req.AccessToken.Name)
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
		tokenHash = hashOrgToken(req.AccessToken.PlainTextToken)
	}

	var expiredAt time.Time
	if req.AccessToken.ExpiredAt == 0 {
		expiredAt = time.Unix(maxOrgTokenExpiredAtUnix, 0)
	} else {
		expiredAt = time.Unix(req.AccessToken.ExpiredAt, 0)
	}

	token := &model.OrganizationAccessToken{
		AccessTokenID:     accessTokenID,
		TokenHash:         tokenHash,
		Name:              req.AccessToken.Name,
		Description:       req.AccessToken.Description,
		OrgID:             req.AccessToken.OrganizationId,
		ExpiredAt:         expiredAt,
		LastUpdatedUserID: req.AccessToken.LastUpdatedUserId,
	}

	registered, err := s.repository.PutOrgAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &organization_iam.PutOrgAccessTokenResponse{AccessToken: convertOrgAccessToken(registered)}, nil
}

func (s *OrganizationIAMService) DeleteOrgAccessToken(ctx context.Context, req *organization_iam.DeleteOrgAccessTokenRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOrgAccessToken(ctx, req.OrganizationId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OrganizationIAMService) AttachOrgAccessTokenRole(ctx context.Context, req *organization_iam.AttachOrgAccessTokenRoleRequest) (*organization_iam.AttachOrgAccessTokenRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := s.repository.AttachOrgAccessTokenRole(ctx, req.OrganizationId, req.RoleId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	return &organization_iam.AttachOrgAccessTokenRoleResponse{AccessTokenRole: convertOrgAccessTokenRole(role)}, nil
}

func (s *OrganizationIAMService) DetachOrgAccessTokenRole(ctx context.Context, req *organization_iam.DetachOrgAccessTokenRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DetachOrgAccessTokenRole(ctx, req.OrganizationId, req.RoleId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertOrgAccessTokenRole(role *model.OrganizationAccessTokenRole) *organization_iam.OrganizationAccessTokenRole {
	return &organization_iam.OrganizationAccessTokenRole{
		AccessTokenId: role.AccessTokenID,
		RoleId:        role.RoleID,
		CreatedAt:     role.CreatedAt.Unix(),
		UpdatedAt:     role.UpdatedAt.Unix(),
	}
}

func hashOrgToken(plainText string) string {
	hash := sha512.Sum512([]byte(plainText))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
