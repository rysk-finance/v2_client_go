package api_client

import (
	"fmt"
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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/joho/godotenv"
)

type ApiClientUnitTestSuite struct {
	suite.Suite
	PrivateKeys     string
	RpcUrl          string
	Go100XApiClient *Go100XAPIClient
	MockHTTPServer  *httptest.Server
}

func (s *ApiClientUnitTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		return
	}

	s.Go100XApiClient = NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})
}

func TestRunSuiteUnit_ApiClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ApiClientUnitTestSuite))
}

func (s *ApiClientUnitTestSuite) TestUnit_NewGo100XAPIClient() {
	apiClient := NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})
	require.Equal(s.T(), constants.API_BASE_URL[constants.ENVIRONMENT_TESTNET], apiClient.baseUrl)
	require.Equal(s.T(), strings.TrimPrefix(string(os.Getenv("PRIVATE_KEYS")), "0x"), apiClient.privateKey)
	require.Equal(s.T(), utils.AddressFromPrivateKey(string(os.Getenv("PRIVATE_KEYS"))), apiClient.address)
	require.Equal(s.T(), int64(1), apiClient.SubAccountId)
	require.Equal(s.T(), constants.CIAO_ADDRESS[constants.ENVIRONMENT_TESTNET], apiClient.verifyingContract)
	require.NotNil(s.T(), apiClient.HttpClient)
	require.NotNil(s.T(), apiClient.domain)
	require.Equal(s.T(), constants.DOMAIN_NAME, apiClient.domain.Name)
	require.Equal(s.T(), constants.DOMAIN_VERSION, apiClient.domain.Version)
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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_BLAST_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_Get24hrPriceChangeStatistics_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

	res, err := s.Go100XApiClient.GetProduct(constants.PRODUCT_BLAST_PERP.Symbol)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetProduct_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

	res, err := s.Go100XApiClient.GetProductById(constants.PRODUCT_BLAST_PERP.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_GetProductById_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

	res, err := s.Go100XApiClient.ListProducts()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ListProducts_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

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
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = s.MockHTTPServer.URL
	defer s.MockHTTPServer.Close()

	res, err := s.Go100XApiClient.ServerTime()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientUnitTestSuite) TestUnit_ServerTime_BadBaseURL() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.MockHTTPServer = httptest.NewServer(http.HandlerFunc(handler))
	s.Go100XApiClient.baseUrl = "://invalid-url"
	defer s.MockHTTPServer.Close()

	res, err := s.Go100XApiClient.ServerTime()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}
