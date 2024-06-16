package ws_client

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
	"github.com/gorilla/websocket"
)

// Go100XWSClientConfiguration represents configuration settings for the 100x WebSocket client.
type Go100XWSClientConfiguration struct {
	Env          types.Environment // Env specifies the environment: `constants.ENVIRONMENT_TESTNET` or `constants.ENVIRONMENT_MAINNET`.
	PrivateKey   string            // PrivateKey is the account private key with or without `0x` prefix.
	RpcUrl       string            // RpcUrl is the URL for the RPC server, e.g., `https://sepolia.blast.io` or `https://rpc.blastblockchain.com`.
	SubAccountId uint8             // SubAccountId is the ID of the subaccount to use.
}

// Go100XWSClient is the WebSocket client for interacting with 100x services.
type Go100XWSClient struct {
	env              types.Environment        // env is the current environment setting.
	rpcUrl           string                   // rpcUrl is the RPC server URL.
	streamUrl        string                   // streamUrl is the WebSocket stream URL.
	privateKeyString string                   // privateKeyString is the private key string.
	addressString    string                   // addressString is the Ethereum address string derived from the private key.
	privateKey       *ecdsa.PrivateKey        // privateKey is the ECDSA private key instance.
	address          common.Address           // address is the Ethereum address derived from the private key.
	ciao             common.Address           // ciao is a common address used in the context.
	usdb             common.Address           // usdb is a common address used in the context.
	domain           apitypes.TypedDataDomain // domain represents the typed data domain for API requests.
	SubAccountId     int64                    // SubAccountId is the ID of the subaccount to use.
	RPCConnection    *websocket.Conn          // RPCConnection is the WebSocket connection for RPC operations.
	StreamConnection *websocket.Conn          // StreamConnection is the WebSocket connection for streaming operations.
	EthClient        types.IEthClient         // EthClient is the Ethereum client interface.
}

