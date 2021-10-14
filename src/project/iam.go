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
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		appLogger.Fatalf("Failed to connect backend gRPC server, addr=%s, err=%+v", addr, err)
	}
	return conn
}

type iamService interface {
	CreateDefaultRole(context.Context, uint32, uint32) error
	DeleteAllProjectRole(context.Context, uint32) error
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
