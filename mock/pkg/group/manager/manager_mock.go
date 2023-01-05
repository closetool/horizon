// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock_manager is a generated GoMock package.
package mock_manager

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/horizoncd/horizon/pkg/environmentregion/models"
	models0 "github.com/horizoncd/horizon/pkg/group/models"
	models1 "github.com/horizoncd/horizon/pkg/region/models"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockManager) Create(ctx context.Context, group *models0.Group) (*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, group)
	ret0, _ := ret[0].(*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockManagerMockRecorder) Create(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManager)(nil).Create), ctx, group)
}

// Delete mocks base method.
func (m *MockManager) Delete(ctx context.Context, id uint) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockManagerMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockManager)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *MockManager) GetAll(ctx context.Context) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockManagerMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockManager)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *MockManager) GetByID(ctx context.Context, id uint) (*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockManagerMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockManager)(nil).GetByID), ctx, id)
}

// GetByIDNameFuzzily mocks base method.
func (m *MockManager) GetByIDNameFuzzily(ctx context.Context, id uint, name string) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDNameFuzzily", ctx, id, name)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDNameFuzzily indicates an expected call of GetByIDNameFuzzily.
func (mr *MockManagerMockRecorder) GetByIDNameFuzzily(ctx, id, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDNameFuzzily", reflect.TypeOf((*MockManager)(nil).GetByIDNameFuzzily), ctx, id, name)
}

// GetByIDs mocks base method.
func (m *MockManager) GetByIDs(ctx context.Context, ids []uint) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDs", ctx, ids)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDs indicates an expected call of GetByIDs.
func (mr *MockManagerMockRecorder) GetByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockManager)(nil).GetByIDs), ctx, ids)
}

// GetByNameFuzzily mocks base method.
func (m *MockManager) GetByNameFuzzily(ctx context.Context, name string) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameFuzzily", ctx, name)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameFuzzily indicates an expected call of GetByNameFuzzily.
func (mr *MockManagerMockRecorder) GetByNameFuzzily(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameFuzzily", reflect.TypeOf((*MockManager)(nil).GetByNameFuzzily), ctx, name)
}

// GetByNameFuzzilyIncludeSoftDelete mocks base method.
func (m *MockManager) GetByNameFuzzilyIncludeSoftDelete(ctx context.Context, name string) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameFuzzilyIncludeSoftDelete", ctx, name)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameFuzzilyIncludeSoftDelete indicates an expected call of GetByNameFuzzilyIncludeSoftDelete.
func (mr *MockManagerMockRecorder) GetByNameFuzzilyIncludeSoftDelete(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameFuzzilyIncludeSoftDelete", reflect.TypeOf((*MockManager)(nil).GetByNameFuzzilyIncludeSoftDelete), ctx, name)
}

// GetByNameOrPathUnderParent mocks base method.
func (m *MockManager) GetByNameOrPathUnderParent(ctx context.Context, name, path string, parentID uint) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameOrPathUnderParent", ctx, name, path, parentID)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameOrPathUnderParent indicates an expected call of GetByNameOrPathUnderParent.
func (mr *MockManagerMockRecorder) GetByNameOrPathUnderParent(ctx, name, path, parentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameOrPathUnderParent", reflect.TypeOf((*MockManager)(nil).GetByNameOrPathUnderParent), ctx, name, path, parentID)
}

// GetByPaths mocks base method.
func (m *MockManager) GetByPaths(ctx context.Context, paths []string) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPaths", ctx, paths)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPaths indicates an expected call of GetByPaths.
func (mr *MockManagerMockRecorder) GetByPaths(ctx, paths interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPaths", reflect.TypeOf((*MockManager)(nil).GetByPaths), ctx, paths)
}

