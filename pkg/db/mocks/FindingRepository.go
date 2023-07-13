// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	db "github.com/ca-risken/core/pkg/db"
	finding "github.com/ca-risken/core/proto/finding"

	mock "github.com/stretchr/testify/mock"

	model "github.com/ca-risken/core/pkg/model"
)

// FindingRepository is an autogenerated mock type for the FindingRepository type
type FindingRepository struct {
	mock.Mock
}

// BatchListFinding provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) BatchListFinding(_a0 context.Context, _a1 *finding.BatchListFindingRequest) (*[]model.Finding, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *[]model.Finding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.BatchListFindingRequest) (*[]model.Finding, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.BatchListFindingRequest) *[]model.Finding); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Finding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.BatchListFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BulkUpsertFinding provides a mock function with given fields: ctx, data
func (_m *FindingRepository) BulkUpsertFinding(ctx context.Context, data []*model.Finding) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Finding) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BulkUpsertFindingTag provides a mock function with given fields: ctx, data
func (_m *FindingRepository) BulkUpsertFindingTag(ctx context.Context, data []*model.FindingTag) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.FindingTag) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BulkUpsertRecommend provides a mock function with given fields: ctx, data
func (_m *FindingRepository) BulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Recommend) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BulkUpsertRecommendFinding provides a mock function with given fields: ctx, data
func (_m *FindingRepository) BulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.RecommendFinding) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BulkUpsertResource provides a mock function with given fields: ctx, data
func (_m *FindingRepository) BulkUpsertResource(ctx context.Context, data []*model.Resource) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Resource) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BulkUpsertResourceTag provides a mock function with given fields: ctx, data
func (_m *FindingRepository) BulkUpsertResourceTag(ctx context.Context, data []*model.ResourceTag) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.ResourceTag) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearScoreFinding provides a mock function with given fields: ctx, req
func (_m *FindingRepository) ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ClearScoreRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFinding provides a mock function with given fields: _a0, _a1, _a2
func (_m *FindingRepository) DeleteFinding(_a0 context.Context, _a1 uint32, _a2 uint64) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFindingSetting provides a mock function with given fields: ctx, projectID, findingSettingID
func (_m *FindingRepository) DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error {
	ret := _m.Called(ctx, projectID, findingSettingID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(ctx, projectID, findingSettingID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePendFinding provides a mock function with given fields: ctx, projectID, findingID
func (_m *FindingRepository) DeletePendFinding(ctx context.Context, projectID uint32, findingID uint64) error {
	ret := _m.Called(ctx, projectID, findingID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) error); ok {
		r0 = rf(ctx, projectID, findingID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteResource provides a mock function with given fields: _a0, _a1, _a2
func (_m *FindingRepository) DeleteResource(_a0 context.Context, _a1 uint32, _a2 uint64) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFinding provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *FindingRepository) GetFinding(_a0 context.Context, _a1 uint32, _a2 uint64, _a3 bool) (*model.Finding, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.Finding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64, bool) (*model.Finding, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64, bool) *model.Finding); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Finding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64, bool) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFindingByDataSource provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *FindingRepository) GetFindingByDataSource(_a0 context.Context, _a1 uint32, _a2 string, _a3 string) (*model.Finding, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.Finding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, string) (*model.Finding, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, string) *model.Finding); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Finding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFindingSetting provides a mock function with given fields: ctx, projectID, findingSettingID
func (_m *FindingRepository) GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	ret := _m.Called(ctx, projectID, findingSettingID)

	var r0 *model.FindingSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.FindingSetting, error)); ok {
		return rf(ctx, projectID, findingSettingID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.FindingSetting); ok {
		r0 = rf(ctx, projectID, findingSettingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FindingSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(ctx, projectID, findingSettingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFindingSettingByResource provides a mock function with given fields: ctx, projectID, resourceName
func (_m *FindingRepository) GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error) {
	ret := _m.Called(ctx, projectID, resourceName)

	var r0 *model.FindingSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*model.FindingSetting, error)); ok {
		return rf(ctx, projectID, resourceName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *model.FindingSetting); ok {
		r0 = rf(ctx, projectID, resourceName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FindingSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(ctx, projectID, resourceName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFindingTagByKey provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *FindingRepository) GetFindingTagByKey(_a0 context.Context, _a1 uint32, _a2 uint64, _a3 string) (*model.FindingTag, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.FindingTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64, string) (*model.FindingTag, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64, string) *model.FindingTag); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FindingTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPendFinding provides a mock function with given fields: ctx, projectID, findingID
func (_m *FindingRepository) GetPendFinding(ctx context.Context, projectID uint32, findingID uint64) (*model.PendFinding, error) {
	ret := _m.Called(ctx, projectID, findingID)

	var r0 *model.PendFinding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) (*model.PendFinding, error)); ok {
		return rf(ctx, projectID, findingID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) *model.PendFinding); ok {
		r0 = rf(ctx, projectID, findingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PendFinding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64) error); ok {
		r1 = rf(ctx, projectID, findingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecommend provides a mock function with given fields: ctx, projectID, findingID
func (_m *FindingRepository) GetRecommend(ctx context.Context, projectID uint32, findingID uint64) (*model.Recommend, error) {
	ret := _m.Called(ctx, projectID, findingID)

	var r0 *model.Recommend
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) (*model.Recommend, error)); ok {
		return rf(ctx, projectID, findingID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) *model.Recommend); ok {
		r0 = rf(ctx, projectID, findingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Recommend)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64) error); ok {
		r1 = rf(ctx, projectID, findingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecommendByDataSourceType provides a mock function with given fields: ctx, dataSource, recommendType
func (_m *FindingRepository) GetRecommendByDataSourceType(ctx context.Context, dataSource string, recommendType string) (*model.Recommend, error) {
	ret := _m.Called(ctx, dataSource, recommendType)

	var r0 *model.Recommend
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*model.Recommend, error)); ok {
		return rf(ctx, dataSource, recommendType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Recommend); ok {
		r0 = rf(ctx, dataSource, recommendType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Recommend)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, dataSource, recommendType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResource provides a mock function with given fields: _a0, _a1, _a2
func (_m *FindingRepository) GetResource(_a0 context.Context, _a1 uint32, _a2 uint64) (*model.Resource, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) (*model.Resource, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) *model.Resource); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResourceByName provides a mock function with given fields: _a0, _a1, _a2
func (_m *FindingRepository) GetResourceByName(_a0 context.Context, _a1 uint32, _a2 string) (*model.Resource, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*model.Resource, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *model.Resource); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResourceTagByKey provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *FindingRepository) GetResourceTagByKey(_a0 context.Context, _a1 uint32, _a2 uint64, _a3 string) (*model.ResourceTag, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.ResourceTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64, string) (*model.ResourceTag, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64, string) *model.ResourceTag); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ResourceTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFinding provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) ListFinding(_a0 context.Context, _a1 *finding.ListFindingRequest) (*[]model.Finding, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *[]model.Finding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingRequest) (*[]model.Finding, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingRequest) *[]model.Finding); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Finding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingCount provides a mock function with given fields: ctx, projectID, alertID, fromScore, toScore, findingID, dataSources, resourceNames, tags, status
func (_m *FindingRepository) ListFindingCount(ctx context.Context, projectID uint32, alertID uint32, fromScore float32, toScore float32, findingID uint64, dataSources []string, resourceNames []string, tags []string, status finding.FindingStatus) (int64, error) {
	ret := _m.Called(ctx, projectID, alertID, fromScore, toScore, findingID, dataSources, resourceNames, tags, status)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, float32, float32, uint64, []string, []string, []string, finding.FindingStatus) (int64, error)); ok {
		return rf(ctx, projectID, alertID, fromScore, toScore, findingID, dataSources, resourceNames, tags, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, float32, float32, uint64, []string, []string, []string, finding.FindingStatus) int64); ok {
		r0 = rf(ctx, projectID, alertID, fromScore, toScore, findingID, dataSources, resourceNames, tags, status)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, float32, float32, uint64, []string, []string, []string, finding.FindingStatus) error); ok {
		r1 = rf(ctx, projectID, alertID, fromScore, toScore, findingID, dataSources, resourceNames, tags, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingSetting provides a mock function with given fields: ctx, req
func (_m *FindingRepository) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	ret := _m.Called(ctx, req)

	var r0 *[]model.FindingSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingSettingRequest) *[]model.FindingSetting); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.FindingSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingSettingRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTag provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error) {
	ret := _m.Called(ctx, param)

	var r0 *[]model.FindingTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagRequest) (*[]model.FindingTag, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagRequest) *[]model.FindingTag); ok {
		r0 = rf(ctx, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.FindingTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingTagRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTagByFindingID provides a mock function with given fields: ctx, projectID, findingID
func (_m *FindingRepository) ListFindingTagByFindingID(ctx context.Context, projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	ret := _m.Called(ctx, projectID, findingID)

	var r0 *[]model.FindingTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) (*[]model.FindingTag, error)); ok {
		return rf(ctx, projectID, findingID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) *[]model.FindingTag); ok {
		r0 = rf(ctx, projectID, findingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.FindingTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64) error); ok {
		r1 = rf(ctx, projectID, findingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTagCount provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error) {
	ret := _m.Called(ctx, param)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagRequest) (int64, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagRequest) int64); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingTagRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTagName provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]db.TagName, error) {
	ret := _m.Called(ctx, param)

	var r0 *[]db.TagName
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagNameRequest) (*[]db.TagName, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagNameRequest) *[]db.TagName); ok {
		r0 = rf(ctx, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]db.TagName)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingTagNameRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTagNameCount provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error) {
	ret := _m.Called(ctx, param)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagNameRequest) (int64, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagNameRequest) int64); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingTagNameRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResource provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) ListResource(_a0 context.Context, _a1 *finding.ListResourceRequest) (*[]model.Resource, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *[]model.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceRequest) (*[]model.Resource, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceRequest) *[]model.Resource); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceCount provides a mock function with given fields: ctx, req
