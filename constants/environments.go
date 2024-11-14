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
	ENVIRONMENT_MAINNET: math.NewHexOrDecimal256(42161),
	ENVIRONMENT_TESTNET: math.NewHexOrDecimal256(421614),
}

var CIAO_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0xD8eB81D7D31b420b435Cb3C61a8B4E7805e12Eff",
	ENVIRONMENT_TESTNET: "0x71728FDDF90233cc35D61bec7858d7c42A310ACe",
}

var USDC_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
	ENVIRONMENT_TESTNET: "0xb8bE1401E65dC08Bfb8f832Fc1A27a16CA821B05",
}

var ORDER_DISPATCHER_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0x6644D5B09EBae015fE4e3a87Eff1A07d33558E59",
	ENVIRONMENT_TESTNET: "0x27809a3Bd3cf44d855f1BE668bFD16D34bcE157C",
}
