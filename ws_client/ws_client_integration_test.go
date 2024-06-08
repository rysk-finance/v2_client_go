package ws_client

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

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

func (s *WsClientIntegrationTestSuite) TestIntegraton_Login() {
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
