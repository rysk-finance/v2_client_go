package ws_client

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/eldief/go100x/constants"
	"github.com/eldief/go100x/types"
	"github.com/eldief/go100x/utils"
	"github.com/gorilla/websocket"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
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
	url               string
	privateKey        string
	address           string
	SubAccountId      int64
	wsConnection      *websocket.Conn
	verifyingContract string
	domain            apitypes.TypedDataDomain
}

// NewGo100XWSClient creates a new `Go100XWSClient` instance.
// Initializes the client with the provided configuration.
func NewGo100XWSClient(config *Go100XWSClientConfiguration) *Go100XWSClient {
	// Remove '0x' from private key.
	privateKey := strings.TrimPrefix(config.PrivateKey, "0x")

	websocket, _, err := websocket.DefaultDialer.DialContext(
		context.Background(),
		constants.WS_URL[config.Env],
		http.Header{},
	)
	if err != nil {
		panic(err)
	}

	// Return a new `go100x.Client`.
	return &Go100XWSClient{
		url:               constants.WS_URL[config.Env],
		privateKey:        privateKey,
		address:           utils.AddressFromPrivateKey(privateKey),
		SubAccountId:      int64(config.SubAccountId),
		wsConnection:      websocket,
		verifyingContract: constants.CIAO_ADDRESS[config.Env],
		domain: apitypes.TypedDataDomain{
			Name:              constants.DOMAIN_NAME,
			Version:           constants.DOMAIN_VERSION,
			ChainId:           constants.CHAIN_ID[config.Env],
			VerifyingContract: constants.VERIFIER_ADDRESS[config.Env],
		},
	}
}

// Login authenticate a websocket connection.
// Authentication using signature is required to create and cancel orders, deposit and withdraw.
func (go100XClient *Go100XWSClient) Login(messageId string) error {
	// Current timestamp in ms, will be rejected if older than 10s, easiest to send in a time in the future.
	timestamp := uint64(time.Now().Add(10 * time.Second).UnixMilli())

	// Generate EIP712 signature.
	signature, err := utils.SignMessage(
		go100XClient.domain,
		go100XClient.privateKey,
		constants.PRIMARY_TYPE_LOGIN_MESSAGE,
		&struct {
			Account   string `json:"account"`
			Message   string `json:"message"`
			Timestamp uint64 `json:"timestamp"`
		}{
			Account:   go100XClient.address,
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
			Account:   go100XClient.address,
			Message:   "I want to log into 100x.finance",
			Timestamp: timestamp,
			Signature: signature,
		},
	}

	// Send RPC request.
	return utils.SendRPCRequest(go100XClient.wsConnection, request)
}
