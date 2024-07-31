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


## Examples

- Look [here](https://github.com/eldief/go100x/tree/master/examples/rest) for REST API Client examples
- Look [here](https://github.com/eldief/go100x/tree/master/examples/websocket) for Websocket Client examples


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
