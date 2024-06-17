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

// Go100XAPIClientConfiguration holds the configuration for the Go100X API client.
type Go100XAPIClientConfiguration struct {
	Env          types.Environment // `constants.ENVIRONMENT_TESTNET` or `constants.ENVIRONMENT_MAINNET`.
	PrivateKey   string            // Private key as a string, e.g., `0x2638b4...` or `2638b4...`.
	RpcUrl       string            // RPC URL of the Ethereum client, e.g., `https://sepolia.blast.io` or `https://rpc.blastblockchain.com`.
	SubAccountId uint8             // ID of the subaccount to use.
}

// Go100XAPIClient is the main client for interacting with the 100X API.
type Go100XAPIClient struct {
	env              types.Environment        // Environment (testnet or mainnet).
	baseUrl          string                   // Base URL for the API.
	privateKeyString string                   // Private key as a string.
	addressString    string                   // Address derived from the private key.
	privateKey       *ecdsa.PrivateKey        // ECDSA private key.
	address          common.Address           // Common address derived from the private key.
	ciao             common.Address           // Address for the CIAO contract.
	usdb             common.Address           // Address for the USDB contract.
	domain           apitypes.TypedDataDomain // Typed data domain for EIP-712.
	SubAccountId     int64                    // Subaccount ID.
	HttpClient       *http.Client             // HTTP client for making requests.
	EthClient        types.IEthClient         // Ethereum client for interacting with the blockchain.
}

// NewGo100XAPIClient creates a new Go100XAPIClient instance.
// Initializes the client with the provided configuration.
//
// Parameters:
//   - config: A pointer to Go100XAPIClientConfiguration containing the configuration settings.
//
// Returns:
//   - A pointer to Go100XAPIClient.
//   - An error if initialization fails.
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
	apiClient := &Go100XAPIClient{
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
	}

	apiClient.addReferee()
	return apiClient, nil
}

// Get24hrPriceChangeStatistics returns 24-hour rolling window price change statistics.
// These statistics do not reflect the UTC day, but rather a 24-hour rolling window for the previous 24 hours.
// If no `Product` is provided, ticker data for all assets will be returned.
//
// Parameters:
//   - product: A pointer to a Product struct for which the statistics are being retrieved.
//     If nil, ticker data for all assets will be returned.
//
// Returns:
//   - A pointer to an http.Response containing the response from the server.
//   - An error if the request fails.
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

// GetProduct returns details for a specific product by its symbol.
//
// Parameters:
//   - symbol: The symbol of the product for which details are being retrieved.
//
// Returns:
//   - A pointer to an http.Response containing the response from the server with product details.
//   - An error if the request fails.
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

// GetProductById retrieves details for a specific product by its unique identifier.
//
// Parameters:
//   - id: The ID of the product.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// GetKlineData retrieves Kline/Candlestick bars for a symbol based on the provided parameters.
//
// Parameters:
//   - params: A pointer to a KlineDataRequest struct containing the parameters for the request,
//     including symbol, interval (timeframe), startTime, and optional endTime.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// ListProducts retrieves a list of products available for trading on the platform.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// OrderBook retrieves the order book (bids and asks) for a specific market.
//
// Parameters:
//   - params: A pointer to an OrderBookRequest struct containing the parameters for the request.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// ServerTime retrieves the current server time from the API.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// ApproveSigner approves a Signer for a SubAccount. This operation allows the specified
// Signer to sign transactions on behalf of the SubAccount.
//
// Params:
//   - params: An instance of types.ApproveRevokeSignerRequest containing the necessary
//     parameters for approving the signer.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
func (go100XClient *Go100XAPIClient) ApproveSigner(params *types.ApproveRevokeSignerRequest) (*http.Response, error) {
	return go100XClient.approveRevokeSigner(params, true)
}

// RevokeSigner revokes a Signer for a SubAccount. This operation disables the specified
// Signer from signing transactions on behalf of the SubAccount.
//
// Params:
//   - params: An instance of types.ApproveRevokeSignerRequest containing the necessary
//     parameters for revoking the signer.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
func (go100XClient *Go100XAPIClient) RevokeSigner(params *types.ApproveRevokeSignerRequest) (*http.Response, error) {
	return go100XClient.approveRevokeSigner(params, false)
}

// approveRevokeSigner approves or revokes a signer for a `SubAccount`.
//
// This function either approves or revokes a signer for a specific `SubAccount`
// based on the value of `isApproved`.
//
// Parameters:
//   - params: The parameters containing the request details, including signer information.
//   - isApproved: Boolean flag indicating whether to approve (true) or revoke (false) the signer.
//
// Returns:
//   - *http.Response: The HTTP response received from the API after the operation.
//   - error: An error if the operation encountered any issues.
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
			Account        string `json:"account"`
			SubAccountId   int64  `json:"subAccountId"`
			Signature      string `json:"signature"`
			ApprovedSigner string `json:"approvedSigner"`
			Nonce          int64  `json:"nonce"`
			IsApproved     bool   `json:"isApproved"`
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

// Withdraw initiates a withdrawal of USDB from the 100x account.
//
// Params:
//   - params: An instance of types.WithdrawRequest containing the withdrawal parameters,
//     including the withdrawal amount and destination address.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			Asset        string `json:"asset"`
			Quantity     string `json:"quantity"`
			Nonce        int64  `json:"nonce"`
			Signature    string `json:"signature"`
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

