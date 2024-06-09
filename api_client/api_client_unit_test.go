package api_client

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

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/eldief/go100x/utils"
	"github.com/eldief/go100x/utils/mocks"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/joho/godotenv"
)

type ApiClientUnitTestSuite struct {
	suite.Suite
	PrivateKey      string
	Address         string
	BaseUrl         string
	RpcUrl          string
	Go100XApiClient *Go100XAPIClient
	EthClient       types.IEthClient
}

func (s *ApiClientUnitTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("ApiClientUnitTestSuite.SetupSuite: Error loading .env file:", err)
		return
	}
	apiClient, _ := NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	s.Go100XApiClient = apiClient
	s.PrivateKey = s.Go100XApiClient.privateKeyString
	s.Address = utils.AddressFromPrivateKey(s.Go100XApiClient.privateKeyString)
	s.BaseUrl = s.Go100XApiClient.baseUrl
	s.EthClient = apiClient.EthClient
}

func (s *ApiClientUnitTestSuite) SetupTest() {
	s.Go100XApiClient.privateKeyString = s.PrivateKey
	s.Go100XApiClient.addressString = s.Address
	s.Go100XApiClient.baseUrl = s.BaseUrl
	s.Go100XApiClient.EthClient = s.EthClient
}

func TestRunSuiteUnit_ApiClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ApiClientUnitTestSuite))
}

func (s *ApiClientUnitTestSuite) TestUnit_NewGo100XAPIClient() {
	apiClient, err := NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), constants.ENVIRONMENT_TESTNET, apiClient.env)
	require.Equal(s.T(), constants.API_BASE_URL[constants.ENVIRONMENT_TESTNET], apiClient.baseUrl)
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), apiClient.privateKeyString)
	require.Equal(s.T(), utils.AddressFromPrivateKey(strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x")), apiClient.addressString)
	require.Equal(s.T(), utils.AddressFromPrivateKey(strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x")), apiClient.address.String())
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), hex.EncodeToString(crypto.FromECDSA(apiClient.privateKey)))
	require.Equal(s.T(), constants.CIAO_ADDRESS[constants.ENVIRONMENT_TESTNET], apiClient.ciao.String())
	require.Equal(s.T(), constants.USDB_ADDRESS[constants.ENVIRONMENT_TESTNET], apiClient.usdb.String())
	require.NotNil(s.T(), apiClient.domain)
	require.Equal(s.T(), constants.DOMAIN_NAME, apiClient.domain.Name)
	require.Equal(s.T(), constants.DOMAIN_VERSION, apiClient.domain.Version)
	require.Equal(s.T(), int64(1), apiClient.SubAccountId)
	require.NotNil(s.T(), apiClient.HttpClient)
	require.NotNil(s.T(), apiClient.EthClient)
}

func (s *ApiClientUnitTestSuite) TestUnit_NewGo100XAPIClient_InvalidPrivateKey() {
	apiClient, err := NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   "0x123",
		RpcUrl:       "",
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *ApiClientUnitTestSuite) TestUnit_NewGo100XAPIClient_InvalidRPCURL() {
	apiClient, err := NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       "",
		SubAccountId: 1,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), apiClient)
}

