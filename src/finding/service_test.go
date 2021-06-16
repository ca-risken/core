package main

import (
	"reflect"
	"testing"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
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

func (m *mockFindingRepository) ListFinding(*finding.ListFindingRequest) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}
func (m *mockFindingRepository) ListFindingCount(req *finding.ListFindingRequest) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}
func (m *mockFindingRepository) GetFinding(uint32, uint64) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *mockFindingRepository) GetFindingByDataSource(uint32, string, string) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *mockFindingRepository) UpsertFinding(*model.Finding) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
func (m *mockFindingRepository) DeleteFinding(uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) ListFindingTag(param *finding.ListFindingTagRequest) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) ListFindingTagCount(param *finding.ListFindingTagRequest) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}
func (m *mockFindingRepository) ListFindingTagName(param *finding.ListFindingTagNameRequest) (*[]tagName, error) {
	args := m.Called()
	return args.Get(0).(*[]tagName), args.Error(1)
}
func (m *mockFindingRepository) ListFindingTagNameCount(param *finding.ListFindingTagNameRequest) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}
func (m *mockFindingRepository) GetFindingTagByKey(uint32, uint64, string) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) GetFindingTagByID(uint32, uint64) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) TagFinding(*model.FindingTag) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}
func (m *mockFindingRepository) UntagFinding(uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}

// Resource

func (m *mockFindingRepository) ListResource(*finding.ListResourceRequest) (*[]model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Resource), args.Error(1)
}
func (m *mockFindingRepository) ListResourceCount(req *finding.ListResourceRequest) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}
func (m *mockFindingRepository) GetResource(uint32, uint64) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *mockFindingRepository) GetResourceByName(uint32, string) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *mockFindingRepository) UpsertResource(*model.Resource) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
func (m *mockFindingRepository) DeleteResource(uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockFindingRepository) ListResourceTag(param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) ListResourceTagCount(param *finding.ListResourceTagRequest) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}
func (m *mockFindingRepository) ListResourceTagName(param *finding.ListResourceTagNameRequest) (*[]tagName, error) {
	args := m.Called()
	return args.Get(0).(*[]tagName), args.Error(1)
}
func (m *mockFindingRepository) ListResourceTagNameCount(param *finding.ListResourceTagNameRequest) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}
func (m *mockFindingRepository) GetResourceTagByKey(uint32, uint64, string) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) GetResourceTagByID(uint32, uint64) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) TagResource(*model.ResourceTag) (*model.ResourceTag, error) {
	args := m.Called()
	return args.Get(0).(*model.ResourceTag), args.Error(1)
}
func (m *mockFindingRepository) UntagResource(uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}

// PendFinding

func (m *mockFindingRepository) GetPendFinding(uint32, uint64) (*model.PendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.PendFinding), args.Error(1)
}
func (m *mockFindingRepository) UpsertPendFinding(*finding.PendFindingForUpsert) (*model.PendFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.PendFinding), args.Error(1)
}
func (m *mockFindingRepository) DeletePendFinding(uint32, uint64) error {
	args := m.Called()
	return args.Error(0)
}

// FindingSetting

func (m *mockFindingRepository) ListFindingSetting(req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) GetFindingSetting(projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) GetFindingSettingByResource(projectID uint32, resourceName string) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) UpsertFindingSetting(data *model.FindingSetting) (*model.FindingSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingSetting), args.Error(1)
}
func (m *mockFindingRepository) DeleteFindingSetting(projectID uint32, findingSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
