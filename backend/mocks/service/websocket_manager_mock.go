// Code generated by MockGen. DO NOT EDIT.
// Source: /home/nakamura/program/webrtc-practice/backend/internal/domain/service/websocket_manager.go
//
// Generated by this command:
//
//	mockgen -source=/home/nakamura/program/webrtc-practice/backend/internal/domain/service/websocket_manager.go -destination=./mocks/service/websocket_manager_mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	entity "example.com/webrtc-practice/internal/domain/entity"
	service "example.com/webrtc-practice/internal/domain/service"
	gomock "go.uber.org/mock/gomock"
)

// MockWebSocketConnection is a mock of WebSocketConnection interface.
type MockWebSocketConnection struct {
	ctrl     *gomock.Controller
	recorder *MockWebSocketConnectionMockRecorder
	isgomock struct{}
}

// MockWebSocketConnectionMockRecorder is the mock recorder for MockWebSocketConnection.
type MockWebSocketConnectionMockRecorder struct {
	mock *MockWebSocketConnection
}

// NewMockWebSocketConnection creates a new mock instance.
func NewMockWebSocketConnection(ctrl *gomock.Controller) *MockWebSocketConnection {
	mock := &MockWebSocketConnection{ctrl: ctrl}
	mock.recorder = &MockWebSocketConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebSocketConnection) EXPECT() *MockWebSocketConnectionMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockWebSocketConnection) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockWebSocketConnectionMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWebSocketConnection)(nil).Close))
}

// ReadMessage mocks base method.
func (m *MockWebSocketConnection) ReadMessage() (int, entity.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(entity.Message)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockWebSocketConnectionMockRecorder) ReadMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockWebSocketConnection)(nil).ReadMessage))
}

// WriteMessage mocks base method.
func (m *MockWebSocketConnection) WriteMessage(arg0 entity.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMessage indicates an expected call of WriteMessage.
func (mr *MockWebSocketConnectionMockRecorder) WriteMessage(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMessage", reflect.TypeOf((*MockWebSocketConnection)(nil).WriteMessage), arg0)
}

// MockWebsocketManager is a mock of WebsocketManager interface.
type MockWebsocketManager struct {
	ctrl     *gomock.Controller
	recorder *MockWebsocketManagerMockRecorder
	isgomock struct{}
}

// MockWebsocketManagerMockRecorder is the mock recorder for MockWebsocketManager.
type MockWebsocketManagerMockRecorder struct {
	mock *MockWebsocketManager
}

// NewMockWebsocketManager creates a new mock instance.
func NewMockWebsocketManager(ctrl *gomock.Controller) *MockWebsocketManager {
	mock := &MockWebsocketManager{ctrl: ctrl}
	mock.recorder = &MockWebsocketManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebsocketManager) EXPECT() *MockWebsocketManagerMockRecorder {
	return m.recorder
}

// DeleteConnection mocks base method.
func (m *MockWebsocketManager) DeleteConnection(conn service.WebSocketConnection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConnection", conn)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConnection indicates an expected call of DeleteConnection.
func (mr *MockWebsocketManagerMockRecorder) DeleteConnection(conn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConnection", reflect.TypeOf((*MockWebsocketManager)(nil).DeleteConnection), conn)
}

// ExistsByID mocks base method.
func (m *MockWebsocketManager) ExistsByID(id string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByID", id)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ExistsByID indicates an expected call of ExistsByID.
func (mr *MockWebsocketManagerMockRecorder) ExistsByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByID", reflect.TypeOf((*MockWebsocketManager)(nil).ExistsByID), id)
}

// GetConnectionByID mocks base method.
func (m *MockWebsocketManager) GetConnectionByID(id string) (service.WebSocketConnection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionByID", id)
	ret0, _ := ret[0].(service.WebSocketConnection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionByID indicates an expected call of GetConnectionByID.
func (mr *MockWebsocketManagerMockRecorder) GetConnectionByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionByID", reflect.TypeOf((*MockWebsocketManager)(nil).GetConnectionByID), id)
}

// RegisterConnection mocks base method.
func (m *MockWebsocketManager) RegisterConnection(conn service.WebSocketConnection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterConnection", conn)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterConnection indicates an expected call of RegisterConnection.
func (mr *MockWebsocketManagerMockRecorder) RegisterConnection(conn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterConnection", reflect.TypeOf((*MockWebsocketManager)(nil).RegisterConnection), conn)
}

// RegisterID mocks base method.
func (m *MockWebsocketManager) RegisterID(conn service.WebSocketConnection, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterID", conn, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterID indicates an expected call of RegisterID.
func (mr *MockWebsocketManagerMockRecorder) RegisterID(conn, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterID", reflect.TypeOf((*MockWebsocketManager)(nil).RegisterID), conn, id)
}
