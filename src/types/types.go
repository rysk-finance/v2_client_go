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

type Product struct {
	Symbol string
	Id     int64
}

type ClientConfiguration struct {
	Env          Environment
	Timeout      time.Duration
	PrivateKey   string
	RpcUrl       string
	SubAccountId uint8
}

type Client struct {
	BaseUri           string
	PrivateKey        string
	Address           string
	SubAccountId      int64
	HttpClient        *http.Client
	VerifyingContract string
	Domain            apitypes.TypedDataDomain
}
