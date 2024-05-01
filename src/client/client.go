package go100x

import (
	"go100x/src/constants"
	"go100x/src/types"
	"go100x/src/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// NewClient creates a new `go100x.Client` instance.
// Initializes the client with the provided configuration.
func NewClient(config *types.ClientConfiguration) *types.Client {
	// Remove '0x' from private key.
	privateKey := strings.TrimPrefix(config.PrivateKey, "0x")

	// Return a new `go100x.Client`.
	return &types.Client{
		BaseUri:           constants.BASE_URI[config.Env],
		PrivateKey:        privateKey,
		Address:           utils.AddressFromPrivateKey(privateKey),
		SubAccountId:      int64(config.SubAccountId),
		HttpClient:        utils.GetHTTPClient(config.Timeout),
		VerifyingContract: constants.CIAO[config.Env],
		Domain: apitypes.TypedDataDomain{
			Name:              constants.NAME,
			Version:           constants.VERSION,
			ChainId:           constants.CHAIN_ID[config.Env],
			VerifyingContract: constants.VERIFIER[config.Env],
		},
	}
}

// Get24hrPriceChangeStatistics returns 24 hour rolling window price change statistics.
// If no `Product` is provided, ticker data for all assets will be returned.
func Get24hrPriceChangeStatistics(c *types.Client, product *types.Product) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		c.BaseUri+string(constants.GET_24H_TICKER_PRICE_CHANGE_STATISTICS),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Add query parameters and URL encode HTTP request.
	if product.Id != 0 && product.Symbol != "" {
		query := req.URL.Query()
		query.Add("symbol", string(product.Symbol))
		req.URL.RawQuery = query.Encode()
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// GetProduct returns details for a specific product by symbol
func GetProduct(c *types.Client, symbol string) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		c.BaseUri+string(constants.GET_PRODUCT)+symbol,
		nil,
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// GetProductById returns details for a specific product by id.
func GetProductById(c *types.Client, id int64) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		c.BaseUri+string(constants.GET_PRODUCT_BY_ID)+strconv.FormatInt(id, 10),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// GetKlineData returns Kline/candlestick bars for a symbol. Klines are uniquely identified by interval(timeframe) and startTime.
func GetKlineData(c *types.Client, params *types.KlineDataRequest) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		string(c.BaseUri)+string(constants.GET_KLINE_DATA),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Add query parameters and URL encode HTTP request.
	query := req.URL.Query()
	query.Add("symbol", string(params.Product.Symbol))
	if params.Interval != "" {
		query.Add("interval", string(params.Interval))
	}
	if params.StartTime != 0 {
		query.Add("startTime", strconv.FormatInt(params.StartTime, 10))
	}
	if params.EndTime != 0 {
		query.Add("endTime", strconv.FormatInt(params.EndTime, 10))
	}
	if params.Limit != 0 {
		query.Add("limit", strconv.FormatInt(params.Limit, 10))
	}
	req.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// ListProducts returns a list of products available to trade.
func ListProducts(c *types.Client) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		string(c.BaseUri)+string(constants.LIST_PRODUCTS),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// OrderBook returns bids and asks for a market.
func OrderBook(c *types.Client, params *types.OrderBookRequest) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		string(c.BaseUri)+string(constants.ORDER_BOOK),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Add query parameters and URL encode HTTP request.
	query := req.URL.Query()
	query.Add("symbol", string(params.Product.Symbol))
	if params.Granularity != 0 {
		query.Add("granularity", strconv.FormatInt(params.Granularity, 10))
	}
	if params.Limit != 0 {
		query.Add("limit", strconv.FormatInt(int64(params.Limit), 10))
	}
	req.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// ServerTime returns current server time.
func ServerTime(c *types.Client) (string, error) {
	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		string(c.BaseUri)+string(constants.SERVER_TIME),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// ApproveSigner approves a Signer for a `SubAccount`.
func ApproveSigner(c *types.Client, params *types.ApproveRevokeSignerRequest) (string, error) {
	return approveRevokeSigner(c, params, true)
}

// RevokeSigner revokes a Signer for a `SubAccount`.
func RevokeSigner(c *types.Client, params *types.ApproveRevokeSignerRequest) (string, error) {
	return approveRevokeSigner(c, params, false)
}

// approveRevokeSigner approves or revoke a signer for a `SubAccount`.
func approveRevokeSigner(c *types.Client, params *types.ApproveRevokeSignerRequest, isApproved bool) (string, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(c, constants.APPROVE_SIGNER, &struct {
		Account        string `json:"account"`
		SubAccountId   string `json:"subAccountId"`
		ApprovedSigner string `json:"approvedSigner"`
		IsApproved     bool   `json:"isApproved"`
		Nonce          string `json:"nonce"`
	}{
		Account:        c.Address,
		SubAccountId:   strconv.FormatInt(c.SubAccountId, 10),
		ApprovedSigner: params.ApprovedSigner,
		IsApproved:     isApproved,
		Nonce:          strconv.FormatInt(params.Nonce, 10),
	})
	if err != nil {
		return "", err
	}

	// Create HTTP request.
	req, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(c.BaseUri)+string(constants.APPROVE_REVOKE_SIGNER),
		&struct {
			Account        string
			SubAccountId   int64
			Signature      string
			ApprovedSigner string
			Nonce          int64
			IsApproved     bool
		}{
			Account:        c.Address,
			SubAccountId:   c.SubAccountId,
			ApprovedSigner: params.ApprovedSigner,
			Nonce:          params.Nonce,
			Signature:      signature,
			IsApproved:     isApproved,
		},
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// Log user in returning a set-cookie header that will need to be attached to authenticated requests to access private endpoints.
func Login(c *types.Client) (string, error) {
	timestamp := time.Now().UnixMilli()

	// Generate EIP712 signature.
	signature, err := utils.SignMessage(c, constants.LOGIN_MESSAGE, &struct {
		Account   string `json:"account"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}{
		Account:   c.Address,
		Message:   "I want to log into 100x.finance",
		Timestamp: strconv.FormatInt(timestamp, 10),
	})
	if err != nil {
		return "", err
	}

	// Create HTTP request.
	req, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(c.BaseUri)+string(constants.LOGIN),
		&struct {
			Account   string
			Message   string
			Timestamp int64
			Signature string
		}{
			Account:   c.Address,
			Message:   "I want to log into 100x.finance",
			Timestamp: timestamp,
			Signature: signature,
		},
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// TODO
// Withdraw USDB from 100x account.
// func Withdraw(c *types.Client100x) (string, error) {

// }

// Returns spot balances for sub account id.
func GetSpotBalances(c *types.Client) (string, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(c, constants.SIGNED_AUTHENTICATION, &struct {
		Account      string `json:"account"`
		SubAccountId string `json:"subAccountId"`
	}{
		Account:      c.Address,
		SubAccountId: strconv.FormatInt(c.SubAccountId, 10),
	})
	if err != nil {
		return "", err
	}

	// Create HTTP request.
	req, err := http.NewRequest(
		http.MethodGet,
		string(c.BaseUri)+string(constants.GET_SPOT_BALANCES),
		nil,
	)
	if err != nil {
		return "", err
	}

	// Add query parameters and URL encode HTTP request.
	query := req.URL.Query()
	query.Add("account", c.Address)
	query.Add("subAccountId", strconv.FormatInt(c.SubAccountId, 10))
	query.Add("signature", signature)
	req.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// Create a new order on the SubAccount.
func NewOrder(c *types.Client, params *types.NewOrderRequest) (string, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(c, constants.ORDER, &struct {
		Account      string `json:"account"`
		SubAccountId string `json:"subAccountId"`
		ProductId    string `json:"productId"`
		IsBuy        bool   `json:"isBuy"`
		OrderType    string `json:"orderType"`
		TimeInForce  string `json:"timeInForce"`
		Expiration   string `json:"expiration"`
		Price        string `json:"price"`
		Quantity     string `json:"quantity"`
		Nonce        string `json:"nonce"`
	}{
		Account:      c.Address,
		SubAccountId: strconv.FormatInt(c.SubAccountId, 10),
		ProductId:    strconv.FormatInt(params.Product.Id, 10),
		IsBuy:        params.IsBuy,
		OrderType:    strconv.FormatInt(int64(params.OrderType), 10),
		TimeInForce:  strconv.FormatInt(int64(params.TimeInForce), 10),
		Expiration:   strconv.FormatInt(params.Expiration, 10),
		Price:        params.Price,
		Quantity:     params.Quantity,
		Nonce:        strconv.FormatInt(params.Nonce, 10),
	})
	if err != nil {
		return "", err
	}

	// Create HTTP request.
	req, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(c.BaseUri)+string(constants.NEW_ORDER),
		&struct {
			Account      string
			SubAccountId int64
			ProductId    int64
			IsBuy        bool
			OrderType    int64
			TimeInForce  int64
			Expiration   int64
			Price        string
			Quantity     string
			Nonce        int64
			Signature    string
		}{
			Account:      c.Address,
			SubAccountId: c.SubAccountId,
			ProductId:    params.Product.Id,
			IsBuy:        params.IsBuy,
			OrderType:    int64(params.OrderType),
			TimeInForce:  int64(params.TimeInForce),
			Expiration:   params.Expiration,
			Price:        params.Price,
			Quantity:     params.Quantity,
			Nonce:        params.Nonce,
			Signature:    signature,
		},
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}

// CancelOrderAndReplace cancel an order and create a new order on the `SubAccount`.
func CancelOrderAndReplace(c *types.Client, params *types.CancelOrderAndReplaceRequest) (string, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(c, constants.ORDER, &struct {
		Account      string `json:"account"`
		SubAccountId string `json:"subAccountId"`
		ProductId    string `json:"productId"`
		IsBuy        bool   `json:"isBuy"`
		OrderType    string `json:"orderType"`
		TimeInForce  string `json:"timeInForce"`
		Expiration   string `json:"expiration"`
		Price        string `json:"price"`
		Quantity     string `json:"quantity"`
		Nonce        string `json:"nonce"`
	}{
		Account:      c.Address,
		SubAccountId: strconv.FormatInt(c.SubAccountId, 10),
		ProductId:    strconv.FormatInt(params.NewOrder.Product.Id, 10),
		IsBuy:        params.NewOrder.IsBuy,
		OrderType:    strconv.FormatInt(int64(params.NewOrder.OrderType), 10),
		TimeInForce:  strconv.FormatInt(int64(params.NewOrder.TimeInForce), 10),
		Expiration:   strconv.FormatInt(params.NewOrder.Expiration, 10),
		Price:        params.NewOrder.Price,
		Quantity:     params.NewOrder.Quantity,
		Nonce:        strconv.FormatInt(params.NewOrder.Nonce, 10),
	})
	if err != nil {
		return "", err
	}

	// Create HTTP request.
	req, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(c.BaseUri)+string(constants.CANCEL_REPLACE_ORDER),
		&struct {
			IdToCancel string
			NewOrder   interface{}
		}{
			IdToCancel: params.IdToCancel,
			NewOrder: struct {
				Account      string
				SubAccountId int64
				ProductId    int64
				IsBuy        bool
				OrderType    int64
				TimeInForce  int64
				Expiration   int64
				Price        string
				Quantity     string
				Nonce        int64
				Signature    string
			}{
				Account:      c.Address,
				SubAccountId: c.SubAccountId,
				ProductId:    params.NewOrder.Product.Id,
				IsBuy:        params.NewOrder.IsBuy,
				OrderType:    int64(params.NewOrder.OrderType),
				TimeInForce:  int64(params.NewOrder.TimeInForce),
				Expiration:   params.NewOrder.Expiration,
				Price:        params.NewOrder.Price,
				Quantity:     params.NewOrder.Quantity,
				Nonce:        params.NewOrder.Nonce,
				Signature:    signature,
			},
		},
	)
	if err != nil {
		return "", err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(c.HttpClient, req)
}
