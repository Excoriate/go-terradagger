// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/daggerio/core.go
//
// Generated by this command:
//
//	mockgen -destination=pkg/mocks/core.go -source=pkg/daggerio/core.go
//

// Package mock_daggerio is a generated GoMock package.
package mock_daggerio

import (
	"io"
	"reflect"

	"dagger.io/dagger"
	"go.uber.org/mock/gomock"
)

// MockBackendClient is a mock of BackendClient interface.
type MockBackendClient struct {
	ctrl     *gomock.Controller
	recorder *MockBackendClientMockRecorder
}

// MockBackendClientMockRecorder is the mock recorder for MockBackendClient.
type MockBackendClientMockRecorder struct {
	mock *MockBackendClient
}

// NewMockBackendClient creates a new mock instance.
func NewMockBackendClient(ctrl *gomock.Controller) *MockBackendClient {
	mock := &MockBackendClient{ctrl: ctrl}
	mock.recorder = &MockBackendClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackendClient) EXPECT() *MockBackendClientMockRecorder {
	return m.recorder
}

// CreateDaggerBackend mocks base method.
func (m *MockBackendClient) CreateDaggerBackend(options ...dagger.ClientOpt) (*dagger.Client, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateDaggerBackend", varargs...)
	ret0, _ := ret[0].(*dagger.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDaggerBackend indicates an expected call of CreateDaggerBackend.
func (mr *MockBackendClientMockRecorder) CreateDaggerBackend(options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDaggerBackend", reflect.TypeOf((*MockBackendClient)(nil).CreateDaggerBackend), options...)
}

// ResolveDaggerLogConfig mocks base method.
func (m *MockBackendClient) ResolveDaggerLogConfig(enableErrorsOnly bool) io.Writer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveDaggerLogConfig", enableErrorsOnly)
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

// ResolveDaggerLogConfig indicates an expected call of ResolveDaggerLogConfig.
func (mr *MockBackendClientMockRecorder) ResolveDaggerLogConfig(enableErrorsOnly any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveDaggerLogConfig", reflect.TypeOf((*MockBackendClient)(nil).ResolveDaggerLogConfig), enableErrorsOnly)
}
