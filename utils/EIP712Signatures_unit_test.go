//go:build !integration
// +build !integration

package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EIP712SignaturesTestSuite struct {
	suite.Suite
	PrivateKey       *ecdsa.PrivateKey
	PrivateKeyString string
	PublicKey        []byte
}

func (suite *EIP712SignaturesTestSuite) SetupSuite() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	suite.PrivateKey = privateKey
	suite.PrivateKeyString = strings.TrimPrefix(hex.EncodeToString(crypto.FromECDSA(privateKey)), "0x")
	suite.PublicKey = crypto.FromECDSAPub(&privateKey.PublicKey)
}

func TestRunSuiteUnit_EIP712SignaturesTestSuite(t *testing.T) {
	suite.Run(t, new(EIP712SignaturesTestSuite))
}

func (suite *EIP712SignaturesTestSuite) TestUnit_MapMessageToTypedData() {
	type SampleMessage struct {
		Name  string `json:"name"`
		Value uint   `json:"value"`
	}
	message := SampleMessage{
		Name:  "test",
		Value: 100,
	}

	typedData, err := mapMessageToTypedData(message)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), typedData)
	require.Equal(suite.T(), "test", typedData["name"])
	require.Equal(suite.T(), float64(100), typedData["value"]) // JSON unmarshalling converts integers to float64
}

