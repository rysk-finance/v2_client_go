package constants

import "github.com/eldief/go100x/types"

const (
	WS_METHOD_LIST_PRODUCTS                   types.WSMethod = "product.list"
	WS_METHOD_GET_PRODUCT                     types.WSMethod = "product.get"
	WS_METHOD_SERVER_TIME                     types.WSMethod = "time"
	WS_METHOD_LOGIN                           types.WSMethod = "session.login"
	WS_METHOD_SESSION_STATUS                  types.WSMethod = "session.status"
	WS_METHOD_SUB_ACCOUNT_LIST                types.WSMethod = "subaccount.list"
	WS_METHOD_WITHDRAW                        types.WSMethod = "withdraw"
	WS_METHOD_APPROVE_REVOKE_SIGNER           types.WSMethod = "signer.set"
	WS_METHOD_NEW_ORDER                       types.WSMethod = "order.place"
	WS_METHOD_ORDER_LIST                      types.WSMethod = "order.list"
	WS_METHOD_CANCEL_ORDER                    types.WSMethod = "order.cancel"
	WS_METHOD_CANCEL_ALL_OPEN_ORDERS          types.WSMethod = "order.cancelOpen"
	WS_METHOD_ORDER_BOOK_DEPTH                types.WSMethod = "depth"
	WS_METHOD_GET_PERPETUAL_POSITION          types.WSMethod = "position.perp.list"
	WS_METHOD_GET_SPOT_BALANCES               types.WSMethod = "position.spot.list"
	WS_METHOD_ACCOUNT_UPDATES                 types.WSMethod = "account.updates"
	WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE   types.WSMethod = "SUBSCRIBE"
	WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE types.WSMethod = "UNSUBSCRIBE"
)
