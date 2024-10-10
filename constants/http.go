package constants

import "github.com/rysk-finance/v2_client_go/types"

var API_BASE_URL = map[types.Environment]string{
	ENVIRONMENT_MAINNET: "https://arbitrum-api.prod.rysk.finance/v1",
	ENVIRONMENT_TESTNET: "https://arbitrum-api.staging.rysk.finance/v1",
}

const (
	API_ENDPOINT_GET_24H_TICKER_PRICE_CHANGE_STATISTICS types.APIEndpoint = "/ticker/24hr"
	API_ENDPOINT_GET_PRODUCT                            types.APIEndpoint = "/products/"
	API_ENDPOINT_GET_PRODUCT_BY_ID                      types.APIEndpoint = "/products/product-by-id/"
	API_ENDPOINT_GET_KLINE_DATA                         types.APIEndpoint = "/uiKlines"
	API_ENDPOINT_LIST_PRODUCTS                          types.APIEndpoint = "/products"
	API_ENDPOINT_ORDER_BOOK                             types.APIEndpoint = "/depth"
	API_ENDPOINT_SERVER_TIME                            types.APIEndpoint = "/time"
	API_ENDPOINT_APPROVE_REVOKE_SIGNER                  types.APIEndpoint = "/approved-signers"
	API_ENDPOINT_WITHDRAW                               types.APIEndpoint = "/withdraw"
	API_ENDPOINT_NEW_ORDER                              types.APIEndpoint = "/order"
	API_ENDPOINT_CANCEL_REPLACE_ORDER                   types.APIEndpoint = "/order/cancel-and-replace"
	API_ENDPOINT_CANCEL_ORDER                           types.APIEndpoint = "/order"
	API_ENDPOINT_CANCEL_ALL_OPEN_ORDERS                 types.APIEndpoint = "/openOrders"
	API_ENDPOINT_GET_SPOT_BALANCES                      types.APIEndpoint = "/balances"
	API_ENDPOINT_GET_PERPETUAL_POSITION                 types.APIEndpoint = "/positionRisk"
	API_ENDPOINT_LIST_APPROVED_SIGNERS                  types.APIEndpoint = "/approved-signers"
	API_ENDPOINT_LIST_OPEN_ORDERS                       types.APIEndpoint = "/openOrders"
	API_ENDPOINT_LIST_ORDERS                            types.APIEndpoint = "/orders"
	API_ENDPOINT_ADD_REFEREE                            types.APIEndpoint = "/referral/add-referee"
)
