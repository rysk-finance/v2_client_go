package utils

import (
	"crypto/ecdsa"

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
