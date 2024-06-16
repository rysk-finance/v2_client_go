# GO 100x

Go100x is a Go SDK for interacting with the 100x API and JSON-RPC Websocket, providing tools and utilities to streamline integration and interaction with 100x's services.

100x: [https://app.100x.finance/](https://app.100x.finance/)

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

Before you start, ensure you have the following installed:

- Go (Golang): Ensure you have Go installed on your machine. You can download it from [https://golang.org/](https://golang.org/) and follow the installation instructions for your operating system.

### Installation

To install `go100x` as a Go module, simply use `go get`:

    $ go get github.com/eldief/go100x
    
This command will download and install go100x and its dependencies.


## Usage

This package follows [100x API documentation](https://100x.readme.io/reference/100x-api-introduction)


Includes:
- REST HTTP client: `Go100XAPIClient` 
- JSON RPC Websocket: `Go100XWSClient`


### Go100XAPIClient

```go
// Import api_client, types and constants packages
import (
  "github.com/eldief/go100x/api_client"
  "github.com/eldief/go100x/constants"
  "github.com/eldief/go100x/types"
)

// Initialize new Go100XAPIClient
client, err := api_client.NewGo100XAPIClient(&api_client.Go100XAPIClientConfiguration{
  Env:          constants.ENVIRONMENT_TESTNET,
  PrivateKey:   privateKey, 
  RpcUrl:       rpcUrl,
  SubAccountId: 0,
})
if err != nil {
    return err
}

// Approve USDB to 100x
transaction, err := client.ApproveUSDB(context.Background(), constants.E20)
if err != nil {
    return err
}
receipt, err := client.WaitTransaction(context.Background(), transaction)
if err != nil {
    return err
}

// Deposit USDB to 100x
transaction, err = client.DepositUSDB(context.Background(), constants.E20)
if err != nil {
    return err
}
receipt, err = client.WaitTransaction(context.Background(), transaction)
if err != nil {
    return err
}

// Create a new order
res, err := client.NewOrder(&types.NewOrderRequest{
  Product:     &constants.PRODUCT_ETH_PERP,
  IsBuy:       true,
  OrderType:   constants.ORDER_TYPE_LIMIT,
  TimeInForce: constants.TIME_IN_FORCE_GTC,
  Price:       price.String(),
  Quantity:    quantity.String(),
  Expiration:  time.Now().Add(time.Minute).UnixMilli(),
  Nonce:       time.Now().UnixMicro(),
})
if err != nil {
    return err
}

// Withdraw funds
res, err = client.Withdraw(&types.WithdrawRequest{
  Quantity: constants.E20.String(),
  Nonce:    time.Now().UnixMicro(),
})
if err != nil {
    return err
}
```

### Go100XWSClient

```go
// Import ws_client package, types and constants packages
import (
  "github.com/eldief/go100x/ws_client"
  "github.com/eldief/go100x/constants"
  "github.com/eldief/go100x/types"
)

// Initialize new Go100XWSClient
client, err := ws_client.NewGo100XWSClient(&ws_client.Go100XWSClientConfiguration{
  Env:          constants.ENVIRONMENT_TESTNET,
  PrivateKey:   privateKeys,
  RpcUrl:       rpcUrl
  SubAccountId: 0,
})
if err != nil {
    return err
}

// Login
err = client.Login("my_login_message_id")
if err != nil {
    return err
}

// Approve USDB to 100x
transaction, err := client.ApproveUSDB(context.Background(), constants.E20)
if err != nil {
    return err
}
receipt, err := client.WaitTransaction(context.Background(), transaction)
if err != nil {
    return err
}


// Deposit USDB to 100x
transaction, err = client.DepositUSDB(context.Background(), constants.E20)
if err != nil {
    return err
}
receipt, err = client.WaitTransaction(context.Background(), transaction)
if err != nil {
    return err
}

// Create a new order
err = client.NewOrder("my_new_order_message_id", &types.NewOrderRequest{
  Product:     &constants.PRODUCT_ETH_PERP,
  IsBuy:       true,
  OrderType:   constants.ORDER_TYPE_LIMIT,
  TimeInForce: constants.TIME_IN_FORCE_GTC,
  Price:       price.String(),
  Quantity:    quanitity.String(),
  Expiration:  time.Now().Add(time.Minute).UnixMilli(),
  Nonce:       time.Now().UnixMicro(),
})
if err != nil {
    return err
}

// Read JSON RPC messages
for {
    _, p, err := client.RPCConnection.ReadMessage()
    if err != nil {
        return err
    }

    var response types.WebsocketResponse
    err = json.Unmarshal(p, &response)
    if err != nil {
        return err
    }
    
    // Elaborate received messages

    break
}

// Subscribe to stream 
err = client.SubscribeAggregateTrades("my_subscription_message_id", []*types.Product{
    &constants.PRODUCT_ETH_PERP
})
if err != nil {
    return err
}

// Read stream messages
for {
    _, p, err := client.StreamConnection.ReadMessage()
    if err != nil {
        return err
    }

    var response types.WebsocketResponse
    err = json.Unmarshal(p, &response)
    if err != nil {
        return err
    }

    // Elaborate received messages

    break
}

// Unsubscribe from stream 
err = client.UnsubscribeAggregateTrades("my_unsubscription_message_id", []*types.Product{
    &constants.PRODUCT_ETH_PERP
})
if err != nil {
    return err
}
```


## Testing

Before running integration tests add a new `.env` file in both `api_client` and `ws_client` folder following both `.env.example` files.

To run tests for GO 100x, you can use the provided Makefile:

```
# Run all tests
$ make test

# Run specific tests
$ make test_utils
$ make test_api_client
$ make test_ws_client

# Run unit tests
$ make test_unit

# Run integration tests
$ make test_integration

# View test coverage
$ make coverage

```

## Contributing

Contributions to GO 100x are welcome! Follow these steps to contribute:

- Fork the repository and create your branch (git checkout -b feature/myfeature).
- Commit your changes (git commit -am 'Add new feature').
- Push to the branch (git push origin feature/myfeature).
- Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
