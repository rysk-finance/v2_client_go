package utils

import (
	"fmt"
	"go100x/src/constants"
	"go100x/src/types"
	"reflect"
	"unicode"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// SignMessage signs a message using EIP-712 and returns the signature.
func SignMessage(c *types.Client, endpoint types.Endpoint, message interface{}) (string, error) {
	// Get the primary type based on the endpoint.
	primaryType := constants.PRIMARY_TYPE[endpoint]

	// Create a pointer to the message value.
	messagePtr := reflect.ValueOf(message)
	if messagePtr.Kind() != reflect.Ptr {
		return "", fmt.Errorf("message must be a pointer to an interface")
	}

	// Generate the EIP-712 message using the provided primary type, client domain, and mapped message.
	unsignedMessage, err := generateEIP712Message(primaryType, c.Domain, mapMessage(message))
	if err != nil {
		return "", err
	}

	// Load the private key from hex
	privateKey, err := crypto.HexToECDSA(c.PrivateKey)
	if err != nil {
		return "", err
	}

	// Sign the message hash with the private key
	signature, err := crypto.Sign(unsignedMessage, privateKey)
	if err != nil {
		return "", err
	}

	// Convert the signature to a hex string
	signatureHex := common.Bytes2Hex(signature)

	// Return the signature
	return "0x" + signatureHex, nil
}

// mapMessage maps any struct to a map[string]interface{}.
func mapMessage(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Get the reflect.Value of the interface.
	value := reflect.ValueOf(data)

	// Dereference the pointer to get the underlying value.
	elemValue := value.Elem()

	// Iterate through fields of the struct.
	for i := 0; i < elemValue.NumField(); i++ {
		field := elemValue.Type().Field(i)
		fieldValue := elemValue.Field(i).Interface()

		// Add field name and value to the result map.
		result[lowercaseFirstRune(field.Name)] = fieldValue
	}

	return result
}

func lowercaseFirstRune(s string) string {
	if len(s) == 0 {
		return s
	}
	rs := []rune(s)
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

// generateEIP712Message generates an EIP-712 compliant message hash.
func generateEIP712Message(primaryType string, domain apitypes.TypedDataDomain, message apitypes.TypedDataMessage) ([]byte, error) {
	// Create a TypedData instance with the provided parameters.
	signerData := apitypes.TypedData{
		Types:       types.EIP712_TYPES,
		PrimaryType: primaryType,
		Domain:      domain,
		Message:     message,
	}

	// Hash the structured data of the message.
	typedDataHash, err := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	if err != nil {
		return nil, err
	}

	// Hash the EIP712 domain separator data.
	domainSeparator, err := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())
	if err != nil {
		return nil, err
	}

	// Construct the raw data for hashing.
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))

	// Hash the raw data using Keccak256 and return.
	return crypto.Keccak256(rawData), nil
}
