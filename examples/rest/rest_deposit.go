package rest_examples

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rysk-finance/v2_client_go/api_client"
	"github.com/rysk-finance/v2_client_go/constants"
)

func Deposit() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Deposit:: %v", err)
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("Deposit:: no PRIVATE_KEYS found")
	}

	// Initialize new RyskV2APIClient
	client, err := api_client.NewRyskV2APIClient(&api_client.RyskV2APIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,                  // Arbistrum sepolia testnet
		PrivateKey:   privateKey,                                     // Your private key
		RpcUrl:       "https://arbitrum-sepolia.gateway.tenderly.co", // Public Arbistrum sepolia testnet RPC url
		SubAccountId: 0,                                              // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("Deposit:: %v", err)
	}

	// Approve USDC to Rysk V2
	transaction, err := client.ApproveUSDC(context.Background(), constants.E9)
	if err != nil {
		log.Fatalf("Deposit:: %v", err)
	}

	receipt, err := client.WaitTransaction(context.Background(), transaction)
	if err != nil {
		log.Fatalf("Deposit:: %v", err)
	}
	log.Printf("Deposit:: approve USDC transaction receipt received: %s", receipt.TxHash)

	// Deposit USDC to Rysk V2
	transaction, err = client.DepositUSDC(context.Background(), constants.E9)
	if err != nil {
		log.Fatalf("Deposit:: %v", err)
	}

	receipt, err = client.WaitTransaction(context.Background(), transaction)
	if err != nil {
		log.Fatalf("Deposit:: %v", err)
	}
	log.Printf("Deposit:: deposit USDC transaction receipt received: %s", receipt.TxHash)
}
