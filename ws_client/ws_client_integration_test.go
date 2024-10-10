//go:build !unit
// +build !unit

package ws_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/rysk-finance/v2_client_go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WsClientIntegrationTestSuite struct {
	suite.Suite
	PrivateKeys    string
	RpcUrl         string
	RyskV2WSClient *RyskV2WSClient
}

func (s *WsClientIntegrationTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	s.RyskV2WSClient, _ = NewRyskV2WSClient(&RyskV2WSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
}

func (s *WsClientIntegrationTestSuite) TearDownSuite() {
	s.RyskV2WSClient.RPCConnection.Close()
	s.RyskV2WSClient.StreamConnection.Close()
}

func TestRunSuiteIntegration_WsClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WsClientIntegrationTestSuite))
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ListProducts() {
	err := s.RyskV2WSClient.ListProducts("LIST_PRODUCTS")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "LIST_PRODUCTS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			resultArray, ok := response.Result.([]interface{})
			require.True(s.T(), ok)
			require.NotEmpty(s.T(), resultArray)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_GetProduct() {
	err := s.RyskV2WSClient.GetProduct("GET_PRODUCT", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PRODUCT" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ServerTime() {
	err := s.RyskV2WSClient.ServerTime("SERVER_TIME")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SERVER_TIME" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_Login() {
	err := s.RyskV2WSClient.Login("LOGIN")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "LOGIN" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SessionStatus() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.SessionStatus("SESSION_STATUS")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SESSION_STATUS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SubAccountList() {
	s.TestIntegration_Login()
	s.TestIntegration_ApproveSigner()

	err := s.RyskV2WSClient.SubAccountList("SUB_ACCOUNT_LIST")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUB_ACCOUNT_LIST" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ApproveSigner() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.ApproveSigner("APPROVE_SIGNER", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.RyskV2WSClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "APPROVE_SIGNER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_RevokeSigner() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.RevokeSigner("REVOKE_SIGNER", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.RyskV2WSClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "REVOKE_SIGNER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_NewOrder() {
	s.TestIntegration_Login()
	s.TestIntegration_ApproveDepositUSDCWaitingTxs()

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
	err = s.RyskV2WSClient.GetProduct("GET_PRODUCT", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	var productUnmarshaled struct {
		Increment string `json:"increment"`
	}

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PRODUCT" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)

			result, ok := response.Result.(map[string]interface{})
			require.True(s.T(), ok)

			resultJSON, err := json.Marshal(result)
			require.NoError(s.T(), err)

			err = json.Unmarshal(resultJSON, &productUnmarshaled)
			require.NoError(s.T(), err)
			break
		}
	}

	increment := new(big.Int)
	_, ok := increment.SetString(productUnmarshaled.Increment, 10)
	require.True(s.T(), ok)

	// adjust price with increment
	remainder := new(big.Int).Mod(price, increment)
	adjustedPrice := new(big.Int).Sub(price, remainder)

	// Limit buy 0.01 ETH for market price, valid for 1 day
	err = s.RyskV2WSClient.NewOrder("NEW_ORDER", &types.NewOrderRequest{
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

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "NEW_ORDER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ListOpenOrders() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.ListOpenOrders("LIST_OPEN_ORDERS", &types.ListOrdersRequest{
		Product:   &constants.PRODUCT_ETH_PERP,
		Ids:       []string{},
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
		Limit:     10,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "LIST_OPEN_ORDERS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_CancelOrder() {
	s.TestIntegration_Login()
	s.TestIntegration_ApproveDepositUSDCWaitingTxs()

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
	price = new(big.Int).Mul(new(big.Int).Div(price, big.NewInt(100)), big.NewInt(110))

	// get product increment
	err = s.RyskV2WSClient.GetProduct("GET_PRODUCT", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	var productUnmarshaled struct {
		Increment string `json:"increment"`
	}

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PRODUCT" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)

			result, ok := response.Result.(map[string]interface{})
			require.True(s.T(), ok)

			resultJSON, err := json.Marshal(result)
			require.NoError(s.T(), err)

			err = json.Unmarshal(resultJSON, &productUnmarshaled)
			require.NoError(s.T(), err)
			break
		}
	}

	increment := new(big.Int)
	_, ok := increment.SetString(productUnmarshaled.Increment, 10)
	require.True(s.T(), ok)

	// adjust price with increment
	remainder := new(big.Int).Mod(price, increment)
	adjustedPrice := new(big.Int).Sub(price, remainder)

	// new order at 10% market premium
	err = s.RyskV2WSClient.NewOrder("NEW_ORDER", &types.NewOrderRequest{
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

	var newOrderUnmarshaled struct {
		ID string `json:"id"`
	}

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "NEW_ORDER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)

			result, ok := response.Result.(map[string]interface{})
			require.True(s.T(), ok)

			resultJSON, err := json.Marshal(result)
			require.NoError(s.T(), err)

			err = json.Unmarshal(resultJSON, &newOrderUnmarshaled)
			require.NoError(s.T(), err)
			break
		}
	}

	// cancel order
	err = s.RyskV2WSClient.CancelOrder("CANCEL_ORDER", &types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: newOrderUnmarshaled.ID,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "CANCEL_ORDER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_CancelAllOpenOrders() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.CancelAllOpenOrders("CANCEL_ALL_OPEN_ORDERS", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "CANCEL_ALL_OPEN_ORDERS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_OrderBook() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.OrderBook("ORDER_BOOK", &types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 10,
		Limit:       constants.LIMIT_FIVE,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "ORDER_BOOK" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_GetPerpetualPosition() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.GetPerpetualPosition("GET_PERPETUAL_POSITION", []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
		&constants.PRODUCT_SOL_PERP,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PERPETUAL_POSITION" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_GetSpotBalances() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.GetSpotBalances("GET_SPOT_BALANCES", []string{})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_SPOT_BALANCES" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			fmt.Println(response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_AccountUpdates() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.AccountUpdates("ACCOUNT_UPDATES")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "ACCOUNT_UPDATES" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Subscribed to updates", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SubscribeAggregateTrades() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.SubscribeAggregateTrades("SUBSCRIBE_AGGREGATE_TRADES", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUBSCRIBE_AGGREGATE_TRADES" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Subscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_UnsubscribeAggregateTrades() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.UnsubscribeAggregateTrades("UNSUBSCRIBE_AGGREGATE_TRADES", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "UNSUBSCRIBE_AGGREGATE_TRADES" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Unsubscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SubscribeSingleTrade() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.SubscribeSingleTrades("SUBSCRIBE_SINGLE_TRADE", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUBSCRIBE_SINGLE_TRADE" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Subscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_UnsubscribeSingleTrade() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.UnsubscribeAggregateTrades("UNSUBSCRIBE_SINGLE_TRADE", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "UNSUBSCRIBE_SINGLE_TRADE" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Unsubscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SubscribeKlineData() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.SubscribeKlineData("SUBSCRIBE_KLINE_DATA", []*types.Product{&constants.PRODUCT_ETH_PERP}, []types.Interval{
		constants.INTERVAL_15M,
		constants.INTERVAL_1H,
		constants.INTERVAL_1M,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUBSCRIBE_KLINE_DATA" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Subscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_UnsubscribeKlineData() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.UnsubscribeKlineData("UNSUBSCRIBE_KLINE_DATA", []*types.Product{&constants.PRODUCT_ETH_PERP}, []types.Interval{
		constants.INTERVAL_15M,
		constants.INTERVAL_1H,
		constants.INTERVAL_1M,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "UNSUBSCRIBE_KLINE_DATA" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Unsubscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SubscribePartialBookDepth() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.SubscribePartialBookDepth("SUBSCRIBE_PARTIAL_BOOK_DEPTH",
		[]*types.Product{
			&constants.PRODUCT_ETH_PERP,
			&constants.PRODUCT_SOL_PERP,
			&constants.PRODUCT_BTC_PERP,
		},
		[]types.Limit{
			constants.LIMIT_FIVE,
			constants.LIMIT_TEN,
			constants.LIMIT_TWENTY,
		},
		[]int64{16, 17, 18},
	)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUBSCRIBE_PARTIAL_BOOK_DEPTH" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Subscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_UnsubscribePartialBookDepth() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.UnsubscribePartialBookDepth("UNSUBSCRIBE_PARTIAL_BOOK_DEPTH",
		[]*types.Product{
			&constants.PRODUCT_ETH_PERP,
			&constants.PRODUCT_SOL_PERP,
			&constants.PRODUCT_BTC_PERP,
		},
		[]types.Limit{
			constants.LIMIT_FIVE,
			constants.LIMIT_TEN,
			constants.LIMIT_TWENTY,
		},
		[]int64{16, 17, 18},
	)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "UNSUBSCRIBE_PARTIAL_BOOK_DEPTH" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Unsubscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_Subscribe24hrPriceChangeStatistics() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.Subscribe24hrPriceChangeStatistics("SUBSCRIBE_24H_PRICE_CHANGE_STATS",
		[]*types.Product{
			&constants.PRODUCT_ETH_PERP,
			&constants.PRODUCT_SOL_PERP,
			&constants.PRODUCT_BTC_PERP,
		},
	)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUBSCRIBE_24H_PRICE_CHANGE_STATS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Subscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_Unsubscribe24hrPriceChangeStatistics() {
	s.TestIntegration_Login()

	err := s.RyskV2WSClient.Unsubscribe24hrPriceChangeStatistics("UNSUBSCRIBE_24H_PRICE_CHANGE_STATS",
		[]*types.Product{
			&constants.PRODUCT_ETH_PERP,
			&constants.PRODUCT_SOL_PERP,
			&constants.PRODUCT_BTC_PERP,
		},
	)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.RyskV2WSClient.StreamConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "UNSUBSCRIBE_24H_PRICE_CHANGE_STATS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "Unsubscribed", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ApproveUSDCWaitingTx() {
	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), constants.E9)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.RyskV2WSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ApproveDepositUSDCWaitingTxs() {
	transaction, err := s.RyskV2WSClient.ApproveUSDC(context.Background(), constants.E9)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.RyskV2WSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)

	transaction, err = s.RyskV2WSClient.DepositUSDC(context.Background(), constants.E9)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err = s.RyskV2WSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)
}

func (s *WsClientIntegrationTestSuite) TestIntegration_addReferee() {
	res, err := s.RyskV2WSClient.addReferee()
	require.NoError(s.T(), err)
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

	fmt.Println(string(bytesBody))
}
