package rest_examples

import (
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rysk-finance/v2_client_go/api_client"
	"github.com/rysk-finance/v2_client_go/constants"
	"github.com/rysk-finance/v2_client_go/types"
)

func PlaceOrder() {
	// Load ".env" file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
		return
	}

	// Get Private Key from environment
	privateKey := os.Getenv("PRIVATE_KEYS")
	if privateKey == "" {
		log.Fatalf("websocketPlaceOrder:: no PRIVATE_KEYS found")
	}

	// Initialize new RyskV2APIClient
	client, err := api_client.NewRyskV2APIClient(&api_client.RyskV2APIClientConfiguration{
		Env:          constants.ENVIRONMENT_TESTNET,                  // Arbistrum sepolia testnet
		PrivateKey:   privateKey,                                     // Your private key
		RpcUrl:       "https://arbitrum-sepolia.gateway.tenderly.co", // Public Arbistrum sepolia testnet RPC url
		SubAccountId: 0,                                              // Default frontend subaccount
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	// get market price
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.coinbase.com/v2/exchange-rates?currency=ETH",
		nil,
	)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	var unmarshaled struct {
		Data struct {
			Currency string
			Rates    struct {
				USD string
			}
		}
	}
	err = json.Unmarshal(body, &unmarshaled)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	priceFloat, err := strconv.ParseFloat(unmarshaled.Data.Rates.USD, 64)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	price := new(big.Int)
	new(big.Float).Mul(big.NewFloat(priceFloat), new(big.Float).SetFloat64(1e18)).Int(price)
	price = new(big.Int).Mul(new(big.Int).Div(price, big.NewInt(100)), big.NewInt(110))

	// get product increment
	res, err = client.GetProductById(constants.PRODUCT_ETH_PERP.Id)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	var productUnmarshaled struct {
		Increment string
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	err = json.Unmarshal(body, &productUnmarshaled)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	increment := new(big.Int)
	_, ok := increment.SetString(productUnmarshaled.Increment, 10)
	if !ok {
		log.Fatalf("PlaceOrder:: !ok")
	}

	// adjust price with increment
	remainder := new(big.Int).Mod(price, increment)
	adjustedPrice := new(big.Int).Sub(price, remainder)

	// Create a new order
	response, err := client.NewOrder(&types.NewOrderRequest{
		Product:     &constants.PRODUCT_ETH_PERP,
		IsBuy:       true,
		OrderType:   constants.ORDER_TYPE_LIMIT,
		TimeInForce: constants.TIME_IN_FORCE_GTC,
		Price:       adjustedPrice.String(),
		Quantity:    constants.E16.String(),
		Expiration:  time.Now().Add(10 * time.Minute).UnixMilli(),
		Nonce:       time.Now().UnixMicro(),
	})
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	body, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("PlaceOrder:: %v", err)
	}

	log.Printf("PlaceOrder:: response received: %d - %s", response.StatusCode, body)
}