func (suite *EIP712SignaturesTestSuite) TestUnit_MapMessageToTypedData_MarshalError() {
	type InvalidMessage struct {
		Channel chan int `json:"channel"`
	}
	invalidMessage := InvalidMessage{make(chan int)}

	typedData, err := mapMessageToTypedData(invalidMessage)
	require.Error(suite.T(), err)
	require.Nil(suite.T(), typedData)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_MapMessageToTypedData_UnmarshalError() {
	invalidJSON := "{invalid_json}"

	typedData, err := mapMessageToTypedData(invalidJSON)
	require.Error(suite.T(), err)
	require.Nil(suite.T(), typedData)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_GenerateEIP712Message() {
	typedDataDomain := apitypes.TypedDataDomain{
		Name:              constants.DOMAIN_NAME,
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	rawDataMessage := &struct {
		Account        string `json:"account"`
		SubAccountId   string `json:"subAccountId"`
		ApprovedSigner string `json:"approvedSigner"`
		IsApproved     bool   `json:"isApproved"`
		Nonce          string `json:"nonce"`
	}{
		Account:        "0x0000000000000000000000000000000000000000",
		SubAccountId:   "1",
		ApprovedSigner: "0x0000000000000000000000000000000000000000",
		IsApproved:     true,
		Nonce:          strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
	typedDataMessage, err := mapMessageToTypedData(rawDataMessage)
	require.NoError(suite.T(), err)

	hash, err := generateEIP712Message(constants.PRIMARY_TYPE_APPROVE_SIGNER, typedDataDomain, typedDataMessage)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), hash)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_GenerateEIP712Message_HashStructError() {
	typedDataDomain := apitypes.TypedDataDomain{
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	typedDataMessage := map[string]interface{}{
		"test": "data",
	}

	hash, err := generateEIP712Message(constants.PRIMARY_TYPE_APPROVE_SIGNER, typedDataDomain, typedDataMessage)
	require.Error(suite.T(), err)
	require.Nil(suite.T(), hash)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_GenerateEIP712Message_DomainHashError() {
	typedDataDomain := apitypes.TypedDataDomain{
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	rawDataMessage := &struct {
		Account        string `json:"account"`
		SubAccountId   string `json:"subAccountId"`
		ApprovedSigner string `json:"approvedSigner"`
		IsApproved     bool   `json:"isApproved"`
		Nonce          string `json:"nonce"`
	}{
		Account:        "0x0000000000000000000000000000000000000000",
		SubAccountId:   "1",
		ApprovedSigner: "0x0000000000000000000000000000000000000000",
		IsApproved:     true,
		Nonce:          strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
	typedDataMessage, err := mapMessageToTypedData(rawDataMessage)
	require.NoError(suite.T(), err)

	hash, err := generateEIP712Message(constants.PRIMARY_TYPE_APPROVE_SIGNER, typedDataDomain, typedDataMessage)
	require.Error(suite.T(), err)
	require.Nil(suite.T(), hash)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_SignEIP712Message() {
	typedDataDomain := apitypes.TypedDataDomain{
		Name:              constants.DOMAIN_NAME,
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	rawDataMessage := &struct {
		Account        string `json:"account"`
		SubAccountId   string `json:"subAccountId"`
		ApprovedSigner string `json:"approvedSigner"`
		IsApproved     bool   `json:"isApproved"`
		Nonce          string `json:"nonce"`
	}{
		Account:        "0x0000000000000000000000000000000000000000",
		SubAccountId:   "1",
		ApprovedSigner: "0x0000000000000000000000000000000000000000",
		IsApproved:     true,
		Nonce:          strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
	typedDataMessage, err := mapMessageToTypedData(rawDataMessage)
	require.NoError(suite.T(), err)
	hash, err := generateEIP712Message(constants.PRIMARY_TYPE_APPROVE_SIGNER, typedDataDomain, typedDataMessage)
	require.NoError(suite.T(), err)

	signature, err := signEIP712Message(hash, suite.PrivateKey)
	require.NoError(suite.T(), err)
	sigBytes := common.Hex2Bytes(signature[2:])
	r := sigBytes[:32]
	s := sigBytes[32:64]
	verified := crypto.VerifySignature(suite.PublicKey, hash, append(r, s...))
	require.True(suite.T(), verified)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_SignEIP712Message_InvalidMessage() {
	hash := new([]byte)
	_, err := signEIP712Message(*hash, suite.PrivateKey)
	require.Error(suite.T(), err)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_SignMessage() {
	typedDataDomain := apitypes.TypedDataDomain{
		Name:              constants.DOMAIN_NAME,
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	signature, err := SignMessage(
		typedDataDomain,
		suite.PrivateKeyString,
		constants.PRIMARY_TYPE_APPROVE_SIGNER,
		&struct {
			Account        string `json:"account"`
			SubAccountId   string `json:"subAccountId"`
			ApprovedSigner string `json:"approvedSigner"`
			IsApproved     bool   `json:"isApproved"`
			Nonce          string `json:"nonce"`
		}{
			Account:        "0x0000000000000000000000000000000000000000",
			SubAccountId:   "1",
			ApprovedSigner: "0x0000000000000000000000000000000000000000",
			IsApproved:     true,
			Nonce:          strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	)
	require.NoError(suite.T(), err)
	require.True(suite.T(), strings.HasPrefix(signature, "0x"))
}

func (suite *EIP712SignaturesTestSuite) TestUnit_SignMessage_MapMessageToTypedDataError() {
	typedDataDomain := apitypes.TypedDataDomain{
		Name:              constants.DOMAIN_NAME,
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	_, err := SignMessage(
		typedDataDomain,
		suite.PrivateKeyString,
		constants.PRIMARY_TYPE_APPROVE_SIGNER,
		&struct {
			Channel chan int `json:"channel"`
		}{
			Channel: make(chan int),
		},
	)
	require.Error(suite.T(), err)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_SignMessage_GenerateEIP712MessageError() {
	typedDataDomain := apitypes.TypedDataDomain{
		Name:              constants.DOMAIN_NAME,
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	_, err := SignMessage(
		typedDataDomain,
		suite.PrivateKeyString,
		"",
		&struct {
			Account        string `json:"account"`
			SubAccountId   string `json:"subAccountId"`
			ApprovedSigner string `json:"approvedSigner"`
			IsApproved     bool   `json:"isApproved"`
			Nonce          string `json:"nonce"`
		}{
			Account:        "0x0000000000000000000000000000000000000000",
			SubAccountId:   "1",
			ApprovedSigner: "0x0000000000000000000000000000000000000000",
			IsApproved:     true,
			Nonce:          strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	)
	require.Error(suite.T(), err)
}

func (suite *EIP712SignaturesTestSuite) TestUnit_SignMessage_HexToECDSAError() {
	typedDataDomain := apitypes.TypedDataDomain{
		Name:              constants.DOMAIN_NAME,
		Version:           constants.DOMAIN_VERSION,
		ChainId:           constants.CHAIN_ID[constants.ENVIRONMENT_TESTNET],
		VerifyingContract: constants.ORDER_DISPATCHER_ADDRESS[constants.ENVIRONMENT_TESTNET],
	}
	_, err := SignMessage(
		typedDataDomain,
		"",
		constants.PRIMARY_TYPE_APPROVE_SIGNER,
		&struct {
			Account        string `json:"account"`
			SubAccountId   string `json:"subAccountId"`
			ApprovedSigner string `json:"approvedSigner"`
			IsApproved     bool   `json:"isApproved"`
			Nonce          string `json:"nonce"`
		}{
			Account:        "0x0000000000000000000000000000000000000000",
			SubAccountId:   "1",
			ApprovedSigner: "0x0000000000000000000000000000000000000000",
			IsApproved:     true,
			Nonce:          strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	)
	require.Error(suite.T(), err)
}
