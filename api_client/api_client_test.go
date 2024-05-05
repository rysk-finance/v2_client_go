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

	"github.com/ethereum/go-ethereum/params"
	"github.com/joho/godotenv"
)

// .env variables to be setup via `TestMain` before running tests.
var (
	PRIVATE_KEYS      string = ""
	RPC_URL           string = ""
	GO100X_API_CLIENT *Go100XAPIClient
)

// Setup .env variables for test suite.
func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		os.Exit(1)
	}

	PRIVATE_KEYS = string(os.Getenv("PRIVATE_KEYS"))
	RPC_URL = os.Getenv("RPC_URL")
	GO100X_API_CLIENT = NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Get24hrPriceChangeStatistics_NoProduct(t *testing.T) {
	res, err := GO100X_API_CLIENT.Get24hrPriceChangeStatistics(&types.Product{})
	require.Nil(t, err, "[Test_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)

	verifyValidJSONResponse("Test_Get24hrPriceChangeStatistics_NoProduct", t, res)
}

func Test_Get24hrPriceChangeStatistics_WithNonExistingProduct(t *testing.T) {
	res, err := GO100X_API_CLIENT.Get24hrPriceChangeStatistics(&types.Product{
		Id:     69420,
		Symbol: "69420",
	})
	require.Nil(t, err, "[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Error: %v", err)

	require.Equal(t, 404, res.StatusCode, "[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Error code")
	require.Equal(t, "404 Not Found", res.Status, "[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Error code")
}

func Test_Get24hrPriceChangeStatistics_WithProduct(t *testing.T) {
	res, err := GO100X_API_CLIENT.Get24hrPriceChangeStatistics(&constants.PRODUCT_ETH_PERP)
	require.Nil(t, err, "[Test_Get24hrPriceChangeStatistics_WithProduct] Error: %v", err)

	verifyResponseStatusCode("Test_Get24hrPriceChangeStatistics_WithProduct", t, 200, res)
	verifyValidJSONResponse("Test_Get24hrPriceChangeStatistics_WithProduct", t, res)
}

func Test_GetProduct(t *testing.T) {
	res, err := GO100X_API_CLIENT.GetProduct(constants.PRODUCT_ETH_PERP.Symbol)
	require.Nil(t, err, "[Test_GetProduct] Error: %v", err)

	verifyResponseStatusCode("Test_GetProduct", t, 200, res)
	verifyValidJSONResponse("Test_GetProduct", t, res)
}

func Test_GetProductById(t *testing.T) {
	res, err := GO100X_API_CLIENT.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	require.Nil(t, err, "[Test_GetProductById] Error: %v", err)

	verifyResponseStatusCode("Test_GetProductById", t, 200, res)
	verifyValidJSONResponse("Test_GetProductById", t, res)
}

func Test_GetKlineData(t *testing.T) {
	res, err := GO100X_API_CLIENT.GetKlineData(&types.KlineDataRequest{
		Product:   &constants.PRODUCT_BTC_PERP,
		Interval:  constants.INTERVAL_D1,
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
	})
	require.Nil(t, err, "[Test_GetKlineData] Error: %v", err)

	verifyResponseStatusCode("Test_GetKlineData", t, 200, res)
	verifyValidJSONResponse("Test_GetKlineData", t, res)
}

func Test_ListProducts(t *testing.T) {
	res, err := GO100X_API_CLIENT.ListProducts()
	require.Nil(t, err, "[Test_ListProducts] Error: %v", err)

	verifyResponseStatusCode("Test_ListProducts", t, 200, res)
	verifyValidJSONResponse("Test_ListProducts", t, res)
}

func Test_OrderBook(t *testing.T) {
	res, err := GO100X_API_CLIENT.OrderBook(&types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 0,
		Limit:       constants.LIMIT_FIVE,
	})
	require.Nil(t, err, "[Test_OrderBook] Error: %v", err)

	verifyResponseStatusCode("Test_OrderBook", t, 200, res)
	verifyValidJSONResponse("Test_OrderBook", t, res)
}

func Test_ServerTime(t *testing.T) {
	res, err := GO100X_API_CLIENT.ServerTime()
	require.Nil(t, err, "[Test_ServerTime] Error: %v", err)

	verifyResponseStatusCode("Test_ServerTime", t, 200, res)
	verifyValidJSONResponse("Test_ServerTime", t, res)
}

