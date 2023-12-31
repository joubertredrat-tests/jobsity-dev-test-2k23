// Code generated by MockGen. DO NOT EDIT.
// Source: chat/domain/event.go
//
// Generated by this command:
//
//	mockgen -package=mock -source=chat/domain/event.go
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	domain "joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMessageEvent is a mock of MessageEvent interface.
type MockMessageEvent struct {
	ctrl     *gomock.Controller
	recorder *MockMessageEventMockRecorder
}

// MockMessageEventMockRecorder is the mock recorder for MockMessageEvent.
type MockMessageEventMockRecorder struct {
	mock *MockMessageEvent
}

// NewMockMessageEvent creates a new mock instance.
func NewMockMessageEvent(ctrl *gomock.Controller) *MockMessageEvent {
	mock := &MockMessageEvent{ctrl: ctrl}
	mock.recorder = &MockMessageEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageEvent) EXPECT() *MockMessageEventMockRecorder {
	return m.recorder
}

// StockCommandReceived mocks base method.
func (m *MockMessageEvent) StockCommandReceived(ctx context.Context, message domain.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StockCommandReceived", ctx, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// StockCommandReceived indicates an expected call of StockCommandReceived.
func (mr *MockMessageEventMockRecorder) StockCommandReceived(ctx, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StockCommandReceived", reflect.TypeOf((*MockMessageEvent)(nil).StockCommandReceived), ctx, message)
}
