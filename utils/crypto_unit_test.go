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

func TestRun_CryptoUnitTestSuite(t *testing.T) {
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
	require.PanicsWithError(s.T(), "invalid hex character 'i' in private key", func() {
		AddressFromPrivateKey(invalidPrivateKeyHex)
	})
}

func (s *CryptoUnitTestSuite) TestUnit_AddressFromPrivateKey_InvalidPublicKey() {
	invalidPrivateKey := make([]byte, 0)
	invalidPrivateKeyHex := hex.EncodeToString(invalidPrivateKey)

	require.PanicsWithError(s.T(), "invalid length, need 256 bits", func() {
		AddressFromPrivateKey(invalidPrivateKeyHex)
	})
}
