// Code generated by MockGen. DO NOT EDIT.
// Source: ./vendor/github.com/serverscom/serverscom-go-client/pkg/kubernetes_clusters.go
//
// Generated by this command:
//
//	mockgen --destination ./internal/mocks/kubernetes_clusters_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/kubernetes_clusters.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	gomock "go.uber.org/mock/gomock"
)

// MockKubernetesClustersService is a mock of KubernetesClustersService interface.
type MockKubernetesClustersService struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesClustersServiceMockRecorder
	isgomock struct{}
}

// MockKubernetesClustersServiceMockRecorder is the mock recorder for MockKubernetesClustersService.
type MockKubernetesClustersServiceMockRecorder struct {
	mock *MockKubernetesClustersService
}

// NewMockKubernetesClustersService creates a new mock instance.
func NewMockKubernetesClustersService(ctrl *gomock.Controller) *MockKubernetesClustersService {
	mock := &MockKubernetesClustersService{ctrl: ctrl}
	mock.recorder = &MockKubernetesClustersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetesClustersService) EXPECT() *MockKubernetesClustersServiceMockRecorder {
	return m.recorder
}

// Collection mocks base method.
func (m *MockKubernetesClustersService) Collection() serverscom.Collection[serverscom.KubernetesCluster] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Collection")
	ret0, _ := ret[0].(serverscom.Collection[serverscom.KubernetesCluster])
	return ret0
}

// Collection indicates an expected call of Collection.
func (mr *MockKubernetesClustersServiceMockRecorder) Collection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collection", reflect.TypeOf((*MockKubernetesClustersService)(nil).Collection))
}

// Get mocks base method.
func (m *MockKubernetesClustersService) Get(ctx context.Context, id string) (*serverscom.KubernetesCluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*serverscom.KubernetesCluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockKubernetesClustersServiceMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKubernetesClustersService)(nil).Get), ctx, id)
}

// GetNode mocks base method.
func (m *MockKubernetesClustersService) GetNode(ctx context.Context, clusterID, nodeID string) (*serverscom.KubernetesClusterNode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode", ctx, clusterID, nodeID)
	ret0, _ := ret[0].(*serverscom.KubernetesClusterNode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNode indicates an expected call of GetNode.
func (mr *MockKubernetesClustersServiceMockRecorder) GetNode(ctx, clusterID, nodeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockKubernetesClustersService)(nil).GetNode), ctx, clusterID, nodeID)
}

// Nodes mocks base method.
func (m *MockKubernetesClustersService) Nodes(id string) serverscom.Collection[serverscom.KubernetesClusterNode] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nodes", id)
	ret0, _ := ret[0].(serverscom.Collection[serverscom.KubernetesClusterNode])
	return ret0
}

// Nodes indicates an expected call of Nodes.
func (mr *MockKubernetesClustersServiceMockRecorder) Nodes(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nodes", reflect.TypeOf((*MockKubernetesClustersService)(nil).Nodes), id)
}

// Update mocks base method.
func (m *MockKubernetesClustersService) Update(ctx context.Context, id string, input serverscom.KubernetesClusterUpdateInput) (*serverscom.KubernetesCluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, input)
	ret0, _ := ret[0].(*serverscom.KubernetesCluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockKubernetesClustersServiceMockRecorder) Update(ctx, id, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockKubernetesClustersService)(nil).Update), ctx, id, input)
}
