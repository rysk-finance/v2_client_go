package go100x_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	go100x "github.com/eldief/go100x/src/client"
	"github.com/eldief/go100x/src/constants"
	"github.com/eldief/go100x/src/types"

	"github.com/ethereum/go-ethereum/params"
	"github.com/joho/godotenv"
)

// .env variables to be setup via `TestMain` before running tests.
var (
	PRIVATE_KEYS string = ""
	RPC_URL      string = ""
	CLIENT_100_X *types.Client
)

// Setup .env variables for test suite.
func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		os.Exit(1)
	}

	PRIVATE_KEYS = string(os.Getenv("PRIVATE_KEYS"))
	RPC_URL = os.Getenv("RPC_URL")
	CLIENT_100_X = go100x.NewClient(&types.ClientConfiguration{
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
	res, err := go100x.Get24hrPriceChangeStatistics(CLIENT_100_X, &types.Product{})

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_Get24hrPriceChangeStatistics_NoProduct", t, res)
}

func Test_Get24hrPriceChangeStatistics_WithNonExistingProduct(t *testing.T) {
	_, err := go100x.Get24hrPriceChangeStatistics(CLIENT_100_X, &types.Product{
		Id:     69420,
		Symbol: "69420",
	})

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Error: %v", err)
		return
	}

	// expectedResponse := "Product not found"
	// if res != expectedResponse {
	// 	t.Errorf("[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Unexpected response. Got: %s, Expected: %s", res, expectedResponse)
	// }
}

func Test_Get24hrPriceChangeStatistics_WithProduct(t *testing.T) {
	res, err := go100x.Get24hrPriceChangeStatistics(CLIENT_100_X, &constants.PRODUCT_ETH_PERP)

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_WithProduct] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_Get24hrPriceChangeStatistics_WithProduct", t, 200, res)
	verifyValidJSONResponse("Test_Get24hrPriceChangeStatistics_WithProduct", t, res)
}

