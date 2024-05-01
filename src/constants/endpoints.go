package constants

import "go100x/src/types"

const (
	GET_24H_TICKER_PRICE_CHANGE_STATISTICS types.Endpoint = "ticker/24hr"
	GET_PRODUCT                            types.Endpoint = "products/"
	GET_PRODUCT_BY_ID                      types.Endpoint = "products/product-by-id"
	GET_KLINE_DATA                         types.Endpoint = "uiKlines"
	LIST_PRODUCTS                          types.Endpoint = "products"
	ORDER_BOOK                             types.Endpoint = "depth"
	SERVER_TIME                            types.Endpoint = "time"
	APPROVE_REVOKE_SIGNER                  types.Endpoint = "approved-signers"
	LOGIN                                  types.Endpoint = "session/login"
	POST_WITHDRAW                          types.Endpoint = "withdraw"
	NEW_ORDER                              types.Endpoint = "order"
	CANCEL_REPLACE_ORDER                   types.Endpoint = "order/cancel-and-replace"
	CANCEL_ORDER                           types.Endpoint = "order"
	DELETE_CANCEL_ALL_OPEN_ORDERS          types.Endpoint = "openOrders"
	GET_SESSION_STATUS                     types.Endpoint = "session/status"
	GET_LOGOUT                             types.Endpoint = "session/logout"
	GET_SPOT_BALANCES                      types.Endpoint = "balances"
	GET_PERPETUAL_POSITION                 types.Endpoint = "positionRisk"
	GET_LIST_APPROVED_SIGNERS              types.Endpoint = "approved-signers"
	GET_LIST_OPEN_ORDERS                   types.Endpoint = "openOrders"
	GET_LIST_ORDERS                        types.Endpoint = "orders"
)
