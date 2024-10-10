//go:build !integration
// +build !integration

package utils

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/rysk-finance/v2_client_go/utils/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WebSocketUnitTestSuite struct {
	suite.Suite
}

func (s *WebSocketUnitTestSuite) SetupSuite() {
}

func TestRunSuiteUnit_WebSocketUnitTestSuite(t *testing.T) {
	suite.Run(t, new(WebSocketUnitTestSuite))
}

func (s *WebSocketUnitTestSuite) TestUnit_SendRPCRequest() {
	mockConnection := new(mocks.MockWebSocketConnection)
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
	mockConnection := new(mocks.MockWebSocketConnection)
	invalidRequest := make(chan int)

	err := SendRPCRequest(mockConnection, invalidRequest)
	require.Error(s.T(), err)
}

func (s *WebSocketUnitTestSuite) TestUnit_SendRPCRequest_WriteError() {
	mockConnection := new(mocks.MockWebSocketConnection)
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
