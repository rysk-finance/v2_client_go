package go100x_test

import (
	"encoding/json"
	"fmt"
	go100x "go100x/src/client"
	"go100x/src/constants"
	"go100x/src/types"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/joho/godotenv"
)

// .env variables to be setup via `TestMain` before running tests.
var (
	PRIVATE_KEYS string = ""
	RPC_URL      string = ""
	CLIENT       *types.Client
)

// Setup .env variables for test suite.
func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		os.Exit(1)
	}

	PRIVATE_KEYS = string(os.Getenv("PRIVATE_KEYS"))
	RPC_URL = os.Getenv("RPC_URL")
	CLIENT = go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Get24hrPriceChangeStatistics_NoProduct(t *testing.T) {
	res, err := go100x.Get24hrPriceChangeStatistics(CLIENT, &types.Product{})

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_Get24hrPriceChangeStatistics_NoProduct", t, res)
}

func Test_Get24hrPriceChangeStatistics_WithNonExistingProduct(t *testing.T) {
	res, err := go100x.Get24hrPriceChangeStatistics(CLIENT, &types.Product{
		Id:     69420,
		Symbol: "69420",
	})

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Error: %v", err)
		return
	}

	expectedResponse := "Product not found"
	if res != expectedResponse {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_WithNonExistingProduct] Unexpected response. Got: %s, Expected: %s", res, expectedResponse)
	}
}

func Test_Get24hrPriceChangeStatistics_WithProduct(t *testing.T) {
	res, err := go100x.Get24hrPriceChangeStatistics(CLIENT, &constants.ETH_PERP)

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_WithProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_Get24hrPriceChangeStatistics_WithProduct", t, res)
}

func Test_GetProduct(t *testing.T) {
	res, err := go100x.GetProduct(CLIENT, constants.ETH_PERP.Symbol)

	if err != nil {
		t.Errorf("[Test_GetProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_GetProduct", t, res)
}

func Test_GetProductById(t *testing.T) {
	res, err := go100x.GetProductById(CLIENT, constants.ETH_PERP.Id)

	if err != nil {
		t.Errorf("[Test_GetProductById] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_GetProductById", t, res)
}

func Test_GetKlineData(t *testing.T) {
	res, err := go100x.GetKlineData(CLIENT, &types.KlineDataRequest{
		Product:   constants.BTC_PERP,
		Interval:  constants.D1,
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_GetKlineData] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_GetKlineData", t, res)
}

func Test_ListProducts(t *testing.T) {
	res, err := go100x.ListProducts(CLIENT)

	if err != nil {
		t.Errorf("[Test_ListProducts] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_ListProducts", t, res)
}

func Test_OrderBook(t *testing.T) {
	res, err := go100x.OrderBook(CLIENT, &types.OrderBookRequest{
		Product:     constants.ETH_PERP,
		Granularity: 0,
		Limit:       constants.FIVE,
	})

	if err != nil {
		t.Errorf("[Test_OrderBook] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_OrderBook", t, res)
}

func Test_ServerTime(t *testing.T) {
	res, err := go100x.ServerTime(CLIENT)

	if err != nil {
		t.Errorf("[Test_ServerTime] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_ServerTime", t, res)
}

func Test_ApproveSigner(t *testing.T) {
	res, err := go100x.ApproveSigner(CLIENT, &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_ApproveRevokeSigner] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_ApproveSigner", t, res)
}

func Test_RevokeSigner(t *testing.T) {
	res, err := go100x.RevokeSigner(CLIENT, &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_RevokeSigner] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_RevokeSigner", t, res)
}

func Test_Login(t *testing.T) {
	res, err := go100x.Login(CLIENT)

	if err != nil {
		t.Errorf("[Test_Login] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_Login", t, res)
}

func Test_NewOrder(t *testing.T) {
	// Limit buy 1 ETH for 3300 USDB, valid for 1 day
	res, err := go100x.NewOrder(CLIENT, &types.NewOrderRequest{
		Product:     constants.ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.LIMIT,
		TimeInForce: constants.GTC,
		Price:       new(big.Int).Mul(big.NewInt(3300), big.NewInt(params.Ether)).String(),
		Quantity:    new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)).String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_NewOrder] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_NewOrder", t, res)
}

func Test_CancelOrderAndReplace(t *testing.T) {
	res, err := go100x.CancelOrderAndReplace(CLIENT, &types.CancelOrderAndReplaceRequest{
		IdToCancel: "1",
		// Limit buy 1 ETH for 3300 USDB, valid for 1 day
		NewOrder: types.NewOrderRequest{
			Product:     constants.ETH_PERP,
			IsBuy:       true,
			OrderType:   constants.LIMIT_MAKER,
			TimeInForce: constants.GTC,
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

	verifyValidJSONResponse("Test_CancelOrderAndReplace", t, res)
}

func Test_GetSpotBalances(t *testing.T) {
	res, err := go100x.GetSpotBalances(CLIENT)
	if err != nil {
		t.Errorf("[Test_GetSpotBalances] Error: %v", err)
		return
	}

	verifyValidJSONResponse("Test_GetSpotBalances", t, res)
}

func verifyValidJSONResponse(testName string, t *testing.T, res string) {
	var data interface{}

	// Check if res is valid JSON by trying to unmarshal it
	if err := json.Unmarshal([]byte(res), &data); err != nil {
		t.Errorf("[%s] Error unmarshalling response: %v", testName, err)
		return
	}

	t.Logf("[%s] %s", testName, res)
}