func (_m *FindingRepository) ListResourceCount(ctx context.Context, req *finding.ListResourceRequest) (int64, error) {
	ret := _m.Called(ctx, req)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceRequest) (int64, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceRequest) int64); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTag provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListResourceTag(ctx context.Context, param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error) {
	ret := _m.Called(ctx, param)

	var r0 *[]model.ResourceTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagRequest) (*[]model.ResourceTag, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagRequest) *[]model.ResourceTag); ok {
		r0 = rf(ctx, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.ResourceTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceTagRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTagByResourceID provides a mock function with given fields: ctx, projectID, resourceID
func (_m *FindingRepository) ListResourceTagByResourceID(ctx context.Context, projectID uint32, resourceID uint64) (*[]model.ResourceTag, error) {
	ret := _m.Called(ctx, projectID, resourceID)

	var r0 *[]model.ResourceTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) (*[]model.ResourceTag, error)); ok {
		return rf(ctx, projectID, resourceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) *[]model.ResourceTag); ok {
		r0 = rf(ctx, projectID, resourceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.ResourceTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint64) error); ok {
		r1 = rf(ctx, projectID, resourceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTagCount provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListResourceTagCount(ctx context.Context, param *finding.ListResourceTagRequest) (int64, error) {
	ret := _m.Called(ctx, param)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagRequest) (int64, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagRequest) int64); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceTagRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTagName provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListResourceTagName(ctx context.Context, param *finding.ListResourceTagNameRequest) (*[]db.TagName, error) {
	ret := _m.Called(ctx, param)

	var r0 *[]db.TagName
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagNameRequest) (*[]db.TagName, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagNameRequest) *[]db.TagName); ok {
		r0 = rf(ctx, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]db.TagName)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceTagNameRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTagNameCount provides a mock function with given fields: ctx, param
func (_m *FindingRepository) ListResourceTagNameCount(ctx context.Context, param *finding.ListResourceTagNameRequest) (int64, error) {
	ret := _m.Called(ctx, param)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagNameRequest) (int64, error)); ok {
		return rf(ctx, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagNameRequest) int64); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceTagNameRequest) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TagFinding provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) TagFinding(_a0 context.Context, _a1 *model.FindingTag) (*model.FindingTag, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.FindingTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FindingTag) (*model.FindingTag, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.FindingTag) *model.FindingTag); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FindingTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.FindingTag) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TagResource provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) TagResource(_a0 context.Context, _a1 *model.ResourceTag) (*model.ResourceTag, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.ResourceTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.ResourceTag) (*model.ResourceTag, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.ResourceTag) *model.ResourceTag); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ResourceTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.ResourceTag) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UntagFinding provides a mock function with given fields: _a0, _a1, _a2
