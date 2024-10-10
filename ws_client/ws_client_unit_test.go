//go:build !integration
// +build !integration

package ws_client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/rysk-finance/v2_client_go/types"
	"github.com/rysk-finance/v2_client_go/utils"
	"github.com/rysk-finance/v2_client_go/utils/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WSClientUnitTestSuite struct {
	suite.Suite
	PrivateKey     string
	Address        string
	BaseUrl        string
	RpcUrl         string
	RyskV2WSClient *RyskV2WSClient
	EthClient      types.IEthClient
}

func (s *WSClientUnitTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("WSClientUnitTestSuite.SetupSuite:  Error loading .env file:", err)
		return
	}

	wsClient, err := NewRyskV2WSClient(&RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.NoError(s.T(), err)

	s.RyskV2WSClient = wsClient
	s.PrivateKey = s.RyskV2WSClient.privateKeyString
	s.Address = utils.AddressFromPrivateKey(s.RyskV2WSClient.privateKeyString)
	s.EthClient = s.RyskV2WSClient.EthClient
}

func (s *WSClientUnitTestSuite) SetupTest() {
	s.RyskV2WSClient.privateKeyString = s.PrivateKey
	s.RyskV2WSClient.addressString = s.Address
	s.RyskV2WSClient.EthClient = s.EthClient
}

func TestRunSuiteUnit_WsClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(WSClientUnitTestSuite))
}

