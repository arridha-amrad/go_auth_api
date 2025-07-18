// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/redis_repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/repositories/redis_repository.go -destination=mocks/mock_repositories/mock_redis_repo.go -package=mockrepositories
//

// Package mockrepositories is a generated GoMock package.
package mockrepositories

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockIRedisRepository is a mock of IRedisRepository interface.
type MockIRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRedisRepositoryMockRecorder
	isgomock struct{}
}

// MockIRedisRepositoryMockRecorder is the mock recorder for MockIRedisRepository.
type MockIRedisRepositoryMockRecorder struct {
	mock *MockIRedisRepository
}

// NewMockIRedisRepository creates a new mock instance.
func NewMockIRedisRepository(ctrl *gomock.Controller) *MockIRedisRepository {
	mock := &MockIRedisRepository{ctrl: ctrl}
	mock.recorder = &MockIRedisRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRedisRepository) EXPECT() *MockIRedisRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockIRedisRepository) Delete(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIRedisRepositoryMockRecorder) Delete(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIRedisRepository)(nil).Delete), key)
}

// HGet mocks base method.
func (m *MockIRedisRepository) HGet(key, field string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HGet", key, field)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HGet indicates an expected call of HGet.
func (mr *MockIRedisRepositoryMockRecorder) HGet(key, field any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HGet", reflect.TypeOf((*MockIRedisRepository)(nil).HGet), key, field)
}

// HGetAll mocks base method.
func (m *MockIRedisRepository) HGetAll(key string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HGetAll", key)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HGetAll indicates an expected call of HGetAll.
func (mr *MockIRedisRepositoryMockRecorder) HGetAll(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HGetAll", reflect.TypeOf((*MockIRedisRepository)(nil).HGetAll), key)
}

// HSet mocks base method.
func (m *MockIRedisRepository) HSet(key string, data map[string]any, expiry time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HSet", key, data, expiry)
	ret0, _ := ret[0].(error)
	return ret0
}

// HSet indicates an expected call of HSet.
func (mr *MockIRedisRepositoryMockRecorder) HSet(key, data, expiry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HSet", reflect.TypeOf((*MockIRedisRepository)(nil).HSet), key, data, expiry)
}