func (s *ApiClientUnitTestSuite) TestUnit_Get24hrPriceChangeStatistics_WithBadRequest() {
	s.Go100XApiClient.baseUrl = "http://\t"

	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_BLAST_PERP)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_Get24hrPriceChangeStatistics_NoProduct() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_24H_TICKER_PRICE_CHANGE_STATISTICS),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&types.Product{})
	require.Nil(s.T(), err, "[TestUnit_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_Get24hrPriceChangeStatistics_WithProduct() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_24H_TICKER_PRICE_CHANGE_STATISTICS),
			req.URL.Path,
		)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), http.MethodGet, req.Method)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_BLAST_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_Get24hrPriceChangeStatistics_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_BLAST_PERP)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetProduct() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_PRODUCT)+constants.PRODUCT_BLAST_PERP.Symbol,
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetProduct(constants.PRODUCT_BLAST_PERP.Symbol)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetProduct_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetProduct(constants.PRODUCT_BLAST_PERP.Symbol)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetProductById() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_PRODUCT_BY_ID)+strconv.FormatInt(constants.PRODUCT_BLAST_PERP.Id, 10),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetProductById(constants.PRODUCT_BLAST_PERP.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetProductById_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetProductById(constants.PRODUCT_BLAST_PERP.Id)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetKlineData() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_KLINE_DATA),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetKlineData_WithInterval() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_KLINE_DATA),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), string(constants.INTERVAL_15M), req.URL.Query().Get("interval"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product:  &constants.PRODUCT_BLAST_PERP,
		Interval: constants.INTERVAL_15M,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetKlineData_WithStartTime() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_KLINE_DATA),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), strconv.FormatInt(123, 10), req.URL.Query().Get("startTime"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product:   &constants.PRODUCT_BLAST_PERP,
		StartTime: 123,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetKlineData_WithEndTime() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_KLINE_DATA),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), strconv.FormatInt(123, 10), req.URL.Query().Get("endTime"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		EndTime: 123,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetKlineData_WithLimit() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_GET_KLINE_DATA),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), strconv.FormatInt(123, 10), req.URL.Query().Get("limit"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		Limit:   123,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetKlineData_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListProducts() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_LIST_PRODUCTS),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListProducts()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListProducts_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListProducts()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_OrderBook() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_ORDER_BOOK),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.OrderBook(&types.OrderBookRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_OrderBook_WithGranularity() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_ORDER_BOOK),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), strconv.FormatInt(1, 10), req.URL.Query().Get("granularity"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.OrderBook(&types.OrderBookRequest{
		Product:     &constants.PRODUCT_BLAST_PERP,
		Granularity: 1,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_OrderBook_WithLimit() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_ORDER_BOOK),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), strconv.FormatInt(5, 10), req.URL.Query().Get("limit"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.OrderBook(&types.OrderBookRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		Limit:   constants.LIMIT_FIVE,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_OrderBook_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.OrderBook(&types.OrderBookRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ServerTime() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(
			s.T(),
			string(constants.API_ENDPOINT_SERVER_TIME),
			req.URL.Path,
		)
		require.Equal(s.T(), http.MethodGet, req.Method)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ServerTime()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ServerTime_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ServerTime()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ApproveSigner() {
	nonce := time.Now().UnixMilli()
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		// Unmarshal the body into a struct for easier comparison
		var requestBody struct {
			Account        string `json:"account"`
			SubAccountId   int64  `json:"subAccountId"`
			Signature      string `json:"signature"`
			ApprovedSigner string `json:"approvedSigner"`
			Nonce          int64  `json:"nonce"`
			IsApproved     bool   `json:"isApproved"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodPost, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_APPROVE_REVOKE_SIGNER), req.URL.Path)
		require.Equal(s.T(), "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", requestBody.ApprovedSigner)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.Account)
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), strconv.FormatInt(requestBody.SubAccountId, 10))
		require.Equal(s.T(), nonce, requestBody.Nonce)
		require.True(s.T(), requestBody.IsApproved)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Nonce:          nonce,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ApproveSigner_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "",
		Nonce:          time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ApproveSigner_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Nonce:          time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_Revokeigner() {
	nonce := time.Now().UnixMicro()
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		// Unmarshal the body into a struct for easier comparison
		var requestBody struct {
			Account        string `json:"account"`
			SubAccountId   int64  `json:"subAccountId"`
			Signature      string `json:"signature"`
			ApprovedSigner string `json:"approvedSigner"`
			Nonce          int64  `json:"nonce"`
			IsApproved     bool   `json:"isApproved"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodPost, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_APPROVE_REVOKE_SIGNER), req.URL.Path)
		require.Equal(s.T(), "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", requestBody.ApprovedSigner)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.Account)
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), strconv.FormatInt(requestBody.SubAccountId, 10))
		require.Equal(s.T(), nonce, requestBody.Nonce)
		require.False(s.T(), requestBody.IsApproved)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Nonce:          nonce,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_RevokeSigner_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "",
		Nonce:          time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_RevokeSigner_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Nonce:          time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_Withdraw() {
	nonce := time.Now().UnixMilli()
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		var requestBody struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			Asset        string `json:"asset"`
			Quantity     string `json:"quantity"`
			Nonce        int64  `json:"nonce"`
			Signature    string `json:"signature"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodPost, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_WITHDRAW), req.URL.Path)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.Account)
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), strconv.FormatInt(requestBody.SubAccountId, 10))
		require.Equal(s.T(), constants.USDB_ADDRESS[constants.ENVIRONMENT_TESTNET], requestBody.Asset)
		require.Equal(s.T(), "456", requestBody.Quantity)
		require.Equal(s.T(), nonce, requestBody.Nonce)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.Withdraw(&types.WithdrawRequest{
		Quantity: "456",
		Nonce:    nonce,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_Withdraw_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.Withdraw(&types.WithdrawRequest{
		Quantity: "456",
		Nonce:    time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_Withdraw_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.Withdraw(&types.WithdrawRequest{
		Quantity: "456",
		Nonce:    time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_NewOrder() {
	nonce := time.Now().UnixMicro()
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		var requestBody struct {
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
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodPost, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_NEW_ORDER), req.URL.Path)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.Account)
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), strconv.FormatInt(requestBody.SubAccountId, 10))
		require.Equal(s.T(), int64(1006), requestBody.ProductId)
		require.True(s.T(), requestBody.IsBuy)
		require.Equal(s.T(), int64(1), requestBody.OrderType)
		require.Equal(s.T(), int64(1), requestBody.TimeInForce)
		require.Equal(s.T(), int64(1627801200), requestBody.Expiration)
		require.Equal(s.T(), "123", requestBody.Price)
		require.Equal(s.T(), "456", requestBody.Quantity)
		require.Equal(s.T(), nonce, requestBody.Nonce)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_BLAST_PERP,
		IsBuy:       true,
		OrderType:   types.OrderType(1),
		TimeInForce: types.TimeInForce(1),
		Expiration:  1627801200,
		Price:       "123",
		Quantity:    "456",
		Nonce:       nonce,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_NewOrder_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_BLAST_PERP,
		IsBuy:       true,
		OrderType:   types.OrderType(1),
		TimeInForce: types.TimeInForce(1),
		Expiration:  1627801200,
		Price:       "123",
		Quantity:    "456",
		Nonce:       time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_NewOrder_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_BLAST_PERP,
		IsBuy:       true,
		OrderType:   types.OrderType(1),
		TimeInForce: types.TimeInForce(1),
		Expiration:  1627801200,
		Price:       "123",
		Quantity:    "456",
		Nonce:       time.Now().UnixMicro(),
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelOrderAndReplace() {
	nonce := time.Now().UnixMicro()
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		var requestBody struct {
			IdToCancel string `json:"idToCancel"`
			NewOrder   struct {
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
			} `json:"newOrder"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodPost, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_CANCEL_REPLACE_ORDER), req.URL.Path)
		require.Equal(s.T(), "order123", requestBody.IdToCancel)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.NewOrder.Account)
		require.Equal(s.T(), s.Go100XApiClient.SubAccountId, requestBody.NewOrder.SubAccountId)
		require.Equal(s.T(), int64(1006), requestBody.NewOrder.ProductId)
		require.True(s.T(), requestBody.NewOrder.IsBuy)
		require.Equal(s.T(), int64(1), requestBody.NewOrder.OrderType)
		require.Equal(s.T(), int64(1), requestBody.NewOrder.TimeInForce)
		require.Equal(s.T(), int64(1627801200), requestBody.NewOrder.Expiration)
		require.Equal(s.T(), "123", requestBody.NewOrder.Price)
		require.Equal(s.T(), "456", requestBody.NewOrder.Quantity)
		require.Equal(s.T(), nonce, requestBody.NewOrder.Nonce)
		require.NotEmpty(s.T(), requestBody.NewOrder.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.CancelOrderAndReplace(&types.CancelOrderAndReplaceRequest{
		IdToCancel: "order123",
		NewOrder: &types.NewOrderRequest{
			Product:     &types.Product{Id: 1006},
			IsBuy:       true,
			OrderType:   types.OrderType(1),
			TimeInForce: types.TimeInForce(1),
			Expiration:  1627801200,
			Price:       "123",
			Quantity:    "456",
			Nonce:       nonce,
		},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelOrderAndReplace_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.CancelOrderAndReplace(&types.CancelOrderAndReplaceRequest{
		IdToCancel: "order123",
		NewOrder: &types.NewOrderRequest{
			Product:     &types.Product{Id: 1006},
			IsBuy:       true,
			OrderType:   types.OrderType(1),
			TimeInForce: types.TimeInForce(1),
			Expiration:  1627801200,
			Price:       "123",
			Quantity:    "456",
			Nonce:       time.Now().UnixMicro(),
		},
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelOrderAndReplace_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.CancelOrderAndReplace(&types.CancelOrderAndReplaceRequest{
		IdToCancel: "order123",
		NewOrder: &types.NewOrderRequest{
			Product:     &types.Product{Id: 1006},
			IsBuy:       true,
			OrderType:   types.OrderType(1),
			TimeInForce: types.TimeInForce(1),
			Expiration:  1627801200,
			Price:       "123",
			Quantity:    "456",
			Nonce:       time.Now().UnixMicro(),
		},
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelOrder() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		var requestBody struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
			OrderId      string `json:"orderId"`
			Signature    string `json:"signature"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodDelete, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_CANCEL_ORDER), req.URL.Path)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.Account)
		require.Equal(s.T(), s.Go100XApiClient.SubAccountId, requestBody.SubAccountId)
		require.Equal(s.T(), int64(1006), requestBody.ProductId)
		require.Equal(s.T(), "order123", requestBody.OrderId)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.CancelOrder(&types.CancelOrderRequest{
		Product:    &types.Product{Id: 1006},
		IdToCancel: "order123",
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelOrder_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.CancelOrder(&types.CancelOrderRequest{
		Product:    &types.Product{Id: 1006},
		IdToCancel: "order123",
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelOrder_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.CancelOrder(&types.CancelOrderRequest{
		Product:    &types.Product{Id: 1006},
		IdToCancel: "order123",
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelAllOpenOrders() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(s.T(), err)
		defer req.Body.Close()

		var requestBody struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
			Signature    string `json:"signature"`
		}
		err = json.Unmarshal(body, &requestBody)
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.MethodDelete, req.Method)
		require.Equal(s.T(), string(constants.API_ENDPOINT_CANCEL_ALL_OPEN_ORDERS), req.URL.Path)
		require.Equal(s.T(), s.Go100XApiClient.addressString, requestBody.Account)
		require.Equal(s.T(), s.Go100XApiClient.SubAccountId, requestBody.SubAccountId)
		require.Equal(s.T(), int64(1006), requestBody.ProductId)
		require.NotEmpty(s.T(), requestBody.Signature)
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	product := &types.Product{Id: 1006}
	res, err := s.Go100XApiClient.CancelAllOpenOrders(product)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelAllOpenOrders_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	product := &types.Product{Id: 1006}
	res, err := s.Go100XApiClient.CancelAllOpenOrders(product)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_CancelAllOpenOrders_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	product := &types.Product{Id: 1006}
	res, err := s.Go100XApiClient.CancelAllOpenOrders(product)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetSpotBalances() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(s.T(), string(constants.API_ENDPOINT_GET_SPOT_BALANCES), req.URL.Path)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), s.Go100XApiClient.addressString, req.URL.Query().Get("account"))
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), req.URL.Query().Get("subAccountId"))
		require.NotEmpty(s.T(), req.URL.Query().Get("signature"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetSpotBalances()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetSpotBalances_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetSpotBalances()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetSpotBalances_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetSpotBalances()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetPerpetualPosition() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(s.T(), string(constants.API_ENDPOINT_GET_PERPETUAL_POSITION), req.URL.Path)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), s.Go100XApiClient.addressString, req.URL.Query().Get("account"))
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), req.URL.Query().Get("subAccountId"))
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.NotEmpty(s.T(), req.URL.Query().Get("signature"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetPerpetualPosition(&constants.PRODUCT_BLAST_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetPerpetualPosition_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetPerpetualPosition(&constants.PRODUCT_BLAST_PERP)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetPerpetualPosition_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.GetPerpetualPosition(&constants.PRODUCT_BLAST_PERP)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListApprovedSigners() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(s.T(), string(constants.API_ENDPOINT_LIST_APPROVED_SIGNERS), req.URL.Path)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), s.Go100XApiClient.addressString, req.URL.Query().Get("account"))
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), req.URL.Query().Get("subAccountId"))
		require.NotEmpty(s.T(), req.URL.Query().Get("signature"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListApprovedSigners()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListApprovedSigners_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListApprovedSigners()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListApprovedSigners_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListApprovedSigners()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOpenOrders() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(s.T(), string(constants.API_ENDPOINT_LIST_OPEN_ORDERS), req.URL.Path)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), s.Go100XApiClient.addressString, req.URL.Query().Get("account"))
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), req.URL.Query().Get("subAccountId"))
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.NotEmpty(s.T(), req.URL.Query().Get("signature"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOpenOrders(&constants.PRODUCT_BLAST_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOpenOrders_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOpenOrders(&constants.PRODUCT_BLAST_PERP)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOpenOrders_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOpenOrders(&constants.PRODUCT_BLAST_PERP)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOrders() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(s.T(), string(constants.API_ENDPOINT_LIST_ORDERS), req.URL.Path)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), s.Go100XApiClient.addressString, req.URL.Query().Get("account"))
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), req.URL.Query().Get("subAccountId"))
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Equal(s.T(), "123", req.URL.Query().Get("ids"))
		require.NotEmpty(s.T(), req.URL.Query().Get("signature"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		Ids:     []string{"123"},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOrders_MultipleIds() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(s.T(), string(constants.API_ENDPOINT_LIST_ORDERS), req.URL.Path)
		require.Equal(s.T(), http.MethodGet, req.Method)
		require.Equal(s.T(), s.Go100XApiClient.addressString, req.URL.Query().Get("account"))
		require.Equal(s.T(), strconv.FormatInt(s.Go100XApiClient.SubAccountId, 10), req.URL.Query().Get("subAccountId"))
		require.Equal(s.T(), constants.PRODUCT_BLAST_PERP.Symbol, req.URL.Query().Get("symbol"))
		require.Contains(s.T(), req.URL.Query()["ids"], "123")
		require.Contains(s.T(), req.URL.Query()["ids"], "456")
		require.NotEmpty(s.T(), req.URL.Query().Get("signature"))
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = mockHttpServer.URL
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		Ids:     []string{"123", "456"},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOrders_BadAddress() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.addressString = ""
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		Ids:     []string{"123"},
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListOrders_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mockHttpServer := httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer mockHttpServer.Close()

	res, err := s.Go100XApiClient.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BLAST_PERP,
		Ids:     []string{"123"},
	})
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *ApiClientUnitTestSuite) TestUnit_ApproveUSDB() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)

	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)
}

func (s *ApiClientUnitTestSuite) TestUnit_ApproveUSDB_ErrorGettingParameters() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), fmt.Errorf("error getting parameters"))

	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *ApiClientUnitTestSuite) TestUnit_ApproveUSDB_ErrorSendTransaction() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to send transaction"))

	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *ApiClientUnitTestSuite) TestUnit_DepositUSDB() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)

	transaction, err := s.Go100XApiClient.DepositUSDB(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)
}

func (s *ApiClientUnitTestSuite) TestUnit_DepositUSDB_ErrorGettingParameters() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), fmt.Errorf("error getting parameters"))

	transaction, err := s.Go100XApiClient.DepositUSDB(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *ApiClientUnitTestSuite) TestUnit_DepositUSDB_ErrorSendTransaction() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to send transaction"))

	transaction, err := s.Go100XApiClient.DepositUSDB(context.Background(), big.NewInt(1000))
	require.Error(s.T(), err)
	require.Nil(s.T(), transaction)
}

func (s *ApiClientUnitTestSuite) TestUnit_WaitTransaction() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)
	mockEthClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(&geth_types.Receipt{}, nil)

	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.Go100XApiClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
}

func (s *ApiClientUnitTestSuite) TestUnit_WaitTransaction_WaitMinedError() {
	mockEthClient := new(mocks.MockEthClient)
	s.Go100XApiClient.EthClient = mockEthClient
	mockEthClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(1), nil)
	mockEthClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockEthClient.On("NetworkID", mock.Anything).Return(big.NewInt(1), nil)
	mockEthClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(21000), nil)
	mockEthClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)
	mockEthClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return((*geth_types.Receipt)(nil), fmt.Errorf("failed to wait transaction"))

	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), big.NewInt(1000))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	receipt, err := s.Go100XApiClient.WaitTransaction(ctx, transaction)
	require.Error(s.T(), err)
	require.Nil(s.T(), receipt)
}
