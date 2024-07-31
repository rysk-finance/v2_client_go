package rest_examples

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/eldief/go100x/api_client"
	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/joho/godotenv"
)

func Withdraw() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Println("Withdraw:: error loading .env file:", err)
		return
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("Withdraw:: no PRIVATE_KEYS found in %q file", ".env")
	}

	// Initialize new Go100XAPIClient
	client, err := api_client.NewGo100XAPIClient(&api_client.Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_MAINNET, // Blast Mainnet
		PrivateKey:   privateKey,                    // Your private key
		RpcUrl:       "https://rpc.blast.io",        // Public Blast Mainnet RPC url
		SubAccountId: 0,                             // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("Withdraw:: error intializing Go100XAPIClient: %v", err)
	}

	response, err := client.Withdraw(&types.WithdrawRequest{
		Quantity: constants.E16.String(),
		Nonce:    time.Now().UnixMicro(),
	})
	if err != nil {
		log.Fatalf("Withdraw:: error performing Withdraw request: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Withdraw:: error during reading response body: %v", err)

	}

	log.Printf("Withdraw:: response received: %d - %s", response.StatusCode, body)
}