// NewOrder creates a new order on the SubAccount.
//
// Params:
//   - params: An instance of types.NewOrderRequest containing the order parameters,
//     including the order type (limit/market), quantity, side (buy/sell), price, and symbol.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
			IsBuy        bool   `json:"isBuy"`
			OrderType    int64  `json:"orderType"`
			TimeInForce  int64  `json:"timeInForce"`
			Expiration   int64  `json:"expiration"`
			Price        string `json:"price"`
			Quantity     string `json:"quantity"`
			Nonce        int64  `json:"nonce"`
			Signature    string `json:"signature"`
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

// CancelOrderAndReplace cancels an order and creates a new order on the SubAccount.
//
// Params:
//   - params: An instance of types.CancelOrderAndReplaceRequest containing the necessary
//     parameters to identify the order to cancel and the new order parameters.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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
			IdToCancel string      `json:"idToCancel"`
			NewOrder   interface{} `json:"newOrder"`
		}{
			IdToCancel: params.IdToCancel,
			NewOrder: &struct {
				Account      string `json:"account"`
				SubAccountId int64  `json:"subAccountId"`
				ProductId    int64  `json:"productId"`
				IsBuy        bool   `json:"isBuy"`
				OrderType    int64  `json:"orderType"`
				TimeInForce  int64  `json:"timeInForce"`
				Expiration   int64  `json:"expiration"`
				Price        string `json:"price"`
				Quantity     string `json:"quantity"`
				Nonce        int64  `json:"nonce"`
				Signature    string `json:"signature"`
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

// CancelOrder cancels an active order on the SubAccount.
//
// Params:
//   - params: An instance of types.CancelOrderRequest containing the necessary
//     parameters to identify the order to cancel.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
			OrderId      string `json:"orderId"`
			Signature    string `json:"signature"`
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

// CancelAllOpenOrders cancels all active orders on a specific product for the SubAccount.
//
// Params:
//   - product: The product for which all active orders should be canceled.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
			Signature    string `json:"signature"`
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

// GetSpotBalances retrieves spot balances for the SubAccount.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// GetPerpetualPosition retrieves the perpetual position for a specific product and SubAccount.
//
// Parameters:
//   - product: The product for which the perpetual position is requested.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// ListApprovedSigners retrieves a list of all approved signers for a specific `SubAccount`.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// ListOpenOrders retrieves all open orders on the `SubAccount` for a specific product.
//
// Parameters:
//   - product: A pointer to a `types.Product` struct representing the product for which open orders are to be retrieved.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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

// ListOrders retrieves all orders on the `SubAccount` for a specific product.
//
// Parameters:
//   - params: A pointer to a `types.ListOrdersRequest` struct containing parameters for listing orders.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
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
//
// Parameters:
//   - ctx: The context.Context for the Ethereum transaction.
//   - amount: The amount of USDB tokens to approve, specified as a *big.Int.
//
// Returns:
//   - A pointer to a geth_types.Transaction representing the Ethereum transaction.
//   - An error if the Ethereum transaction fails or encounters an issue.
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
//
// Parameters:
//   - ctx: The context.Context for the Ethereum transaction.
//   - amount: The amount of USDB tokens to deposit, specified as a *big.Int.
//
// Returns:
//   - A pointer to a geth_types.Transaction representing the Ethereum transaction.
//   - An error if the Ethereum transaction fails or encounters an issue.
func (go100XClient *Go100XAPIClient) DepositUSDB(ctx context.Context, amount *big.Int) (*geth_types.Transaction, error) {
	// Approve self as signer
	_, err := go100XClient.ApproveSigner(&types.ApproveRevokeSignerRequest{
		ApprovedSigner: go100XClient.addressString,
		Nonce:          time.Now().UnixMicro(),
	})
	if err != nil {
		return nil, err
	}

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
//
// Parameters:
//   - ctx: The context.Context for the Ethereum transaction.
//   - transaction: The Ethereum transaction (*geth_types.Transaction) to monitor.
//
// Returns:
//   - A pointer to a geth_types.Receipt containing the transaction receipt once the transaction is mined.
//   - An error if the transaction fails to be mined or encounters an issue.
func (go100XClient *Go100XAPIClient) WaitTransaction(ctx context.Context, transaction *geth_types.Transaction) (*geth_types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, go100XClient.EthClient, transaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// addReferee adds a referee to author referral code.
//
// Returns:
//   - A pointer to an http.Response containing the response from the API call.
//   - An error if the API call fails or if the response is not as expected.
func (go100XClient *Go100XAPIClient) addReferee() (*http.Response, error) {
	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_REFERRAL,
		&struct {
			Account string `json:"account"`
			Code    string `json:"code"`
		}{
			Account: go100XClient.addressString,
			Code:    "eldief",
		},
	)
	if err != nil {
		return nil, err
	}

	// Create HTTP request.
	request, err := utils.CreateHTTPRequestWithBody(
		http.MethodPost,
		string(go100XClient.baseUrl)+string(constants.API_ENDPOINT_ADD_REFEREE),
		&struct {
			Account   string `json:"account"`
			Code      string `json:"code"`
			Signature string `json:"signature"`
		}{
			Account:   go100XClient.addressString,
			Code:      "eldief",
			Signature: signature,
		},
	)
	if err != nil {
		return nil, err
	}

	// Send HTTP request and return result.
	return utils.SendHTTPRequest(go100XClient.HttpClient, request)
}
