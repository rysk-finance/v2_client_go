package websocket_examples

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/eldief/go100x/ws_client"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func PlaceOrder() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Println("PlaceOrder:: Error loading .env file:", err)
		return
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("PlaceOrder:: no PRIVATE_KEYS found in %q file", ".env")
	}

	// Initialize Go100XWSClient
	client, err := ws_client.NewGo100XWSClient(&ws_client.Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_MAINNET, // Blast Mainnet
		PrivateKey:   privateKey,                    // Your private key
		RpcUrl:       "https://rpc.blast.io",        // Public Blast Mainnet RPC url
		SubAccountId: 0,                             // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: error during %q: %v", "ws_client.NewGo100XWSClient", err)
	}

	// Start sending pings periodically
	go startPinging(client)

	// Get response channel
	responseChan := make(chan types.WebsocketResponse)
	defer close(responseChan)

	// Start listening for RPC messages
	go listenForRPCResponses(client, responseChan)

	// Start listening for Stream messages
	go listenForStreamResponses(client, responseChan)

	// Print responses
	go printResponses(responseChan)

	// Login
	err = client.Login("LOGIN")
	if err != nil {
		log.Fatalf("PlaceOrder:: error performing Login: %v", err)
	}

	// Send order
	err = client.NewOrder("NEW_ORDER", &types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       new(big.Int).Mul(big.NewInt(3150), constants.E18).String(),
		Quantity:    constants.E16.String(),
		Expiration:  time.Now().Add(10 * time.Minute).UnixMilli(),
		Nonce:       time.Now().UnixMicro(),
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: error performing NewOrder: %v", err)
	}

	// Keep alive
	select {}
}

func startPinging(client *ws_client.Go100XWSClient) {
	log.Println("startPinging:: start sending Pings...")

	for {
		err := client.RPCConnection.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			log.Fatalf("startPinging:: error sending ping: %v", err)
		}
		time.Sleep(30 * time.Second)
	}
}

func listenForRPCResponses(client *ws_client.Go100XWSClient, responseChan chan types.WebsocketResponse) {
	log.Println("listenForStreamResponses:: start listening for RPC responses...")

	for {
		// Wait for messages
		_, data, err := client.RPCConnection.ReadMessage()
		if err != nil {
			log.Fatalf("listenForRPCResponses:: error during reading message: %v", err)
		}

		// Unmarshal response
		var response types.WebsocketResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Printf("listenForRPCResponses:: error during unmarshaling response: %v", err)
			continue
		}

		// Dispatch response to channel
		responseChan <- response
	}
}

func listenForStreamResponses(client *ws_client.Go100XWSClient, responseChan chan types.WebsocketResponse) {
	log.Println("listenForStreamResponses:: start listening for Stream responses...")

	for {
		// Wait for messages
		_, data, err := client.StreamConnection.ReadMessage()
		if err != nil {
			log.Fatalf("listenForStreamResponses:: error during reading %q message: %v", "Login", err)
		}

		// Unmarshal response
		var response types.WebsocketResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Printf("listenForStreamResponses:: error during unmarshaling response: %v", err)
			continue
		}

		// Dispatch response to channel
		responseChan <- response
	}
}

func printResponses(responseChan chan types.WebsocketResponse) {
	for {
		response, ok := <-responseChan
		if !ok {
			log.Println("printResponses:: channel closed, exiting")
			return
		}

		if response.Error != nil {
			log.Printf("printResponses:: error response received for message ID %q: %d -  %s", response.ID, response.Error.Code, response.Error.Message)
			break
		}

		log.Printf("printResponses:: response received for message ID %q: %v", response.ID, response.Result)
	}
}
