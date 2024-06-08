package ws_client

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WsClientUnitTestSuite struct {
	suite.Suite
	PrivateKeys    string
	RpcUrl         string
	Go100XWSClient *Go100XWSClient
}

func (s *WsClientUnitTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		return
	}

	s.Go100XWSClient, _ = NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
}

func TestRunSuiteUnit_WsClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(WsClientUnitTestSuite))
}

func (s *WsClientUnitTestSuite) TestUnit_NewGo100XWSClient() {
	wsClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), constants.ENVIRONMENT_TESTNET, wsClient.env)
	require.Equal(s.T(), constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET], wsClient.rpcUrl)
	require.Equal(s.T(), constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET], wsClient.streamUrl)
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), wsClient.privateKeyString)
	require.Equal(s.T(), utils.AddressFromPrivateKey(strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x")), wsClient.addressString)
	require.Equal(s.T(), utils.AddressFromPrivateKey(strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x")), wsClient.address.String())
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), hex.EncodeToString(crypto.FromECDSA(wsClient.privateKey)))
	require.Equal(s.T(), constants.CIAO_ADDRESS[constants.ENVIRONMENT_TESTNET], wsClient.ciao.String())
	require.Equal(s.T(), constants.USDB_ADDRESS[constants.ENVIRONMENT_TESTNET], wsClient.usdb.String())
	require.NotNil(s.T(), wsClient.domain)
	require.Equal(s.T(), constants.DOMAIN_NAME, wsClient.domain.Name)
	require.Equal(s.T(), constants.DOMAIN_VERSION, wsClient.domain.Version)
	require.Equal(s.T(), int64(1), wsClient.SubAccountId)
	require.NotNil(s.T(), wsClient.RPCConnection)
	require.NotNil(s.T(), wsClient.StreamConnection)
	require.NotNil(s.T(), wsClient.EthClient)
}

func (s *WsClientUnitTestSuite) TestUnitNewGo100XWSClient_InvalidPrivateKey() {
	apiClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   "0x123",
		RpcUrl:       "",
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *WsClientUnitTestSuite) TestUnit_NewGo100XWSClient_InvalidRPCURL() {
	apiClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       "",
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}
