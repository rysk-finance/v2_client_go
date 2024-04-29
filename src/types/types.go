package types

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type Environment string
type MarginAssetKey string
type Endpoint string
type Interval string
type Limit int64

type Product struct {
	Symbol string
	Id     int64
}

type Client100xConfiguration struct {
	Env          Environment
	Timeout      time.Duration
	PrivateKey   string
	RpcUrl       string
	SubAccountId uint8
}

type Client100x struct {
	BaseUri           string
	PrivateKey        string
	Address           string
	SubAccountId      int64
	HttpClient        *http.Client
	VerifyingContract string
	Domain            apitypes.TypedDataDomain
}
