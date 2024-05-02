package constants

import (
	"github.com/eldief/go100x/types"
)

const (
	DOMAIN_NAME    string = "100x"
	DOMAIN_VERSION string = "0.0.0"
)

const (
	PRIMARY_TYPE_LOGIN_MESSAGE         types.PrimaryType = "LoginMessage"
	PRIMARY_TYPE_ORDER                 types.PrimaryType = "Order"
	PRIMARY_TYPE_CANCEL_ORDER          types.PrimaryType = "CancelOrder"
	PRIMARY_TYPE_CANCEL_ORDERS         types.PrimaryType = "CancelOrders"
	PRIMARY_TYPE_APPROVE_SIGNER        types.PrimaryType = "ApproveSigner"
	PRIMARY_TYPE_WITHDRAW              types.PrimaryType = "Withdraw"
	PRIMARY_TYPE_SIGNED_AUTHENTICATION types.PrimaryType = "SignedAuthentication"
)
