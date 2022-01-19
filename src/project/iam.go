package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/core/proto/iam"
	"github.com/gassara-kys/envconfig"
	"google.golang.org/grpc"
)

type iamConfig struct {
	IAMSvcAddr string `required:"true" split_words:"true" default:"iam.core.svc.cluster.local:8002"`
}

func newIAMService() iamService {
	var conf iamConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("project config load error: err=%+v", err)
	}
	ctx := context.Background()
	return &iamServiceImpl{
		client: iam.NewIAMServiceClient(getGRPCConn(ctx, conf.IAMSvcAddr)),
	}

}

func getGRPCConn(ctx context.Context, addr string) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()), grpc.WithInsecure())
	if err != nil {
		appLogger.Fatalf("Failed to connect backend gRPC server, addr=%s, err=%+v", addr, err)
	}
	return conn
}

type iamService interface {
	CreateDefaultRole(ctx context.Context, ownerUserID, projectID uint32) error
	DeleteAllProjectRole(ctx context.Context, projectID uint32) error
	IsActiveProject(ctx context.Context, projectID uint32) (bool, error)
}

type iamServiceImpl struct {
	client iam.IAMServiceClient
}

func (i *iamServiceImpl) CreateDefaultRole(ctx context.Context, ownerUserID, projectID uint32) error {
	policy, err := i.client.PutPolicy(ctx, &iam.PutPolicyRequest{
		ProjectId: projectID,
		Policy: &iam.PolicyForUpsert{
			Name:        "project-admin",
			ProjectId:   projectID,
			ActionPtn:   ".*",
			ResourcePtn: ".*",
		},
	})
	if err != nil {
		return fmt.Errorf("Could not put default policy, err=%+v", err)
	}
	role, err := i.client.PutRole(ctx, &iam.PutRoleRequest{
		ProjectId: projectID,
		Role: &iam.RoleForUpsert{
			Name:      "project-admin-role",
			ProjectId: projectID,
		},
	})
	if err != nil {
		return fmt.Errorf("Could not put project-admin-role, err=%+v", err)
	}
	if _, err := i.client.AttachPolicy(ctx, &iam.AttachPolicyRequest{
		ProjectId: projectID,
		RoleId:    role.Role.RoleId,
		PolicyId:  policy.Policy.PolicyId,
	}); err != nil {
		return fmt.Errorf("Could not attach default policy, err=%+v", err)
	}
	if _, err := i.client.AttachRole(ctx, &iam.AttachRoleRequest{
		ProjectId: projectID,
		UserId:    ownerUserID,
		RoleId:    role.Role.RoleId,
	}); err != nil {
		return fmt.Errorf("Could not attach default role, err=%+v", err)
	}
	return nil
}

func (i *iamServiceImpl) DeleteAllProjectRole(ctx context.Context, projectID uint32) error {
	list, err := i.client.ListRole(ctx, &iam.ListRoleRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	for _, roleID := range list.RoleId {
		if _, err := i.client.DeleteRole(ctx, &iam.DeleteRoleRequest{
			ProjectId: projectID,
			RoleId:    roleID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (i *iamServiceImpl) IsActiveProject(ctx context.Context, projectID uint32) (bool, error) {
	resp, err := i.client.ListUser(ctx, &iam.ListUserRequest{
		ProjectId: projectID,
		Activated: true,
	})
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, nil
	}
	return len(resp.UserId) > 0, nil
}
