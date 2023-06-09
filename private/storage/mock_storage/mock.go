// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/private/storage (interfaces: TrustDB)

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	x509 "crypto/x509"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cppki "github.com/scionproto/scion/pkg/scrypto/cppki"
	trust "github.com/scionproto/scion/private/storage/trust"
	trust0 "github.com/scionproto/scion/private/trust"
)

// MockTrustDB is a mock of TrustDB interface.
type MockTrustDB struct {
	ctrl     *gomock.Controller
	recorder *MockTrustDBMockRecorder
}

// MockTrustDBMockRecorder is the mock recorder for MockTrustDB.
type MockTrustDBMockRecorder struct {
	mock *MockTrustDB
}

// NewMockTrustDB creates a new mock instance.
func NewMockTrustDB(ctrl *gomock.Controller) *MockTrustDB {
	mock := &MockTrustDB{ctrl: ctrl}
	mock.recorder = &MockTrustDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrustDB) EXPECT() *MockTrustDBMockRecorder {
	return m.recorder
}

// Chain mocks base method.
func (m *MockTrustDB) Chain(arg0 context.Context, arg1 []byte) ([]*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chain", arg0, arg1)
	ret0, _ := ret[0].([]*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Chain indicates an expected call of Chain.
func (mr *MockTrustDBMockRecorder) Chain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chain", reflect.TypeOf((*MockTrustDB)(nil).Chain), arg0, arg1)
}

// Chains mocks base method.
func (m *MockTrustDB) Chains(arg0 context.Context, arg1 trust0.ChainQuery) ([][]*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chains", arg0, arg1)
	ret0, _ := ret[0].([][]*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Chains indicates an expected call of Chains.
func (mr *MockTrustDBMockRecorder) Chains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chains", reflect.TypeOf((*MockTrustDB)(nil).Chains), arg0, arg1)
}

// Close mocks base method.
func (m *MockTrustDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTrustDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTrustDB)(nil).Close))
}

// InsertChain mocks base method.
func (m *MockTrustDB) InsertChain(arg0 context.Context, arg1 []*x509.Certificate) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertChain", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertChain indicates an expected call of InsertChain.
func (mr *MockTrustDBMockRecorder) InsertChain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertChain", reflect.TypeOf((*MockTrustDB)(nil).InsertChain), arg0, arg1)
}

// InsertTRC mocks base method.
func (m *MockTrustDB) InsertTRC(arg0 context.Context, arg1 cppki.SignedTRC) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTRC", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTRC indicates an expected call of InsertTRC.
func (mr *MockTrustDBMockRecorder) InsertTRC(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTRC", reflect.TypeOf((*MockTrustDB)(nil).InsertTRC), arg0, arg1)
}

// SignedTRC mocks base method.
func (m *MockTrustDB) SignedTRC(arg0 context.Context, arg1 cppki.TRCID) (cppki.SignedTRC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignedTRC", arg0, arg1)
	ret0, _ := ret[0].(cppki.SignedTRC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignedTRC indicates an expected call of SignedTRC.
func (mr *MockTrustDBMockRecorder) SignedTRC(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignedTRC", reflect.TypeOf((*MockTrustDB)(nil).SignedTRC), arg0, arg1)
}

// SignedTRCs mocks base method.
func (m *MockTrustDB) SignedTRCs(arg0 context.Context, arg1 trust.TRCsQuery) (cppki.SignedTRCs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignedTRCs", arg0, arg1)
	ret0, _ := ret[0].(cppki.SignedTRCs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignedTRCs indicates an expected call of SignedTRCs.
func (mr *MockTrustDBMockRecorder) SignedTRCs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignedTRCs", reflect.TypeOf((*MockTrustDB)(nil).SignedTRCs), arg0, arg1)
}
