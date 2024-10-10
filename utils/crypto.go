package utils

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rysk-finance/v2_client_go/types"
)

// AddressFromPrivateKey derives the Ethereum address from a hexadecimal private key string.
//
// This function takes a hexadecimal string representing an Ethereum private key
// (with or without the '0x' prefix) and computes the corresponding Ethereum address.
//
// Parameters:
//   - privateKeyHex: Hexadecimal string representing the Ethereum private key.
//
// Returns:
//   - string: The Ethereum address derived from the private key.
func AddressFromPrivateKey(privateKeyHex string) string {
	// Convert private key hex string to an ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}

	// Derive the Ethereum address (EOA) from the public key
	publicKey := privateKey.Public()

	// No need to assert if cast went fine since 'crypto.HexToECDSA' is already validating 'privateKeyHex'
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}

// GetTransactionParams retrieves transaction parameters required for sending a transaction.
//
// This function queries the Ethereum network using the provided Ethereum client (`ethClient`)
// to fetch the following transaction parameters:
// - Nonce: The transaction count of the sender's address.
// - GasPrice: The current gas price for the transaction.
// - ChainID: The ID of the Ethereum chain the transaction will be sent on.
// - GasLimit: The maximum amount of gas that can be used for the transaction.
//
// Parameters:
//   - ctx: The context for the Ethereum client operations.
//   - ethClient: Interface for interacting with the Ethereum blockchain.
//   - privateKey: The sender's private key for signing the transaction.
//   - from: The sender's Ethereum address.
//   - to: The recipient's Ethereum address (optional for contract creation).
//   - data: The data payload for the transaction (optional).
//
// Returns:
//   - nonce: The current nonce (transaction count) of the sender's address.
//   - gasPrice: The current gas price in Wei.
//   - chainID: The ID of the Ethereum chain.
//   - gasLimit: The maximum gas limit for the transaction.
//   - err: Any error encountered during the retrieval of parameters.
func GetTransactionParams(
	ctx context.Context,
	ethClient types.IEthClient,
	privateKey *ecdsa.PrivateKey,
	from *common.Address,
	to *common.Address,
	data *[]byte,
) (nonce uint64, gasPrice *big.Int, chainID *big.Int, gasLimit uint64, err error) {
	// Get next nonce
	nonce, err = ethClient.PendingNonceAt(ctx, *from)
	if err != nil {
		return
	}
	// Get gas price
	gasPrice, err = ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return
	}
	// Get Chain ID
	chainID, err = ethClient.NetworkID(ctx)
	if err != nil {
		return
	}
	// Compute estimated gas usage
	gasLimit, err = ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From: *from,
		To:   to,
		Data: *data,
	})
	if err != nil {
		return
	}
	return
}
