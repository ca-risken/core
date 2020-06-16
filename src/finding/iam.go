package main

import (
	"context"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
	"google.golang.org/grpc"
)

type iamConf struct {
	IamSvcAddr string `required:"true" split_words:"true"`
}

func newIAMServiceClient() iam.IAMServiceClient {
	var conf iamConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Unexpected load error to IAM config, err=%+v", err)
		return nil
	}
	ctx := context.Background()
	conn, err := getGRPCClientConn(ctx, conf.IamSvcAddr)
	if err != nil {
		appLogger.Fatalf("Could not connected IAM service, err=%+v", err)
		return nil
	}
	appLogger.Infof("Connected to IAM Service... %s", conf.IamSvcAddr)
	return iam.NewIAMServiceClient(conn)
}

func getGRPCClientConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}
	return conn, err
}

/**
 * Finding
 */

const prefixFinding = "finding/"

func (f *findingService) callIsAuthorizedForFinding(ctx context.Context, userID uint32, apiName string, data *model.Finding) bool {
	authzRequest := iam.IsAuthorizedRequest{
		UserId:       userID,
		ProjectId:    data.ProjectID,
		ActionName:   prefixFinding + apiName,                                   // e.g. `finding/GetFinding`
		ResourceName: prefixFinding + data.DataSource + "/" + data.ResourceName, // e.g. `finding/aws:guardduty/bucket-name`
	}
	resp, err := f.iamClient.IsAuthorized(ctx, &authzRequest)
	if err != nil {
		appLogger.Warnf("[callIsAuthorizedForFinding]Something wrong, err=%+v", err)
		return false
	}
	return resp.Ok
}

func (f *findingService) isAuthorizedWithFindingID(ctx context.Context, userID uint32, apiName string, findingID uint64) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	data, err := f.repository.GetFinding(findingID)
	if err != nil {
		appLogger.Warnf("[isAuthorizedToFindingID]Could not get finding data, err=%+v", err)
		return false
	}
	return f.callIsAuthorizedForFinding(ctx, userID, apiName, data)
}

func (f *findingService) isAuthorizedWithFinding(ctx context.Context, userID uint32, apiName string, data *model.Finding) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	return f.callIsAuthorizedForFinding(ctx, userID, apiName, data)
}

/**
 * FindingTag
 */

const prefixFindingTag = "finding-tag/"

func (f *findingService) callIsAuthorizedForFindingTag(ctx context.Context, userID uint32, apiName string, data *model.FindingTag) bool {
	authzRequest := iam.IsAuthorizedRequest{
		UserId:       userID,
		ProjectId:    data.ProjectID,
		ActionName:   prefixFindingTag + apiName,     // e.g. `finding-tag/GetFindingTag`
		ResourceName: prefixFindingTag + data.TagKey, // e.g. `finding-tag/key-1`
	}
	resp, err := f.iamClient.IsAuthorized(ctx, &authzRequest)
	if err != nil {
		appLogger.Warnf("[callIsAuthorizedForFindingTag]Something wrong, err=%+v", err)
		return false
	}
	return resp.Ok
}
func (f *findingService) isAuthorizedWithFindingTagID(ctx context.Context, userID uint32, apiName string, findingTagID uint64) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	data, err := f.repository.GetFindingTagByID(findingTagID)
	if err != nil {
		appLogger.Warnf("[isAuthorizedWithFindingTagID]Could not get finding data, err=%+v", err)
		return false
	}
	return f.callIsAuthorizedForFindingTag(ctx, userID, apiName, data)
}

func (f *findingService) isAuthorizedWithFindingTag(ctx context.Context, userID uint32, apiName string, data *model.FindingTag) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	return f.callIsAuthorizedForFindingTag(ctx, userID, apiName, data)
}

/**
 * Resource
 */

const prefixResource = "resource/"

func (f *findingService) callIsAuthorizedForResource(ctx context.Context, userID uint32, apiName string, data *model.Resource) bool {
	authzRequest := iam.IsAuthorizedRequest{
		UserId:       userID,
		ProjectId:    data.ProjectID,
		ActionName:   prefixResource + apiName,           // e.g. `resource/GetResource`
		ResourceName: prefixResource + data.ResourceName, // e.g. `resource/ec2-instance-id`
	}
	resp, err := f.iamClient.IsAuthorized(ctx, &authzRequest)
	if err != nil {
		appLogger.Warnf("[callIsAuthorizedForResource]Something wrong, err=%+v", err)
		return false
	}
	return resp.Ok
}
func (f *findingService) isAuthorizedWithResourceID(ctx context.Context, userID uint32, apiName string, resourceID uint64) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	data, err := f.repository.GetResource(resourceID)
	if err != nil {
		appLogger.Warnf("[isAuthorizedWithResourceID]Could not get finding data, err=%+v", err)
		return false
	}
	return f.callIsAuthorizedForResource(ctx, userID, apiName, data)
}

func (f *findingService) isAuthorizedWithResource(ctx context.Context, userID uint32, apiName string, data *model.Resource) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	return f.callIsAuthorizedForResource(ctx, userID, apiName, data)
}

/**
 * ResourceTag
 */

const prefixResourceTag = "resource-tag/"

func (f *findingService) callIsAuthorizedForResourceTag(ctx context.Context, userID uint32, apiName string, data *model.ResourceTag) bool {
	authzRequest := iam.IsAuthorizedRequest{
		UserId:       userID,
		ProjectId:    data.ProjectID,
		ActionName:   prefixResourceTag + apiName,     // e.g. `resource/GetResource`
		ResourceName: prefixResourceTag + data.TagKey, // e.g. `resource/ec2-instance-id`
	}
	resp, err := f.iamClient.IsAuthorized(ctx, &authzRequest)
	if err != nil {
		appLogger.Warnf("[callIsAuthorizedForResourceTag]Something wrong, err=%+v", err)
		return false
	}
	return resp.Ok
}
func (f *findingService) isAuthorizedWithResourceTagID(ctx context.Context, userID uint32, apiName string, resourceTagID uint64) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	data, err := f.repository.GetResourceTagByID(resourceTagID)
	if err != nil {
		appLogger.Warnf("[isAuthorizedWithResourceTagID]Could not get finding data, err=%+v", err)
		return false
	}
	return f.callIsAuthorizedForResourceTag(ctx, userID, apiName, data)
}

func (f *findingService) isAuthorizedWithResourceTag(ctx context.Context, userID uint32, apiName string, data *model.ResourceTag) bool {
	if zero.IsZeroVal(userID) {
		return true // no userID(internal service call)
	}
	return f.callIsAuthorizedForResourceTag(ctx, userID, apiName, data)
}
