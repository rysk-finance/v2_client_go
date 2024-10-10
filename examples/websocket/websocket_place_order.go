package websocket_examples

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/rysk-finance/v2_client_go/types"
	"github.com/rysk-finance/v2_client_go/ws_client"
)

func PlaceOrder() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("PlaceOrder:: no PRIVATE_KEYS found")
	}

	// Initialize RyskV2WSClient
	client, err := ws_client.NewRyskV2WSClient(&ws_client.RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,                  // Arbistrum sepolia testnet
		PrivateKey:   privateKey,                                     // Your private key
		RpcUrl:       "https://arbitrum-sepolia.gateway.tenderly.co", // Public Arbistrum sepolia testnet RPC url
		SubAccountId: 0,                                              // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
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
		log.Fatalf("PlaceOrder:: %v", err)
	}

	// Send order
	err = client.NewOrder("NEW_ORDER", &types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       new(big.Int).Mul(big.NewInt(3150), constants.E17).String(),
		Quantity:    constants.E4.String(),
		Expiration:  time.Now().Add(10 * time.Minute).UnixMilli(),
		Nonce:       time.Now().UnixMicro(),
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	// Keep alive
	select {}
}

func startPinging(client *ws_client.RyskV2WSClient) {
	log.Println("startPinging:: start sending Pings...")

	for {
		err := client.RPCConnection.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			log.Fatalf("startPinging:: %v", err)
		}
		time.Sleep(30 * time.Second)
	}
}

func listenForRPCResponses(client *ws_client.RyskV2WSClient, responseChan chan types.WebsocketResponse) {
	log.Println("listenForStreamResponses:: start listening for RPC responses...")

	for {
		// Wait for messages
		_, data, err := client.RPCConnection.ReadMessage()
		if err != nil {
			log.Fatalf("listenForRPCResponses:: %v", err)
		}

		// Unmarshal response
		var response types.WebsocketResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Printf("listenForRPCResponses:: %v", err)
			continue
		}

		// Dispatch response to channel
		responseChan <- response
	}
}

func listenForStreamResponses(client *ws_client.RyskV2WSClient, responseChan chan types.WebsocketResponse) {
	log.Println("listenForStreamResponses:: start listening for Stream responses...")

	for {
		// Wait for messages
		_, data, err := client.StreamConnection.ReadMessage()
		if err != nil {
			log.Fatalf("listenForStreamResponses:: %v", err)
		}

		// Unmarshal response
		var response types.WebsocketResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Printf("listenForStreamResponses:: %v", err)
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
