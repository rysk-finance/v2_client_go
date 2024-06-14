package constants

import (
	"github.com/eldief/go100x/types"
	"github.com/ethereum/go-ethereum/common/math"
)

const (
	ENVIRONMENT_MAINNET types.Environment = "mainnet"
	ENVIRONMENT_TESTNET types.Environment = "testnet"
)

var CHAIN_ID = map[types.Environment]*math.HexOrDecimal256{
	ENVIRONMENT_MAINNET: math.NewHexOrDecimal256(81457),
	ENVIRONMENT_TESTNET: math.NewHexOrDecimal256(168587773),
}

var CIAO_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0x1BaEbEE6B00B3f559B0Ff0719B47E0aF22A6bfC4",
	ENVIRONMENT_TESTNET: "0x0c3b9472b3923CfE199bAE24B5f5bD75FAD2bae9",
}

var USDB_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0x4300000000000000000000000000000000000003",
	ENVIRONMENT_TESTNET: "0x79A59c326C715AC2d31C169C85d1232319E341ce",
}

var VERIFIER_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0x65CbB566D1A6E60107c0c7888761de1AdFa1ccC0",
	ENVIRONMENT_TESTNET: "0x02Ca4fcB63E2D3C89fa20D86ccDcfc540c683545",
}
