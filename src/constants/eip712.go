package constants

import (
	"go100x/src/types"
)

const (
	NAME    string = "100x"
	VERSION string = "0.0.0"
)

var PRIMARY_TYPE = map[types.Endpoint]string{
	POST_LOGIN:                    "LoginMessage",
	POST_NEW_ORDER:                "Order",
	POST_CANCEL_REPLACE_ORDER:     "Order",
	DELETE_CANCEL_ALL_OPEN_ORDERS: "CancelOrders",
	DELETE_CANCEL_ORDER:           "CancelOrder",
	POST_APPROVE_REVOKE_SIGNER:    "ApproveSigner",
	POST_WITHDRAW:                 "Withdraw",
	GET_SPOT_BALANCES:             "SignedAuthentication",
	GET_PERPETUAL_POSITION:        "SignedAuthentication",
	GET_LIST_APPROVED_SIGNERS:     "SignedAuthentication",
	GET_LIST_OPEN_ORDERS:          "SignedAuthentication",
	GET_LIST_ORDERS:               "SignedAuthentication",
}
