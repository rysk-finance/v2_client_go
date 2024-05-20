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
		fmt.Println("Error loading .env file:", err)
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
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_Get24hrPriceChangeStatistics_WithNonExistingProduct() {
	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&types.Product{
		Id:     69420,
		Symbol: "69420",
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 404, res.StatusCode)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_Get24hrPriceChangeStatistics_WithProduct() {
	res, err := s.Go100XApiClient.Get24hrPriceChangeStatistics(&constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetProduct() {
	res, err := s.Go100XApiClient.GetProduct(constants.PRODUCT_ETH_PERP.Symbol)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetProductById() {
	res, err := s.Go100XApiClient.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetKlineData() {
	res, err := s.Go100XApiClient.GetKlineData(&types.KlineDataRequest{
		Product:   &constants.PRODUCT_BTC_PERP,
		Interval:  constants.INTERVAL_D1,
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListProducts() {
	res, err := s.Go100XApiClient.ListProducts()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_OrderBook() {
	res, err := s.Go100XApiClient.OrderBook(&types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 0,
		Limit:       constants.LIMIT_FIVE,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ServerTime() {
	res, err := s.Go100XApiClient.ServerTime()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ApproveSigner() {
	res, err := s.Go100XApiClient.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Nonce:          time.Now().UnixMilli(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_RevokeSigner() {
	res, err := s.Go100XApiClient.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Nonce:          time.Now().UnixMilli(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
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
	require.NoError(s.T(), err)
	require.Equal(s.T(), 400, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
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
	require.NoError(s.T(), err)
	require.Equal(s.T(), 404, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelOrder() {
	res, err := s.Go100XApiClient.CancelOrder(&types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "1",
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 404, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelAllOpenOrders() {
	res, err := s.Go100XApiClient.CancelAllOpenOrders(&constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetSpotBalances() {
	res, err := s.Go100XApiClient.GetSpotBalances()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_GetPerpetualPosition() {
	res, err := s.Go100XApiClient.GetPerpetualPosition(&constants.PRODUCT_BLAST_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListApproveSigners() {
	res, err := s.Go100XApiClient.ListApprovedSigners()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListOpenOrders() {
	res, err := s.Go100XApiClient.ListOpenOrders(&constants.PRODUCT_BLAST_PERP)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ListOrders() {
	res, err := s.Go100XApiClient.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BTC_PERP,
		Ids:     []string{"A", "B"},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func verifyValidJSONResponse(t *testing.T, response *http.Response) {
	// Read response
	defer response.Body.Close()
	bytesBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	// Check if res is valid JSON by trying to unmarshal it
	var data interface{}
	err = json.Unmarshal([]byte(bytesBody), &data)
	require.NoError(t, err)
}
