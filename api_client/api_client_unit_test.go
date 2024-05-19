package api_client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
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

func TestRun_ApiClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ApiClientUnitTestSuite))
}

func (s *ApiClientUnitTestSuite) TestUnit_Get24hrPriceChangeStatistics_WithBadRequest() {
	s.Go100XApiClient.baseUrl = "http://\t"

	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_BLAST_PERP)
	require.NotNil(s.T(), err)
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
	require.Nil(s.T(), err, "[TestUnit_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)
	require.Equal(s.T(), 200, res.StatusCode)
}
