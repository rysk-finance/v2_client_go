package constants

import (
	"github.com/eldief/go100x/types"
)

const (
	MARGIN_ASSET_USDB types.MarginAsset = "USDB"
)

var MARGIN_ASSET = map[types.Environment]map[types.MarginAsset]string{
	ENVIRONMENT_MAINNET: {
		MARGIN_ASSET_USDB: "0x79a59c326c715ac2d31c169c85d1232319e341ce",
	},
	ENVIRONMENT_TESTNET: {
		MARGIN_ASSET_USDB: "0x79a59c326c715ac2d31c169c85d1232319e341ce",
	},
}
