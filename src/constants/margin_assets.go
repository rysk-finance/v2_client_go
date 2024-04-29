package constants

import (
	"go100x/src/types"
)

const (
	USDB types.MarginAssetKey = "USDB"
)

var MARGIN_ASSET = map[types.Environment]map[types.MarginAssetKey]string{
	MAINNET: {
		USDB: "0x79a59c326c715ac2d31c169c85d1232319e341ce",
	},
	TESTNET: {
		USDB: "0x79a59c326c715ac2d31c169c85d1232319e341ce",
	},
}