func Test_ApproveSigner(t *testing.T) {
	res, err := GO100X_API_CLIENT.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})
	require.Nil(t, err, "[Test_ApproveSigner] Error: %v", err)

	verifyResponseStatusCode("Test_ApproveSigner", t, 200, res)
	verifyValidJSONResponse("Test_ApproveSigner", t, res)
}

func Test_RevokeSigner(t *testing.T) {
	res, err := GO100X_API_CLIENT.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})
	require.Nil(t, err, "[Test_RevokeSigner] Error: %v", err)

	verifyResponseStatusCode("Test_RevokeSigner", t, 500, res)
	verifyValidJSONResponse("Test_RevokeSigner", t, res)
}

func Test_NewOrder(t *testing.T) {
	// Limit buy 1 ETH for 3300 USDB, valid for 1 day
	res, err := GO100X_API_CLIENT.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       new(big.Int).Mul(big.NewInt(3300), big.NewInt(params.Ether)).String(),
		Quantity:    new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)).String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMilli(),
	})
	require.Nil(t, err, "[Test_NewOrder] Error: %v", err)

	verifyResponseStatusCode("Test_NewOrder", t, 400, res)
	verifyValidJSONResponse("Test_NewOrder", t, res)
}

func Test_CancelOrderAndReplace(t *testing.T) {
	res, err := GO100X_API_CLIENT.CancelOrderAndReplace(&types.CancelOrderAndReplaceRequest{
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
	require.Nil(t, err, "[Test_CancelOrderAndReplace] Error: %v", err)

	verifyResponseStatusCode("Test_CancelOrderAndReplace", t, 404, res)
	verifyValidJSONResponse("Test_CancelOrderAndReplace", t, res)
}

func Test_CancelOrder(t *testing.T) {
	res, err := GO100X_API_CLIENT.CancelOrder(&types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "1",
	})
	require.Nil(t, err, "[Test_CancelOrder] Error: %v", err)

	verifyResponseStatusCode("Test_CancelOrder", t, 404, res)
	verifyValidJSONResponse("Test_CancelOrder", t, res)
}

func Test_CancelAllOpenOrders(t *testing.T) {
	res, err := GO100X_API_CLIENT.CancelAllOpenOrders(&constants.PRODUCT_ETH_PERP)
	require.Nil(t, err, "[Test_CancelAllOpenOrders] Error: %v", err)

	verifyResponseStatusCode("Test_CancelAllOpenOrders", t, 200, res)
	verifyValidJSONResponse("Test_CancelAllOpenOrders", t, res)
}

func Test_GetSpotBalances(t *testing.T) {
	res, err := GO100X_API_CLIENT.GetSpotBalances()
	require.Nil(t, err, "[Test_GetSpotBalances] Error: %v", err)

	verifyResponseStatusCode("Test_GetSpotBalances", t, 200, res)
	verifyValidJSONResponse("Test_GetSpotBalances", t, res)
}

func Test_GetPerpetualPosition(t *testing.T) {
	res, err := GO100X_API_CLIENT.GetPerpetualPosition(&constants.PRODUCT_BLAST_PERP)
	require.Nil(t, err, "[Test_GetPerpetualPosition] Error: %v", err)

	verifyResponseStatusCode("Test_GetPerpetualPosition", t, 200, res)
	verifyValidJSONResponse("Test_GetPerpetualPosition", t, res)
}

func Test_ListApproveSigners(t *testing.T) {
	res, err := GO100X_API_CLIENT.ListApprovedSigners()
	require.Nil(t, err, "[Test_ListApproveSigners] Error: %v", err)

	verifyResponseStatusCode("Test_ListApproveSigners", t, 200, res)
	verifyValidJSONResponse("Test_ListApproveSigners", t, res)
}

func Test_ListOpenOrders(t *testing.T) {
	res, err := GO100X_API_CLIENT.ListOpenOrders(&constants.PRODUCT_BLAST_PERP)
	require.Nil(t, err, "[Test_ListOpenOrders] Error: %v", err)

	verifyResponseStatusCode("Test_ListOpenOrders", t, 200, res)
	verifyValidJSONResponse("Test_ListOpenOrders", t, res)
}

func Test_ListOrders(t *testing.T) {
	res, err := GO100X_API_CLIENT.ListOrders(&types.ListOrdersRequest{
		Product: &constants.PRODUCT_BTC_PERP,
		Ids:     []string{"A", "B"},
	})
	require.Nil(t, err, "[Test_ListOrders] Error: %v", err)

	verifyResponseStatusCode("Test_ListOrders", t, 200, res)
	verifyValidJSONResponse("Test_ListOrders", t, res)
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
