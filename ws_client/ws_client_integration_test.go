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

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WsClientIntegrationTestSuite struct {
	suite.Suite
	PrivateKeys    string
	RpcUrl         string
	Go100XWSClient *Go100XWSClient
}

func (s *WsClientIntegrationTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	s.Go100XWSClient, _ = NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
}

func TestRunSuiteIntegration_WsClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WsClientIntegrationTestSuite))
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ListProducts() {
	err := s.Go100XWSClient.ListProducts("LIST_PRODUCTS")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
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
	err := s.Go100XWSClient.GetProduct("GET_PRODUCT", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PRODUCT" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ServerTime() {
	err := s.Go100XWSClient.ServerTime("SERVER_TIME")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SERVER_TIME" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_Login() {
	err := s.Go100XWSClient.Login("LOGIN")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "LOGIN" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SessionStatus() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.SessionStatus("SESSION_STATUS")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SESSION_STATUS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_SubAccountList() {
	s.TestIntegration_Login()
	s.TestIntegration_ApproveSigner()

	err := s.Go100XWSClient.SubAccountList("SUB_ACCOUNT_LIST")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "SUB_ACCOUNT_LIST" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ApproveSigner() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.ApproveSigner("APPROVE_SIGNER", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.Go100XWSClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "APPROVE_SIGNER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_RevokeSigner() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.RevokeSigner("REVOKE_SIGNER", &types.ApproveRevokeSignerRequest{
		ApprovedSigner: s.Go100XWSClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "REVOKE_SIGNER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_NewOrder() {
	s.TestIntegration_Login()
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
	err = s.Go100XWSClient.GetProduct("GET_PRODUCT", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	var productUnmarshaled struct {
		Increment string `json:"increment"`
	}

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PRODUCT" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))

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
	err = s.Go100XWSClient.NewOrder("NEW_ORDER", &types.NewOrderRequest{
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
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "NEW_ORDER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ListOpenOrders() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.ListOpenOrders("LIST_OPEN_ORDERS", &types.ListOrdersRequest{
		Product:   &constants.PRODUCT_ETH_PERP,
		Ids:       []string{},
		StartTime: time.Now().Add(-24 * time.Hour).UnixMilli(),
		EndTime:   time.Now().UnixMilli(),
		Limit:     10,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "LIST_OPEN_ORDERS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Empty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_CancelOrder() {
	s.TestIntegration_Login()
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
	price = new(big.Int).Mul(new(big.Int).Div(price, big.NewInt(100)), big.NewInt(110))

	// get product increment
	err = s.Go100XWSClient.GetProduct("GET_PRODUCT", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	var productUnmarshaled struct {
		Increment string `json:"increment"`
	}

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PRODUCT" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))

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
	err = s.Go100XWSClient.NewOrder("NEW_ORDER", &types.NewOrderRequest{
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
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "NEW_ORDER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))

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
	err = s.Go100XWSClient.CancelOrder("CANCEL_ORDER", &types.CancelOrderRequest{
		Product:    &constants.PRODUCT_ETH_PERP,
		IdToCancel: newOrderUnmarshaled.ID,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "CANCEL_ORDER" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_CancelAllOpenOrders() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.CancelAllOpenOrders("CANCEL_ALL_OPEN_ORDERS", &constants.PRODUCT_ETH_PERP)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "CANCEL_ALL_OPEN_ORDERS" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Equal(s.T(), "OK", response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_OrderBook() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.OrderBook("ORDER_BOOK", &types.OrderBookRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		Granularity: 10,
		Limit:       constants.LIMIT_FIVE,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "ORDER_BOOK" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.NotEmpty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_GetPerpetualPosition() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.GetPerpetualPosition("GET_PERPETUAL_POSITION", []*types.Product{
		&constants.PRODUCT_ETH_PERP,
		&constants.PRODUCT_BTC_PERP,
		&constants.PRODUCT_SOL_PERP,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_PERPETUAL_POSITION" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Empty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_GetSpotBalances() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.GetSpotBalances("GET_SPOT_BALANCES", []string{})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
		require.NoError(s.T(), err)

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.NoError(s.T(), err)

		if response.ID == "GET_SPOT_BALANCES" {
			require.True(s.T(), response.Success)
			require.Nil(s.T(), response.Error)
			require.Equal(s.T(), "2.0", string(response.JsonRPC))
			require.Empty(s.T(), response.Result)
			break
		}
	}
}

func (s *WsClientIntegrationTestSuite) TestIntegration_AccountUpdates() {
	s.TestIntegration_Login()

	err := s.Go100XWSClient.AccountUpdates("ACCOUNT_UPDATES")
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.RPCConnection.ReadMessage()
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

	err := s.Go100XWSClient.SubscribeAggregateTrades("SUBSCRIBE_AGGREGATE_TRADES", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.UnsubscribeAggregateTrades("UNSUBSCRIBE_AGGREGATE_TRADES", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.SubscribeSingleTrades("SUBSCRIBE_SINGLE_TRADE", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.UnsubscribeAggregateTrades("UNSUBSCRIBE_SINGLE_TRADE", []*types.Product{&constants.PRODUCT_ETH_PERP})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.SubscribeKlineData("SUBSCRIBE_KLINE_DATA", []*types.Product{&constants.PRODUCT_ETH_PERP}, []types.Interval{
		constants.INTERVAL_15M,
		constants.INTERVAL_1H,
		constants.INTERVAL_1M,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.UnsubscribeKlineData("UNSUBSCRIBE_KLINE_DATA", []*types.Product{&constants.PRODUCT_ETH_PERP}, []types.Interval{
		constants.INTERVAL_15M,
		constants.INTERVAL_1H,
		constants.INTERVAL_1M,
	})
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.SubscribePartialBookDepth("SUBSCRIBE_PARTIAL_BOOK_DEPTH",
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
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.UnsubscribePartialBookDepth("UNSUBSCRIBE_PARTIAL_BOOK_DEPTH",
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
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.Subscribe24hrPriceChangeStatistics("SUBSCRIBE_24H_PRICE_CHANGE_STATS",
		[]*types.Product{
			&constants.PRODUCT_ETH_PERP,
			&constants.PRODUCT_SOL_PERP,
			&constants.PRODUCT_BTC_PERP,
		},
	)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

	err := s.Go100XWSClient.Unsubscribe24hrPriceChangeStatistics("UNSUBSCRIBE_24H_PRICE_CHANGE_STATS",
		[]*types.Product{
			&constants.PRODUCT_ETH_PERP,
			&constants.PRODUCT_SOL_PERP,
			&constants.PRODUCT_BTC_PERP,
		},
	)
	require.NoError(s.T(), err)

	for {
		_, p, err := s.Go100XWSClient.StreamConnection.ReadMessage()
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

func (s *WsClientIntegrationTestSuite) TestIntegration_ApproveUSDBWaitingTx() {
	transaction, err := s.Go100XWSClient.ApproveUSDB(context.Background(), constants.E22)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.Go100XWSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)
}

func (s *WsClientIntegrationTestSuite) TestIntegration_ApproveDepositUSDBWaitingTxs() {
	transaction, err := s.Go100XWSClient.ApproveUSDB(context.Background(), constants.E20)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err := s.Go100XWSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)

	transaction, err = s.Go100XWSClient.DepositUSDB(context.Background(), constants.E20)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), transaction)

	receipt, err = s.Go100XWSClient.WaitTransaction(context.Background(), transaction)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), receipt)
	require.Equal(s.T(), uint64(1), receipt.Status)
}

func (s *WsClientIntegrationTestSuite) TestIntegration_addReferee() {
	res, err := s.Go100XWSClient.addReferee()
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
