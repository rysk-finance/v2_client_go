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

// Go100XAPIClient configuration
type Go100XWSClientConfiguration struct {
	Env          types.Environment // `constants.ENVIRONMENT_TESTNET` or `constants.ENVIRONMENT_MAINNET`.
	PrivateKey   string            // Account private key with or without `0x` prefix.
	RpcUrl       string            // e.g. `https://sepolia.blast.io` or `https://rpc.blastblockchain.com`.
	SubAccountId uint8             // ID of the subaccount to use.
}

// 100x Websocket client.
type Go100XWSClient struct {
	env              types.Environment
	rpcUrl           string
	streamUrl        string
	privateKeyString string
	addressString    string
	privateKey       *ecdsa.PrivateKey
	address          common.Address
	ciao             common.Address
	usdb             common.Address
	domain           apitypes.TypedDataDomain
	SubAccountId     int64
	RPCConnection    *websocket.Conn
	StreamConnection *websocket.Conn
	EthClient        types.IEthClient
}

// NewGo100XWSClient creates a new `Go100XWSClient` instance.
// Initializes the client with the provided configuration.
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

// ListProducts returns the list of products.
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

// GetProduct returns details for a specific product.
func (go100XClient *Go100XWSClient) GetProduct(messageId string, product *types.Product) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_GET_PRODUCT,
		Params: &struct {
			Id string `json:"id"`
		}{
			Id: strconv.FormatInt(product.Id, 10),
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// ServerTime tests connectivity and get the current server time.
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

// Login authenticate a websocket connection.
// Authentication using signature is required to create and cancel orders, deposit and withdraw.
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

// SessionStatus checks active session and return the address currently authenticated.
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

// SubAccountList lists all sub accounts.
func (go100XClient *Go100XWSClient) SubAccountList(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_SUB_ACCOUNT_LIST,
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// Withdraw withdraws USDB from 100x account.
func (go100XClient *Go100XWSClient) Withdraw(messageId string, params *types.WithdrawRequest) error {
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
		return err
	}

	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_WITHDRAW,
		Params: &struct {
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
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// ApproveSigner approves a Signer for a `SubAccount`.
func (go100XClient *Go100XWSClient) ApproveSigner(messageId string, params *types.ApproveRevokeSignerRequest) error {
	return go100XClient.approveRevokeSigner(messageId, params, true)
}

// RevokeSigner revokes a Signer for a `SubAccount`.
func (go100XClient *Go100XWSClient) RevokeSigner(messageId string, params *types.ApproveRevokeSignerRequest) error {
	return go100XClient.approveRevokeSigner(messageId, params, false)
}

// approveRevokeSigner approves or revoke a signer for a `SubAccount`.
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
			Account      string
			SubAccountId int64
			Signer       string
			Approved     bool
			Nonce        int64
			Signature    string
		}{
			Account:      go100XClient.addressString,
			SubAccountId: go100XClient.SubAccountId,
			Signer:       params.ApprovedSigner,
			Approved:     isApproved,
			Nonce:        params.Nonce,
			Signature:    signature,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// NewOrder creates a new order on the SubAccount.
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
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// ListOpenOrders returns all open orders on the `SubAccount` per product.
func (go100XClient *Go100XWSClient) ListOpenOrders(messageId string, params *types.ListOrdersRequest) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ORDER_LIST,
		Params: &struct {
			Account      string
			SubAccountId int64
			ProductId    int64
			OrderIds     []string
			StartTime    int64
			EndTime      int64
			Limit        int64
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

// CancelOrder cancel an active order on the `SubAccount`.
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
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// CancelAllOpenOrders cancel all active orders on a product.
// Returns number of deleted orders.
func (go100XClient *Go100XWSClient) CancelAllOpenOrders(messageId string, product *types.Product) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_CANCEL_ALL_OPEN_ORDERS,
		Params: &struct {
			Account      string
			SubAccountId int64
			ProductId    int64
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
func (go100XClient *Go100XWSClient) OrderBook(messageId string, params *types.OrderBookRequest) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ORDER_BOOK_DEPTH,
		Params: &struct {
			Symbol      string
			Granularity int64
			Limit       int64
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
		Method:  constants.WS_METHOD_ORDER_BOOK_DEPTH,
		Params: &struct {
			Account    string
			SubAccount int64
			ProductIds []int64
		}{
			Account:    go100XClient.addressString,
			SubAccount: go100XClient.SubAccountId,
			ProductIds: productIds,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// GetPerpetualPosition returns perpetual position for sub account id.
func (go100XClient *Go100XWSClient) GetSpotBalances(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ORDER_BOOK_DEPTH,
		Params: &struct {
			Account    string
			SubAccount int64
			Assets     []string
		}{
			Account:    go100XClient.addressString,
			SubAccount: go100XClient.SubAccountId,
			Assets:     []string{constants.USDB_ADDRESS[go100XClient.env]},
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// AccountUpdates returns immediate order updates on placement, execution, cancellation,
// up to date spot balances and perp positions pushed out every 5s.
func (go100XClient *Go100XWSClient) AccountUpdates(messageId string) error {
	// Generate RPC request.
	request := &types.WebsocketRequest{
		JsonRPC: constants.WS_JSON_RPC,
		ID:      messageId,
		Method:  constants.WS_METHOD_ACCOUNT_UPDATES,
		Params: &struct {
			Account    string
			SubAccount int64
		}{
			Account:    go100XClient.addressString,
			SubAccount: go100XClient.SubAccountId,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.RPCConnection, request)
}

// SubscribeAggregateTrades subscribes to aggregate trade (aggTrade) that represents one or more individual trades.
// Trades that fill at the same time, from the same taker order.
func (go100XClient *Go100XWSClient) SubscribeAggregateTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeAggregateTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products)
}

// UnubscribeAggregateTrades unsubscribes from aggregate trade (aggTrade) that represents one or more individual trades.
// Trades that fill at the same time, from the same taker order.
func (go100XClient *Go100XWSClient) UnsubscribeAggregateTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeAggregateTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products)
}

// subscribeUnsubscribeAggregateTrades subscribe or unsubscribe to/from aggregate trade (aggTrade).
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
func (go100XClient *Go100XWSClient) SubscribeSingleTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeSingleTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products)
}

// UnubscribeSingleTrades unsubscribes from Trade Streams that push raw trade information; each trade has a unique buyer and seller.
func (go100XClient *Go100XWSClient) UnubscribeSingleTrades(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribeSingleTrades(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products)
}

// subscribeUnsubscribeSingleTrades subscribe or unsubscribe to/from Trade Streams.
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
func (go100XClient *Go100XWSClient) SubscribeKlineData(messageId string, products []*types.Product, intervals []types.Interval) error {
	return go100XClient.subscribeUnsubscribeKlineData(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products, intervals)
}

// SubscribeKlineData unsubscribes from Kline/Candlestick Stream that push updates to the current klines/candlestick every second.
func (go100XClient *Go100XWSClient) UnsubscribeKlineData(messageId string, products []*types.Product, intervals []types.Interval) error {
	return go100XClient.subscribeUnsubscribeKlineData(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products, intervals)
}

// subscribeUnsubscribeKlineData subscribe or unsubscribe to/from Kline/Candlestick Stream.
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
func (go100XClient *Go100XWSClient) SubscribePartialBookDepth(messageId string, products []*types.Product, limits []types.Limit, granularities []int64) error {
	return go100XClient.subscribeUnsubscribePartialBookDepth(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products, limits, granularities)
}

// UnubscribePartialBookDepth unsubscribes from top {limit} bids and asks, pushed every second.
// Prices are rounded by 1e{granularity}.
func (go100XClient *Go100XWSClient) UnubscribePartialBookDepth(messageId string, products []*types.Product, limits []types.Limit, granularities []int64) error {
	return go100XClient.subscribeUnsubscribePartialBookDepth(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products, limits, granularities)
}

// subscribeUnsubscribeKlineData subscribe or unsubscribe to/from Kline/Candlestick Stream.
func (go100XClient *Go100XWSClient) subscribeUnsubscribePartialBookDepth(messageId string, method types.WSMethod, products []*types.Product, limits []types.Limit, granularities []int64) error {
	// Create @depth params.
	var params []string
	for _, product := range products {
		for _, limit := range limits {
			for _, granularity := range granularities {
				params = append(params, product.Symbol+"@depth_"+strconv.FormatInt(int64(limit), 10)+"_"+strconv.FormatInt(granularity, 10))
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

// SubscribeSingleTrades subscribes to 24hr rolling window mini-ticker statistics.
// These are NOT the statistics of the UTC day, but a 24hr rolling window for the previous 24hrs.
// Pushed out every 5s.
func (go100XClient *Go100XWSClient) Subscribe24hrPriceChangeStatistics(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribe24hrPriceChangeStatistics(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_SUBSCRIBE, products)
}

// SubscribeSingleTrUnubscribe24hrPriceChangeStatisticsades unsubscribes from 24hr rolling window mini-ticker statistics.
// These are NOT the statistics of the UTC day, but a 24hr rolling window for the previous 24hrs.
// Pushed out every 5s.
func (go100XClient *Go100XWSClient) Unsubscribe24hrPriceChangeStatistics(messageId string, products []*types.Product) error {
	return go100XClient.subscribeUnsubscribe24hrPriceChangeStatistics(messageId, constants.WS_METHOD_MARKET_DATA_STREAMS_UNSUBSCRIBE, products)
}

// subscribeUnsubscribeSingleTrades subscribe or unsubscribe to/from 24hr rolling window mini-ticker statistics.
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
func (go100XClient *Go100XWSClient) DepositUSDB(ctx context.Context, amount *big.Int) (*geth_types.Transaction, error) {
	// Parse ABI
	parsedABI, _ := abi.JSON(strings.NewReader(constants.CIAO_ABI))

	// Pack transaction data
	data, _ := parsedABI.Pack("deposit", go100XClient.addressString, uint8(go100XClient.SubAccountId), amount, go100XClient.usdb)

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
func (go100XClient *Go100XWSClient) WaitTransaction(ctx context.Context, transaction *geth_types.Transaction) (*geth_types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, go100XClient.EthClient, transaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
