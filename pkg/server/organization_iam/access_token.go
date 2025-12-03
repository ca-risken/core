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

func convertOrgAccessToken(token *model.OrganizationAccessToken) *organization_iam.OrganizationAccessToken {
	return &organization_iam.OrganizationAccessToken{
		AccessTokenId:     token.AccessTokenID,
		Name:              token.Name,
		Description:       token.Description,
		OrganizationId:    token.OrganizationID,
		ExpiredAt:         token.ExpiredAt.Unix(),
		LastUpdatedUserId: token.LastUpdatedUserID,
		CreatedAt:         token.CreatedAt.Unix(),
		UpdatedAt:         token.UpdatedAt.Unix(),
	}
}

func convertOrgAccessTokenRole(role *model.OrganizationAccessTokenRole) *organization_iam.OrganizationAccessTokenRole {
	return &organization_iam.OrganizationAccessTokenRole{
		AccessTokenId: role.AccessTokenID,
		RoleId:        role.RoleID,
		CreatedAt:     role.CreatedAt.Unix(),
		UpdatedAt:     role.UpdatedAt.Unix(),
	}
}

func (s *OrganizationIAMService) PutOrganizationAccessToken(ctx context.Context, req *organization_iam.PutOrganizationAccessTokenRequest) (*organization_iam.PutOrganizationAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
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

	token := &model.OrganizationAccessToken{
		AccessTokenID:     accessTokenID,
		TokenHash:         tokenHash,
		Name:              req.Name,
		Description:       req.Description,
		OrganizationID:    req.OrganizationId,
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

func (s *OrganizationIAMService) AuthenticateOrganizationAccessToken(ctx context.Context, req *organization_iam.AuthenticateOrganizationAccessTokenRequest) (*organization_iam.AuthenticateOrganizationAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	token, err := s.repository.GetActiveOrgAccessTokenHash(ctx, req.OrganizationId, req.AccessTokenId, hashOrgToken(req.PlainTextToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.AuthenticateOrganizationAccessTokenResponse{}, nil
		}
		return nil, err
	}
	return &organization_iam.AuthenticateOrganizationAccessTokenResponse{AccessToken: convertOrgAccessToken(token)}, nil
}

func (s *OrganizationIAMService) AttachOrganizationAccessTokenRole(ctx context.Context, req *organization_iam.AttachOrganizationAccessTokenRoleRequest) (*organization_iam.AttachOrganizationAccessTokenRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := s.repository.AttachOrgAccessTokenRole(ctx, req.OrganizationId, req.RoleId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	return &organization_iam.AttachOrganizationAccessTokenRoleResponse{AccessTokenRole: convertOrgAccessTokenRole(role)}, nil
}

func (s *OrganizationIAMService) DetachOrganizationAccessTokenRole(ctx context.Context, req *organization_iam.DetachOrganizationAccessTokenRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DetachOrgAccessTokenRole(ctx, req.OrganizationId, req.RoleId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func hashOrgToken(plainText string) string {
	hash := sha512.Sum512([]byte(plainText))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
