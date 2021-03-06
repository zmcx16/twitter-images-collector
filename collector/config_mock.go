// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package collector is a generated GoMock package.
package collector

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIConfig is a mock of IConfig interface.
type MockIConfig struct {
	ctrl     *gomock.Controller
	recorder *MockIConfigMockRecorder
}

// MockIConfigMockRecorder is the mock recorder for MockIConfig.
type MockIConfigMockRecorder struct {
	mock *MockIConfig
}

// NewMockIConfig creates a new mock instance.
func NewMockIConfig(ctrl *gomock.Controller) *MockIConfig {
	mock := &MockIConfig{ctrl: ctrl}
	mock.recorder = &MockIConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConfig) EXPECT() *MockIConfigMockRecorder {
	return m.recorder
}

// LoadConfigFromPath mocks base method.
func (m *MockIConfig) LoadConfigFromPath(configPath string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfigFromPath", configPath)
	ret0, _ := ret[0].(bool)
	return ret0
}

// LoadConfigFromPath indicates an expected call of LoadConfigFromPath.
func (mr *MockIConfigMockRecorder) LoadConfigFromPath(configPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfigFromPath", reflect.TypeOf((*MockIConfig)(nil).LoadConfigFromPath), configPath)
}

// LoadConfigFromReader mocks base method.
func (m *MockIConfig) LoadConfigFromReader(r io.Reader) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfigFromReader", r)
	ret0, _ := ret[0].(bool)
	return ret0
}

// LoadConfigFromReader indicates an expected call of LoadConfigFromReader.
func (mr *MockIConfigMockRecorder) LoadConfigFromReader(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfigFromReader", reflect.TypeOf((*MockIConfig)(nil).LoadConfigFromReader), r)
}
