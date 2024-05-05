package constants

import "github.com/eldief/go100x/types"

var WS_URL = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "wss://api.staging.100x.finance/v1/ws/operate",
	ENVIRONMENT_TESTNET: "wss://api.staging.100x.finance/v1/ws/operate",
}
