package types

type KlineDataRequest struct {
	Product   Product  // The product. Can be `constants.ETH_PERP`, `constants.BTC_PERP` or `constants.BLAST_PERP`.
	Interval  Interval // The interval. Can be `constants.M1`, `constants.M5`, `constants.M15`, `constants.M30`, `constants.H1`, `constants.H2`, `constants.H4`, `constants.H8`, `constants.D1`, `constants.D3` or `constants.W1`.
	StartTime int64    // Start timestamp in ms.
	EndTime   int64    // End timestamp in ms.
	Limit     int64    // Number of values to retrieve (max 1000).
}

type OrderBookRequest struct {
	Product     Product // The product. Can be `constants.ETH_PERP`, `constants.BTC_PERP` or `constants.BLAST_PERP`.
	Granularity int64   // The number of decimals to remove from prices.
	Limit       Limit   // The number of bids and asks to retrieve. Can be `constants.FIVE`, `constants.TEN` or `constants.TWENTY`.
}

type ApproveRevokeSignerRequest struct {
	ApprovedSigner string // The address of the account that will be an approvedSigner on the given subaccount.
	Nonce          int64  // The nonce. Suggest using the current UNIX timestamp in milliseconds.
}

type NewOrderRequest struct {
	Product     Product     // The product. Can be `constants.ETH_PERP`, `constants.BTC_PERP` or `constants.BLAST_PERP`.
	IsBuy       bool        // Whether the account is buying or selling.
	OrderType   OrderType   // The order type. Can be `constants.LIMIT`, `constants.LIMIT_MAKER`, `constants.MARKET`, `constants.STOP_LOSS`, `constants.STOP_LOSS_LIMIT`, `constants.TAKE_PROFIT` or `constants.TAKE_PROFIT_LIMIT`.
	TimeInForce TimeInForce // Order time in force. Can be `constants.GTC`, `constants.FOK` or `constants.IOC`.
	Price       string      // Price in wei (e18).
	Quantity    string      // Quantity in wei (e18).
	Expiration  int64       // UNIX timestamp (in ms) after which the order is no longer active.
	Nonce       int64       // The nonce. Suggest using the current UNIX timestamp in milliseconds.
}

type CancelOrderAndReplaceRequest struct {
	IdToCancel string          // ID of the order to be replaced.
	NewOrder   NewOrderRequest // The new order details to be used in the replacement.
}

type CancelOrderRequest struct {
	Product    Product // The product. Can be `constants.ETH_PERP`, `constants.BTC_PERP` or `constants.BLAST_PERP`.
	IdToCancel string  // The unique ID of the order you wish to cancel.
}
