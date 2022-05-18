// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock_manager is a generated GoMock package.
package mock_manager

import (
	context "context"
	reflect "reflect"

	q "g.hz.netease.com/horizon/lib/q"
	models "g.hz.netease.com/horizon/pkg/cluster/models"
	models0 "g.hz.netease.com/horizon/pkg/clustertag/models"
	gomock "github.com/golang/mock/gomock"
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

// CheckClusterExists mocks base method.
func (m *MockManager) CheckClusterExists(ctx context.Context, cluster string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckClusterExists", ctx, cluster)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckClusterExists indicates an expected call of CheckClusterExists.
func (mr *MockManagerMockRecorder) CheckClusterExists(ctx, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckClusterExists", reflect.TypeOf((*MockManager)(nil).CheckClusterExists), ctx, cluster)
}

// Create mocks base method.
func (m *MockManager) Create(ctx context.Context, cluster *models.Cluster, clusterTags []*models0.ClusterTag, extraMembers map[string]string) (*models.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, cluster, clusterTags, extraMembers)
	ret0, _ := ret[0].(*models.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockManagerMockRecorder) Create(ctx, cluster, clusterTags, extraMembers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManager)(nil).Create), ctx, cluster, clusterTags, extraMembers)
}

// DeleteByID mocks base method.
func (m *MockManager) DeleteByID(ctx context.Context, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockManagerMockRecorder) DeleteByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockManager)(nil).DeleteByID), ctx, id)
}

// GetByID mocks base method.
func (m *MockManager) GetByID(ctx context.Context, id uint) (*models.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockManagerMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockManager)(nil).GetByID), ctx, id)
}

// GetByName mocks base method.
func (m *MockManager) GetByName(ctx context.Context, clusterName string) (*models.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, clusterName)
	ret0, _ := ret[0].(*models.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockManagerMockRecorder) GetByName(ctx, clusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockManager)(nil).GetByName), ctx, clusterName)
}

// ListByApplicationAndEnvs mocks base method.
func (m *MockManager) ListByApplicationAndEnvs(ctx context.Context, applicationID uint, environments []string, filter string, query *q.Query) (int, []*models.ClusterWithEnvAndRegion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByApplicationAndEnvs", ctx, applicationID, environments, filter, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]*models.ClusterWithEnvAndRegion)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListByApplicationAndEnvs indicates an expected call of ListByApplicationAndEnvs.
func (mr *MockManagerMockRecorder) ListByApplicationAndEnvs(ctx, applicationID, environments, filter, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByApplicationAndEnvs", reflect.TypeOf((*MockManager)(nil).ListByApplicationAndEnvs), ctx, applicationID, environments, filter, query)
}

// ListByApplicationID mocks base method.
func (m *MockManager) ListByApplicationID(ctx context.Context, applicationID uint) ([]*models.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByApplicationID", ctx, applicationID)
	ret0, _ := ret[0].([]*models.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByApplicationID indicates an expected call of ListByApplicationID.
func (mr *MockManagerMockRecorder) ListByApplicationID(ctx, applicationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByApplicationID", reflect.TypeOf((*MockManager)(nil).ListByApplicationID), ctx, applicationID)
}

// ListByNameFuzzily mocks base method.
func (m *MockManager) ListByNameFuzzily(ctx context.Context, environment, name string, query *q.Query) (int, []*models.ClusterWithEnvAndRegion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByNameFuzzily", ctx, environment, name, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]*models.ClusterWithEnvAndRegion)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListByNameFuzzily indicates an expected call of ListByNameFuzzily.
func (mr *MockManagerMockRecorder) ListByNameFuzzily(ctx, environment, name, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByNameFuzzily", reflect.TypeOf((*MockManager)(nil).ListByNameFuzzily), ctx, environment, name, query)
}

// ListUserAuthorizedByNameFuzzily mocks base method.
func (m *MockManager) ListUserAuthorizedByNameFuzzily(ctx context.Context, environment, name string, applicationIDs []uint, userInfo uint, query *q.Query) (int, []*models.ClusterWithEnvAndRegion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserAuthorizedByNameFuzzily", ctx, environment, name, applicationIDs, userInfo, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]*models.ClusterWithEnvAndRegion)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUserAuthorizedByNameFuzzily indicates an expected call of ListUserAuthorizedByNameFuzzily.
func (mr *MockManagerMockRecorder) ListUserAuthorizedByNameFuzzily(ctx, environment, name, applicationIDs, userInfo, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserAuthorizedByNameFuzzily", reflect.TypeOf((*MockManager)(nil).ListUserAuthorizedByNameFuzzily), ctx, environment, name, applicationIDs, userInfo, query)
}

// UpdateByID mocks base method.
func (m *MockManager) UpdateByID(ctx context.Context, id uint, cluster *models.Cluster) (*models.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, cluster)
	ret0, _ := ret[0].(*models.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateByID indicates an expected call of UpdateByID.
func (mr *MockManagerMockRecorder) UpdateByID(ctx, id, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockManager)(nil).UpdateByID), ctx, id, cluster)
}
