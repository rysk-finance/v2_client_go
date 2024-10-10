package utils

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/rysk-finance/v2_client_go/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// SignMessage signs a message using EIP-712 and returns the signature.
//
// This function signs a message using the EIP-712 standard, which defines structured data hashing
// and signing for Ethereum. It takes the domain parameters (like name, version, chainId, verifyingContract),
// private key of the signer, primary type (the structure of the message), and the message itself.
//
// Parameters:
//   - domain: The domain parameters required for EIP-712 signing.
//   - privateKey: The private key of the signer in hexadecimal format (without '0x' prefix).
//   - primaryType: The primary type describing the structure of the message being signed.
//   - message: The message payload to be signed. It should conform to the primaryType structure.
//
// Returns:
//   - string: The signature of the message in hexadecimal format (with '0x' prefix).
//   - error: An error if the signing process fails, nil otherwise.
func SignMessage(domain apitypes.TypedDataDomain, privateKey string, primaryType types.PrimaryType, message interface{}) (string, error) {
	// Map message to `TypedDataMessage` interface.
	typedDataMessage, err := mapMessageToTypedData(message)
	if err != nil {
		return "", err
	}

	// Generate the EIP-712 message using the provided primary type, client `TypedDataDomain`, and `TypedDataMessage` message.
	unsignedMessage, err := generateEIP712Message(primaryType, domain, typedDataMessage)
	if err != nil {
		return "", err
	}

	// Load the private key from hex.
	hexPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	// Sign EIP-712 message and return the signature.
	return signEIP712Message(unsignedMessage, hexPrivateKey)
}

// mapMessageToTypedData maps any struct to `TypedDataMessage`.
//
// This function takes an input `message` of any struct type and converts it into
// a map[string]interface{} representation suitable for EIP-712 signing (TypedDataMessage).
// The struct fields are mapped to corresponding key-value pairs in the map.
//
// Parameters:
//   - message: The input struct message to be converted into `TypedDataMessage`.
//
// Returns:
//   - map[string]interface{}: The mapped representation of the struct as a map.
//   - error: An error if the mapping process fails, nil otherwise.
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
//
// This function computes the EIP-712 message hash using the provided parameters:
//   - primaryType: The primary type of the structured data schema.
//   - typedDataDomain: The domain separator for the structured data schema.
//   - typedDataMessage: The structured data message containing the data to hash.
//
// Parameters:
//   - primaryType: The primary type of the structured data schema.
//   - typedDataDomain: The domain separator for the structured data schema.
//   - typedDataMessage: The structured data message containing the data to hash.
//
// Returns:
//   - []byte: The EIP-712 compliant message hash.
//   - error: An error if the message hash computation fails, nil otherwise.
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
//
// This function signs the provided EIP-712 compliant message using the provided private key.
//
// Parameters:
//   - unsignedMessage: The unsigned EIP-712 message hash to sign.
//   - privateKey: The private key used for signing the message.
//
// Returns:
//   - string: The hexadecimal representation of the signed EIP-712 signature.
//   - error: An error if signing fails, nil otherwise.
func signEIP712Message(unsignedMessage []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	// Sign EIP-712 message.
	signature, err := crypto.Sign(unsignedMessage, privateKey)
	if err != nil {
		return "", err
	}
	signature[64] += 27

	// Convert the signature to a hex string.
	signatureHex := common.Bytes2Hex(signature)

	// Return the signature prepending `0x`.
	return "0x" + signatureHex, nil
}
