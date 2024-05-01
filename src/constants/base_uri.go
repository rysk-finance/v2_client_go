package constants

import "go100x/src/types"

var BASE_URI = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "https://api.100x.finance/v1/",
	ENVIRONMENT_TESTNET: "https://api.staging.100x.finance/v1/",
}
