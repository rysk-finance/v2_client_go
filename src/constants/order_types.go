package constants

import "go100x/src/types"

const (
	LIMIT             types.OrderType = 0
	LIMIT_MAKER       types.OrderType = 1
	MARKET            types.OrderType = 2
	STOP_LOSS         types.OrderType = 3
	STOP_LOSS_LIMIT   types.OrderType = 4
	TAKE_PROFIT       types.OrderType = 5
	TAKE_PROFIT_LIMIT types.OrderType = 6
)
