package api_client

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/params"
	"github.com/joho/godotenv"
)

type ApiClientIntegrationTestSuite struct {
	suite.Suite
	PrivateKeys     string
	RpcUrl          string
	Go100XApiClient *Go100XAPIClient
}

func (s *ApiClientIntegrationTestSuite) SetupSuite() {
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

func TestRunSuiteIntegration_ApiClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ApiClientIntegrationTestSuite))
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_Get24hrPriceChangeStatistics_NoProduct() {
	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&types.Product{})
	require.Nil(s.T(), err, "[TestIntegration_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_Get24hrPriceChangeStatistics_NoProduct", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_Get24hrPriceChangeStatistics_NoProduct", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_Get24hrPriceChangeStatistics_WithNonExistingProduct() {
	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&types.Product{
		Id:     69420,
		Symbol: "69420",
	})
	require.Nil(s.T(), err, "[TestIntegration_Get24hrPriceChangeStatistics_WithNonExistingProduct] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_Get24hrPriceChangeStatistics_WithNonExistingProduct", s.T(), 404, res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_Get24hrPriceChangeStatistics_WithProduct() {
	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_ETH_PERP)
	require.Nil(s.T(), err, "[TestIntegration_Get24hrPriceChangeStatistics_WithProduct] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_Get24hrPriceChangeStatistics_WithProduct", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_Get24hrPriceChangeStatistics_WithProduct", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetProduct() {
	res, err := s.Go100XApiClient.GetProduct(constants.PRODUCT_ETH_PERP.Symbol)
	require.Nil(s.T(), err, "[TestIntegration_GetProduct] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_GetProduct", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_GetProduct", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetProductById() {
	res, err := s.Go100XApiClient.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	require.Nil(s.T(), err, "[TestIntegration_GetProductById] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_GetProductById", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_GetProductById", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetKlineData() {
	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product:   &constants.PRODUCT_BTC_PERP,
		Interval:  constants.INTERVAL_D1,
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
	})
	require.Nil(s.T(), err, "[TestIntegration_GetKlineData] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_GetKlineData", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_GetKlineData", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListProducts() {
	res, err := s.Go100XApiClient.ListProducts()
	require.Nil(s.T(), err, "[TestIntegration_ListProducts] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_ListProducts", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_ListProducts", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_OrderBook() {
	res, err := s.Go100XApiClient.OrderBook(&types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 0,
		Limit:       constants.LIMIT_FIVE,
	})
	require.Nil(s.T(), err, "[TestIntegration_OrderBook] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_OrderBook", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_OrderBook", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ServerTime() {
	res, err := s.Go100XApiClient.ServerTime()
	require.Nil(s.T(), err, "[TestIntegration_ServerTime] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_ServerTime", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_ServerTime", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ApproveSigner() {
	res, err := s.Go100XApiClient.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})
	require.Nil(s.T(), err, "[TestIntegration_ApproveSigner] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_ApproveSigner", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_ApproveSigner", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_RevokeSigner() {
	res, err := s.Go100XApiClient.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})
	require.Nil(s.T(), err, "[TestIntegration_RevokeSigner] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_RevokeSigner", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_RevokeSigner", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_NewOrder() {
	// Limit buy 1 ETH for 3300 USDB, valid for 1 day
	res, err := s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       new(big.Int).Mul(big.NewInt(3300), big.NewInt(params.Ether)).String(),
		Quantity:    new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)).String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMilli(),
	})
	require.Nil(s.T(), err, "[TestIntegration_NewOrder] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_NewOrder", s.T(), 400, res)
	verifyValidJSONResponse("TestIntegration_NewOrder", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelOrderAndReplace() {
	res, err := s.Go100XApiClient.CancelOrderAndReplace(&types.CancelOrderAndReplaceRequest{
		IdToCancel: "1",
		// Limit buy 1 ETH for 3300 USDB, valid for 1 day
		NewOrder: &types.NewOrderRequest{
			Product:     &constants.PRODUCT_ETH_PERP,
			IsBuy:       true,
			OrderType:   constants.ORDER_TYPE_LIMIT_MAKER,
			TimeInForce: constants.TIME_IN_FORCE_GTC,
			Price:       new(big.Int).Mul(big.NewInt(3300), big.NewInt(params.Ether)).String(),
			Quantity:    new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)).String(),
			Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
			Nonce:       time.Now().UnixMilli(),
		},
	})
	require.Nil(s.T(), err, "[TestIntegration_CancelOrderAndReplace] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_CancelOrderAndReplace", s.T(), 404, res)
	verifyValidJSONResponse("TestIntegration_CancelOrderAndReplace", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelOrder() {
	res, err := s.Go100XApiClient.CancelOrder(&types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "1",
	})
	require.Nil(s.T(), err, "[TestIntegration_CancelOrder] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_CancelOrder", s.T(), 404, res)
	verifyValidJSONResponse("TestIntegration_CancelOrder", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelAllOpenOrders() {
	res, err := s.Go100XApiClient.CancelAllOpenOrders(&constants.PRODUCT_ETH_PERP)
	require.Nil(s.T(), err, "[TestIntegration_CancelAllOpenOrders] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_CancelAllOpenOrders", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_CancelAllOpenOrders", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetSpotBalances() {
	res, err := s.Go100XApiClient.GetSpotBalances()
	require.Nil(s.T(), err, "[TestIntegration_GetSpotBalances] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_GetSpotBalances", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_GetSpotBalances", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetPerpetualPosition() {
	res, err := s.Go100XApiClient.GetPerpetualPosition(&constants.PRODUCT_BLAST_PERP)
	require.Nil(s.T(), err, "[TestIntegration_GetPerpetualPosition] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_GetPerpetualPosition", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_GetPerpetualPosition", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListApproveSigners() {
	res, err := s.Go100XApiClient.ListApprovedSigners()
	require.Nil(s.T(), err, "[TestIntegration_ListApproveSigners] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_ListApproveSigners", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_ListApproveSigners", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListOpenOrders() {
	res, err := s.Go100XApiClient.ListOpenOrders(&constants.PRODUCT_BLAST_PERP)
	require.Nil(s.T(), err, "[TestIntegration_ListOpenOrders] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_ListOpenOrders", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_ListOpenOrders", s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListOrders() {
	res, err := s.Go100XApiClient.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BTC_PERP,
		Ids:     []string{"A", "B"},
	})
	require.Nil(s.T(), err, "[TestIntegration_ListOrders] Error: %v", err)
	verifyResponseStatusCode("TestIntegration_ListOrders", s.T(), 200, res)
	verifyValidJSONResponse("TestIntegration_ListOrders", s.T(), res)
}

func verifyResponseStatusCode(testName string, t *testing.T, expectedStatusCode int, response *http.Response) {
	require.Equal(t, expectedStatusCode, response.StatusCode, "[%s] Expected status code %v but got %v", testName, expectedStatusCode, response.StatusCode)
}

func verifyValidJSONResponse(testName string, t *testing.T, response *http.Response) {
	// Read response
	defer response.Body.Close()
	bytesBody, err := io.ReadAll(response.Body)
	require.Nil(t, err, "[%s] Error reading response body: %v", testName, err)

	// Check if res is valid JSON by trying to unmarshal it
	var data interface{}
	err = json.Unmarshal([]byte(bytesBody), &data)
	require.Nil(t, err, "[%s] Error unmarshalling response: %v", testName, bytesBody)
	t.Logf("[%s] Response is a valid JSON: %s", testName, bytesBody)
}
