package utils

import (
	"encoding/json"

	"github.com/eldief/go100x/types"
	"github.com/gorilla/websocket"
)

// SendRPCRequest sends a RPC request via a WebSocket connection and returns any error encountered.
//
// Parameters:
//   - connection: WebSocket connection implementing `types.IWSConnection` interface.
//   - request: JSON-RPC request payload to be sent over the WebSocket connection.
//
// Returns:
//   - error: Returns an error if there was any issue sending the RPC request.
func SendRPCRequest(connection types.IWSConnection, request interface{}) error {
	// Marshal request into JSON.
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Send RPC request.
	return connection.WriteMessage(websocket.TextMessage, body)
}
