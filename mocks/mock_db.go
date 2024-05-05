// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/storage/db.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDataBaseInterface is a mock of DataBaseInterface interface.
type MockDataBaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDataBaseInterfaceMockRecorder
}

// MockDataBaseInterfaceMockRecorder is the mock recorder for MockDataBaseInterface.
type MockDataBaseInterfaceMockRecorder struct {
	mock *MockDataBaseInterface
}

// NewMockDataBaseInterface creates a new mock instance.
func NewMockDataBaseInterface(ctrl *gomock.Controller) *MockDataBaseInterface {
	mock := &MockDataBaseInterface{ctrl: ctrl}
	mock.recorder = &MockDataBaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataBaseInterface) EXPECT() *MockDataBaseInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDataBaseInterface) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockDataBaseInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDataBaseInterface)(nil).Close))
}

// Ping mocks base method.
func (m *MockDataBaseInterface) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockDataBaseInterfaceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDataBaseInterface)(nil).Ping))
}

// SaveURLDB mocks base method.
func (m *MockDataBaseInterface) SaveURLDB(id, shortURL, originalURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURLDB", id, shortURL, originalURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveURLDB indicates an expected call of SaveURLDB.
func (mr *MockDataBaseInterfaceMockRecorder) SaveURLDB(id, shortURL, originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURLDB", reflect.TypeOf((*MockDataBaseInterface)(nil).SaveURLDB), id, shortURL, originalURL)
}
