package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockWebSocketConnection struct {
	mock.Mock
}

func (m *MockWebSocketConnection) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}
