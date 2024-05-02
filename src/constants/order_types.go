package constants

import "github.com/eldief/go100x/src/types"

const (
	ORDER_TYPE_LIMIT             types.OrderType = 0
	ORDER_TYPE_LIMIT_MAKER       types.OrderType = 1
	ORDER_TYPE_MARKET            types.OrderType = 2
	ORDER_TYPE_STOP_LOSS         types.OrderType = 3
	ORDER_TYPE_STOP_LOSS_LIMIT   types.OrderType = 4
	ORDER_TYPE_TAKE_PROFIT       types.OrderType = 5
	ORDER_TYPE_TAKE_PROFIT_LIMIT types.OrderType = 6
)
