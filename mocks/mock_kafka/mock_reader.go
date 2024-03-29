// Code generated by MockGen. DO NOT EDIT.
// Source: src/server/kafka/kafka.go

// Package mock_kafka is a generated GoMock package.
package mock_kafka

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	kafka "github.com/segmentio/kafka-go"
	kafka0 "github.com/walmartdigital/katalog/server/kafka"
)

// MockReader is a mock of Reader interface
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockReader) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockReaderMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReader)(nil).Close))
}

// ReadMessage mocks base method
func (m *MockReader) ReadMessage(arg0 context.Context) (kafka.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage", arg0)
	ret0, _ := ret[0].(kafka.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMessage indicates an expected call of ReadMessage
func (mr *MockReaderMockRecorder) ReadMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockReader)(nil).ReadMessage), arg0)
}

// MockReaderFactory is a mock of ReaderFactory interface
type MockReaderFactory struct {
	ctrl     *gomock.Controller
	recorder *MockReaderFactoryMockRecorder
}

// MockReaderFactoryMockRecorder is the mock recorder for MockReaderFactory
type MockReaderFactoryMockRecorder struct {
	mock *MockReaderFactory
}

// NewMockReaderFactory creates a new mock instance
func NewMockReaderFactory(ctrl *gomock.Controller) *MockReaderFactory {
	mock := &MockReaderFactory{ctrl: ctrl}
	mock.recorder = &MockReaderFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReaderFactory) EXPECT() *MockReaderFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockReaderFactory) Create(arg0, arg1 string) kafka0.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(kafka0.Reader)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockReaderFactoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockReaderFactory)(nil).Create), arg0, arg1)
}
