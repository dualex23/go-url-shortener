// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/storage/db.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	storage "github.com/dualex23/go-url-shortener/internal/app/storage"
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

// BatchSaveUrls mocks base method.
func (m *MockDataBaseInterface) BatchSaveUrls(urls []storage.URLData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSaveUrls", urls)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSaveUrls indicates an expected call of BatchSaveUrls.
func (mr *MockDataBaseInterfaceMockRecorder) BatchSaveUrls(urls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSaveUrls", reflect.TypeOf((*MockDataBaseInterface)(nil).BatchSaveUrls), urls)
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

// LoadURLByID mocks base method.
func (m *MockDataBaseInterface) LoadURLByID(id string) (*storage.URLData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadURLByID", id)
	ret0, _ := ret[0].(*storage.URLData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadURLByID indicates an expected call of LoadURLByID.
func (mr *MockDataBaseInterfaceMockRecorder) LoadURLByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadURLByID", reflect.TypeOf((*MockDataBaseInterface)(nil).LoadURLByID), id)
}

// LoadUrls mocks base method.
func (m *MockDataBaseInterface) LoadUrls() (map[string]storage.URLData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUrls")
	ret0, _ := ret[0].(map[string]storage.URLData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUrls indicates an expected call of LoadUrls.
func (mr *MockDataBaseInterfaceMockRecorder) LoadUrls() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUrls", reflect.TypeOf((*MockDataBaseInterface)(nil).LoadUrls))
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

// SaveUrls mocks base method.
func (m *MockDataBaseInterface) SaveUrls(id, shortURL, originalURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUrls", id, shortURL, originalURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUrls indicates an expected call of SaveUrls.
func (mr *MockDataBaseInterfaceMockRecorder) SaveUrls(id, shortURL, originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUrls", reflect.TypeOf((*MockDataBaseInterface)(nil).SaveUrls), id, shortURL, originalURL)
}
