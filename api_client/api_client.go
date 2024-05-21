package api_client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/eldief/go100x/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// Go100XAPIClient configuration
type Go100XAPIClientConfiguration struct {
	Env          types.Environment // `constants.ENVIRONMENT_TESTNET` or `constants.ENVIRONMENT_MAINNET`.
	PrivateKey   string            // e.g. `0x2638b4...` or `2638b4...`.
	RpcUrl       string            // e.g. `https://sepolia.blast.io` or `https://rpc.blastblockchain.com`.
	SubAccountId uint8             // ID of the subaccount to use.
}

// 100x API client.
type Go100XAPIClient struct {
	env              types.Environment
	baseUrl          string
	privateKeyString string
	addressString    string
	privateKey       *ecdsa.PrivateKey
	address          common.Address
	ciao             common.Address
	usdb             common.Address
	domain           apitypes.TypedDataDomain
	SubAccountId     int64
	HttpClient       *http.Client
	EthClient        types.IEthClient
}

// NewGo100XAPIClient creates a new `Go100XAPIClient` instance.
// Initializes the client with the provided configuration.
func NewGo100XAPIClient(config *Go100XAPIClientConfiguration) (*Go100XAPIClient, error) {
	// Remove '0x' from private key.
	privateKeyString := strings.TrimPrefix(config.PrivateKey, "0x")

	// Get ecdsa.PrivateKey.
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Instanciate new Ethereum Client.
	client, err := ethclient.Dial(config.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	// Return a new `go100x.Client`.
	return &Go100XAPIClient{
		env:              config.Env,
		baseUrl:          constants.API_BASE_URL[config.Env],
		privateKeyString: privateKeyString,
		addressString:    utils.AddressFromPrivateKey(privateKeyString),
		privateKey:       privateKey,
		address:          common.HexToAddress(utils.AddressFromPrivateKey(privateKeyString)),
		ciao:             common.HexToAddress(constants.CIAO_ADDRESS[config.Env]),
		usdb:             common.HexToAddress(constants.USDB_ADDRESS[config.Env]),
		domain: apitypes.TypedDataDomain{
			Name:              constants.DOMAIN_NAME,
			Version:           constants.DOMAIN_VERSION,
			ChainId:           constants.CHAIN_ID[config.Env],
			VerifyingContract: constants.VERIFIER_ADDRESS[config.Env],
		},
		SubAccountId: int64(config.SubAccountId),
		HttpClient:   utils.GetHTTPClient(10 * time.Second),
		EthClient:    client,
	}, nil
}

// Get24hrPriceChangeStatistics returns 24 hour rolling window price change statistics.
// If no `Product` is provided, ticker data for all assets will be returned.
func (go100XClient *Go100XAPIClient) Get24hrPriceChangeStatistics(product *types.Product) (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		go100XClient.baseUrl+string(constants.API_ENDPOINT_GET_24H_TICKER_PRICE_CHANGE_STATISTICS),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	if product.Id != 0 && product.Symbol != "" {
		query := request.URL.Query()
		query.Add("symbol", string(product.Symbol))
		request.URL.RawQuery = query.Encode()
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// GetProduct returns details for a specific product by symbol
func (go100XClient *Go100XAPIClient) GetProduct(symbol string) (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		go100XClient.baseUrl+string(constants.API_ENDPOINT_GET_PRODUCT)+symbol,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// GetProductById returns details for a specific product by id.
func (go100XClient *Go100XAPIClient) GetProductById(id int64) (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		go100XClient.baseUrl+string(constants.API_ENDPOINT_GET_PRODUCT_BY_ID)+strconv.FormatInt(id, 10),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// GetKlineData returns Kline/Candlestick bars for a symbol. Klines are uniquely identified by interval(timeframe) and startTime.
func (go100XClient *Go100XAPIClient) GetKlineData(params *types.KlineDataRequest) (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_GET_KLINE_DATA),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
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
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ListProducts returns a list of products available to trade.
func (go100XClient *Go100XAPIClient) ListProducts() (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_LIST_PRODUCTS),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// OrderBook returns bids and asks for a market.
func (go100XClient *Go100XAPIClient) OrderBook(params *types.OrderBookRequest) (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_ORDER_BOOK),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
	query.Add("symbol", string(params.Product.Symbol))
	if params.Granularity != 0 {
		query.Add("granularity", strconv.FormatInt(params.Granularity, 10))
	}
	if params.Limit != 0 {
		query.Add("limit", strconv.FormatInt(int64(params.Limit), 10))
	}
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ServerTime returns current server time.
func (go100XClient *Go100XAPIClient) ServerTime() (*http.Response, error) {
	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_SERVER_TIME),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ApproveSigner approves a Signer for a `SubAccount`.
func (go100XClient *Go100XAPIClient) ApproveSigner(params *types.ApproveRevokeSignerRequest) (*http.Response, error) {
	return go100XClient.approveRevokeSigner(params, true)
}

// RevokeSigner revokes a Signer for a `SubAccount`.
func (go100XClient *Go100XAPIClient) RevokeSigner(params *types.ApproveRevokeSignerRequest) (*http.Response, error) {
	return go100XClient.approveRevokeSigner(params, false)
}

// approveRevokeSigner approves or revoke a signer for a `SubAccount`.
func (go100XClient *Go100XAPIClient) approveRevokeSigner(params *types.ApproveRevokeSignerRequest, isApproved bool) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_APPROVE_SIGNER,
		&struct {
			Account        string `json:"account"`
			SubAccountId   string `json:"subAccountId"`
			ApprovedSigner string `json:"approvedSigner"`
			IsApproved     bool   `json:"isApproved"`
			Nonce          string `json:"nonce"`
		}{
			Account:        go100XClient.addressString,
			SubAccountId:   strconv.FormatInt(go100XClient.SubAccountId, 10),
			ApprovedSigner: params.ApprovedSigner,
			IsApproved:     isApproved,
			Nonce:          strconv.FormatInt(params.Nonce, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_APPROVE_REVOKE_SIGNER),
		&struct {
			Account        string
			SubAccountId   int64
			Signature      string
			ApprovedSigner string
			Nonce          int64
			IsApproved     bool
		}{
			Account:        go100XClient.addressString,
			SubAccountId:   go100XClient.SubAccountId,
			ApprovedSigner: params.ApprovedSigner,
			Nonce:          params.Nonce,
			Signature:      signature,
			IsApproved:     isApproved,
		},
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// Withdraw USDB from 100x account.
func (go100XClient *Go100XAPIClient) Withdraw(params *types.WithdrawRequest) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_WITHDRAW,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
			Asset        string `json:"asset"`
			Quantity     string `json:"quantity"`
			Nonce        string `json:"nonce"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
			Asset:        constants.USDB_ADDRESS[go100XClient.env],
			Quantity:     params.Quantity,
			Nonce:        strconv.FormatInt(params.Nonce, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_WITHDRAW),
		&struct {
			Account      string
			SubAccountId int64
			Asset        string
			Quantity     string
			Nonce        int64
			Signature    string
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			Asset:        constants.USDB_ADDRESS[go100XClient.env],
			Quantity:     params.Quantity,
			Nonce:        params.Nonce,
			Signature:    signature,
		},
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// NewOrder creates a new order on the `SubAccount`.
func (go100XClient *Go100XAPIClient) NewOrder(params *types.NewOrderRequest) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_ORDER,
		&struct {
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
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
			ProductId:    strconv.FormatInt(params.Product.Id, 10),
			IsBuy:        params.IsBuy,
			OrderType:    strconv.FormatInt(int64(params.OrderType), 10),
			TimeInForce:  strconv.FormatInt(int64(params.TimeInForce), 10),
			Expiration:   strconv.FormatInt(params.Expiration, 10),
			Price:        params.Price,
			Quantity:     params.Quantity,
			Nonce:        strconv.FormatInt(params.Nonce, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_NEW_ORDER),
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
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
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
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// CancelOrderAndReplace cancel an order and create a new order on the `SubAccount`.
func (go100XClient *Go100XAPIClient) CancelOrderAndReplace(params *types.CancelOrderAndReplaceRequest) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_ORDER,
		&struct {
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
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
			ProductId:    strconv.FormatInt(params.NewOrder.Product.Id, 10),
			IsBuy:        params.NewOrder.IsBuy,
			OrderType:    strconv.FormatInt(int64(params.NewOrder.OrderType), 10),
			TimeInForce:  strconv.FormatInt(int64(params.NewOrder.TimeInForce), 10),
			Expiration:   strconv.FormatInt(params.NewOrder.Expiration, 10),
			Price:        params.NewOrder.Price,
			Quantity:     params.NewOrder.Quantity,
			Nonce:        strconv.FormatInt(params.NewOrder.Nonce, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_CANCEL_REPLACE_ORDER),
		&struct {
			IdToCancel string
			NewOrder   interface{}
		}{
			IdToCancel: params.IdToCancel,
			NewOrder: &struct {
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
				Account:      go100XClient.addressString,
				SubAccountId: go100XClient.SubAccountId,
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
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// CancelOrder cancel an active order on the `SubAccount`.
func (go100XClient *Go100XAPIClient) CancelOrder(params *types.CancelOrderRequest) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_CANCEL_ORDER,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
			ProductId    string `json:"productId"`
			OrderId      string `json:"orderId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
			ProductId:    strconv.FormatInt(params.Product.Id, 10),
			OrderId:      params.IdToCancel,
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodDelete,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_CANCEL_ORDER),
		&struct {
			Account      string
			SubAccountId int64
			ProductId    int64
			OrderId      string
			Signature    string
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			ProductId:    params.Product.Id,
			OrderId:      params.IdToCancel,
			Signature:    signature,
		},
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// CancelAllOpenOrders cancel all active orders on a product.
func (go100XClient *Go100XAPIClient) CancelAllOpenOrders(product *types.Product) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_CANCEL_ORDERS,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
			ProductId    string `json:"productId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
			ProductId:    strconv.FormatInt(product.Id, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodDelete,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_CANCEL_ALL_OPEN_ORDERS),
		&struct {
			Account      string
			SubAccountId int64
			ProductId    int64
			Signature    string
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			ProductId:    product.Id,
			Signature:    signature,
		},
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// GetSpotBalances returns spot balances for sub account id.
func (go100XClient *Go100XAPIClient) GetSpotBalances() (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_SIGNED_AUTHENTICATION,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_GET_SPOT_BALANCES),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
	query.Add("account", go100XClient.addressString)
	query.Add("subAccountId", strconv.FormatInt(go100XClient.SubAccountId, 10))
	query.Add("signature", signature)
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// GetPerpetualPosition returns perpetual position for sub account id.
func (go100XClient *Go100XAPIClient) GetPerpetualPosition(product *types.Product) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_SIGNED_AUTHENTICATION,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_GET_PERPETUAL_POSITION),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
	query.Add("account", go100XClient.addressString)
	query.Add("subAccountId", strconv.FormatInt(go100XClient.SubAccountId, 10))
	query.Add("symbol", product.Symbol)
	query.Add("signature", signature)
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ListApprovedSigners returns a list of all approved signers for a `SubAccount`.
func (go100XClient *Go100XAPIClient) ListApprovedSigners() (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_SIGNED_AUTHENTICATION,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_LIST_APPROVED_SIGNERS),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
	query.Add("account", go100XClient.addressString)
	query.Add("subAccountId", strconv.FormatInt(go100XClient.SubAccountId, 10))
	query.Add("signature", signature)
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ListOpenOrders returns all open orders on the `SubAccount` per product.
func (go100XClient *Go100XAPIClient) ListOpenOrders(product *types.Product) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_SIGNED_AUTHENTICATION,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_LIST_OPEN_ORDERS),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
	query.Add("account", go100XClient.addressString)
	query.Add("subAccountId", strconv.FormatInt(go100XClient.SubAccountId, 10))
	query.Add("symbol", product.Symbol)
	query.Add("signature", signature)
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ListOrders returns all orders on the `SubAccount` per product.
func (go100XClient *Go100XAPIClient) ListOrders(params *types.ListOrdersRequest) (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_SIGNED_AUTHENTICATION,
		&struct {
			Account      string `json:"account"`
			SubAccountId string `json:"subAccountId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: strconv.FormatInt(go100XClient.SubAccountId, 10),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := http.NewRequest(
		http.MethodGet,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_LIST_ORDERS),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add query parameters and URL encode HTTP request.
	query := request.URL.Query()
	query.Add("account", go100XClient.addressString)
	query.Add("subAccountId", strconv.FormatInt(go100XClient.SubAccountId, 10))
	query.Add("symbol", params.Product.Symbol)
	for _, id := range params.Ids {
		query.Add("ids", id)
	}
	query.Add("signature", signature)
	request.URL.RawQuery = query.Encode()

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}

// ApproveUSDB approves 100x to spend USDB on your behalf.
func (go100XClient *Go100XAPIClient) ApproveUSDB(ctx context.Context, amount *big.Int) (*geth_types.Transaction, error) {
	// Parse ABI
	parsedABI, _ := abi.JSON(strings.NewReader(constants.ERC20_ABI))

	// Pack transaction data
	data, _ := parsedABI.Pack("approve", go100XClient.ciao, amount)

	// Get transaction parameters
	nonce, gasPrice, chainID, gasLimit, err := utils.GetTransactionParams(ctx, go100XClient.EthClient, go100XClient.privateKey, &go100XClient.address, &go100XClient.usdb, &data)
	if err != nil {
		return nil, err
	}

	// Create a new transaction
	tx := geth_types.NewTransaction(nonce, go100XClient.usdb, big.NewInt(0), gasLimit, gasPrice, data)

	// Sign transaction
	signedTx, _ := geth_types.SignTx(tx, geth_types.NewEIP155Signer(chainID), go100XClient.privateKey)

	// Send transaction
	err = go100XClient.EthClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx, nil
}

// DepositUSDB sends USDB to 100x.
func (go100XClient *Go100XAPIClient) DepositUSDB(ctx context.Context, amount *big.Int) (*geth_types.Transaction, error) {
	// Parse ABI
	parsedABI, _ := abi.JSON(strings.NewReader(constants.CIAO_ABI))

	// Pack transaction data
	data, _ := parsedABI.Pack("deposit", go100XClient.address, uint8(go100XClient.SubAccountId), amount, go100XClient.usdb)

	// Get transaction parameters
	nonce, gasPrice, chainID, gasLimit, err := utils.GetTransactionParams(ctx, go100XClient.EthClient, go100XClient.privateKey, &go100XClient.address, &go100XClient.ciao, &data)
	if err != nil {
		return nil, err
	}

	// Create a new transaction
	tx := geth_types.NewTransaction(nonce, go100XClient.ciao, big.NewInt(0), gasLimit, gasPrice, data)

	// Sign transaction
	signedTx, _ := geth_types.SignTx(tx, geth_types.NewEIP155Signer(chainID), go100XClient.privateKey)

	// Send transaction
	err = go100XClient.EthClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx, nil
}

// WaitTransaction waits for a transaction to be mined and returns its receipt.
func (go100XClient *Go100XAPIClient) WaitTransaction(ctx context.Context, transaction *geth_types.Transaction) (*geth_types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, go100XClient.EthClient, transaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
