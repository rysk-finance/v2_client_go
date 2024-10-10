package main

import (
	rest_examples "github.com/rysk-finance/v2_client_go/examples/rest"
	websocket_examples "github.com/rysk-finance/v2_client_go/examples/websocket"
)

// Comment not needed examples and run 'go run examples/main.go'
func main() {
	// REST
	rest_examples.Deposit()
	rest_examples.PlaceOrder()
	rest_examples.Withdraw()

	// WEBSOCKET
	websocket_examples.Deposit()
	websocket_examples.PlaceOrder()
}
