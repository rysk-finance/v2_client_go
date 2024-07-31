package rest_examples

import (
	"context"
	"log"
	"os"

	"github.com/eldief/go100x/api_client"
	"github.com/eldief/go100x/constants"
	"github.com/joho/godotenv"
)

func Deposit() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Println("Deposit:: error loading .env file:", err)
		return
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("Deposit:: no PRIVATE_KEYS found in %q file", ".env")
	}

	// Initialize new Go100XAPIClient
	client, err := api_client.NewGo100XAPIClient(&api_client.Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_MAINNET, // Blast Mainnet
		PrivateKey:   privateKey,                    // Your private key
		RpcUrl:       "https://rpc.blast.io",        // Public Blast Mainnet RPC url
		SubAccountId: 0,                             // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("Deposit:: error intializing Go100XAPIClient: %v", err)
	}

	// Approve USDB to 100x
	transaction, err := client.ApproveUSDB(context.Background(), constants.E20)
	if err != nil {
		log.Fatalf("Deposit:: error dispatching approve USDB transaction: %v", err)
	}
	receipt, err := client.WaitTransaction(context.Background(), transaction)
	if err != nil {
		log.Fatalf("Deposit:: error waiting for approve USDB transaction: %v", err)
	}
	log.Printf("Deposit:: approve USDB transaction receipt received: %s", receipt.TxHash)

	// Deposit USDB to 100x
	transaction, err = client.DepositUSDB(context.Background(), constants.E20)
	if err != nil {
		log.Fatalf("Deposit:: error dispatching deposit USDB transaction: %v", err)
	}
	receipt, err = client.WaitTransaction(context.Background(), transaction)
	if err != nil {
		log.Fatalf("Deposit:: error waiting for deposit USDB transaction: %v", err)
	}
	log.Printf("Deposit:: deposit USDB transaction receipt received: %s", receipt.TxHash)
}
