// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/metinorak/envloader (interfaces: EnvReader)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEnvReader is a mock of EnvReader interface.
type MockEnvReader struct {
	ctrl     *gomock.Controller
	recorder *MockEnvReaderMockRecorder
}

// MockEnvReaderMockRecorder is the mock recorder for MockEnvReader.
type MockEnvReaderMockRecorder struct {
	mock *MockEnvReader
}

// NewMockEnvReader creates a new mock instance.
func NewMockEnvReader(ctrl *gomock.Controller) *MockEnvReader {
	mock := &MockEnvReader{ctrl: ctrl}
	mock.recorder = &MockEnvReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnvReader) EXPECT() *MockEnvReaderMockRecorder {
	return m.recorder
}

// LookupEnv mocks base method.
func (m *MockEnvReader) LookupEnv(arg0 string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupEnv", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// LookupEnv indicates an expected call of LookupEnv.
func (mr *MockEnvReaderMockRecorder) LookupEnv(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupEnv", reflect.TypeOf((*MockEnvReader)(nil).LookupEnv), arg0)
}
