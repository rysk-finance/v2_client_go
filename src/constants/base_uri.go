package constants

import "go100x/src/types"

var BASE_URI = map[types.Environment]string{
	MAINNET: "https://api.100x.finance/v1/",
	TESTNET: "https://api.staging.100x.finance/v1/",
}
