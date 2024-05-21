package utils

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/eldief/go100x/utils/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
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

func (s *CryptoUnitTestSuite) TestUnit_GetTransactionParams() {
	mockEthClient := new(mocks.MockEthClient)
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(1), nil)

	nonce, gasPrice, chainID, gasLimit, err := GetTransactionParams(
		context.Background(),
		mockEthClient,
		s.privateKey,
		&common.MaxAddress,
		&common.MaxAddress,
		new([]byte),
	)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), nonce)
	require.NotEmpty(s.T(), gasPrice)
	require.NotEmpty(s.T(), chainID)
	require.NotEmpty(s.T(), gasLimit)
}

func (s *CryptoUnitTestSuite) TestUnit_GetTransactionParams_ErrorPendingNonceAt() {
	mockEthClient := new(mocks.MockEthClient)
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), fmt.Errorf("failed to get nonce"))

	nonce, gasPrice, chainID, gasLimit, err := GetTransactionParams(
		context.Background(),
		mockEthClient,
		s.privateKey,
		&common.MaxAddress,
		&common.MaxAddress,
		new([]byte),
	)
	require.Error(s.T(), err)
	require.Empty(s.T(), nonce)
	require.Empty(s.T(), gasPrice)
	require.Empty(s.T(), chainID)
	require.Empty(s.T(), gasLimit)
}

func (s *CryptoUnitTestSuite) TestUnit_GetTransactionParams_ErrorSuggestGasPrice() {
	mockEthClient := new(mocks.MockEthClient)
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return((*big.Int)(nil), fmt.Errorf("failed to get gas price"))

	nonce, gasPrice, chainID, gasLimit, err := GetTransactionParams(
		context.Background(),
		mockEthClient,
		s.privateKey,
		&common.MaxAddress,
		&common.MaxAddress,
		new([]byte),
	)
	require.Error(s.T(), err)
	require.NotEmpty(s.T(), nonce)
	require.Empty(s.T(), gasPrice)
	require.Empty(s.T(), chainID)
	require.Empty(s.T(), gasLimit)
}

func (s *CryptoUnitTestSuite) TestUnit_GetTransactionParams_ErrorNetworkID() {
	mockEthClient := new(mocks.MockEthClient)
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return((*big.Int)(nil), fmt.Errorf("failed to get network ID"))

	nonce, gasPrice, chainID, gasLimit, err := GetTransactionParams(
		context.Background(),
		mockEthClient,
		s.privateKey,
		&common.MaxAddress,
		&common.MaxAddress,
		new([]byte),
	)
	require.Error(s.T(), err)
	require.NotEmpty(s.T(), nonce)
	require.NotEmpty(s.T(), gasPrice)
	require.Empty(s.T(), chainID)
	require.Empty(s.T(), gasLimit)
}

func (s *CryptoUnitTestSuite) TestUnit_GetTransactionParams_ErrorEstimateGas() {
	mockEthClient := new(mocks.MockEthClient)
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(0), fmt.Errorf("failed to estimate gas"))

	nonce, gasPrice, chainID, gasLimit, err := GetTransactionParams(
		context.Background(),
		mockEthClient,
		s.privateKey,
		&common.MaxAddress,
		&common.MaxAddress,
		new([]byte),
	)
	require.Error(s.T(), err)
	require.NotEmpty(s.T(), nonce)
	require.NotEmpty(s.T(), gasPrice)
	require.NotEmpty(s.T(), chainID)
	require.Empty(s.T(), gasLimit)
}