func (s *WSClientUnitTestSuite) TestUnit_NewRyskV2WSClient() {
	wsClient, err := NewRyskV2WSClient(&RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), constants.ENVIRONMENT_TESTNET, wsClient.env)
	require.Equal(s.T(), constants.API_BASE_URL[constants.ENVIRONMENT_TESTNET], wsClient.baseUrl)
	require.Equal(s.T(), constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET], wsClient.rpcUrl)
	require.Equal(s.T(), constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET], wsClient.streamUrl)
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), wsClient.privateKeyString)
	require.Equal(s.T(), utils.AddressFromPrivateKey(strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x")), wsClient.addressString)
	require.Equal(s.T(), utils.AddressFromPrivateKey(strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x")), wsClient.address.String())
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), hex.EncodeToString(crypto.FromECDSA(wsClient.privateKey)))
	require.Equal(s.T(), constants.CIAO_ADDRESS[constants.ENVIRONMENT_TESTNET], wsClient.ciao.String())
	require.Equal(s.T(), constants.USDC_ADDRESS[constants.ENVIRONMENT_TESTNET], wsClient.usdc.String())
	require.NotNil(s.T(), wsClient.domain)
	require.Equal(s.T(), constants.DOMAIN_NAME, wsClient.domain.Name)
	require.Equal(s.T(), constants.DOMAIN_VERSION, wsClient.domain.Version)
	require.Equal(s.T(), int64(1), wsClient.SubAccountId)
	require.NotNil(s.T(), wsClient.RPCConnection)
	require.NotNil(s.T(), wsClient.StreamConnection)
	require.NotNil(s.T(), wsClient.EthClient)
}

func (s *WSClientUnitTestSuite) TestUnitNewRyskV2WSClient_InvalidPrivateKey() {
	apiClient, err := NewRyskV2WSClient(&RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   "0x123",
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *WSClientUnitTestSuite) TestUnit_NewRyskV2WSClient_InvalidRPCURL() {
	apiClient, err := NewRyskV2WSClient(&RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       "",
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *WSClientUnitTestSuite) TestUnit_NewRyskV2WSClient_InvalidRPCWebsocketURL() {
	ogRPCURL := constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET]
	constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET] = "invalid_rpc_url"
	apiClient, err := NewRyskV2WSClient(&RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
	constants.WS_RPC_URL[constants.ENVIRONMENT_TESTNET] = ogRPCURL
}

func (s *WSClientUnitTestSuite) TestUnit_NewRyskV2WSClient_InvalidStreamWebsocketURL() {
	ogStreamURL := constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET]
	constants.WS_STREAM_URL[constants.ENVIRONMENT_TESTNET] = "invalid_stream_url"
	apiClient, err := NewRyskV2WSClient(&RyskV2WSClientConfiguration{
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
	done := make(chan struct{})
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
			require.Equal(s.T(), "null", string(requestBody.Params))
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.ListProducts("69420")
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_GetProduct() {
	done := make(chan struct{})
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
				Symbol string `json:"symbol"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), constants.PRODUCT_ETH_PERP.Symbol, params.Symbol)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.GetProduct("69420", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_ServerTime() {
	done := make(chan struct{})
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
			require.Equal(s.T(), string(constants.WS_METHOD_SERVER_TIME), requestBody.Method)
			require.Equal(s.T(), "null", string(requestBody.Params))
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.ServerTime("69420")
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_Login() {
	done := make(chan struct{})
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
			require.Equal(s.T(), string(constants.WS_METHOD_LOGIN), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account   string `json:"account"`
				Message   string `json:"message"`
				Timestamp uint64 `json:"timestamp"`
				Signature string `json:"signature"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), "I want to log into rysk.finance", params.Message)
			require.Greater(s.T(), params.Timestamp, uint64(0))
			require.NotEmpty(s.T(), params.Signature)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.Login("69420")
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_Login_BadAddress() {
	s.RyskV2WSClient.addressString = ""
	err := s.RyskV2WSClient.Login("69420")
	require.Error(s.T(), err)
}

func (s *WSClientUnitTestSuite) TestUnit_SessionStatus() {
	done := make(chan struct{})
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
			require.Equal(s.T(), string(constants.WS_METHOD_SESSION_STATUS), requestBody.Method)
			require.Equal(s.T(), "null", string(requestBody.Params))
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.SessionStatus("69420")
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubAccountList() {
	done := make(chan struct{})
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
			require.Equal(s.T(), string(constants.WS_METHOD_SUB_ACCOUNT_LIST), requestBody.Method)
			require.Equal(s.T(), string(requestBody.Params), "{\"account\":\""+s.RyskV2WSClient.addressString+"\"}")
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.SubAccountList("69420")
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_ApproveSigner() {
	done := make(chan struct{})
	nonce := time.Now().UnixMicro()
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
			require.Equal(s.T(), string(constants.WS_METHOD_APPROVE_REVOKE_SIGNER), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account        string `json:"account"`
				SubAccountId   int64  `json:"subAccountId"`
				ApprovedSigner string `json:"approvedSigner"`
				IsApproved     bool   `json:"isApproved"`
				Nonce          int64  `json:"nonce"`
				Signature      string `json:"signature"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), int64(1), params.SubAccountId)
			require.Equal(s.T(), s.Address, params.ApprovedSigner)
			require.True(s.T(), params.IsApproved)
			require.Equal(s.T(), nonce, params.Nonce)
			require.NotEmpty(s.T(), params.Signature)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.ApproveSigner("69420", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.Address,
		Nonce:          nonce,
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_ApproveSigner_BadAddress() {
	err := s.RyskV2WSClient.ApproveSigner("69420", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "",
		Nonce:          time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
}

func (s *WSClientUnitTestSuite) TestUnit_RevokeSigner() {
	done := make(chan struct{})
	nonce := time.Now().UnixMicro()
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
			require.Equal(s.T(), string(constants.WS_METHOD_APPROVE_REVOKE_SIGNER), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account        string `json:"account"`
				SubAccountId   int64  `json:"subAccountId"`
				ApprovedSigner string `json:"approvedSigner"`
				IsApproved     bool   `json:"isApproved"`
				Nonce          int64  `json:"nonce"`
				Signature      string `json:"signature"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), int64(1), params.SubAccountId)
			require.Equal(s.T(), s.Address, params.ApprovedSigner)
			require.False(s.T(), params.IsApproved)
			require.Equal(s.T(), nonce, params.Nonce)
			require.NotEmpty(s.T(), params.Signature)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.RevokeSigner("69420", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.Address,
		Nonce:          nonce,
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_RevokeSigner_BadAddress() {
	err := s.RyskV2WSClient.RevokeSigner("69420", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "",
		Nonce:          time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
}

func (s *WSClientUnitTestSuite) TestUnit_NewOrder() {
	done := make(chan struct{})
	nonce := time.Now().UnixMicro()
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
			require.Equal(s.T(), string(constants.WS_METHOD_NEW_ORDER), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account      string `json:"account"`
				SubAccountId int64  `json:"subAccountId"`
				ProductId    int64  `json:"productId"`
				IsBuy        bool   `json:"isBuy"`
				OrderType    int64  `json:"orderType"`
				TimeInForce  int64  `json:"timeInForce"`
				Expiration   int64  `json:"expiration"`
				Price        string `json:"price"`
				Quantity     string `json:"quantity"`
				Nonce        int64  `json:"nonce"`
				Signature    string `json:"signature"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), int64(1), params.SubAccountId)
			require.Equal(s.T(), constants.PRODUCT_ETH_PERP.Id, params.ProductId)
			require.True(s.T(), params.IsBuy)

			require.Equal(s.T(), int64(1), params.OrderType)
			require.Equal(s.T(), int64(1), params.TimeInForce)
			require.Equal(s.T(), int64(1627801200), params.Expiration)
			require.Equal(s.T(), "123", params.Price)
			require.Equal(s.T(), "456", params.Quantity)

			require.Equal(s.T(), nonce, params.Nonce)
			require.NotEmpty(s.T(), params.Signature)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.NewOrder("69420", &types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   types.OrderType(1),
		TimeInForce: types.TimeInForce(1),
		Expiration:  1627801200,
		Price:       "123",
		Quantity:    "456",
		Nonce:       nonce,
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_NewOrder_BadAddress() {
	s.RyskV2WSClient.addressString = ""
	nonce := time.Now().UnixMicro()
	err := s.RyskV2WSClient.NewOrder("69420", &types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   types.OrderType(1),
		TimeInForce: types.TimeInForce(1),
		Expiration:  1627801200,
		Price:       "123",
		Quantity:    "456",
		Nonce:       nonce,
	})
	require.Error(s.T(), err)
}

func (s *WSClientUnitTestSuite) TestUnit_ListOpenOrders() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		_, message, err := conn.ReadMessage()
		require.NoError(s.T(), err)

		var requestBody struct {
			JsonRPC string          `json:"jsonrpc"`
			Id      string          `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		err = json.Unmarshal(message, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "2.0", requestBody.JsonRPC)
		require.Equal(s.T(), "69420", requestBody.Id)
		require.Equal(s.T(), string(constants.WS_METHOD_ORDER_LIST), requestBody.Method)
		require.NotEmpty(s.T(), requestBody.Params)

		var params struct {
			Account      string   `json:"account"`
			SubAccountId int64    `json:"subAccountId"`
			ProductId    int64    `json:"productId"`
			OrderIds     []string `json:"orderIds"`
			StartTime    int64    `json:"startTime"`
			EndTime      int64    `json:"endTime"`
			Limit        int64    `json:"limit"`
		}
		err = json.Unmarshal(requestBody.Params, &params)
		require.NoError(s.T(), err)
		require.Equal(s.T(), s.Address, params.Account)
		require.Equal(s.T(), int64(1), params.SubAccountId)
		require.Equal(s.T(), constants.PRODUCT_ETH_PERP.Id, params.ProductId)
		require.Equal(s.T(), []string{"order1", "order2"}, params.OrderIds)
		require.Equal(s.T(), int64(1627801200), params.StartTime)
		require.Equal(s.T(), int64(1627801800), params.EndTime)
		require.Equal(s.T(), int64(50), params.Limit)
		done <- struct{}{}
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.ListOpenOrders("69420", &types.ListOrdersRequest{
		Product:   &constants.PRODUCT_ETH_PERP,
		Ids:       []string{"order1", "order2"},
		StartTime: 1627801200,
		EndTime:   1627801800,
		Limit:     50,
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_CancelOrder() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		_, message, err := conn.ReadMessage()
		require.NoError(s.T(), err)

		var requestBody struct {
			JsonRPC string          `json:"jsonrpc"`
			Id      string          `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		err = json.Unmarshal(message, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "2.0", requestBody.JsonRPC)
		require.Equal(s.T(), "69420", requestBody.Id)
		require.Equal(s.T(), string(constants.WS_METHOD_CANCEL_ORDER), requestBody.Method)
		require.NotEmpty(s.T(), requestBody.Params)

		var params struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
			OrderId      string `json:"orderId"`
			Signature    string `json:"signature"`
		}
		err = json.Unmarshal(requestBody.Params, &params)
		require.NoError(s.T(), err)
		require.Equal(s.T(), s.Address, params.Account)
		require.Equal(s.T(), int64(1), params.SubAccountId)
		require.Equal(s.T(), constants.PRODUCT_ETH_PERP.Id, params.ProductId)
		require.Equal(s.T(), "order123", params.OrderId)
		require.NotEmpty(s.T(), params.Signature)
		done <- struct{}{}
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.CancelOrder("69420", &types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "order123",
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_CancelOrder_BadAddress() {
	s.RyskV2WSClient.addressString = ""
	err := s.RyskV2WSClient.CancelOrder("69420", &types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "order123",
	})
	require.Error(s.T(), err)
}
func (s *WSClientUnitTestSuite) TestUnit_CancelAllOpenOrders() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		_, message, err := conn.ReadMessage()
		require.NoError(s.T(), err)

		var requestBody struct {
			JsonRPC string          `json:"jsonrpc"`
			Id      string          `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		err = json.Unmarshal(message, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "2.0", requestBody.JsonRPC)
		require.Equal(s.T(), "69420", requestBody.Id)
		require.Equal(s.T(), string(constants.WS_METHOD_CANCEL_ALL_OPEN_ORDERS), requestBody.Method)
		require.NotEmpty(s.T(), requestBody.Params)

		var params struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
		}
		err = json.Unmarshal(requestBody.Params, &params)
		require.NoError(s.T(), err)
		require.Equal(s.T(), s.Address, params.Account)
		require.Equal(s.T(), int64(1), params.SubAccountId)
		require.Equal(s.T(), constants.PRODUCT_ETH_PERP.Id, params.ProductId)
		done <- struct{}{}
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.CancelAllOpenOrders("69420", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_OrderBook() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		_, message, err := conn.ReadMessage()
		require.NoError(s.T(), err)

		var requestBody struct {
			JsonRPC string          `json:"jsonrpc"`
			Id      string          `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		err = json.Unmarshal(message, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "2.0", requestBody.JsonRPC)
		require.Equal(s.T(), "69420", requestBody.Id)
		require.Equal(s.T(), string(constants.WS_METHOD_ORDER_BOOK_DEPTH), requestBody.Method)
		require.NotEmpty(s.T(), requestBody.Params)

		var params struct {
			Symbol      string `json:"symbol"`
			Granularity int64  `json:"granularity"`
			Limit       int64  `json:"limit"`
		}
		err = json.Unmarshal(requestBody.Params, &params)
		require.NoError(s.T(), err)
		require.Equal(s.T(), constants.PRODUCT_ETH_PERP.Symbol, params.Symbol)
		require.Equal(s.T(), int64(1), params.Granularity)
		require.Equal(s.T(), int64(100), params.Limit)
		done <- struct{}{}
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.OrderBook("69420", &types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 1,
		Limit:       100,
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_GetPerpetualPosition() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_GET_PERPETUAL_POSITION), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account      string  `json:"account"`
				SubAccountId int64   `json:"subAccountId"`
				ProductIds   []int64 `json:"productIds"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), int64(1), params.SubAccountId)
			require.ElementsMatch(s.T(), []int64{constants.PRODUCT_ETH_PERP.Id, constants.PRODUCT_ETH_PERP.Id}, params.ProductIds)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.GetPerpetualPosition("69420", []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_ETH_PERP,
	})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_GetSpotBalances() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_GET_SPOT_BALANCES), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account      string   `json:"account"`
				SubAccountId int64    `json:"subAccountId"`
				Assets       []string `json:"assets"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), int64(1), params.SubAccountId)
			require.ElementsMatch(s.T(), []string{}, params.Assets)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.GetSpotBalances("69420", []string{})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_AccountUpdates() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_ACCOUNT_UPDATES), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params struct {
				Account      string `json:"account"`
				SubAccountId int64  `json:"subAccountId"`
			}
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Equal(s.T(), s.Address, params.Account)
			require.Equal(s.T(), int64(1), params.SubAccountId)
			done <- struct{}{}
			break
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
	s.RyskV2WSClient.RPCConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.AccountUpdates("69420")
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubscribeAggregateTrades_ZeroProducts() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 0)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.SubscribeAggregateTrades("69420", []*types.Product{})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubscribeAggregateTrades() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 2)

			expectedSymbol1 := constants.PRODUCT_ETH_PERP.Symbol + "@aggTrade"
			expectedSymbol2 := constants.PRODUCT_BTC_PERP.Symbol + "@aggTrade"
			require.Contains(s.T(), params, expectedSymbol1)
			require.Contains(s.T(), params, expectedSymbol2)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}
	err = s.RyskV2WSClient.SubscribeAggregateTrades("69420", products)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_UnsubscribeAggregateTrades_ZeroProducts() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 0)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.UnsubscribeAggregateTrades("69420", []*types.Product{})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_UnsubscribeAggregateTrades() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 2)

			expectedSymbol1 := constants.PRODUCT_ETH_PERP.Symbol + "@aggTrade"
			expectedSymbol2 := constants.PRODUCT_BTC_PERP.Symbol + "@aggTrade"
			require.Contains(s.T(), params, expectedSymbol1)
			require.Contains(s.T(), params, expectedSymbol2)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}
	err = s.RyskV2WSClient.UnsubscribeAggregateTrades("69420", products)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubscribeSingleTrades_ZeroProducts() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 0)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.SubscribeSingleTrades("69420", []*types.Product{})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubscribeSingleTrades() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 2)

			expectedSymbol1 := constants.PRODUCT_ETH_PERP.Symbol + "@trade"
			expectedSymbol2 := constants.PRODUCT_BTC_PERP.Symbol + "@trade"
			require.Contains(s.T(), params, expectedSymbol1)
			require.Contains(s.T(), params, expectedSymbol2)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}
	err = s.RyskV2WSClient.SubscribeSingleTrades("69420", products)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_UnsubscribeSingleTrades_ZeroProducts() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 0)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.UnubscribeSingleTrades("69420", []*types.Product{})
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_UnsubscribeSingleTrades() {
	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, 2)

			expectedSymbol1 := constants.PRODUCT_ETH_PERP.Symbol + "@trade"
			expectedSymbol2 := constants.PRODUCT_BTC_PERP.Symbol + "@trade"
			require.Contains(s.T(), params, expectedSymbol1)
			require.Contains(s.T(), params, expectedSymbol2)

			done <- struct{}{}
			break
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}
	err = s.RyskV2WSClient.UnubscribeSingleTrades("69420", products)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubscribeKlineData() {
	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}

	intervals := []types.Interval{
		constants.INTERVAL_1M,
		constants.INTERVAL_5M,
		constants.INTERVAL_1H,
	}

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, len(intervals)*len(products))

			for _, product := range products {
				for _, interval := range intervals {
					expectedSymbol := product.Symbol + "@klines_" + string(interval)
					require.Contains(s.T(), params, expectedSymbol)
				}
			}

			done <- struct{}{}
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.SubscribeKlineData("69420", products, intervals)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_UnsubscribeKlineData() {
	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}

	intervals := []types.Interval{
		constants.INTERVAL_1M,
		constants.INTERVAL_5M,
		constants.INTERVAL_1H,
	}

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, len(intervals)*len(products))

			for _, product := range products {
				for _, interval := range intervals {
					expectedSymbol := product.Symbol + "@klines_" + string(interval)
					require.Contains(s.T(), params, expectedSymbol)
				}
			}

			done <- struct{}{}
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.UnsubscribeKlineData("69420", products, intervals)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_SubscribePartialBookDepth() {
	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}

	limits := []types.Limit{
		constants.LIMIT_FIVE,
		constants.LIMIT_TEN,
		constants.LIMIT_TWENTY,
	}

	granularities := []int64{
		1,
		10,
		100,
	}

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, len(products)*len(limits)*len(granularities))

			for _, product := range products {
				for _, limit := range limits {
					for _, granularity := range granularities {
						expectedSymbol := product.Symbol + "@depth" + strconv.FormatInt(int64(limit), 10) + "_" + strconv.FormatInt(granularity, 10)
						require.Contains(s.T(), params, expectedSymbol)
					}
				}
			}

			done <- struct{}{}
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.SubscribePartialBookDepth("69420", products, limits, granularities)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_UnubscribePartialBookDepth() {
	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}

	limits := []types.Limit{
		constants.LIMIT_FIVE,
		constants.LIMIT_TEN,
		constants.LIMIT_TWENTY,
	}

	granularities := []int64{
		1,
		10,
		100,
	}

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, len(products)*len(limits)*len(granularities))

			for _, product := range products {
				for _, limit := range limits {
					for _, granularity := range granularities {
						expectedSymbol := product.Symbol + "@depth" + strconv.FormatInt(int64(limit), 10) + "_" + strconv.FormatInt(granularity, 10)
						require.Contains(s.T(), params, expectedSymbol)
					}
				}
			}

			done <- struct{}{}
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.UnsubscribePartialBookDepth("69420", products, limits, granularities)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_Subscribe24hrPriceChangeStatistics() {
	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, len(products))

			for _, product := range products {
				expectedSymbol := product.Symbol + "@ticker"
				require.Contains(s.T(), params, expectedSymbol)
			}

			done <- struct{}{}
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.Subscribe24hrPriceChangeStatistics("69420", products)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_Unsubscribe24hrPriceChangeStatistics() {
	products := []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
	}

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			require.NoError(s.T(), err)

			var requestBody struct {
				JsonRPC string          `json:"jsonrpc"`
				Id      string          `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}
			err = json.Unmarshal(message, &requestBody)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "2.0", requestBody.JsonRPC)
			require.Equal(s.T(), "69420", requestBody.Id)
			require.Equal(s.T(), string(constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE), requestBody.Method)
			require.NotEmpty(s.T(), requestBody.Params)

			var params []string
			err = json.Unmarshal(requestBody.Params, &params)
			require.NoError(s.T(), err)
			require.Len(s.T(), params, len(products))

			for _, product := range products {
				expectedSymbol := product.Symbol + "@ticker"
				require.Contains(s.T(), params, expectedSymbol)
			}

			done <- struct{}{}
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
	s.RyskV2WSClient.StreamConnection = rpcWebsocket
	s.RyskV2WSClient.rpcUrl = url

	err = s.RyskV2WSClient.Unsubscribe24hrPriceChangeStatistics("69420", products)
	require.NoError(s.T(), err)
	<-done
}

func (s *WSClientUnitTestSuite) TestUnit_ApproveUSDC() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)

	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_ApproveUSDC_ErrorGettingParameters() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), fmt.Errorf("error getting parameters"))

	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_ApproveUSDC_ErrorSendTransaction() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to send transaction"))

	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_DepositUSDC() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)

	transaction, err := s.RyskV2WSClient.DepositUSDC(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_DepositUSDC_ErrorApproveSigner() {
	s.RyskV2WSClient.addressString = ""
	transaction, err := s.RyskV2WSClient.DepositUSDC(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_DepositUSDC_ErrorGettingParameters() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), fmt.Errorf("error getting parameters"))

	transaction, err := s.RyskV2WSClient.DepositUSDC(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_DepositUSDC_ErrorSendTransaction() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to send transaction"))

	transaction, err := s.RyskV2WSClient.DepositUSDC(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *WSClientUnitTestSuite) TestUnit_WaitTransaction() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)
	mockEthClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(&geth_types.Receipt{}, nil)

	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.RyskV2WSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
}

func (s *WSClientUnitTestSuite) TestUnit_WaitTransaction_WaitMinedError() {
	mockEthClient := new(mocks.MockEthClient)
	s.RyskV2WSClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)
	mockEthClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return((*geth_types.Receipt)(nil), fmt.Errorf("failed to wait transaction"))

	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	receipt, err := s.RyskV2WSClient.WaitTransaction(ctx, transaction)
	require.Error(s.T(), err)
	require.Nil(s.T(), receipt)
}

func (s *WSClientUnitTestSuite) TestUnit_addReferee() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		var requestBody struct {
			Account   string `json:"account"`
			Code      string `json:"code"`
			Signature string `json:"signature"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodPost, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_ADD_REFEREE), req.URL.Path)
		require.Equal(s.T(), s.RyskV2WSClient.addressString, requestBody.Account)
		require.Equal(s.T(), "eldief", requestBody.Code)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.RyskV2WSClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.RyskV2WSClient.addReferee()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *WSClientUnitTestSuite) TestUnit_addReferee_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.RyskV2WSClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.RyskV2WSClient.addReferee()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *WSClientUnitTestSuite) TestUnit_addReferee_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.RyskV2WSClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.RyskV2WSClient.addReferee()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}
