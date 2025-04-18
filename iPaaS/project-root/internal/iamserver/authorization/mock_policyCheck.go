// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/yuansuan/ticp/common/project-root-iam/internal/iamserver/authorization (interfaces: PolicyCheck)

// Package authorization is a generated GoMock package.
package authorization

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ladon "github.com/ory/ladon"
)

// MockPolicyCheck is a mock of PolicyCheck interface.
type MockPolicyCheck struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyCheckMockRecorder
}

// MockPolicyCheckMockRecorder is the mock recorder for MockPolicyCheck.
type MockPolicyCheckMockRecorder struct {
	mock *MockPolicyCheck
}

// NewMockPolicyCheck creates a new mock instance.
func NewMockPolicyCheck(ctrl *gomock.Controller) *MockPolicyCheck {
	mock := &MockPolicyCheck{ctrl: ctrl}
	mock.recorder = &MockPolicyCheckMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPolicyCheck) EXPECT() *MockPolicyCheckMockRecorder {
	return m.recorder
}

// DoPoliciesAllow mocks base method.
func (m *MockPolicyCheck) DoPoliciesAllow(arg0 *ladon.Request, arg1 []ladon.DefaultPolicy) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoPoliciesAllow", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoPoliciesAllow indicates an expected call of DoPoliciesAllow.
func (mr *MockPolicyCheckMockRecorder) DoPoliciesAllow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoPoliciesAllow", reflect.TypeOf((*MockPolicyCheck)(nil).DoPoliciesAllow), arg0, arg1)
}
