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

// `constants.ETH_PERP`, `constants.BTC_PERP` or `constants.BLAST_PERP`.
type Product struct {
	Symbol string // `constants.ETH_PERP.Symbol`, `constants.BTC_PERP.Symbol` or `constants.BLAST_PERP.Symbol`.
	Id     int64  // `constants.ETH_PERP.Id`, `constants.BTC_PERP.Id` or `constants.BLAST_PERP.Id`.
}

// Client configuration
type ClientConfiguration struct {
	Env          Environment   // `constants.TESTNET` or `constants.MAINNET`.
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
