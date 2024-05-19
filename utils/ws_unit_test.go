package utils

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WebSocketUnitTestSuite struct {
	suite.Suite
}

func (s *WebSocketUnitTestSuite) SetupSuite() {
}

func TestRun_WebSocketUnitTestSuite(t *testing.T) {
	suite.Run(t, new(WebSocketUnitTestSuite))
}

type MockWebSocketConnection struct {
	mock.Mock
}

func (m *MockWebSocketConnection) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

func (s *WebSocketUnitTestSuite) TestUnit_SendRPCRequest() {
	mockConnection := new(MockWebSocketConnection)
	request := map[string]interface{}{
		"method": "example",
		"params": nil,
	}
	requestJSON, err := json.Marshal(request)
	require.NoError(s.T(), err)
	mockConnection.On("WriteMessage", websocket.TextMessage, requestJSON).Return(nil)

	err = SendRPCRequest(mockConnection, request)
	require.NoError(s.T(), err)
	mockConnection.AssertExpectations(s.T())
}

func (s *WebSocketUnitTestSuite) TestUnit_SendRPCRequest_MarshalError() {
	mockConnection := new(MockWebSocketConnection)
	invalidRequest := make(chan int)

	err := SendRPCRequest(mockConnection, invalidRequest)
	require.Error(s.T(), err)
}

func (s *WebSocketUnitTestSuite) TestUnit_SendRPCRequest_WriteError() {
	mockConnection := new(MockWebSocketConnection)
	request := map[string]interface{}{
		"method": "example",
		"params": nil,
	}
	requestJSON, _ := json.Marshal(request)
	mockConnection.On("WriteMessage", websocket.TextMessage, requestJSON).Return(errors.New(""))

	err := SendRPCRequest(mockConnection, request)
	require.Error(s.T(), err)
	mockConnection.AssertExpectations(s.T())
}
