//go:build !integration
// +build !integration

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HttpUnitTestSuite struct {
	suite.Suite
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	var res *http.Response
	if args.Get(0) != nil {
		res = args.Get(0).(*http.Response)
	}
	return res, args.Error(1)
}

func (s *HttpUnitTestSuite) SetupSuite() {
}

func TestRunSuiteUnit_HttpUnitTestSuite(t *testing.T) {
	suite.Run(t, new(HttpUnitTestSuite))
}

func (s *HttpUnitTestSuite) TestUnit_GetHTTPClient() {
	timeout := 5 * time.Second
	client := GetHTTPClient(timeout)
	require.NotNil(s.T(), client)
	require.Equal(s.T(), timeout, client.Timeout)

	transport, ok := client.Transport.(*http.Transport)
	require.True(s.T(), ok)
	require.Equal(s.T(), 100, transport.MaxIdleConns)
	require.Equal(s.T(), 100, transport.MaxConnsPerHost)
	require.Equal(s.T(), 100, transport.MaxIdleConnsPerHost)
}

func (s *HttpUnitTestSuite) TestUnit_CreateHTTPRequestWithBody() {
	body := map[string]interface{}{
		"key": "value",
	}

	req, err := CreateHTTPRequestWithBody(http.MethodPost, "http://example.com", body)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), req)
	require.Equal(s.T(), http.MethodPost, req.Method)
	require.Equal(s.T(), "http://example.com", req.URL.String())

	expectedBody, _ := json.Marshal(body)
	actualBody := make([]byte, req.ContentLength)
	req.Body.Read(actualBody)
	require.Equal(s.T(), expectedBody, actualBody)
}

func (s *HttpUnitTestSuite) TestUnit_CreateHTTPRequestWithBody_MarshalError() {
	unsupportedType := make(chan int)
	req, err := CreateHTTPRequestWithBody(http.MethodPost, "http://example.com", unsupportedType)
	require.Error(s.T(), err)
	require.Nil(s.T(), req)
}

func (s *HttpUnitTestSuite) TestUnit_CreateHTTPRequestWithBody_NilBody() {
	req, err := CreateHTTPRequestWithBody(http.MethodGet, "http://example.com", nil)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), req)
	require.Equal(s.T(), http.MethodGet, req.Method)
	require.Equal(s.T(), "http://example.com", req.URL.String())
	require.NotNil(s.T(), req.Body)
}

func (s *HttpUnitTestSuite) TestUnit_SendHTTPRequest() {
	mockClient := new(MockHTTPClient)
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(s.T(), err)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("success")),
	}
	mockClient.On("Do", req).Return(mockResponse, nil)

	res, err := SendHTTPRequest(mockClient, req)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	require.Equal(s.T(), http.StatusOK, res.StatusCode)
	mockClient.AssertExpectations(s.T())
}

func (s *HttpUnitTestSuite) TestUnit_SendHTTPRequest_ClientError() {
	mockClient := new(MockHTTPClient)
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(s.T(), err)
	mockClient.On("Do", req).Return(nil, errors.New("mock client error"))

	res, err := SendHTTPRequest(mockClient, req)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
	require.Equal(s.T(), "mock client error", err.Error())
	mockClient.AssertExpectations(s.T())
}

func (s *HttpUnitTestSuite) TestUnit_SendHTTPRequest_ServerError() {
	mockClient := new(MockHTTPClient)
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(s.T(), err)
	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("success")),
	}
	mockClient.On("Do", req).Return(mockResponse, nil)

	res, err := SendHTTPRequest(mockClient, req)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	require.Equal(s.T(), http.StatusInternalServerError, res.StatusCode)
	mockClient.AssertExpectations(s.T())
}
