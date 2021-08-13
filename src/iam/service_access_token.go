package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (i *iamService) ListAccessToken(ctx context.Context, req *iam.ListAccessTokenRequest) (*iam.ListAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	resp, err := i.repository.ListAccessToken(ctx, req.ProjectId, req.Name, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.ListAccessTokenResponse{}, nil
		}
		return nil, err
	}
	var list []*iam.AccessToken
	for _, a := range *resp {
		list = append(list, convertAccessToken(&a))
	}
	return &iam.ListAccessTokenResponse{AccessToken: list}, nil
}

func convertAccessToken(a *model.AccessToken) *iam.AccessToken {
	return &iam.AccessToken{
		AccessTokenId:     a.AccessTokenID,
		Name:              a.Name,
		Description:       a.Description,
		ProjectId:         a.ProjectID,
		ExpiredAt:         a.ExpiredAt.Unix(),
		LastUpdatedUserId: a.LastUpdatedUserID,
		CreatedAt:         a.CreatedAt.Unix(),
		UpdatedAt:         a.UpdatedAt.Unix(),
	}
}

func (i *iamService) AuthenticateAccessToken(ctx context.Context, req *iam.AuthenticateAccessTokenRequest) (*iam.AuthenticateAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	token, err := i.repository.GetActiveAccessTokenHash(ctx, req.ProjectId, req.AccessTokenId, hash(req.PlainTextToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.AuthenticateAccessTokenResponse{}, nil
		}
		return nil, err
	}
	return &iam.AuthenticateAccessTokenResponse{AccessToken: convertAccessToken(token)}, nil
}

func hash(plainText string) string {
	hash := sha512.Sum512([]byte(plainText))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

const (
	maxExpiredAtUnix int64 = 253402268399 // 9999-12-31T23:59:59
)

func (i *iamService) PutAccessToken(ctx context.Context, req *iam.PutAccessTokenRequest) (*iam.PutAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetAccessTokenByUniqueKey(ctx, req.AccessToken.ProjectId, req.AccessToken.Name)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}

	var accessTokenID uint32
	var tokenHash string
	if !noRecord {
		accessTokenID = savedData.AccessTokenID
		tokenHash = savedData.TokenHash
	} else {
		tokenHash = hash(req.AccessToken.PlainTextToken) // Create new token
	}
	var expiredAt time.Time
	if zero.IsZeroVal(req.AccessToken.ExpiredAt) {
		expiredAt = time.Unix(maxExpiredAtUnix, 0) // Update max expired_at
	} else {
		expiredAt = time.Unix(req.AccessToken.ExpiredAt, 0)
	}
	at := &model.AccessToken{
		AccessTokenID:     accessTokenID,
		TokenHash:         tokenHash,
		Name:              req.AccessToken.Name,
		Description:       req.AccessToken.Description,
		ProjectID:         req.AccessToken.ProjectId,
		ExpiredAt:         expiredAt,
		LastUpdatedUserID: req.AccessToken.LastUpdatedUserId,
	}

	// upsert
	registerdData, err := i.repository.PutAccessToken(ctx, at)
	if err != nil {
		return nil, err
	}
	return &iam.PutAccessTokenResponse{AccessToken: convertAccessToken(registerdData)}, nil
}

func (i *iamService) DeleteAccessToken(ctx context.Context, req *iam.DeleteAccessTokenRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteAccessToken(ctx, req.ProjectId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (i *iamService) AttachAccessTokenRole(ctx context.Context, req *iam.AttachAccessTokenRoleRequest) (*iam.AttachAccessTokenRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	tokenRole, err := i.repository.AttachAccessTokenRole(ctx, req.ProjectId, req.RoleId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	return &iam.AttachAccessTokenRoleResponse{AccessTokenRole: convertAccessTokenRole(tokenRole)}, nil
}

func convertAccessTokenRole(a *model.AccessTokenRole) *iam.AccessTokenRole {
	return &iam.AccessTokenRole{
		AccessTokenId: a.AccessTokenID,
		RoleId:        a.RoleID,
		CreatedAt:     a.CreatedAt.Unix(),
		UpdatedAt:     a.UpdatedAt.Unix(),
	}
}

func (i *iamService) DetachAccessTokenRole(ctx context.Context, req *iam.DetachAccessTokenRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DetachAccessTokenRole(ctx, req.ProjectId, req.RoleId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
