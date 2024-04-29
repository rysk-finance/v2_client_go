package go100x

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"go100x/src/constants"
	"go100x/src/types"
	"go100x/src/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	once       sync.Once
	httpClient *http.Client
	instance   *types.Client100x
)

// New100xClient creates a new Client100x instance.
// It initializes the client with the provided configuration.
func New100xClient(config *types.Client100xConfiguration) *types.Client100x {
	once.Do(func() {
		// Remove '0x' from private key
		privateKey := strings.TrimPrefix(config.PrivateKey, "0x")

		// Compute address from private key
		address, err := addressFromPrivateKey(privateKey)
		if err != nil {
			panic(err)
		}

		// Instance Client100x
		instance = &types.Client100x{
			BaseUri:           constants.BASE_URI[config.Env],
			PrivateKey:        privateKey,
			Address:           address,
			SubAccountId:      int64(config.SubAccountId),
			HttpClient:        getHTTPClient(config),
			VerifyingContract: constants.CIAO[config.Env],
			Domain: apitypes.TypedDataDomain{
				Name:              constants.NAME,
				Version:           constants.VERSION,
				ChainId:           constants.CHAIN_ID[config.Env],
				VerifyingContract: constants.VERIFIER[config.Env],
			},
		}
	})

	return instance
}

// Get24hrPriceChangeStatistics returns 24 hour rolling window price change statistics.
func Get24hrPriceChangeStatistics(c *types.Client100x, product types.Product) (string, error) {
	uri := c.BaseUri + string(constants.GET_24H_TICKER_PRICE_CHANGE_STATISTICS)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	if product.Id != 0 && product.Symbol != "" {
		query := req.URL.Query()
		query.Add("symbol", string(product.Symbol))
		req.URL.RawQuery = query.Encode()
	}

	return sendRequest(c.HttpClient, req)
}

// GetProduct returns details for a specific product by symbol
func GetProduct(c *types.Client100x, symbol string) (string, error) {
	uri := c.BaseUri + string(constants.GET_PRODUCT) + symbol

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	return sendRequest(c.HttpClient, req)
}

// GetProductById returns details for a specific product by id.
func GetProductById(c *types.Client100x, id int64) (string, error) {
	uri := c.BaseUri + string(constants.GET_PRODUCT_BY_ID) + strconv.FormatInt(id, 10)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	return sendRequest(c.HttpClient, req)
}

// GetKlineData returns Kline/candlestick bars for a symbol. Klines are uniquely identified by interval(timeframe) and startTime.
func GetKlineData(c *types.Client100x, params types.KlineDataRequest) (string, error) {
	uri := string(c.BaseUri) + string(constants.GET_KLINE_DATA)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")

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

	return sendRequest(c.HttpClient, req)
}

// ListProducts returns a list of products available to trade.
func ListProducts(c *types.Client100x) (string, error) {
	uri := string(c.BaseUri) + string(constants.GET_LIST_PRODUCTS)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")

	return sendRequest(c.HttpClient, req)
}

// OrderBook returns bids and asks for a market.
func OrderBook(c *types.Client100x, params types.OrderBookRequest) (string, error) {
	uri := string(c.BaseUri) + string(constants.GET_ORDER_BOOK)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()
	query.Add("symbol", string(params.Product.Symbol))

	if params.Granularity != 0 {
		query.Add("granularity", strconv.FormatInt(params.Granularity, 10))
	}
	if params.Limit != 0 {
		query.Add("limit", strconv.FormatInt(int64(params.Limit), 10))
	}

	req.URL.RawQuery = query.Encode()
	req.Header.Add("accept", "application/json")

	return sendRequest(c.HttpClient, req)
}

// ServerTime returns current server time.
func ServerTime(c *types.Client100x) (string, error) {
	uri := string(c.BaseUri) + string(constants.GET_SERVER_TIME)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")

	return sendRequest(c.HttpClient, req)
}

// Approve or revoke a Signer for a SubAccount
func ApproveRevokeSigner(c *types.Client100x, params *types.ApproveRevokeSignerRequest) (string, error) {
	// Generate EIP712 signature
	signature, err := utils.SignMessage(c, constants.POST_APPROVE_REVOKE_SIGNER, struct {
		Account        string
		SubAccountId   string
		ApprovedSigner string
		IsApproved     bool
		Nonce          string
	}{
		Account:        c.Address,
		SubAccountId:   strconv.FormatInt(c.SubAccountId, 10),
		ApprovedSigner: params.ApprovedSigner,
		IsApproved:     params.IsApproved,
		Nonce:          strconv.FormatInt(params.Nonce, 10),
	})
	if err != nil {
		return "", err
	}

	// Set signature in typed request
	params.Account = c.Address
	params.SubAccountId = c.SubAccountId
	params.Signature = signature

	// Marshal the login request into JSON
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// Create API request
	uri := string(c.BaseUri) + string(constants.POST_APPROVE_REVOKE_SIGNER)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	return sendRequest(c.HttpClient, req)
}

// Login log user in returning a set-cookie header that will need to be attached to authenticated requests to access private endpoints.
func Login(c *types.Client100x) (string, error) {
	timestamp := time.Now().UnixMilli()

	// Generate EIP712 signature
	signature, err := utils.SignMessage(c, constants.POST_LOGIN, &struct {
		Account   string
		Message   string
		Timestamp string
	}{
		Account:   c.Address,
		Message:   "I want to log into 100x.finance",
		Timestamp: strconv.FormatInt(timestamp, 10),
	})
	if err != nil {
		return "", err
	}

	// Marshal the login request into JSON
	body, err := json.Marshal(types.LoginRequest{
		Account:   c.Address,
		Message:   "I want to log into 100x.finance",
		Timestamp: timestamp,
		Signature: signature,
	})
	if err != nil {
		return "", err
	}

	// Create API request
	uri := string(c.BaseUri) + string(constants.POST_LOGIN)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	return sendRequest(c.HttpClient, req)
}

// TODO
// Withdraw USDB from 100x account.
// func Withdraw(c *types.Client100x) (string, error) {

// }

// GetSpotBalances returns spot balances for sub account id.
func GetSpotBalances(c *types.Client100x) (string, error) {
	// Generate EIP712 signature
	signature, err := utils.SignMessage(c, constants.GET_SPOT_BALANCES, &struct {
		Account      string
		SubAccountId string
	}{
		Account:      c.Address,
		SubAccountId: strconv.FormatInt(c.SubAccountId, 10),
	})
	if err != nil {
		return "", err
	}

	// Create API request
	uri := string(c.BaseUri) + string(constants.GET_SPOT_BALANCES)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")

	query := req.URL.Query()
	query.Add("account", c.Address)
	query.Add("subAccountId", strconv.FormatInt(c.SubAccountId, 10))
	query.Add("signature", signature)
	req.URL.RawQuery = query.Encode()

	return sendRequest(c.HttpClient, req)
}

// getHTTPClient returns a singleton instance of http.Client.
// It ensures that only one instance of http.Client is created and reused.
func getHTTPClient(config *types.Client100xConfiguration) *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}
	return httpClient
}

// sendRequest send HTTP request using Client100x.HttpClient and returns response as string.
func sendRequest(c *http.Client, req *http.Request) (string, error) {
	// Send request
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read response
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func addressFromPrivateKey(privateKeyHex string) (string, error) {
	// Convert private key hex string to an ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}

	// Derive the Ethereum address (EOA) from the public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error converting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address, nil
}