// GetChildren mocks base method.
func (m *MockManager) GetChildren(ctx context.Context, parentID uint, pageNumber, pageSize int) ([]*models0.GroupOrApplication, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChildren", ctx, parentID, pageNumber, pageSize)
	ret0, _ := ret[0].([]*models0.GroupOrApplication)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetChildren indicates an expected call of GetChildren.
func (mr *MockManagerMockRecorder) GetChildren(ctx, parentID, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChildren", reflect.TypeOf((*MockManager)(nil).GetChildren), ctx, parentID, pageNumber, pageSize)
}

// GetDefaultRegions mocks base method.
func (m *MockManager) GetDefaultRegions(ctx context.Context, id uint) ([]*models.EnvironmentRegion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultRegions", ctx, id)
	ret0, _ := ret[0].([]*models.EnvironmentRegion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultRegions indicates an expected call of GetDefaultRegions.
func (mr *MockManagerMockRecorder) GetDefaultRegions(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultRegions", reflect.TypeOf((*MockManager)(nil).GetDefaultRegions), ctx, id)
}

// GetSelectableRegions mocks base method.
func (m *MockManager) GetSelectableRegions(ctx context.Context, id uint) (models1.RegionParts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelectableRegions", ctx, id)
	ret0, _ := ret[0].(models1.RegionParts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSelectableRegions indicates an expected call of GetSelectableRegions.
func (mr *MockManagerMockRecorder) GetSelectableRegions(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelectableRegions", reflect.TypeOf((*MockManager)(nil).GetSelectableRegions), ctx, id)
}

// GetSelectableRegionsByEnv mocks base method.
func (m *MockManager) GetSelectableRegionsByEnv(ctx context.Context, id uint, env string) (models1.RegionParts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelectableRegionsByEnv", ctx, id, env)
	ret0, _ := ret[0].(models1.RegionParts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSelectableRegionsByEnv indicates an expected call of GetSelectableRegionsByEnv.
func (mr *MockManagerMockRecorder) GetSelectableRegionsByEnv(ctx, id, env interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelectableRegionsByEnv", reflect.TypeOf((*MockManager)(nil).GetSelectableRegionsByEnv), ctx, id, env)
}

// GetSubGroups mocks base method.
func (m *MockManager) GetSubGroups(ctx context.Context, id uint, pageNumber, pageSize int) ([]*models0.Group, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubGroups", ctx, id, pageNumber, pageSize)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSubGroups indicates an expected call of GetSubGroups.
func (mr *MockManagerMockRecorder) GetSubGroups(ctx, id, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubGroups", reflect.TypeOf((*MockManager)(nil).GetSubGroups), ctx, id, pageNumber, pageSize)
}

// GetSubGroupsByGroupIDs mocks base method.
func (m *MockManager) GetSubGroupsByGroupIDs(ctx context.Context, groupIDs []uint) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubGroupsByGroupIDs", ctx, groupIDs)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubGroupsByGroupIDs indicates an expected call of GetSubGroupsByGroupIDs.
func (mr *MockManagerMockRecorder) GetSubGroupsByGroupIDs(ctx, groupIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubGroupsByGroupIDs", reflect.TypeOf((*MockManager)(nil).GetSubGroupsByGroupIDs), ctx, groupIDs)
}

// GetSubGroupsUnderParentIDs mocks base method.
func (m *MockManager) GetSubGroupsUnderParentIDs(ctx context.Context, parentIDs []uint) ([]*models0.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubGroupsUnderParentIDs", ctx, parentIDs)
	ret0, _ := ret[0].([]*models0.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubGroupsUnderParentIDs indicates an expected call of GetSubGroupsUnderParentIDs.
func (mr *MockManagerMockRecorder) GetSubGroupsUnderParentIDs(ctx, parentIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubGroupsUnderParentIDs", reflect.TypeOf((*MockManager)(nil).GetSubGroupsUnderParentIDs), ctx, parentIDs)
}

// GroupExist mocks base method.
func (m *MockManager) GroupExist(ctx context.Context, groupID uint) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GroupExist", ctx, groupID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// GroupExist indicates an expected call of GroupExist.
func (mr *MockManagerMockRecorder) GroupExist(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GroupExist", reflect.TypeOf((*MockManager)(nil).GroupExist), ctx, groupID)
}

// IsRootGroup mocks base method.
func (m *MockManager) IsRootGroup(ctx context.Context, groupID uint) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRootGroup", ctx, groupID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRootGroup indicates an expected call of IsRootGroup.
func (mr *MockManagerMockRecorder) IsRootGroup(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRootGroup", reflect.TypeOf((*MockManager)(nil).IsRootGroup), ctx, groupID)
}

// Transfer mocks base method.
func (m *MockManager) Transfer(ctx context.Context, id, newParentID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", ctx, id, newParentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transfer indicates an expected call of Transfer.
func (mr *MockManagerMockRecorder) Transfer(ctx, id, newParentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockManager)(nil).Transfer), ctx, id, newParentID)
}

// UpdateBasic mocks base method.
func (m *MockManager) UpdateBasic(ctx context.Context, group *models0.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBasic", ctx, group)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBasic indicates an expected call of UpdateBasic.
func (mr *MockManagerMockRecorder) UpdateBasic(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBasic", reflect.TypeOf((*MockManager)(nil).UpdateBasic), ctx, group)
}

// UpdateRegionSelector mocks base method.
func (m *MockManager) UpdateRegionSelector(ctx context.Context, id uint, regionSelector string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRegionSelector", ctx, id, regionSelector)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRegionSelector indicates an expected call of UpdateRegionSelector.
func (mr *MockManagerMockRecorder) UpdateRegionSelector(ctx, id, regionSelector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRegionSelector", reflect.TypeOf((*MockManager)(nil).UpdateRegionSelector), ctx, id, regionSelector)
}