// NewGo100XWSClient creates a new `Go100XWSClient` instance based on the provided configuration.
// It initializes and returns a new client that connects to the 100x WebSocket API.
//
// Parameters:
//   - config: A pointer to a `Go100XWSClientConfiguration` struct that contains configuration parameters
//     such as environment (`types.Environment`), private key (`string`), RPC URL (`string`), and subaccount ID (`uint8`).
//
// Returns:
//   - *Go100XWSClient: A pointer to the initialized `Go100XWSClient` instance.
//   - error: An error if the client initialization fails.
func NewGo100XWSClient(config *Go100XWSClientConfiguration) (*Go100XWSClient, error) {
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

	// Create RPC websocket connection.
	rpcWebsocket, _, err := websocket.DefaultDialer.DialContext(
		context.Background(),
		constants.WS_RPC_URL[config.Env],
		http.Header{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC websocket: %v", err)
	}

	// Create streamWebsocket websocket connection.
	streamWebsocket, _, err := websocket.DefaultDialer.DialContext(
		context.Background(),
		constants.WS_STREAM_URL[config.Env],
		http.Header{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to stream websocket: %v", err)
	}

	// Return a new `go100x.Client`.
	return &Go100XWSClient{
		env:              config.Env,
		rpcUrl:           constants.WS_RPC_URL[config.Env],
		streamUrl:        constants.WS_STREAM_URL[config.Env],
		privateKeyString: privateKeyString,
		addressString:    utils.AddressFromPrivateKey(privateKeyString),
		address:          common.HexToAddress(utils.AddressFromPrivateKey(privateKeyString)),
		privateKey:       privateKey,
		ciao:             common.HexToAddress(constants.CIAO_ADDRESS[config.Env]),
		usdb:             common.HexToAddress(constants.USDB_ADDRESS[config.Env]),
		domain: apitypes.TypedDataDomain{
			Name:              constants.DOMAIN_NAME,
			Version:           constants.DOMAIN_VERSION,
			ChainId:           constants.CHAIN_ID[config.Env],
			VerifyingContract: constants.VERIFIER_ADDRESS[config.Env],
		},
		SubAccountId:     int64(config.SubAccountId),
		RPCConnection:    rpcWebsocket,
		StreamConnection: streamWebsocket,
		EthClient:        client,
	}, nil
}

// ListProducts sends a request to retrieve the list of products available on the 100x WebSocket API.
// It subscribes to the `LIST_PRODUCTS` message identifier to fetch the products.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//
// Returns:
//   - error: An error if the request to fetch the products fails.
func (go100XClient *Go100XWSClient) ListProducts(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_LIST_PRODUCTS,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// GetProduct sends a request to retrieve details for a specific product using the 100x WebSocket API.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - product: A pointer to the product details structure (types.Product) where the retrieved data will be stored.
//
// Returns:
//   - error: An error if the request to fetch the product details fails.
func (go100XClient *Go100XWSClient) GetProduct(messageId string, product *types.Product) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_GET_PRODUCT,
		Params: &struct {
			Symbol string `json:"symbol"`
		}{
			Symbol: product.Symbol,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// ServerTime sends a request to test connectivity and retrieve the current server time
// using the 100x WebSocket API.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//
// Returns:
//   - error: An error if the request to fetch the server time fails.
func (go100XClient *Go100XWSClient) ServerTime(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_SERVER_TIME,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// Login performs authentication for the WebSocket connection.
// Authentication using signature is required to create and cancel orders, deposit and withdraw.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//
// Returns:
//   - error: An error if the authentication fails.
func (go100XClient *Go100XWSClient) Login(messageId string) error {
	// Current timestamp in ms, will be rejected if older than 10s, easiest to send in a time in the future.
	timestamp := uint64(time.Now().Add(10 * time.Second).UnixMilli())

	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKeyString,
		constants.PRIMARY_TYPE_LOGIN_MESSAGE,
		&struct {
			Account   string `json:"account"`
			Message   string `json:"message"`
			Timestamp uint64 `json:"timestamp"`
		}{
			Account:   go100XClient.addressString,
			Message:   "I want to log into 100x.finance",
			Timestamp: timestamp,
		},
	)
	if err != nil {
		return err
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_LOGIN,
		Params: &struct {
			Account   string `json:"account"`
			Message   string `json:"message"`
			Timestamp uint64 `json:"timestamp"`
			Signature string `json:"signature"`
		}{
			Account:   go100XClient.addressString,
			Message:   "I want to log into 100x.finance",
			Timestamp: timestamp,
			Signature: signature,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// SessionStatus checks the active session and returns the address currently authenticated.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//
// Returns:
//   - error: An error if the session status retrieval fails.
func (go100XClient *Go100XWSClient) SessionStatus(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_SESSION_STATUS,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// SubAccountList retrieves a list of all sub-accounts associated with the authenticated account.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//
// Returns:
//   - error: An error if the sub-account list retrieval fails.
func (go100XClient *Go100XWSClient) SubAccountList(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_SUB_ACCOUNT_LIST,
		Params: &struct {
			Account string `json:"account"`
		}{
			Account: go100XClient.addressString,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// ApproveSigner approves a signer for a sub-account.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - params: Approval parameters including the signer details.
//
// Returns:
//   - error: An error if the approval process fails.
func (go100XClient *Go100XWSClient) ApproveSigner(messageId string, params *types.ApproveRevokeSignerRequest) error {
	return go100XClient.approveRevokeSigner(messageId, params, true)
}

// RevokeSigner revokes a signer for a sub-account.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - params: Revocation parameters including the signer details.
//
// Returns:
//   - error: An error if the revocation process fails.
func (go100XClient *Go100XWSClient) RevokeSigner(messageId string, params *types.ApproveRevokeSignerRequest) error {
	return go100XClient.approveRevokeSigner(messageId, params, false)
}

// approveRevokeSigner approves or revokes a signer for a sub-account.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - params: Approval or revocation parameters, including signer details.
//   - isApproved: Boolean flag indicating whether to approve or revoke the signer.
//
// Returns:
//   - error: An error if the operation fails.
func (go100XClient *Go100XWSClient) approveRevokeSigner(messageId string, params *types.ApproveRevokeSignerRequest, isApproved bool) error {
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
		return err
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_APPROVE_REVOKE_SIGNER,
		Params: &struct {
			Account        string `json:"account"`
			SubAccountId   int64  `json:"subAccountId"`
			ApprovedSigner string `json:"approvedSigner"`
			IsApproved     bool   `json:"isApproved"`
			Nonce          int64  `json:"nonce"`
			Signature      string `json:"signature"`
		}{
			Account:        go100XClient.addressString,
			SubAccountId:   go100XClient.SubAccountId,
			ApprovedSigner: params.ApprovedSigner,
			IsApproved:     isApproved,
			Nonce:          params.Nonce,
			Signature:      signature,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// NewOrder creates a new order on the SubAccount.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - params: A struct containing details for the new order, such as product symbol, order type, quantity, price, etc.
//
// Returns:
//   - error: An error if the operation fails.
func (go100XClient *Go100XWSClient) NewOrder(messageId string, params *types.NewOrderRequest) error {
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
		return err
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_NEW_ORDER,
		Params: &struct {
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
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// ListOpenOrders returns all open orders on the `SubAccount` per product.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - params: A struct containing parameters to specify the product and additional filtering criteria for the orders.
//
// Returns:
//   - error: An error if the operation fails.
func (go100XClient *Go100XWSClient) ListOpenOrders(messageId string, params *types.ListOrdersRequest) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ORDER_LIST,
		Params: &struct {
			Account      string   `json:"account"`
			SubAccountId int64    `json:"subAccountId"`
			ProductId    int64    `json:"productId"`
			OrderIds     []string `json:"orderIds"`
			StartTime    int64    `json:"startTime"`
			EndTime      int64    `json:"endTime"`
			Limit        int64    `json:"limit"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			ProductId:    params.Product.Id,
			OrderIds:     params.Ids,
			StartTime:    params.StartTime,
			EndTime:      params.EndTime,
			Limit:        params.Limit,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// CancelOrder cancels an active order on the `SubAccount`.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - params: A struct containing parameters to specify the order to be canceled.
//
// Returns:
//   - error: An error if the operation fails.
func (go100XClient *Go100XWSClient) CancelOrder(messageId string, params *types.CancelOrderRequest) error {
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
		return err
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_CANCEL_ORDER,
		Params: &struct {
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
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// CancelAllOpenOrders cancels all active orders on a product for the `SubAccount`.
//
// Parameters:
//   - messageId: The unique identifier for the message.
//   - product: The product for which all active orders should be canceled.
//
// Returns:
//   - error: An error if the operation fails.
//
// Returns number of deleted orders.
func (go100XClient *Go100XWSClient) CancelAllOpenOrders(messageId string, product *types.Product) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_CANCEL_ALL_OPEN_ORDERS,
		Params: &struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
			ProductId    int64  `json:"productId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			ProductId:    product.Id,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// OrderBook returns bids and asks for a market.
//
// It retrieves the order book data for the specified market based on the provided parameters.
// The order book includes bids and asks, which represent buy and sell orders respectively.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - params: An OrderBookRequest struct pointer containing parameters such as market ID.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) OrderBook(messageId string, params *types.OrderBookRequest) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ORDER_BOOK_DEPTH,
		Params: &struct {
			Symbol      string `json:"symbol"`
			Granularity int64  `json:"granularity"`
			Limit       int64  `json:"limit"`
		}{
			Symbol:      params.Product.Symbol,
			Granularity: params.Granularity,
			Limit:       int64(params.Limit),
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// GetPerpetualPosition returns perpetual position for sub account id.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products for which to retrieve perpetual positions.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) GetPerpetualPosition(messageId string, products []*types.Product) error {
	// Create ProductIds slice.
	var productIds []int64
	for _, product := range products {
		productIds = append(productIds, product.Id)
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_GET_PERPETUAL_POSITION,
		Params: &struct {
			Account      string  `json:"account"`
			SubAccountId int64   `json:"subAccountId"`
			ProductIds   []int64 `json:"productIds"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			ProductIds:   productIds,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// GetSpotBalances returns spot balances for sub account id.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - assets: A slice of strings representing the assets for which to retrieve spot balances.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) GetSpotBalances(messageId string, assets []string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_GET_SPOT_BALANCES,
		Params: &struct {
			Account      string   `json:"account"`
			SubAccountId int64    `json:"subAccountId"`
			Assets       []string `json:"assets"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			Assets:       assets,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// AccountUpdates returns immediate order updates on placement, execution, cancellation,
// up to date spot balances and perp positions pushed out every 5s.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) AccountUpdates(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ACCOUNT_UPDATES,
		Params: &struct {
			Account      string `json:"account"`
			SubAccountId int64  `json:"subAccountId"`
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// SubscribeAggregateTrades subscribes to aggregate trade (aggTrade) that represents one or more individual trades.
// Trades that fill at the same time, from the same taker order.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products to subscribe to for aggregate trades.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) SubscribeAggregateTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeAggregateTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products)
}

// UnsubscribeAggregateTrades unsubscribes from aggregate trade (aggTrade) that represents one or more individual trades.
// Trades that fill at the same time, from the same taker order.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products to unsubscribe from for aggregate trades.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) UnsubscribeAggregateTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeAggregateTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products)
}

// subscribeUnsubscribeAggregateTrades subscribes or unsubscribes to/from aggregate trade (aggTrade).
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - method: The WebSocket method (subscribe or unsubscribe).
//   - products: A slice of Product pointers representing the products to subscribe or unsubscribe for aggregate trades.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) subscribeUnsubscribeAggregateTrades(messageId string, method types.WSMethod, products []*types.Product) error {
	// Create @aggTrade params.
	var params []string
	for _, product := range products {
		params = append(params, product.Symbol+"@aggTrade")
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  method,
		Params:  params,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.StreamConnection, request)
}

// SubscribeSingleTrades subscribes to Trade Streams that push raw trade information; each trade has a unique buyer and seller.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products to subscribe to for single trades.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) SubscribeSingleTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeSingleTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products)
}

// UnsubscribeSingleTrades unsubscribes from Trade Streams that push raw trade information; each trade has a unique buyer and seller.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products to unsubscribe from for single trades.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) UnubscribeSingleTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeSingleTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products)
}

// subscribeUnsubscribeSingleTrades subscribes or unsubscribes to/from Trade Streams.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - method: The WebSocket method (subscribe or unsubscribe).
//   - products: A slice of Product pointers representing the products to subscribe or unsubscribe for single trades.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) subscribeUnsubscribeSingleTrades(messageId string, method types.WSMethod, products []*types.Product) error {
	// Create @trade params.
	var params []string
	for _, product := range products {
		params = append(params, product.Symbol+"@trade")
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  method,
		Params:  params,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.StreamConnection, request)
}

// SubscribeKlineData subscribes to Kline/Candlestick Stream that push updates to the current klines/candlestick every second.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products to subscribe to for Kline/Candlestick data.
//   - intervals: A slice of Interval values representing the time intervals for the Kline/Candlestick data.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) SubscribeKlineData(messageId string, products []*types.Product, intervals []types.Interval) error {
	return go100XClient.subscribeUnsubscribeKlineData(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products, intervals)
}

// UnsubscribeKlineData unsubscribes from Kline/Candlestick Stream that pushes updates to the current klines/candlestick every second.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - products: A slice of Product pointers representing the products to unsubscribe from for Kline/Candlestick data.
//   - intervals: A slice of Interval values representing the time intervals for the Kline/Candlestick data.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) UnsubscribeKlineData(messageId string, products []*types.Product, intervals []types.Interval) error {
	return go100XClient.subscribeUnsubscribeKlineData(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products, intervals)
}

// subscribeUnsubscribeKlineData subscribes or unsubscribes to/from Kline/Candlestick Stream.
//
// Parameters:
//   - messageId: A unique identifier for the message.
//   - method: The WebSocket method (subscribe or unsubscribe).
//   - products: A slice of Product pointers representing the products to subscribe or unsubscribe for Kline/Candlestick data.
//   - intervals: A slice of Interval values representing the time intervals for the Kline/Candlestick data.
//
// Returns:
//   - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) subscribeUnsubscribeKlineData(messageId string, method types.WSMethod, products []*types.Product, intervals []types.Interval) error {
	// Create @klines params.
	var params []string
	for _, product := range products {
		for _, interval := range intervals {
			params = append(params, product.Symbol+"@klines_"+string(interval))
		}
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  method,
		Params:  params,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.StreamConnection, request)
}

// SubscribePartialBookDepth subscribes to top {limit} bids and asks, pushed every second.
// Prices are rounded by 1e{granularity}.
//
// Parameters:
// - messageId: A unique identifier for the message.
// - products: A slice of Product pointers representing the products to subscribe to for partial book depth.
// - limits: A slice of Limit values representing the depth limits for the book.
// - granularities: A slice of int64 values representing the price rounding granularity.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) SubscribePartialBookDepth(messageId string, products []*types.Product, limits []types.Limit, granularities []int64) error {
	return go100XClient.subscribeUnsubscribePartialBookDepth(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products, limits, granularities)
}

// UnsubscribePartialBookDepth unsubscribes from top {limit} bids and asks, pushed every second.
// Prices are rounded by 1e{granularity}.
//
// Parameters:
// - messageId: A unique identifier for the message.
// - products: A slice of Product pointers representing the products to unsubscribe from for partial book depth.
// - limits: A slice of Limit values representing the depth limits for the book.
// - granularities: A slice of int64 values representing the price rounding granularity.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) UnsubscribePartialBookDepth(messageId string, products []*types.Product, limits []types.Limit, granularities []int64) error {
	return go100XClient.subscribeUnsubscribePartialBookDepth(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products, limits, granularities)
}

// subscribeUnsubscribePartialBookDepth subscribes or unsubscribes to/from partial book depth updates.
//
// Parameters:
// - messageId: A unique identifier for the message.
// - method: The WebSocket method (subscribe or unsubscribe).
// - products: A slice of Product pointers representing the products to subscribe or unsubscribe for partial book depth updates.
// - limits: A slice of Limit values representing the depth limits for the book.
// - granularities: A slice of int64 values representing the price rounding granularities.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) subscribeUnsubscribePartialBookDepth(messageId string, method types.WSMethod, products []*types.Product, limits []types.Limit, granularities []int64) error {
	// Create @depth params.
	var params []string
	for _, product := range products {
		for _, limit := range limits {
			for _, granularity := range granularities {
				params = append(params, product.Symbol+"@depth"+strconv.FormatInt(int64(limit), 10)+"_"+strconv.FormatInt(granularity, 10))
			}
		}
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  method,
		Params:  params,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.StreamConnection, request)
}

// Subscribe24hrPriceChangeStatistics subscribes to 24hr rolling window mini-ticker statistics.
// These are NOT the statistics of the UTC day, but a 24hr rolling window for the previous 24hrs.
// Pushed out every 5s.
//
// Parameters:
// - messageId: A unique identifier for the message.
// - products: A slice of Product pointers representing the products to subscribe to for 24hr rolling window mini-ticker statistics.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) Subscribe24hrPriceChangeStatistics(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribe24hrPriceChangeStatistics(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products)
}

// Unsubscribe24hrPriceChangeStatistics unsubscribes from 24hr rolling window mini-ticker statistics.
// These are NOT the statistics of the UTC day, but a 24hr rolling window for the previous 24hrs.
// Pushed out every 5s.
//
// Parameters:
// - messageId: A unique identifier for the message.
// - products: A slice of Product pointers representing the products to unsubscribe from for 24hr rolling window mini-ticker statistics.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) Unsubscribe24hrPriceChangeStatistics(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribe24hrPriceChangeStatistics(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products)
}

// subscribeUnsubscribe24hrPriceChangeStatistics subscribes or unsubscribes to/from 24hr rolling window mini-ticker statistics.
//
// Parameters:
// - messageId: A unique identifier for the message.
// - method: The WebSocket method (subscribe or unsubscribe).
// - products: A slice of Product pointers representing the products to subscribe or unsubscribe for 24hr rolling window mini-ticker statistics.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
func (go100XClient *Go100XWSClient) subscribeUnsubscribe24hrPriceChangeStatistics(messageId string, method types.WSMethod, products []*types.Product) error {
	// Create @ticker params.
	var params []string
	for _, product := range products {
		params = append(params, product.Symbol+"@ticker")
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  method,
		Params:  params,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.StreamConnection, request)
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
func (go100XClient *Go100XWSClient) ApproveUSDB(ctx context.Context, amount *big.Int) (*geth_types.Transaction, error) {
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
func (go100XClient *Go100XWSClient) DepositUSDB(ctx context.Context, amount *big.Int) (*geth_types.Transaction, error) {
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
func (go100XClient *Go100XWSClient) WaitTransaction(ctx context.Context, transaction *geth_types.Transaction) (*geth_types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, go100XClient.EthClient, transaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
