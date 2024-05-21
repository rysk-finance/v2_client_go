package utils

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/eldief/go100x/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

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
