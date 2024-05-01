package constants

import (
	"go100x/src/types"
)

const (
	NAME                  string            = "100x"
	VERSION               string            = "0.0.0"
	LOGIN_MESSAGE         types.PrimaryType = "LoginMessage"
	ORDER                 types.PrimaryType = "Order"
	CANCEL_ORDER          types.PrimaryType = "CancelOrder"
	CANCEL_ORDERS         types.PrimaryType = "CancelOrders"
	APPROVE_SIGNER        types.PrimaryType = "ApproveSigner"
	WITHDRAW              types.PrimaryType = "Withdraw"
	SIGNED_AUTHENTICATION types.PrimaryType = "SignedAuthentication"
)
