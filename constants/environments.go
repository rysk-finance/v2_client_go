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
	ENVIRONMENT_TESTNET: "0x9Bb24CBd60b649a490f35466FDA65b86e986Ee9b",
}

var USDB_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0x4300000000000000000000000000000000000003",
	ENVIRONMENT_TESTNET: "0x79A59c326C715AC2d31C169C85d1232319E341ce",
}

var VERIFIER_ADDRESS = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "0x691a5fc3a81a144e36c6C4fBCa1fC82843c80d0d",
	ENVIRONMENT_TESTNET: "0x888F59A023c82C7FEFaC49CbEFcBEb62db111E4c",
}
