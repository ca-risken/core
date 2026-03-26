package org_iam

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/org_iam"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

const (
	maxOrgTokenExpiredAtUnix int64 = 253402268399 // 9999-12-31T23:59:59
)

func (s *OrgIAMService) ListOrgAccessToken(ctx context.Context, req *org_iam.ListOrgAccessTokenRequest) (*org_iam.ListOrgAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	resp, err := s.repository.ListOrgAccessToken(ctx, req.OrganizationId, req.Name, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.ListOrgAccessTokenResponse{}, nil
		}
		return nil, err
	}
	var list []*org_iam.OrgAccessToken
	for _, token := range *resp {
		list = append(list, convertOrgAccessToken(&token))
	}
	return &org_iam.ListOrgAccessTokenResponse{AccessToken: list}, nil
}

func convertOrgAccessToken(token *model.OrganizationAccessToken) *org_iam.OrgAccessToken {
	return &org_iam.OrgAccessToken{
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

func convertOrgAccessTokenRole(role *model.OrganizationAccessTokenRole) *org_iam.OrgAccessTokenRole {
	return &org_iam.OrgAccessTokenRole{
		AccessTokenId: role.AccessTokenID,
		RoleId:        role.RoleID,
		CreatedAt:     role.CreatedAt.Unix(),
		UpdatedAt:     role.UpdatedAt.Unix(),
	}
}

func (s *OrgIAMService) PutOrgAccessToken(ctx context.Context, req *org_iam.PutOrgAccessTokenRequest) (*org_iam.PutOrgAccessTokenResponse, error) {
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
	return &org_iam.PutOrgAccessTokenResponse{AccessToken: convertOrgAccessToken(registered)}, nil
}

func (s *OrgIAMService) DeleteOrgAccessToken(ctx context.Context, req *org_iam.DeleteOrgAccessTokenRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOrgAccessToken(ctx, req.OrganizationId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OrgIAMService) AuthenticateOrgAccessToken(ctx context.Context, req *org_iam.AuthenticateOrgAccessTokenRequest) (*org_iam.AuthenticateOrgAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	token, err := s.repository.GetActiveOrgAccessTokenHash(ctx, req.OrganizationId, req.AccessTokenId, hashOrgToken(req.PlainTextToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.AuthenticateOrgAccessTokenResponse{}, nil
		}
		return nil, err
	}
	return &org_iam.AuthenticateOrgAccessTokenResponse{AccessToken: convertOrgAccessToken(token)}, nil
}

func (s *OrgIAMService) AttachOrgAccessTokenRole(ctx context.Context, req *org_iam.AttachOrgAccessTokenRoleRequest) (*org_iam.AttachOrgAccessTokenRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := s.repository.AttachOrgAccessTokenRole(ctx, req.OrganizationId, req.RoleId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	return &org_iam.AttachOrgAccessTokenRoleResponse{AccessTokenRole: convertOrgAccessTokenRole(role)}, nil
}

func (s *OrgIAMService) DetachOrgAccessTokenRole(ctx context.Context, req *org_iam.DetachOrgAccessTokenRoleRequest) (*empty.Empty, error) {
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
