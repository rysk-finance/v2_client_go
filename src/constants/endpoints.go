package constants

import "go100x/src/types"

const (
	ENDPOINT_GET_24H_TICKER_PRICE_CHANGE_STATISTICS types.Endpoint = "ticker/24hr"
	ENDPOINT_GET_PRODUCT                            types.Endpoint = "products/"
	ENDPOINT_GET_PRODUCT_BY_ID                      types.Endpoint = "products/product-by-id"
	ENDPOINT_GET_KLINE_DATA                         types.Endpoint = "uiKlines"
	ENDPOINT_LIST_PRODUCTS                          types.Endpoint = "products"
	ENDPOINT_ORDER_BOOK                             types.Endpoint = "depth"
	ENDPOINT_SERVER_TIME                            types.Endpoint = "time"
	ENDPOINT_APPROVE_REVOKE_SIGNER                  types.Endpoint = "approved-signers"
	ENDPOINT_POST_WITHDRAW                          types.Endpoint = "withdraw"
	ENDPOINT_NEW_ORDER                              types.Endpoint = "order"
	ENDPOINT_CANCEL_REPLACE_ORDER                   types.Endpoint = "order/cancel-and-replace"
	ENDPOINT_CANCEL_ORDER                           types.Endpoint = "order"
	ENDPOINT_CANCEL_ALL_OPEN_ORDERS                 types.Endpoint = "openOrders"
	ENDPOINT_GET_SPOT_BALANCES                      types.Endpoint = "balances"
	ENDPOINT_GET_PERPETUAL_POSITION                 types.Endpoint = "positionRisk"
	ENDPOINT_GET_LIST_APPROVED_SIGNERS              types.Endpoint = "approved-signers"
	ENDPOINT_GET_LIST_OPEN_ORDERS                   types.Endpoint = "openOrders"
	ENDPOINT_GET_LIST_ORDERS                        types.Endpoint = "orders"
)