func (_m *FindingRepository) UntagFinding(_a0 context.Context, _a1 uint32, _a2 uint64) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UntagResource provides a mock function with given fields: _a0, _a1, _a2
func (_m *FindingRepository) UntagResource(_a0 context.Context, _a1 uint32, _a2 uint64) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint64) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertFinding provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) UpsertFinding(_a0 context.Context, _a1 *model.Finding) (*model.Finding, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.Finding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Finding) (*model.Finding, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Finding) *model.Finding); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Finding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Finding) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertFindingSetting provides a mock function with given fields: ctx, data
func (_m *FindingRepository) UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error) {
	ret := _m.Called(ctx, data)

	var r0 *model.FindingSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FindingSetting) (*model.FindingSetting, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.FindingSetting) *model.FindingSetting); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FindingSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.FindingSetting) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertPendFinding provides a mock function with given fields: ctx, pend
func (_m *FindingRepository) UpsertPendFinding(ctx context.Context, pend *finding.PendFindingForUpsert) (*model.PendFinding, error) {
	ret := _m.Called(ctx, pend)

	var r0 *model.PendFinding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PendFindingForUpsert) (*model.PendFinding, error)); ok {
		return rf(ctx, pend)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PendFindingForUpsert) *model.PendFinding); ok {
		r0 = rf(ctx, pend)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PendFinding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *finding.PendFindingForUpsert) error); ok {
		r1 = rf(ctx, pend)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertRecommend provides a mock function with given fields: ctx, data
func (_m *FindingRepository) UpsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error) {
	ret := _m.Called(ctx, data)

	var r0 *model.Recommend
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Recommend) (*model.Recommend, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Recommend) *model.Recommend); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Recommend)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Recommend) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertRecommendFinding provides a mock function with given fields: ctx, data
func (_m *FindingRepository) UpsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error) {
	ret := _m.Called(ctx, data)

	var r0 *model.RecommendFinding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.RecommendFinding) (*model.RecommendFinding, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.RecommendFinding) *model.RecommendFinding); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RecommendFinding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.RecommendFinding) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertResource provides a mock function with given fields: _a0, _a1
func (_m *FindingRepository) UpsertResource(_a0 context.Context, _a1 *model.Resource) (*model.Resource, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Resource) (*model.Resource, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Resource) *model.Resource); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Resource) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFindingRepository creates a new instance of FindingRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFindingRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FindingRepository {
	mock := &FindingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
