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
)

// Setup .env variables for test suite.
func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		os.Exit(1)
	}

	PRIVATE_KEYS = string(os.Getenv("PRIVATE_KEYS"))
	RPC_URL = os.Getenv("RPC_URL")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Get24hrPriceChangeStatistics_NoProduct(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.Get24hrPriceChangeStatistics(client, &types.Product{})

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_NoProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_Get24hrPriceChangeStatistics_WithNonExistingProduct(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.Get24hrPriceChangeStatistics(client, &types.Product{
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
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.Get24hrPriceChangeStatistics(client, &constants.ETH_PERP)

	if err != nil {
		t.Errorf("[Test_Get24hrPriceChangeStatistics_WithProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_GetProduct(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.GetProduct(client, constants.ETH_PERP.Symbol)

	if err != nil {
		t.Errorf("[Test_GetProduct] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_GetProductById(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.GetProductById(client, constants.ETH_PERP.Id)

	if err != nil {
		t.Errorf("[Test_GetProductById] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_GetKlineData(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.GetKlineData(client, &types.KlineDataRequest{
		Product:   constants.BTC_PERP,
		Interval:  constants.D1,
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
	})

	if err != nil {
		t.Errorf("[Test_GetKlineData] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_ListProducts(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.ListProducts(client)

	if err != nil {
		t.Errorf("[Test_ListProducts] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_OrderBook(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.OrderBook(client, &types.OrderBookRequest{
		Product:     constants.ETH_PERP,
		Granularity: 0,
		Limit:       constants.FIVE,
	})

	if err != nil {
		t.Errorf("[Test_OrderBook] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_ServerTime(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.ServerTime(client)

	if err != nil {
		t.Errorf("[Test_ServerTime] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_ApproveRevokeSigner(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.ApproveRevokeSigner(client, &types.ApproveRevokeSignerRequest{
		ApprovedSigner: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // vitalik.eth
		Nonce:          time.Now().UnixMilli(),
		IsApproved:     false,
	})

	if err != nil {
		t.Errorf("[Test_ApproveRevokeSigner] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_Login(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.Login(client)

	if err != nil {
		t.Errorf("[Test_Login] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func Test_NewOrder(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	// Limit buy 1 ETH for 3300 USDB, valid for 1 day
	res, err := go100x.NewOrder(client, &types.NewOrderRequest{
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

	verifyValidJSONResponse(t, res)
}

func Test_GetSpotBalances(t *testing.T) {
	client := go100x.NewClient(&types.ClientConfiguration{
		Env:          constants.TESTNET,
		PrivateKey:   PRIVATE_KEYS,
		RpcUrl:       RPC_URL,
		Timeout:      10 * time.Second,
		SubAccountId: 1,
	})

	res, err := go100x.GetSpotBalances(client)
	if err != nil {
		t.Errorf("[Test_GetSpotBalances] Error: %v", err)
		return
	}

	verifyValidJSONResponse(t, res)
}

func verifyValidJSONResponse(t *testing.T, res string) {
	var data interface{}

	// Check if res is valid JSON by trying to unmarshal it
	if err := json.Unmarshal([]byte(res), &data); err != nil {
		t.Errorf("[verifyValidJSONResponse] Error unmarshalling response: %v", err)
		return
	}

	t.Logf("[verifyValidJSONResponse] Response is valid JSON:\n%s", res)
}
