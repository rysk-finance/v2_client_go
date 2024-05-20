package ws_client

import (
	"fmt"
	"os"
	"testing"

	"github.com/eldief/go100x/constants"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type WsClientUnitTestSuite struct {
	suite.Suite
	PrivateKeys    string
	RpcUrl         string
	Go100XWSClient *Go100XWSClient
}

func (s *WsClientUnitTestSuite) SetupSuite() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[TestMain] Error loading .env file:", err)
		return
	}

	s.Go100XWSClient = NewGo100XWSClient(&Go100XWSClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,
		PrivateKey:   string(os.Getenv("PRIVATE_KEYS")),
		RpcUrl:       os.Getenv("RPC_URL"),
		SubAccountId: 1,
	})
}

func TestRunSuiteUnit_WsClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(WsClientUnitTestSuite))
}
