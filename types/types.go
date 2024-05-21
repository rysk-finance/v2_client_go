package types

type Environment string
type MarginAsset string
type APIEndpoint string
type Interval string
type Limit int64
type OrderType int64
type TimeInForce int64

type Product struct {
	Symbol string // `constants.PRODUCT_ETH_PERP.Symbol`, `constants.PRODUCT_BTC_PERP.Symbol` or `constants.PRODUCT_BLAST_PERP.Symbol`.
	Id     int64  // `constants.PRODUCT_ETH_PERP.Id`, `constants.PRODUCT_BTC_PERP.Id` or `constants.PRODUCT_BLAST_PERP.Id`.
}
