package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eldief/go100x/types"
)

// GetHTTPClient initializes a new `http.Client` with custom `http.Transport` and timeout.
// The `http.Transport` is configured with:
//   - 100 `MaxIdleConns`
//   - 100 `MaxConnsPerHost`
//   - 100 `MaxIdleConnsPerHost`
//
// The timeout for the `http.Client` can be configured via the `time.Duration` parameter.
//
// Parameters:
//   - timeout: Timeout duration for the `http.Client`.
//
// Returns:
//   - *http.Client: Configured HTTP client instance.
func GetHTTPClient(timeout time.Duration) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   timeout,
		Transport: t,
	}
}

// CreateHTTPRequestWithBody creates a new HTTP request with a request body.
//
// Parameters:
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.).
//   - uri: Request URI.
//   - body: Request body to be included in the HTTP request. It can be a string,
//     []byte, or any other type that can be marshaled into a valid HTTP request body.
//
// Returns:
//   - *http.Request: Created HTTP request instance.
//   - error: Any error encountered during request creation.
func CreateHTTPRequestWithBody(method string, uri string, body interface{}) (*http.Request, error) {
	// Marshal body into JSON.
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create API request.
	return http.NewRequest(method, uri, bytes.NewBuffer(bodyJSON))
}

// SendHTTPRequest sends an HTTP request using a provided `http.Client` and returns the response.
//
// Parameters:
//   - c: Custom HTTP client implementing `types.IHTTPClient`.
//   - req: HTTP request instance to be sent.
//
// Returns:
//   - *http.Response: HTTP response received from the server.
//   - error: Any error encountered during the HTTP request or response handling.
func SendHTTPRequest(c types.IHTTPClient, req *http.Request) (*http.Response, error) {
	// Send request
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
