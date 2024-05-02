package utils

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// SignMessage signs a message using EIP-712 and returns the signature.
func SignMessage(c *types.Go100XClient, primaryType types.PrimaryType, message interface{}) (string, error) {
	// Map message to `TypedDataMessage` interface.
	typedDataMessage, err := mapMessageToTypedData(message)
	if err != nil {
		return "", err
	}

	// Generate the EIP-712 message using the provided primary type, client `TypedDataDomain`, and `TypedDataMessage` message.
	unsignedMessage, err := generateEIP712Message(primaryType, c.Domain, typedDataMessage)
	if err != nil {
		return "", err
	}

	// Load the private key from hex.
	privateKey, err := crypto.HexToECDSA(c.PrivateKey)
	if err != nil {
		return "", err
	}

	// Sign EIP-712 message and return the signature.
	return signEIP712Message(unsignedMessage, privateKey)
}

// mapMessageToTypedData maps any struct to `TypedDataMessage`.
func mapMessageToTypedData(message interface{}) (map[string]interface{}, error) {
	// Convert message to JSON.
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// Instance empty `TypedDataMessage`.
	var typedData map[string]interface{}

	// Populate the empty `TypedDataMessage` with message values.
	err = json.Unmarshal(messageJSON, &typedData)
	if err != nil {
		return nil, err
	}

	return typedData, nil
}

// generateEIP712Message generates an EIP-712 compliant message hash.
func generateEIP712Message(primaryType types.PrimaryType, typedDataDomain apitypes.TypedDataDomain, typedDataMessage apitypes.TypedDataMessage) ([]byte, error) {
	// Create a TypedData instance with the provided parameters.
	signerData := apitypes.TypedData{
		Types:       types.EIP712_TYPES,
		PrimaryType: string(primaryType),
		Domain:      typedDataDomain,
		Message:     typedDataMessage,
	}

	// Hash the structured data of the message.
	typedDataHash, err := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	if err != nil {
		return nil, err
	}

	// Hash the EIP-712 domain separator data.
	domainSeparator, err := signerData.HashStruct(constants.EIP_712_DOMAIN, signerData.Domain.Map())
	if err != nil {
		return nil, err
	}

	// Construct the raw data for hashing.
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))

	// Hash the raw data using Keccak256 and return.
	return crypto.Keccak256(rawData), nil
}

// signEIP712Message returns a signed EIP-712 signature.
func signEIP712Message(unsignedMessage []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	// Sign EIP-712 message.
	signature, err := crypto.Sign(unsignedMessage, privateKey)
	if err != nil {
		return "", err
	}

	// Convert the signature to a hex string.
	signatureHex := common.Bytes2Hex(signature)

	// Return the signature prepending `0x`.
	return "0x" + signatureHex, nil
}