func Test_GetProduct(t *testing.T) {
	res, err := go100x.GetProduct(CLIENT_100_X, constants.PRODUCT_ETH_PERP.Symbol)

	if err != nil {
		t.Errorf("[Test_GetProduct] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_GetProduct", t, 200, res)
	verifyValidJSONResponse("Test_GetProduct", t, res)
}

func Test_GetProductById(t *testing.T) {
	res, err := go100x.GetProductById(CLIENT_100_X, constants.PRODUCT_ETH_PERP.Id)

	if err != nil {
		t.Errorf("[Test_GetProductById] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_GetProductById", t, 200, res)
	verifyValidJSONResponse("Test_GetProductById", t, res)
}

func Test_GetKlineData(t *testing.T) {
	res, err := go100x.GetKlineData(CLIENT_100_X, &types.KlineDataRequest{
		Product:   &constants.PRODUCT_BTC_PERP,
		Interval:  constants.INTERVAL_D1,
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_GetKlineData] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_GetKlineData", t, 200, res)
	verifyValidJSONResponse("Test_GetKlineData", t, res)
}

func Test_ListProducts(t *testing.T) {
	res, err := go100x.ListProducts(CLIENT_100_X)

	if err != nil {
		t.Errorf("[Test_ListProducts] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_ListProducts", t, 200, res)
	verifyValidJSONResponse("Test_ListProducts", t, res)
}

func Test_OrderBook(t *testing.T) {
	res, err := go100x.OrderBook(CLIENT_100_X, &types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 0,
		Limit:       constants.LIMIT_FIVE,
	})

	if err != nil {
		t.Errorf("[Test_OrderBook] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_OrderBook", t, 200, res)
	verifyValidJSONResponse("Test_OrderBook", t, res)
}

func Test_ServerTime(t *testing.T) {
	res, err := go100x.ServerTime(CLIENT_100_X)

	if err != nil {
		t.Errorf("[Test_ServerTime] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_ServerTime", t, 200, res)
	verifyValidJSONResponse("Test_ServerTime", t, res)
}

func Test_ApproveSigner(t *testing.T) {
	res, err := go100x.ApproveSigner(CLIENT_100_X, &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_ApproveRevokeSigner] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_ApproveSigner", t, 200, res)
	verifyValidJSONResponse("Test_ApproveSigner", t, res)
}

func Test_RevokeSigner(t *testing.T) {
	res, err := go100x.RevokeSigner(CLIENT_100_X, &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_RevokeSigner] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_RevokeSigner", t, 500, res)
	verifyValidJSONResponse("Test_RevokeSigner", t, res)
}

func Test_NewOrder(t *testing.T) {
	// Limit buy 1 ETH for 3300 USDB, valid for 1 day
	res, err := go100x.NewOrder(CLIENT_100_X, &types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       new(big.Int).Mul(big.NewInt(3300), big.NewInt(params.Ether)).String(),
		Quantity:    new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)).String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_NewOrder] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_NewOrder", t, 400, res)
	verifyValidJSONResponse("Test_NewOrder", t, res)
}

func Test_CancelOrderAndReplace(t *testing.T) {
	res, err := go100x.CancelOrderAndReplace(CLIENT_100_X, &types.CancelOrderAndReplaceRequest{
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

	if err != nil {
		t.Errorf("[Test_CancelOrderAndReplace] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_CancelOrderAndReplace", t, 404, res)
	verifyValidJSONResponse("Test_CancelOrderAndReplace", t, res)
}

func Test_CancelOrder(t *testing.T) {
	res, err := go100x.CancelOrder(CLIENT_100_X, &types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "1",
	})

	if err != nil {
		t.Errorf("[Test_CancelOrder] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_CancelOrder", t, 404, res)
	verifyValidJSONResponse("Test_CancelOrder", t, res)
}

func Test_CancelAllOpenOrders(t *testing.T) {
	res, err := go100x.CancelAllOpenOrders(CLIENT_100_X, &constants.PRODUCT_ETH_PERP)

	if err != nil {
		t.Errorf("[Test_CancelAllOpenOrders] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_CancelAllOpenOrders", t, 200, res)
	verifyValidJSONResponse("Test_CancelAllOpenOrders", t, res)
}

func Test_GetSpotBalances(t *testing.T) {
	res, err := go100x.GetSpotBalances(CLIENT_100_X)
	if err != nil {
		t.Errorf("[Test_GetSpotBalances] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_GetSpotBalances", t, 200, res)
	verifyValidJSONResponse("Test_GetSpotBalances", t, res)
}

func Test_GetPerpetualPosition(t *testing.T) {
	res, err := go100x.GetPerpetualPosition(CLIENT_100_X, &constants.PRODUCT_BLAST_PERP)
	if err != nil {
		t.Errorf("[Test_GetPerpetualPosition] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_GetPerpetualPosition", t, 200, res)
	verifyValidJSONResponse("Test_GetPerpetualPosition", t, res)
}

func Test_ListApproveSigners(t *testing.T) {
	res, err := go100x.ListApprovedSigners(CLIENT_100_X)
	if err != nil {
		t.Errorf("[Test_ListApproveSigners] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_ListApproveSigners", t, 200, res)
	verifyValidJSONResponse("Test_ListApproveSigners", t, res)
}

func Test_ListOpenOrders(t *testing.T) {
	res, err := go100x.ListOpenOrders(CLIENT_100_X, &constants.PRODUCT_BLAST_PERP)
	if err != nil {
		t.Errorf("[Test_ListOpenOrders] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_ListOpenOrders", t, 200, res)
	verifyValidJSONResponse("Test_ListOpenOrders", t, res)
}

func Test_ListOrders(t *testing.T) {
	res, err := go100x.ListOrders(CLIENT_100_X, &types.ListOrdersRequest{
		Product: &constants.PRODUCT_BTC_PERP,
		Ids:     []string{"A", "B"},
	})
	if err != nil {
		t.Errorf("[Test_ListOrders] Error: %v", err)
		return
	}

	verifyResponseStatusCode("Test_ListOrders", t, 200, res)
	verifyValidJSONResponse("Test_ListOrders", t, res)
}

func verifyResponseStatusCode(testName string, t *testing.T, expectedStatusCode int, response *http.Response) {
	if response.StatusCode != expectedStatusCode {
		t.Errorf("[%s] Expected status code %v but got %v", testName, expectedStatusCode, response.StatusCode)
	}
}

func verifyValidJSONResponse(testName string, t *testing.T, response *http.Response) {
	// Read response
	defer response.Body.Close()
	bytesBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("[%s] Error reading response body: %v", testName, err)
		return
	}

	// Check if res is valid JSON by trying to unmarshal it
	var data interface{}
	if err := json.Unmarshal([]byte(bytesBody), &data); err != nil {
		t.Errorf("[%s] Error unmarshalling response: %v", testName, bytesBody)
		return
	}

	t.Logf("[%s] Response is a valid JSON: %s", testName, bytesBody)
}
