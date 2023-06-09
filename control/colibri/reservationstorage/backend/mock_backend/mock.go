// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/control/colibri/reservationstorage/backend (interfaces: DB,Transaction)

// Package mock_backend is a generated GoMock package.
package mock_backend

import (
	context "context"
	sql "database/sql"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	e2e "github.com/scionproto/scion/control/colibri/reservation/e2e"
	segment "github.com/scionproto/scion/control/colibri/reservation/segment"
	backend "github.com/scionproto/scion/control/colibri/reservationstorage/backend"
	addr "github.com/scionproto/scion/pkg/addr"
	reservation "github.com/scionproto/scion/pkg/experimental/colibri/reservation"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// BeginTransaction mocks base method.
func (m *MockDB) BeginTransaction(arg0 context.Context, arg1 *sql.TxOptions) (backend.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTransaction", arg0, arg1)
	ret0, _ := ret[0].(backend.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTransaction indicates an expected call of BeginTransaction.
func (mr *MockDBMockRecorder) BeginTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTransaction", reflect.TypeOf((*MockDB)(nil).BeginTransaction), arg0, arg1)
}

// Close mocks base method.
func (m *MockDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDB)(nil).Close))
}

// DeleteExpiredIndices mocks base method.
func (m *MockDB) DeleteExpiredIndices(arg0 context.Context, arg1 time.Time) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredIndices", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteExpiredIndices indicates an expected call of DeleteExpiredIndices.
func (mr *MockDBMockRecorder) DeleteExpiredIndices(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredIndices", reflect.TypeOf((*MockDB)(nil).DeleteExpiredIndices), arg0, arg1)
}

// DeleteSegmentRsv mocks base method.
func (m *MockDB) DeleteSegmentRsv(arg0 context.Context, arg1 *reservation.SegmentID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSegmentRsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSegmentRsv indicates an expected call of DeleteSegmentRsv.
func (mr *MockDBMockRecorder) DeleteSegmentRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegmentRsv", reflect.TypeOf((*MockDB)(nil).DeleteSegmentRsv), arg0, arg1)
}

