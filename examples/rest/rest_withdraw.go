package rest_examples

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rysk-finance/v2_client_go/api_client"
	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/rysk-finance/v2_client_go/types"
)

func Withdraw() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Withdraw:: %v", err)
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("Withdraw:: no PRIVATE_KEYS found")
	}

	// Initialize new RyskV2APIClient
	client, err := api_client.NewRyskV2APIClient(&api_client.RyskV2APIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,                  // Arbistrum sepolia testnet
		PrivateKey:   privateKey,                                     // Your private key
		RpcUrl:       "https://arbitrum-sepolia.gateway.tenderly.co", // Public Arbistrum sepolia testnet RPC url
		SubAccountId: 0,                                              // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("Withdraw:: %v", err)
	}

	response, err := client.Withdraw(&types.WithdrawRequest{
		Quantity: constants.E9.String(),
		Nonce:    time.Now().UnixMicro(),
	})
	if err != nil {
		log.Fatalf("Withdraw:: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Withdraw:: %v", err)
	}

	log.Printf("Withdraw:: response received: %d - %s", response.StatusCode, body)
}
