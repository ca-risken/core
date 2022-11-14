package mocks

import (
	"context"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/mock"
)

/*
 * Mock Repository
 */
type MockFindingRepository struct {
	mock.Mock
}

// Finding

func (m *MockFindingRepository) ListFinding(context.Context, *finding.ListFindingRequest) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}
func (m *MockFindingRepository) BatchListFinding(context.Context, *finding.BatchListFindingRequest) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}
func (m *MockFindingRepository) ListFindingCount(
	ctx context.Context,
	projectID, alertID uint32,
	fromScore, toScore float32,
	fromAt, toAt int64,
	findingID uint64,
	dataSources, resourceNames, tags []string,
	status finding.FindingStatus) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockFindingRepository) GetFinding(context.Context, uint32, uint64, bool) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *MockFindingRepository) GetFindingByDataSource(context.Context, uint32, string, string) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *MockFindingRepository) UpsertFinding(context.Context, *model.Finding) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *MockFindingRepository) DeleteFinding(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}
func (m *MockFindingRepository) ListFindingTagByFindingID(ctx context.Context, projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}
func (m *MockFindingRepository) ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockFindingRepository) ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]db.TagName, error) {
	args := m.Called()
	return args.Get(0).(*[]db.TagName), args.Error(1)
}
func (m *MockFindingRepository) ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockFindingRepository) GetFindingTagByKey(context.Context, uint32, uint64, string) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *MockFindingRepository) GetFindingTagByID(context.Context, uint32, uint64) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *MockFindingRepository) TagFinding(context.Context, *model.FindingTag) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *MockFindingRepository) UntagFinding(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) BulkUpsertFinding(ctx context.Context, data []*model.Finding) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) BulkUpsertFindingTag(ctx context.Context, data []*model.FindingTag) error {
	args := m.Called()
	return args.Error(0)
}

// Resource

func (m *MockFindingRepository) ListResource(context.Context, *finding.ListResourceRequest) (*[]model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Resource), args.Error(1)
}
func (m *MockFindingRepository) ListResourceCount(ctx context.Context, req *finding.ListResourceRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockFindingRepository) GetResource(context.Context, uint32, uint64) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *MockFindingRepository) GetResourceByName(context.Context, uint32, string) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *MockFindingRepository) UpsertResource(context.Context, *model.Resource) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *MockFindingRepository) DeleteResource(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) ListResourceTag(ctx context.Context, param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ResourceTag), args.Error(1)
}
func (m *MockFindingRepository) ListResourceTagByResourceID(ctx context.Context, projectID uint32, resourceID uint64) (*[]model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ResourceTag), args.Error(1)
}
func (m *MockFindingRepository) ListResourceTagCount(ctx context.Context, param *finding.ListResourceTagRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockFindingRepository) ListResourceTagName(ctx context.Context, param *finding.ListResourceTagNameRequest) (*[]db.TagName, error) {
	args := m.Called()
	return args.Get(0).(*[]db.TagName), args.Error(1)
}
func (m *MockFindingRepository) ListResourceTagNameCount(ctx context.Context, param *finding.ListResourceTagNameRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockFindingRepository) GetResourceTagByKey(context.Context, uint32, uint64, string) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *MockFindingRepository) GetResourceTagByID(context.Context, uint32, uint64) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *MockFindingRepository) TagResource(context.Context, *model.ResourceTag) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *MockFindingRepository) UntagResource(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) BulkUpsertResource(ctx context.Context, data []*model.Resource) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) BulkUpsertResourceTag(ctx context.Context, data []*model.ResourceTag) error {
	args := m.Called()
	return args.Error(0)
}

// PendFinding

func (m *MockFindingRepository) GetPendFinding(context.Context, uint32, uint64) (*model.PendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.PendFinding), args.Error(1)
}
func (m *MockFindingRepository) UpsertPendFinding(context.Context, *finding.PendFindingForUpsert) (*model.PendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.PendFinding), args.Error(1)
}
func (m *MockFindingRepository) DeletePendFinding(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}

// FindingSetting

func (m *MockFindingRepository) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingSetting), args.Error(1)
}
func (m *MockFindingRepository) GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *MockFindingRepository) GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *MockFindingRepository) UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *MockFindingRepository) DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}

// Recommend

func (m *MockFindingRepository) GetRecommend(ctx context.Context, projectID uint32, findingID uint64) (*model.Recommend, error) {
	args := m.Called()
	return args.Get(0).(*model.Recommend), args.Error(1)
}
func (m *MockFindingRepository) UpsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error) {
	args := m.Called()
	return args.Get(0).(*model.Recommend), args.Error(1)
}
func (m *MockFindingRepository) UpsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RecommendFinding), args.Error(1)
}
func (m *MockFindingRepository) GetRecommendByDataSourceType(ctx context.Context, dataSource, recommendType string) (*model.Recommend, error) {
	args := m.Called()
	return args.Get(0).(*model.Recommend), args.Error(1)
}
func (m *MockFindingRepository) BulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFindingRepository) BulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error {
	args := m.Called()
	return args.Error(0)
}
