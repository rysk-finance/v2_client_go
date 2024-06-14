//go:build !unit
// +build !unit

package api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	apiClient, _ := NewGo100XAPIClient(&Go100XAPIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 0,
	})
	s.Go100XApiClient = apiClient
}

func (s *ApiClientIntegrationTestSuite) SetupTest() {
	time.Sleep(100 * time.Millisecond)
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
		ApprovedSigner: s.Go100XApiClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)

	// verify approval granted
	res, err = s.Go100XApiClient.ListApprovedSigners()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	var unmarshaled []struct {
		Account    string
		Subaccount uint8
		Signer     string
		Approved   bool
	}
	err = json.Unmarshal(body, &unmarshaled)
	require.NoError(s.T(), err)

	for _, signer := range unmarshaled {
		if signer.Account == strings.ToLower(s.Go100XApiClient.addressString) &&
			signer.Subaccount == uint8(s.Go100XApiClient.SubAccountId) &&
			signer.Signer == strings.ToLower(s.Go100XApiClient.addressString) {
			require.True(s.T(), signer.Approved)
			return
		}
	}
	require.FailNow(s.T(), "Signer not found")
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_RevokeSigner() {
	res, err := s.Go100XApiClient.RevokeSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.Go100XApiClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)

	// verify approval revoked
	res, err = s.Go100XApiClient.ListApprovedSigners()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	var unmarshaled []struct {
		Account    string
		Subaccount uint8
		Signer     string
		Approved   bool
	}
	err = json.Unmarshal(body, &unmarshaled)
	require.NoError(s.T(), err)

	for _, signer := range unmarshaled {
		if signer.Account == strings.ToLower(s.Go100XApiClient.addressString) &&
			signer.Subaccount == uint8(s.Go100XApiClient.SubAccountId) &&
			signer.Signer == strings.ToLower(s.Go100XApiClient.addressString) {
			require.False(s.T(), signer.Approved)
			return
		}
	}
	require.FailNow(s.T(), "Signer not found")
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_Withdraw() {
	s.TestIntegration_ApproveDepositUSDBWaitingTxs()

	// withdraw
	res, err := s.Go100XApiClient.Withdraw(&types.WithdrawRequest{
		Quantity: constants.E20.String(),
		Nonce:    time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_NewOrder() {
	s.TestIntegration_ApproveDepositUSDBWaitingTxs()

	// get market price
	request, err := http.NewRequest(
		http.MethodGet,
		"https://api.coinbase.com/v2/exchange-rates?currency=ETH",
		nil,
	)
	require.NoError(s.T(), err)

	res, err := http.DefaultClient.Do(request)
	require.NoError(s.T(), err)

	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	var priceUnmarshaled struct {
		Data struct {
			Currency string
			Rates    struct {
				USD string
			}
		}
	}
	err = json.Unmarshal(body, &priceUnmarshaled)
	require.NoError(s.T(), err)

	priceFloat, err := strconv.ParseFloat(priceUnmarshaled.Data.Rates.USD, 64)
	require.NoError(s.T(), err)

	price := new(big.Int)
	new(big.Float).Mul(big.NewFloat(priceFloat), new(big.Float).SetFloat64(1e18)).Int(price)

	// get product increment
	res, err = s.Go100XApiClient.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	var productUnmarshaled struct {
		Increment string
	}
	body, err = io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	err = json.Unmarshal(body, &productUnmarshaled)
	require.NoError(s.T(), err)

	increment := new(big.Int)
	_, ok := increment.SetString(productUnmarshaled.Increment, 10)
	require.True(s.T(), ok)

	// adjust price with increment
	remainder := new(big.Int).Mod(price, increment)
	adjustedPrice := new(big.Int).Sub(price, remainder)

	// Limit buy 0.01 ETH for market price, valid for 1 day
	res, err = s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       adjustedPrice.String(),
		Quantity:    constants.E16.String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMicro(),
	})
	verifyValidJSONResponse(s.T(), res)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelOrderAndReplace() {
	s.TestIntegration_ApproveDepositUSDBWaitingTxs()

	// get market price
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.coinbase.com/v2/exchange-rates?currency=ETH",
		nil,
	)
	require.NoError(s.T(), err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)

	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)
	defer res.Body.Close()

	var unmarshaled struct {
		Data struct {
			Currency string
			Rates    struct {
				USD string
			}
		}
	}
	err = json.Unmarshal(body, &unmarshaled)
	require.NoError(s.T(), err)

	priceFloat, err := strconv.ParseFloat(unmarshaled.Data.Rates.USD, 64)
	require.NoError(s.T(), err)

	price := new(big.Int)
	new(big.Float).Mul(big.NewFloat(priceFloat), new(big.Float).SetFloat64(1e18)).Int(price)
	price = new(big.Int).Mul(new(big.Int).Div(price, big.NewInt(100)), big.NewInt(110))

	// get product increment
	res, err = s.Go100XApiClient.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	var productUnmarshaled struct {
		Increment string
	}
	body, err = io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	err = json.Unmarshal(body, &productUnmarshaled)
	require.NoError(s.T(), err)

	increment := new(big.Int)
	_, ok := increment.SetString(productUnmarshaled.Increment, 10)
	require.True(s.T(), ok)

	// adjust price with increment
	remainder := new(big.Int).Mod(price, increment)
	adjustedPrice := new(big.Int).Sub(price, remainder)

	// new order at 10% market premium
	res, err = s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       adjustedPrice.String(),
		Quantity:    constants.E16.String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	body, err = io.ReadAll(res.Body)
	require.NoError(s.T(), err)
	fmt.Println(string(body))

	var order struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(body, &order)
	require.NoError(s.T(), err)

	// cancel and replace
	res, err = s.Go100XApiClient.CancelOrderAndReplace(&types.CancelOrderAndReplaceRequest{
		IdToCancel: order.ID,
		NewOrder: &types.NewOrderRequest{
			Product:     &constants.PRODUCT_ETH_PERP,
			IsBuy:       true,
			OrderType:   constants.ORDER_TYPE_LIMIT_MAKER,
			TimeInForce: constants.TIME_IN_FORCE_GTC,
			Price:       adjustedPrice.String(),
			Quantity:    constants.E16.String(),
			Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
			Nonce:       time.Now().UnixMicro(),
		},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_CancelOrder() {
	s.TestIntegration_ApproveDepositUSDBWaitingTxs()

	// get market price
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.coinbase.com/v2/exchange-rates?currency=ETH",
		nil,
	)
	require.NoError(s.T(), err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)

	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)
	defer res.Body.Close()

	var unmarshaled struct {
		Data struct {
			Currency string
			Rates    struct {
				USD string
			}
		}
	}
	err = json.Unmarshal(body, &unmarshaled)
	require.NoError(s.T(), err)

	priceFloat, err := strconv.ParseFloat(unmarshaled.Data.Rates.USD, 64)
	require.NoError(s.T(), err)

	price := new(big.Int)
	new(big.Float).Mul(big.NewFloat(priceFloat), new(big.Float).SetFloat64(1e18)).Int(price)
	price = new(big.Int).Mul(new(big.Int).Div(price, big.NewInt(100)), big.NewInt(110))

	// get product increment
	res, err = s.Go100XApiClient.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	var productUnmarshaled struct {
		Increment string
	}
	body, err = io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	err = json.Unmarshal(body, &productUnmarshaled)
	require.NoError(s.T(), err)

	increment := new(big.Int)
	_, ok := increment.SetString(productUnmarshaled.Increment, 10)
	require.True(s.T(), ok)

	// adjust price with increment
	remainder := new(big.Int).Mod(price, increment)
	adjustedPrice := new(big.Int).Sub(price, remainder)

	// new order at 10% market premium
	res, err = s.Go100XApiClient.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       adjustedPrice.String(),
		Quantity:    constants.E16.String(),
		Expiration:  time.Now().Add(24 * time.Hour).UnixMilli(),
		Nonce:       time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	body, err = io.ReadAll(res.Body)
	require.NoError(s.T(), err)
	fmt.Println(string(body))

	// cancel order
	res, err = s.Go100XApiClient.CancelOrder(&types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: "1",
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
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
		Ids:     []string{},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)
	verifyValidJSONResponse(s.T(), res)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ApproveUSDBWaitingTx() {
	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), constants.E22)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.Go100XApiClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)
}

func (s *ApiClientIntegrationTestSuite) TestIntegration_ApproveDepositUSDBWaitingTxs() {
	transaction, err := s.Go100XApiClient.ApproveUSDB(context.Background(), constants.E20)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.Go100XApiClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)

	transaction, err = s.Go100XApiClient.DepositUSDB(context.Background(), constants.E20)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err = s.Go100XApiClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)

	res, err := s.Go100XApiClient.GetSpotBalances()
	require.NoError(s.T(), err)
	require.Equal(s.T(), 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)

	var unmarshaled []struct {
		Account    string
		SubAccount int64
		Quantity   string
	}
	err = json.Unmarshal(body, &unmarshaled)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), unmarshaled)

	for _, balance := range unmarshaled {
		if balance.Account == strings.ToLower(s.Go100XApiClient.addressString) &&
			balance.SubAccount == s.Go100XApiClient.SubAccountId {
			quantity := new(big.Int)
			quantity, ok := quantity.SetString(balance.Quantity, 10)
			require.True(s.T(), ok)
			require.GreaterOrEqual(s.T(), quantity.Cmp(constants.E20), 0)
			return
		}
	}
	require.FailNow(s.T(), "Balance not found")
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

	fmt.Println(string(bytesBody))
}
