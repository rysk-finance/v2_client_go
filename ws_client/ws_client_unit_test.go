package ws_client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/eldief/go100x/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WSClientUnitTestSuite struct {
	suite.Suite
	PrivateKeys    string
	PrivateKey     string
	Address        string
	BaseUrl        string
	RpcUrl         string
	Go100XWSClient *Go100XWSClient
	EthClient      types.IEthClient
}

func (s *WSClientUnitTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("WSClientUnitTestSuite.SetupSuite:  Error loading .env file:", err)
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
	suite.Run(t, new(WSClientUnitTestSuite))
}

func (s *WSClientUnitTestSuite) TestUnit_NewGo100XWSClient() {
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

func (s *WSClientUnitTestSuite) TestUnitNewGo100XWSClient_InvalidPrivateKey() {
	apiClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   "0x123",
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *WSClientUnitTestSuite) TestUnit_NewGo100XWSClient_InvalidRPCURL() {
	apiClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       "",
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *WSClientUnitTestSuite) TestUnit_NewGo100XWSClient_InvalidRPCWebsocketURL() {
	ogRPCURL := constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET]
	constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET] = "invalid_rpc_url"
	apiClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
	constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET] = ogRPCURL
}

func (s *WSClientUnitTestSuite) TestUnit_NewGo100XWSClient_InvalidStreamWebsocketURL() {
	ogStreamURL := constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET]
	constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET] = "invalid_stream_url"
	apiClient, err := NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
	constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET] = ogStreamURL
}

func (s *WSClientUnitTestSuite) TestUnit_ListProducts() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		for {
			_, message, _ := conn.ReadMessage()

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err := json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_LIST_PRODUCTS), requestBody.Method)
			require.Empty(s.T(), requestBody.Params)
		}
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	url := strings.Replace(mockHttpServer.URL, "http", "ws", 1)
	defer mockHttpServer.Close()
	rpcWebsocket, _, err := websocket.DefaultDialer.DialContext(
		context.Background(),
		url,
		http.Header{},
	)
	require.NoError(s.T(), err)
	s.Go100XWSClient.RPCConnection = rpcWebsocket
	s.Go100XWSClient.rpcUrl = url

	s.Go100XWSClient.ListProducts("69420")
}

func (s *WSClientUnitTestSuite) TestUnit_GetProduct() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		for {
			_, message, _ := conn.ReadMessage()

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err := json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_GET_PRODUCT), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Id string `json:"id"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), strconv.FormatInt(constants.PRODUCT_ETH_PERP.Id, 10), params.Id)
		}
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	url := strings.Replace(mockHttpServer.URL, "http", "ws", 1)
	defer mockHttpServer.Close()
	rpcWebsocket, _, err := websocket.DefaultDialer.DialContext(
		context.Background(),
		url,
		http.Header{},
	)
	require.NoError(s.T(), err)
	s.Go100XWSClient.RPCConnection = rpcWebsocket
	s.Go100XWSClient.rpcUrl = url

	s.Go100XWSClient.GetProduct("69420", &constants.PRODUCT_ETH_PERP)
}
