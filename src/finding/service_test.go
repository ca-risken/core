package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/src/finding/model"
	"github.com/stretchr/testify/mock"
)

func TestCalculateScore(t *testing.T) {
	cases := []struct {
		name    string
		input   [2]float32
		setting *findingSetting
		want    float32
	}{
		{
			name:  "OK Score 1%",
			input: [2]float32{1.0, 100.0},
			want:  0.01,
		},
		{
			name:  "OK Score 100%",
			input: [2]float32{100.0, 100.0},
			want:  1.00,
		},
		{
			name:  "OK Score 0%",
			input: [2]float32{0, 100.0},
			want:  0.00,
		},
		{
			name:    "OK Setting x1",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: 1.0},
			want:    0.1,
		},
		{
			name:    "OK Setting x1.5",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: 1.5},
			want:    0.15,
		},
		{
			name:    "OK Setting x100",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: 100},
			want:    1.0,
		},
		{
			name:    "OK Setting x-1",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: -1},
			want:    0.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := calculateScore(c.input[0], c.input[1], c.setting)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * Mock Repository
 */
type mockFindingRepository struct {
	mock.Mock
}

// Finding

func (m *mockFindingRepository) ListFinding(context.Context, *finding.ListFindingRequest) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}
func (m *mockFindingRepository) BatchListFinding(context.Context, *finding.BatchListFindingRequest) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}
func (m *mockFindingRepository) ListFindingCount(
	ctx context.Context,
	projectID uint32,
	fromScore, toScore float32,
	fromAt, toAt int64,
	findingID uint64,
	dataSources, resourceNames, tags []string,
	status finding.FindingStatus) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockFindingRepository) GetFinding(context.Context, uint32, uint64, bool) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *mockFindingRepository) GetFindingByDataSource(context.Context, uint32, string, string) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *mockFindingRepository) UpsertFinding(context.Context, *model.Finding) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *mockFindingRepository) DeleteFinding(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockFindingRepository) ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]tagName, error) {
	args := m.Called()
	return args.Get(0).(*[]tagName), args.Error(1)
}
func (m *mockFindingRepository) ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockFindingRepository) GetFindingTagByKey(context.Context, uint32, uint64, string) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) GetFindingTagByID(context.Context, uint32, uint64) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) TagFinding(context.Context, *model.FindingTag) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) UntagFinding(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) BulkUpsertFinding(ctx context.Context, data []*model.Finding) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) BulkUpsertFindingTag(ctx context.Context, data []*model.FindingTag) error {
	args := m.Called()
	return args.Error(0)
}

// Resource

func (m *mockFindingRepository) ListResource(context.Context, *finding.ListResourceRequest) (*[]model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Resource), args.Error(1)
}
func (m *mockFindingRepository) ListResourceCount(ctx context.Context, req *finding.ListResourceRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockFindingRepository) GetResource(context.Context, uint32, uint64) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *mockFindingRepository) GetResourceByName(context.Context, uint32, string) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *mockFindingRepository) UpsertResource(context.Context, *model.Resource) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *mockFindingRepository) DeleteResource(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) ListResourceTag(ctx context.Context, param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) ListResourceTagCount(ctx context.Context, param *finding.ListResourceTagRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockFindingRepository) ListResourceTagName(ctx context.Context, param *finding.ListResourceTagNameRequest) (*[]tagName, error) {
	args := m.Called()
	return args.Get(0).(*[]tagName), args.Error(1)
}
func (m *mockFindingRepository) ListResourceTagNameCount(ctx context.Context, param *finding.ListResourceTagNameRequest) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockFindingRepository) GetResourceTagByKey(context.Context, uint32, uint64, string) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) GetResourceTagByID(context.Context, uint32, uint64) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) TagResource(context.Context, *model.ResourceTag) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) UntagResource(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) BulkUpsertResource(ctx context.Context, data []*model.Resource) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) BulkUpsertResourceTag(ctx context.Context, data []*model.ResourceTag) error {
	args := m.Called()
	return args.Error(0)
}

// PendFinding

func (m *mockFindingRepository) GetPendFinding(context.Context, uint32, uint64) (*model.PendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.PendFinding), args.Error(1)
}
func (m *mockFindingRepository) UpsertPendFinding(context.Context, *finding.PendFindingForUpsert) (*model.PendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.PendFinding), args.Error(1)
}
func (m *mockFindingRepository) DeletePendFinding(context.Context, uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}

// FindingSetting

func (m *mockFindingRepository) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}

// Recommend

func (m *mockFindingRepository) GetRecommend(ctx context.Context, projectID uint32, findingID uint64) (*model.Recommend, error) {
	args := m.Called()
	return args.Get(0).(*model.Recommend), args.Error(1)
}
func (m *mockFindingRepository) UpsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error) {
	args := m.Called()
	return args.Get(0).(*model.Recommend), args.Error(1)
}
func (m *mockFindingRepository) UpsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RecommendFinding), args.Error(1)
}
func (m *mockFindingRepository) GetRecommendByDataSourceType(ctx context.Context, dataSource, recommendType string) (*model.Recommend, error) {
	args := m.Called()
	return args.Get(0).(*model.Recommend), args.Error(1)
}
func (m *mockFindingRepository) BulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) BulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error {
	args := m.Called()
	return args.Error(0)
}