// GetAllSegmentRsvs mocks base method.
func (m *MockDB) GetAllSegmentRsvs(arg0 context.Context) ([]*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSegmentRsvs", arg0)
	ret0, _ := ret[0].([]*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSegmentRsvs indicates an expected call of GetAllSegmentRsvs.
func (mr *MockDBMockRecorder) GetAllSegmentRsvs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSegmentRsvs", reflect.TypeOf((*MockDB)(nil).GetAllSegmentRsvs), arg0)
}

// GetE2ERsvFromID mocks base method.
func (m *MockDB) GetE2ERsvFromID(arg0 context.Context, arg1 *reservation.E2EID) (*e2e.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetE2ERsvFromID", arg0, arg1)
	ret0, _ := ret[0].(*e2e.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetE2ERsvFromID indicates an expected call of GetE2ERsvFromID.
func (mr *MockDBMockRecorder) GetE2ERsvFromID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetE2ERsvFromID", reflect.TypeOf((*MockDB)(nil).GetE2ERsvFromID), arg0, arg1)
}

// GetE2ERsvsOnSegRsv mocks base method.
func (m *MockDB) GetE2ERsvsOnSegRsv(arg0 context.Context, arg1 *reservation.SegmentID) ([]*e2e.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetE2ERsvsOnSegRsv", arg0, arg1)
	ret0, _ := ret[0].([]*e2e.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetE2ERsvsOnSegRsv indicates an expected call of GetE2ERsvsOnSegRsv.
func (mr *MockDBMockRecorder) GetE2ERsvsOnSegRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetE2ERsvsOnSegRsv", reflect.TypeOf((*MockDB)(nil).GetE2ERsvsOnSegRsv), arg0, arg1)
}

// GetSegmentRsvFromID mocks base method.
func (m *MockDB) GetSegmentRsvFromID(arg0 context.Context, arg1 *reservation.SegmentID) (*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvFromID", arg0, arg1)
	ret0, _ := ret[0].(*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvFromID indicates an expected call of GetSegmentRsvFromID.
func (mr *MockDBMockRecorder) GetSegmentRsvFromID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvFromID", reflect.TypeOf((*MockDB)(nil).GetSegmentRsvFromID), arg0, arg1)
}

// GetSegmentRsvFromPath mocks base method.
func (m *MockDB) GetSegmentRsvFromPath(arg0 context.Context, arg1 segment.ReservationTransparentPath) (*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvFromPath", arg0, arg1)
	ret0, _ := ret[0].(*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvFromPath indicates an expected call of GetSegmentRsvFromPath.
func (mr *MockDBMockRecorder) GetSegmentRsvFromPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvFromPath", reflect.TypeOf((*MockDB)(nil).GetSegmentRsvFromPath), arg0, arg1)
}

// GetSegmentRsvsFromIFPair mocks base method.
func (m *MockDB) GetSegmentRsvsFromIFPair(arg0 context.Context, arg1, arg2 *uint16) ([]*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvsFromIFPair", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvsFromIFPair indicates an expected call of GetSegmentRsvsFromIFPair.
func (mr *MockDBMockRecorder) GetSegmentRsvsFromIFPair(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvsFromIFPair", reflect.TypeOf((*MockDB)(nil).GetSegmentRsvsFromIFPair), arg0, arg1, arg2)
}

// GetSegmentRsvsFromSrcDstIA mocks base method.
func (m *MockDB) GetSegmentRsvsFromSrcDstIA(arg0 context.Context, arg1, arg2 addr.IA) ([]*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvsFromSrcDstIA", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvsFromSrcDstIA indicates an expected call of GetSegmentRsvsFromSrcDstIA.
func (mr *MockDBMockRecorder) GetSegmentRsvsFromSrcDstIA(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvsFromSrcDstIA", reflect.TypeOf((*MockDB)(nil).GetSegmentRsvsFromSrcDstIA), arg0, arg1, arg2)
}

// NewSegmentRsv mocks base method.
func (m *MockDB) NewSegmentRsv(arg0 context.Context, arg1 *segment.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSegmentRsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewSegmentRsv indicates an expected call of NewSegmentRsv.
func (mr *MockDBMockRecorder) NewSegmentRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSegmentRsv", reflect.TypeOf((*MockDB)(nil).NewSegmentRsv), arg0, arg1)
}

// PersistE2ERsv mocks base method.
func (m *MockDB) PersistE2ERsv(arg0 context.Context, arg1 *e2e.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PersistE2ERsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PersistE2ERsv indicates an expected call of PersistE2ERsv.
func (mr *MockDBMockRecorder) PersistE2ERsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersistE2ERsv", reflect.TypeOf((*MockDB)(nil).PersistE2ERsv), arg0, arg1)
}

// PersistSegmentRsv mocks base method.
func (m *MockDB) PersistSegmentRsv(arg0 context.Context, arg1 *segment.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PersistSegmentRsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PersistSegmentRsv indicates an expected call of PersistSegmentRsv.
func (mr *MockDBMockRecorder) PersistSegmentRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersistSegmentRsv", reflect.TypeOf((*MockDB)(nil).PersistSegmentRsv), arg0, arg1)
}

// SetMaxIdleConns mocks base method.
func (m *MockDB) SetMaxIdleConns(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxIdleConns", arg0)
}

// SetMaxIdleConns indicates an expected call of SetMaxIdleConns.
func (mr *MockDBMockRecorder) SetMaxIdleConns(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxIdleConns", reflect.TypeOf((*MockDB)(nil).SetMaxIdleConns), arg0)
}

// SetMaxOpenConns mocks base method.
func (m *MockDB) SetMaxOpenConns(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxOpenConns", arg0)
}

// SetMaxOpenConns indicates an expected call of SetMaxOpenConns.
func (mr *MockDBMockRecorder) SetMaxOpenConns(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxOpenConns", reflect.TypeOf((*MockDB)(nil).SetMaxOpenConns), arg0)
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockTransaction) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTransactionMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTransaction)(nil).Commit))
}

// DeleteExpiredIndices mocks base method.
func (m *MockTransaction) DeleteExpiredIndices(arg0 context.Context, arg1 time.Time) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredIndices", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteExpiredIndices indicates an expected call of DeleteExpiredIndices.
func (mr *MockTransactionMockRecorder) DeleteExpiredIndices(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredIndices", reflect.TypeOf((*MockTransaction)(nil).DeleteExpiredIndices), arg0, arg1)
}

// DeleteSegmentRsv mocks base method.
func (m *MockTransaction) DeleteSegmentRsv(arg0 context.Context, arg1 *reservation.SegmentID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSegmentRsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSegmentRsv indicates an expected call of DeleteSegmentRsv.
func (mr *MockTransactionMockRecorder) DeleteSegmentRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegmentRsv", reflect.TypeOf((*MockTransaction)(nil).DeleteSegmentRsv), arg0, arg1)
}

// GetAllSegmentRsvs mocks base method.
func (m *MockTransaction) GetAllSegmentRsvs(arg0 context.Context) ([]*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSegmentRsvs", arg0)
	ret0, _ := ret[0].([]*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSegmentRsvs indicates an expected call of GetAllSegmentRsvs.
func (mr *MockTransactionMockRecorder) GetAllSegmentRsvs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSegmentRsvs", reflect.TypeOf((*MockTransaction)(nil).GetAllSegmentRsvs), arg0)
}

// GetE2ERsvFromID mocks base method.
func (m *MockTransaction) GetE2ERsvFromID(arg0 context.Context, arg1 *reservation.E2EID) (*e2e.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetE2ERsvFromID", arg0, arg1)
	ret0, _ := ret[0].(*e2e.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetE2ERsvFromID indicates an expected call of GetE2ERsvFromID.
func (mr *MockTransactionMockRecorder) GetE2ERsvFromID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetE2ERsvFromID", reflect.TypeOf((*MockTransaction)(nil).GetE2ERsvFromID), arg0, arg1)
}

// GetE2ERsvsOnSegRsv mocks base method.
func (m *MockTransaction) GetE2ERsvsOnSegRsv(arg0 context.Context, arg1 *reservation.SegmentID) ([]*e2e.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetE2ERsvsOnSegRsv", arg0, arg1)
	ret0, _ := ret[0].([]*e2e.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetE2ERsvsOnSegRsv indicates an expected call of GetE2ERsvsOnSegRsv.
func (mr *MockTransactionMockRecorder) GetE2ERsvsOnSegRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetE2ERsvsOnSegRsv", reflect.TypeOf((*MockTransaction)(nil).GetE2ERsvsOnSegRsv), arg0, arg1)
}

// GetSegmentRsvFromID mocks base method.
func (m *MockTransaction) GetSegmentRsvFromID(arg0 context.Context, arg1 *reservation.SegmentID) (*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvFromID", arg0, arg1)
	ret0, _ := ret[0].(*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvFromID indicates an expected call of GetSegmentRsvFromID.
func (mr *MockTransactionMockRecorder) GetSegmentRsvFromID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvFromID", reflect.TypeOf((*MockTransaction)(nil).GetSegmentRsvFromID), arg0, arg1)
}

// GetSegmentRsvFromPath mocks base method.
func (m *MockTransaction) GetSegmentRsvFromPath(arg0 context.Context, arg1 segment.ReservationTransparentPath) (*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvFromPath", arg0, arg1)
	ret0, _ := ret[0].(*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvFromPath indicates an expected call of GetSegmentRsvFromPath.
func (mr *MockTransactionMockRecorder) GetSegmentRsvFromPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvFromPath", reflect.TypeOf((*MockTransaction)(nil).GetSegmentRsvFromPath), arg0, arg1)
}

// GetSegmentRsvsFromIFPair mocks base method.
func (m *MockTransaction) GetSegmentRsvsFromIFPair(arg0 context.Context, arg1, arg2 *uint16) ([]*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvsFromIFPair", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvsFromIFPair indicates an expected call of GetSegmentRsvsFromIFPair.
func (mr *MockTransactionMockRecorder) GetSegmentRsvsFromIFPair(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvsFromIFPair", reflect.TypeOf((*MockTransaction)(nil).GetSegmentRsvsFromIFPair), arg0, arg1, arg2)
}

// GetSegmentRsvsFromSrcDstIA mocks base method.
func (m *MockTransaction) GetSegmentRsvsFromSrcDstIA(arg0 context.Context, arg1, arg2 addr.IA) ([]*segment.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegmentRsvsFromSrcDstIA", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*segment.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentRsvsFromSrcDstIA indicates an expected call of GetSegmentRsvsFromSrcDstIA.
func (mr *MockTransactionMockRecorder) GetSegmentRsvsFromSrcDstIA(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentRsvsFromSrcDstIA", reflect.TypeOf((*MockTransaction)(nil).GetSegmentRsvsFromSrcDstIA), arg0, arg1, arg2)
}

// NewSegmentRsv mocks base method.
func (m *MockTransaction) NewSegmentRsv(arg0 context.Context, arg1 *segment.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSegmentRsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewSegmentRsv indicates an expected call of NewSegmentRsv.
func (mr *MockTransactionMockRecorder) NewSegmentRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSegmentRsv", reflect.TypeOf((*MockTransaction)(nil).NewSegmentRsv), arg0, arg1)
}

// PersistE2ERsv mocks base method.
func (m *MockTransaction) PersistE2ERsv(arg0 context.Context, arg1 *e2e.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PersistE2ERsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PersistE2ERsv indicates an expected call of PersistE2ERsv.
func (mr *MockTransactionMockRecorder) PersistE2ERsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersistE2ERsv", reflect.TypeOf((*MockTransaction)(nil).PersistE2ERsv), arg0, arg1)
}

// PersistSegmentRsv mocks base method.
func (m *MockTransaction) PersistSegmentRsv(arg0 context.Context, arg1 *segment.Reservation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PersistSegmentRsv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PersistSegmentRsv indicates an expected call of PersistSegmentRsv.
func (mr *MockTransactionMockRecorder) PersistSegmentRsv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersistSegmentRsv", reflect.TypeOf((*MockTransaction)(nil).PersistSegmentRsv), arg0, arg1)
}

// Rollback mocks base method.
func (m *MockTransaction) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTransactionMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTransaction)(nil).Rollback))
}
