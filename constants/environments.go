package constants

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rysk-finance/v2_client_go/types"
)

const (
	ENVIRONMENT_MAINNET types.Environment = "mainnet"
	ENVIRONMENT_TESTNET types.Environment = "testnet"
)

var CHAIN_ID = map[types.Environment]*math.HexOrDecimal256{
	ENVIRONMENT_MAINNET: nil,
	ENVIRONMENT_TESTNET: math.NewHexOrDecimal256(421614),
}

var CIAO_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "",
	ENVIRONMENT_TESTNET: "0x71728FDDF90233cc35D61bec7858d7c42A310ACe",
}

var USDC_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "",
	ENVIRONMENT_TESTNET: "0xb8bE1401E65dC08Bfb8f832Fc1A27a16CA821B05",
}

var ORDER_DISPATCHER_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "",
	ENVIRONMENT_TESTNET: "0x27809a3Bd3cf44d855f1BE668bFD16D34bcE157C",
}
