package constants

import "github.com/eldief/go100x/types"

var BASE_URI = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "https://api.100x.finance/v1/",
	ENVIRONMENT_TESTNET: "https://api.staging.100x.finance/v1/",
}
