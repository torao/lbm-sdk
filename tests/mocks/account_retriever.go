// Code generated by MockGen. DO NOT EDIT.
// Source: client/account_retriever.go

// Package mocks is a generated GoMock package.
package mocks

import (
	client "github.com/cosmos/cosmos-sdk/client"
	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAccountRetriever is a mock of AccountRetriever interface
type MockAccountRetriever struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRetrieverMockRecorder
}

// MockAccountRetrieverMockRecorder is the mock recorder for MockAccountRetriever
type MockAccountRetrieverMockRecorder struct {
	mock *MockAccountRetriever
}

// NewMockAccountRetriever creates a new mock instance
func NewMockAccountRetriever(ctrl *gomock.Controller) *MockAccountRetriever {
	mock := &MockAccountRetriever{ctrl: ctrl}
	mock.recorder = &MockAccountRetrieverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountRetriever) EXPECT() *MockAccountRetrieverMockRecorder {
	return m.recorder
}

// EnsureExists mocks base method
func (m *MockAccountRetriever) EnsureExists(nodeQuerier client.NodeQuerier, addr types.AccAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureExists", nodeQuerier, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureExists indicates an expected call of EnsureExists
func (mr *MockAccountRetrieverMockRecorder) EnsureExists(nodeQuerier, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureExists", reflect.TypeOf((*MockAccountRetriever)(nil).EnsureExists), nodeQuerier, addr)
}

// GetAccountNumberSequence mocks base method
func (m *MockAccountRetriever) GetAccountNumberSequence(nodeQuerier client.NodeQuerier, addr types.AccAddress) (uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountNumberSequence", nodeQuerier, addr)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAccountNumberSequence indicates an expected call of GetAccountNumberSequence
func (mr *MockAccountRetrieverMockRecorder) GetAccountNumberSequence(nodeQuerier, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountNumberSequence", reflect.TypeOf((*MockAccountRetriever)(nil).GetAccountNumberSequence), nodeQuerier, addr)
}

// MockNodeQuerier is a mock of NodeQuerier interface
type MockNodeQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockNodeQuerierMockRecorder
}

// MockNodeQuerierMockRecorder is the mock recorder for MockNodeQuerier
type MockNodeQuerierMockRecorder struct {
	mock *MockNodeQuerier
}

// NewMockNodeQuerier creates a new mock instance
func NewMockNodeQuerier(ctrl *gomock.Controller) *MockNodeQuerier {
	mock := &MockNodeQuerier{ctrl: ctrl}
	mock.recorder = &MockNodeQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodeQuerier) EXPECT() *MockNodeQuerierMockRecorder {
	return m.recorder
}

// QueryWithData mocks base method
func (m *MockNodeQuerier) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryWithData", path, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryWithData indicates an expected call of QueryWithData
func (mr *MockNodeQuerierMockRecorder) QueryWithData(path, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryWithData", reflect.TypeOf((*MockNodeQuerier)(nil).QueryWithData), path, data)
}