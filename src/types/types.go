package types

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type Environment string
type MarginAsset string
type Endpoint string
type Interval string
type Limit int64
type OrderType int64
type TimeInForce int64

// `constants.PRODUCT_ETH_PERP`, `constants.PRODUCT_BTC_PERP` or `constants.PRODUCT_BLAST_PERP`.
type Product struct {
	Symbol string // `constants.PRODUCT_ETH_PERP.Symbol`, `constants.PRODUCT_BTC_PERP.Symbol` or `constants.PRODUCT_BLAST_PERP.Symbol`.
	Id     int64  // `constants.PRODUCT_ETH_PERP.Id`, `constants.PRODUCT_BTC_PERP.Id` or `constants.PRODUCT_BLAST_PERP.Id`.
}

// Client configuration
type ClientConfiguration struct {
	Env          Environment   // `constants.ENVIRONMENT_TESTNET` or `constants.ENVIRONMENT_MAINNET`.
	Timeout      time.Duration // e.g. `10 * time.Second`.
	PrivateKey   string        // e.g. `0x2638b4...` or `2638b4...`.
	RpcUrl       string        // e.g. `https://sepolia.blast.io` or `https://rpc.blastblockchain.com`.
	SubAccountId uint8         // ID of the subaccount to use.
}

// 100x API Client.
type Client struct {
	BaseUri           string
	PrivateKey        string
	Address           string
	SubAccountId      int64
	HttpClient        *http.Client
	VerifyingContract string
	Domain            apitypes.TypedDataDomain
}
