//go:build !unit
// +build !unit

package ws_client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
		fmt.Println("[TestMain] Error loading .env file:", err)
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
