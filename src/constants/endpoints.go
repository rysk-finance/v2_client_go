package constants

import "go100x/src/types"

const (
	GET_24H_TICKER_PRICE_CHANGE_STATISTICS types.Endpoint = "ticker/24hr"
	GET_PRODUCT                            types.Endpoint = "products/"
	GET_PRODUCT_BY_ID                      types.Endpoint = "products/product-by-id/"
	GET_KLINE_DATA                         types.Endpoint = "uiKlines"
	GET_LIST_PRODUCTS                      types.Endpoint = "products"
	GET_ORDER_BOOK                         types.Endpoint = "depth"
	GET_SERVER_TIME                        types.Endpoint = "time"
	POST_APPROVE_REVOKE_SIGNER             types.Endpoint = "approved-signers"
	POST_LOGIN                             types.Endpoint = "session/login"
	POST_WITHDRAW                          types.Endpoint = "withdraw"
	POST_NEW_ORDER                         types.Endpoint = "order"
	POST_CANCEL_REPLACE_ORDER              types.Endpoint = "order/cancel-and-replace"
	DELETE_CANCEL_ORDER                    types.Endpoint = "order/"
	DELETE_CANCEL_ALL_OPEN_ORDERS          types.Endpoint = "openOrders"
	GET_SESSION_STATUS                     types.Endpoint = "session/status"
	GET_LOGOUT                             types.Endpoint = "session/logout"
	GET_SPOT_BALANCES                      types.Endpoint = "balances"
	GET_PERPETUAL_POSITION                 types.Endpoint = "positionRisk"
	GET_LIST_APPROVED_SIGNERS              types.Endpoint = "approved-signers/"
	GET_LIST_OPEN_ORDERS                   types.Endpoint = "openOrders/"
	GET_LIST_ORDERS                        types.Endpoint = "orders"
)
