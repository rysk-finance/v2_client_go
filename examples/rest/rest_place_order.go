package rest_examples

import (
	"io"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/eldief/go100x/api_client"
	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/joho/godotenv"
)

func PlaceOrder() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Println("websocketPlaceOrder:: error loading .env file:", err)
		return
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("websocketPlaceOrder:: no PRIVATE_KEYS found in %q file", ".env")
	}

	// Initialize new Go100XAPIClient
	client, err := api_client.NewGo100XAPIClient(&api_client.Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_MAINNET, // Blast Mainnet
		PrivateKey:   privateKey,                    // Your private key
		RpcUrl:       "https://rpc.blast.io",        // Public Blast Mainnet RPC url
		SubAccountId: 0,                             // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: error initializing new Go100XAPIClient: %v", err)
	}

	// Create a new order
	response, err := client.NewOrder(&types.NewOrderRequest{
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
		log.Fatalf("PlaceOrder:: error performing NewOrder request: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("PlaceOrder:: error during reading response body: %v", err)

	}

	log.Printf("PlaceOrder:: response received: %d - %s", response.StatusCode, body)
}
