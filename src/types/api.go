package types

type Get24hrPriceChangeStatisticsResponse struct {
	Results []PriceChangeStatisticsResponse `json:"results"`
}

type PriceChangeStatisticsResponse struct {
	FundingRateHourly  string
	FundingRateYearly  string
	High               string
	Low                string
	MarkPrice          string
	NextFundingTime    string
	OpenInterest       string
	OraclePrice        string
	PriceChange        string
	PriceChangePercent string
	ProductId          int
	ProductSymbol      string
	Volume             string
}

type KlineDataRequest struct {
	Product   Product
	Interval  Interval
	StartTime int64
	EndTime   int64
	Limit     int64
}

type OrderBookRequest struct {
	Product     Product
	Granularity int64
	Limit       Limit
}

type ApproveRevokeSignerRequest struct {
	Account        string
	SubAccountId   int64
	ApprovedSigner string
	Nonce          int64
	IsApproved     bool
	Signature      string
}

type LoginRequest struct {
	Account   string
	Message   string
	Timestamp int64
	Signature string
}
