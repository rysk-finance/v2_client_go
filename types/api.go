package types

type KlineDataRequest struct {
	Product   *Product // The product. Can be `constants.PRODUCT_ETH_PERP`, `constants.PRODUCT_BTC_PERP` or `constants.PRODUCT_BLAST_PERP`.
	Interval  Interval // The interval. Can be `constants.INTERVAL_M1`, `constants.INTERVAL_5M`, `constants.INTERVAL_15M`, `constants.INTERVAL_30M`, `constants.INTERVAL_1H`, `constants.INTERVAL_2H`, `constants.INTERVAL_4H`, `constants.INTERVAL_8H`, `constants.INTERVAL_1D`, `constants.INTERVAL_3D` or `constants.INTERVAL_1W`.
	StartTime int64    // Start timestamp in ms.
	EndTime   int64    // End timestamp in ms.
	Limit     int64    // Number of values to retrieve (max 1000).
}

type OrderBookRequest struct {
	Product     *Product // The product. Can be `constants.PRODUCT_ETH_PERP`, `constants.PRODUCT_BTC_PERP` or `constants.PRODUCT_BLAST_PERP`.
	Granularity int64    // The number of decimals to remove from prices.
	Limit       Limit    // The number of bids and asks to retrieve. Can be `constants.FIVE`, `constants.TEN` or `constants.TWENTY`.
}

type ApproveRevokeSignerRequest struct {
	ApprovedSigner string // The address of the account that will be an approvedSigner on the given subaccount.
	Nonce          int64  // The nonce. Suggest using the current UNIX timestamp in milliseconds.
}

type NewOrderRequest struct {
	Product     *Product    // The product. Can be `constants.PRODUCT_ETH_PERP`, `constants.PRODUCT_BTC_PERP` or `constants.PRODUCT_BLAST_PERP`.
	IsBuy       bool        // Whether the account is buying or selling.
	OrderType   OrderType   // The order type. Can be `constants.ORDER_TYPE_LIMIT`, `constants.ORDER_TYPE_LIMIT_MAKER`, `constants.ORDER_TYPE_MARKET`, `constants.ORDER_TYPE_STOP_LOSS`, `constants.ORDER_TYPE_STOP_LOSS_LIMIT`, `constants.ORDER_TYPE_TAKE_PROFIT` or `constants.ORDER_TYPE_TAKE_PROFIT_LIMIT`.
	TimeInForce TimeInForce // Order time in force. Can be `constants.TIME_IN_FORCE_GTC`, `constants.TIME_IN_FORCE_FOK` or `constants.TIME_IN_FORCE_IOC`.
	Price       string      // Price in wei (e18).
	Quantity    string      // Quantity in wei (e18).
	Expiration  int64       // UNIX timestamp (in ms) after which the order is no longer active.
	Nonce       int64       // The nonce. Suggest using the current UNIX timestamp in milliseconds.
}

type CancelOrderAndReplaceRequest struct {
	IdToCancel string           // ID of the order to be replaced.
	NewOrder   *NewOrderRequest // The new order details to be used in the replacement.
}

type CancelOrderRequest struct {
	Product    *Product // The product. Can be `constants.PRODUCT_ETH_PERP`, `constants.PRODUCT_BTC_PERP` or `constants.PRODUCT_BLAST_PERP`.
	IdToCancel string   // The unique ID of the order you wish to cancel.
}

type ListOrdersRequest struct {
	Product *Product // The product. Can be `constants.PRODUCT_ETH_PERP`, `constants.PRODUCT_BTC_PERP` or `constants.PRODUCT_BLAST_PERP`.
	Ids     []string // IDs of specific orders you would like to retrieve.
}
