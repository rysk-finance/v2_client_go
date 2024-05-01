package constants

import (
	"go100x/src/types"

	"github.com/ethereum/go-ethereum/common/math"
)

var CHAIN_ID = map[types.Environment]*math.HexOrDecimal256{
	ENVIRONMENT_MAINNET: math.NewHexOrDecimal256(81457),
	ENVIRONMENT_TESTNET: math.NewHexOrDecimal256(168587773),
}
