package types

import (
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type PrimaryType string

var EIP712_TYPES = apitypes.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"LoginMessage": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "message",
			Type: "string",
		},
		{
			Name: "timestamp",
			Type: "uint64",
		},
	},
	"Order": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
		{
			Name: "productId",
			Type: "uint32",
		},
		{
			Name: "isBuy",
			Type: "bool",
		},
		{
			Name: "orderType",
			Type: "uint8",
		},
		{
			Name: "timeInForce",
			Type: "uint8",
		},
		{
			Name: "expiration",
			Type: "uint64",
		},
		{
			Name: "price",
			Type: "uint128",
		},
		{
			Name: "quantity",
			Type: "uint128",
		},
		{
			Name: "nonce",
			Type: "uint64",
		},
	},
	"CancelOrders": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
		{
			Name: "productId",
			Type: "uint32",
		},
	},
	"CancelOrder": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
		{
			Name: "productId",
			Type: "uint32",
		},
		{
			Name: "orderId",
			Type: "string",
		},
	},
	"ApproveSigner": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
		{
			Name: "approvedSigner",
			Type: "address",
		},
		{
			Name: "isApproved",
			Type: "bool",
		},
		{
			Name: "nonce",
			Type: "uint64",
		},
	},
	"Deposit": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
		{
			Name: "asset",
			Type: "address",
		},
		{
			Name: "quantity",
			Type: "uint256",
		},
		{
			Name: "nonce",
			Type: "uint64",
		},
	},
	"Withdraw": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
		{
			Name: "asset",
			Type: "address",
		},
		{
			Name: "quantity",
			Type: "uint128",
		},
		{
			Name: "nonce",
			Type: "uint64",
		},
	},
	"SignedAuthentication": {
		{
			Name: "account",
			Type: "address",
		},
		{
			Name: "subAccountId",
			Type: "uint8",
		},
	},
}
