package types

type WSMethod string

type WebsocketRequest struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  WSMethod    `json:"method"`
	Params  interface{} `json:"params"`
}

type WebsocketResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      string          `json:"id" binding:"required"`
	Success bool            `json:"success" binding:"required"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *WebsocketError `json:"error,omitempty"`
}

type WebsocketError struct {
	Code    int         `json:"code,omitempty" binding:"required"`
	Message string      `json:"message,omitempty" binding:"required"`
	Data    interface{} `json:"data,omitempty"`
}
