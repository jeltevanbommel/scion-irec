// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/daemon/fetcher (interfaces: Fetcher)

// Package mock_fetcher is a generated GoMock package.
package mock_fetcher

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/pkg/addr"
	snet "github.com/scionproto/scion/pkg/snet"
)

// MockFetcher is a mock of Fetcher interface.
type MockFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockFetcherMockRecorder
}

// MockFetcherMockRecorder is the mock recorder for MockFetcher.
type MockFetcherMockRecorder struct {
	mock *MockFetcher
}

// NewMockFetcher creates a new mock instance.
func NewMockFetcher(ctrl *gomock.Controller) *MockFetcher {
	mock := &MockFetcher{ctrl: ctrl}
	mock.recorder = &MockFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFetcher) EXPECT() *MockFetcherMockRecorder {
	return m.recorder
}

// GetPaths mocks base method.
func (m *MockFetcher) GetPaths(arg0 context.Context, arg1, arg2 addr.IA, arg3 bool) ([]snet.Path, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaths", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]snet.Path)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaths indicates an expected call of GetPaths.
func (mr *MockFetcherMockRecorder) GetPaths(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaths", reflect.TypeOf((*MockFetcher)(nil).GetPaths), arg0, arg1, arg2, arg3)
}
