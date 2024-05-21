package utils

import (
	"encoding/json"

	"github.com/eldief/go100x/types"
	"github.com/gorilla/websocket"
)

// SendRPCRequest send a RPC request via Websocket connection and returns any error.
func SendRPCRequest(connection types.IWSConnection, request interface{}) error {
	// Marshal request into JSON.
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Send RPC request.
	return connection.WriteMessage(websocket.TextMessage, body)
}
