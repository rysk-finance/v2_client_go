package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// GetHTTPClient initialize a new `http.Client` with custom `http.Transport` and timeout.
// Transport is setup with 100 `MaxIdleConns`, 100 `MaxConnsPerHost` and 100 `MaxIdleConnsPerHost`.
// Timeout can be configured via a `time.Duration` parameter.
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
func CreateHTTPRequestWithBody(method string, uri string, body interface{}) (*http.Request, error) {
	// Marshal body into JSON.
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create API request.
	return http.NewRequest(method, uri, bytes.NewBuffer(bodyJSON))
}

// SendHTTPRequest send HTTP request using a `http.Client` and returns response as string.
func SendHTTPRequest(c *http.Client, req *http.Request) (*http.Response, error) {
	// Send request
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil

	// defer res.Body.Close()

	// // Read response
	// response, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	return "", err
	// }

	// return string(response), nil
}
