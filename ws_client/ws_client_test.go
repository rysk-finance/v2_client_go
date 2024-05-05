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
)

// .env variables to be setup via `TestMain` before running tests.
var (
	WS_PRIVATE_KEYS  string = ""
	WS_RPC_URL       string = ""
	GO100X_WS_CLIENT *Go100XWSClient
)

// Setup .env variables for test suite.
func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		os.Exit(1)
	}

	WS_PRIVATE_KEYS = string(os.Getenv("PRIVATE_KEYS"))
	WS_RPC_URL = os.Getenv("RPC_URL")
	GO100X_WS_CLIENT = NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   WS_PRIVATE_KEYS,
		RpcUrl:       WS_RPC_URL,
		SubAccountId: 1,
	})

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Login(t *testing.T) {
	err := GO100X_WS_CLIENT.Login("LOGIN")
	if err != nil {
		t.Errorf("[Test_Login] Error during login process: %v", err)
	}

	for {
		_, p, err := GO100X_WS_CLIENT.wsConnection.ReadMessage()
		require.Nil(t, err, "[Test_Login] Error reading message from WebSocket")

		var response types.WebsocketResponse
		err = json.Unmarshal(p, &response)
		require.Nil(t, err, "[Test_Login] Error unmarshalling JSON response")

		if response.ID == "LOGIN" {
			require.True(t, response.Success, "[Test_Login] Expected success response")
			require.Nil(t, response.Error, "[Test_Login] Expected no error in response")
			require.Equal(t, "2.0", string(response.JsonRPC), "[Test_Login] Unexpected JSON-RPC version")
			require.Equal(t, "OK", response.Result, "[Test_Login] Expected 'response.Result' to be 'OK'")

			break
		}
	}
}
