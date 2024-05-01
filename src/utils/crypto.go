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
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error converting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}
