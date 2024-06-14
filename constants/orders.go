package constants

import (
	"math/big"

	"github.com/eldief/go100x/types"
)

const (
	ORDER_TYPE_LIMIT             types.OrderType = 0
	ORDER_TYPE_LIMIT_MAKER       types.OrderType = 1
	ORDER_TYPE_MARKET            types.OrderType = 2
	ORDER_TYPE_STOP_LOSS         types.OrderType = 3
	ORDER_TYPE_STOP_LOSS_LIMIT   types.OrderType = 4
	ORDER_TYPE_TAKE_PROFIT       types.OrderType = 5
	ORDER_TYPE_TAKE_PROFIT_LIMIT types.OrderType = 6
)

const (
	INTERVAL_1M  types.Interval = "1m"  // 1 minute
	INTERVAL_5M  types.Interval = "5m"  // 5 minutes
	INTERVAL_15M types.Interval = "15m" // 15 minutes
	INTERVAL_30M types.Interval = "30m" // 30 minutes
	INTERVAL_1H  types.Interval = "1h"  // 1 hour
	INTERVAL_2H  types.Interval = "2h"  // 2 hours
	INTERVAL_4H  types.Interval = "4h"  // 4 hours
	INTERVAL_8H  types.Interval = "8h"  // 8 hours
	INTERVAL_D1  types.Interval = "1d"  // 1 day
	INTERVAL_D3  types.Interval = "3d"  // 3 days
	INTERVAL_1W  types.Interval = "1w"  // 1 week
)

const (
	LIMIT_FIVE   types.Limit = 5
	LIMIT_TEN    types.Limit = 10
	LIMIT_TWENTY types.Limit = 20
)

const (
	TIME_IN_FORCE_GTC types.TimeInForce = 0
	TIME_IN_FORCE_FOK types.TimeInForce = 1
	TIME_IN_FORCE_IOC types.TimeInForce = 2
)

var (
	E27 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1_000_000_000))
	E26 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100_000_000))
	E25 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10_000_000))
	E24 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1_000_000))
	E23 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100_000))
	E22 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10_000))
	E21 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1_000))
	E20 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100))
	E19 = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10))
	E18 = big.NewInt(1e18)
	E17 = big.NewInt(1e17)
	E16 = big.NewInt(1e16)
	E15 = big.NewInt(1e15)
	E14 = big.NewInt(1e14)
	E13 = big.NewInt(1e13)
	E12 = big.NewInt(1e12)
	E11 = big.NewInt(1e11)
	E10 = big.NewInt(1e10)
)
