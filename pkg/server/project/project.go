package project

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func convertProjectWithTag(p *db.ProjectWithTag) *project.Project {
	if p == nil {
		return &project.Project{}
	}
	tags := []*project.ProjectTag{}
	if p.Tag != nil {
		for _, t := range *p.Tag {
			tags = append(tags, &project.ProjectTag{
				ProjectId: t.ProjectID,
				Tag:       t.Tag,
				Color:     t.Color,
			})
		}
	}
	return &project.Project{
		ProjectId: p.ProjectID,
		Name:      p.Name,
		Tag:       tags,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
}

func (p *ProjectService) ListProject(ctx context.Context, req *project.ListProjectRequest) (*project.ListProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := p.repository.ListProject(ctx, req.UserId, req.ProjectId, req.OrganizationId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &project.ListProjectResponse{}, nil
		}
		return nil, err
	}
	var prs []*project.Project
	for _, pr := range *list {
		prs = append(prs, convertProjectWithTag(&pr))
	}
	return &project.ListProjectResponse{Project: prs}, nil
}

func convertProject(p *model.Project) *project.Project {
	return &project.Project{
		ProjectId: p.ProjectID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
}

func (p *ProjectService) CreateProject(ctx context.Context, req *project.CreateProjectRequest) (*project.CreateProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	ownerUserID, err := resolveCreateProjectOwnerUserID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	req.UserId = ownerUserID

	pr, err := p.repository.CreateProject(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if err := p.createDefaultRole(ctx, req.UserId, pr.ProjectID); err != nil {
		return nil, err
	}
	p.logger.Infof(ctx, "Project created: owner=%d, project=%+v", req.UserId, pr)

	return &project.CreateProjectResponse{Project: convertProject(pr)}, nil
}

func resolveCreateProjectOwnerUserID(ctx context.Context, requestUserID uint32) (uint32, error) {
	if ctx == nil {
		return requestUserID, nil
	}
	pr, ok := peer.FromContext(ctx)
	if !ok || pr == nil || pr.Addr == nil {
		return requestUserID, nil
	}
	if isLoopbackPeerAddr(pr.Addr) {
		return requestUserID, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.PermissionDenied, "caller user_id metadata is required")
	}
	var callerUserIDRaw string
	for _, key := range []string{"x-risken-user-id", "x-user-id"} {
		values := md.Get(key)
		if len(values) > 0 {
			callerUserIDRaw = values[0]
			break
		}
	}
	if callerUserIDRaw == "" {
		return 0, status.Error(codes.PermissionDenied, "caller user_id metadata is required")
	}
	callerUserID, err := strconv.ParseUint(callerUserIDRaw, 10, 32)
	if err != nil || callerUserID == 0 {
		return 0, status.Error(codes.InvalidArgument, "invalid caller user_id metadata")
	}
	return uint32(callerUserID), nil
}

func isLoopbackPeerAddr(addr net.Addr) bool {
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		host = addr.String()
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func (p *ProjectService) UpdateProject(ctx context.Context, req *project.UpdateProjectRequest) (*project.UpdateProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	pr, err := p.repository.UpdateProject(ctx, req.ProjectId, req.Name)
	if err != nil {
		return nil, err
	}
	return &project.UpdateProjectResponse{Project: convertProject(pr)}, nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, req *project.DeleteProjectRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := p.repository.DeleteProject(ctx, req.ProjectId); err != nil {
		return nil, err
	}
	p.logger.Infof(ctx, "Project deleted: project=%+v", req.ProjectId)

	return &empty.Empty{}, nil
}

func (p *ProjectService) IsActive(ctx context.Context, req *project.IsActiveRequest) (*project.IsActiveResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	active, err := p.isActiveProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return &project.IsActiveResponse{Active: active}, nil
}

func (p *ProjectService) createDefaultRole(ctx context.Context, ownerUserID, projectID uint32) error {
	projectAdmin := "project-admin"
	projectViewer := "project-viewer"
	findingEditor := "finding-editor"
	viewerActionPtn := "get|list"

	for name, actionPtn := range map[string]string{
		projectAdmin:  ".*",
		projectViewer: viewerActionPtn,
		findingEditor: viewerActionPtn + "|^finding/.+|^alert/.+",
	} {
		policy, err := p.iamClient.PutPolicy(ctx, &iam.PutPolicyRequest{
			ProjectId: projectID,
			Policy: &iam.PolicyForUpsert{
				Name:        name,
				ProjectId:   projectID,
				ActionPtn:   actionPtn,
				ResourcePtn: ".*",
			},
		})
		if err != nil {
			return fmt.Errorf("could not put %s-policy, err=%w", name, err)
		}
		role, err := p.iamClient.PutRole(ctx, &iam.PutRoleRequest{
			ProjectId: projectID,
			Role: &iam.RoleForUpsert{
				Name:      name + "-role",
				ProjectId: projectID,
			},
		})
		if err != nil {
			return fmt.Errorf("could not put %s-role, err=%w", name, err)
		}
		if _, err := p.iamClient.AttachPolicy(ctx, &iam.AttachPolicyRequest{
			ProjectId: projectID,
			RoleId:    role.Role.RoleId,
			PolicyId:  policy.Policy.PolicyId,
		}); err != nil {
			return fmt.Errorf("could not attach %s-policy to %s-role, err=%w", name, name, err)
		}
		if name == projectAdmin {
			if _, err := p.iamClient.AttachRole(ctx, &iam.AttachRoleRequest{
				ProjectId: projectID,
				UserId:    ownerUserID,
				RoleId:    role.Role.RoleId,
			}); err != nil {
				return fmt.Errorf("could not attach default %s-role to project owner, err=%w", name, err)
			}
		}
	}
	return nil
}

func (p *ProjectService) isActiveProject(ctx context.Context, projectID uint32) (bool, error) {
	projects, err := p.repository.ListProject(ctx, 0, projectID, 0, "")
	if err != nil {
		return false, err
	}
	if len(*projects) == 0 {
		return false, nil
	}
	resp, err := p.iamClient.ListUser(ctx, &iam.ListUserRequest{
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

func (p *ProjectService) CleanProject(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	if err := p.repository.CleanWithNoProject(ctx); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
