package iam

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (i *IAMService) ListAccessToken(ctx context.Context, req *iam.ListAccessTokenRequest) (*iam.ListAccessTokenResponse, error) {
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

func (i *IAMService) AuthenticateAccessToken(ctx context.Context, req *iam.AuthenticateAccessTokenRequest) (*iam.AuthenticateAccessTokenResponse, error) {
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

func (i *IAMService) PutAccessToken(ctx context.Context, req *iam.PutAccessTokenRequest) (*iam.PutAccessTokenResponse, error) {
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

func (i *IAMService) DeleteAccessToken(ctx context.Context, req *iam.DeleteAccessTokenRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteAccessToken(ctx, req.ProjectId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (i *IAMService) AttachAccessTokenRole(ctx context.Context, req *iam.AttachAccessTokenRoleRequest) (*iam.AttachAccessTokenRoleResponse, error) {
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

func (i *IAMService) DetachAccessTokenRole(ctx context.Context, req *iam.DetachAccessTokenRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DetachAccessTokenRole(ctx, req.ProjectId, req.RoleId, req.AccessTokenId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

const (
	riskenDataSource = "RISKEN"
	accessTokenTag   = "access-token"
)

func (i *IAMService) AnalyzeTokenExpiration(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	tokens, err := i.repository.ListExpiredAccessToken(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		return nil, err
	}
	// Clear score
	if _, err := i.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: riskenDataSource,
		Tag:        []string{accessTokenTag},
	}); err != nil {
		return nil, err
	}

	// Put finding
	for _, token := range *tokens {
		token.TokenHash = "xxx" // mask credentials
		buf, err := json.Marshal(token)
		if err != nil {
			appLogger.Errorf(ctx, "Failed to encoding json, accessToken=%+v, err=%+v", token, err)
			return nil, err
		}
		resp, err := i.findingClient.PutFinding(ctx, &finding.PutFindingRequest{
			ProjectId: token.ProjectID,
			Finding: &finding.FindingForUpsert{
				Description:      "RISKEN AccessToken expired",
				DataSource:       riskenDataSource,
				DataSourceId:     fmt.Sprintf("risken-access-token-id-%d", token.AccessTokenID),
				ResourceName:     token.Name,
				ProjectId:        token.ProjectID,
				OriginalScore:    0.8,
				OriginalMaxScore: 1.0,
				Data:             string(buf),
			},
		})
		if err != nil {
			return nil, err
		}
		if _, err = i.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
			ProjectId: token.ProjectID,
			Tag: &finding.FindingTagForUpsert{
				FindingId: resp.Finding.FindingId,
				ProjectId: token.ProjectID,
				Tag:       accessTokenTag, // tag
			},
		}); err != nil {
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}
