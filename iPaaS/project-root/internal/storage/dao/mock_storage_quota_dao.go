// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao (interfaces: StorageQuotaDao)

// Package dao is a generated GoMock package.
package dao

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	snowflake "github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	model "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	xorm "xorm.io/xorm"
)

// MockStorageQuotaDao is a mock of StorageQuotaDao interface.
type MockStorageQuotaDao struct {
	ctrl     *gomock.Controller
	recorder *MockStorageQuotaDaoMockRecorder
}

// MockStorageQuotaDaoMockRecorder is the mock recorder for MockStorageQuotaDao.
type MockStorageQuotaDaoMockRecorder struct {
	mock *MockStorageQuotaDao
}

// NewMockStorageQuotaDao creates a new mock instance.
func NewMockStorageQuotaDao(ctrl *gomock.Controller) *MockStorageQuotaDao {
	mock := &MockStorageQuotaDao{ctrl: ctrl}
	mock.recorder = &MockStorageQuotaDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageQuotaDao) EXPECT() *MockStorageQuotaDaoMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockStorageQuotaDao) Get(arg0 *xorm.Session, arg1 *model.StorageQuota) (bool, *model.StorageQuota, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuota", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*model.StorageQuota)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockStorageQuotaDaoMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuota", reflect.TypeOf((*MockStorageQuotaDao)(nil).Get), arg0, arg1)
}

// Insert mocks base method.
func (m *MockStorageQuotaDao) Insert(arg0 *xorm.Session, arg1 *model.StorageQuota) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockStorageQuotaDaoMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockStorageQuotaDao)(nil).Insert), arg0, arg1)
}

// List mocks base method.
func (m *MockStorageQuotaDao) List(arg0 *xorm.Session, arg1, arg2 int) ([]*model.StorageQuota, error, int64, int64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*model.StorageQuota)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].(int64)
	return ret0, ret1, ret2, ret3
}

// List indicates an expected call of List.
func (mr *MockStorageQuotaDaoMockRecorder) List(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockStorageQuotaDao)(nil).List), arg0, arg1, arg2)
}

// Total mocks base method.
func (m *MockStorageQuotaDao) Total(arg0 *xorm.Session) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Total", arg0)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Total indicates an expected call of Total.
func (mr *MockStorageQuotaDaoMockRecorder) Total(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Total", reflect.TypeOf((*MockStorageQuotaDao)(nil).Total), arg0)
}

// Update mocks base method.
func (m *MockStorageQuotaDao) Update(arg0 *xorm.Session, arg1 snowflake.ID, arg2 *model.StorageQuota) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockStorageQuotaDaoMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStorageQuotaDao)(nil).Update), arg0, arg1, arg2)
}
