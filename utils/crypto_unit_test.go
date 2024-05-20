package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CryptoUnitTestSuite struct {
	suite.Suite
	privateKey *ecdsa.PrivateKey
}

func (s *CryptoUnitTestSuite) SetupSuite() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	s.privateKey = privateKey
}

func TestRunSuiteUnit_CryptoUnitTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoUnitTestSuite))
}

func (s *CryptoUnitTestSuite) TestUnit_AddressFromPrivateKey() {
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(s.privateKey))
	address := AddressFromPrivateKey(privateKeyHex)
	expectedAddress := crypto.PubkeyToAddress(s.privateKey.PublicKey).Hex()
	require.Equal(s.T(), expectedAddress, address)
}

func (s *CryptoUnitTestSuite) TestUnit_AddressFromPrivateKey_InvalidHex() {
	invalidPrivateKeyHex := "invalid_hex_string"
	require.Panics(s.T(), func() {
		AddressFromPrivateKey(invalidPrivateKeyHex)
	})
}