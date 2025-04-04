package usd0pp

import (
	"errors"
)

const (
	DexType = "usd0pp"
	USD0PP  = "0x35d8949372d46b7a3d5a56006ae77b215fc69bc0"
	USD0    = "0x73a15fed60bf67631dc6cd7bc5b6e8da8190acf5"
)

const (
	defaultReserves = "1000000000000000000000000"
)

var (
	defaultGas = Gas{
		Mint: 200000,
	}
)

const (
	usd0ppMethodPaused       = "paused"
	usd0ppMethodGetEndTime   = "getEndTime"
	usd0ppMethodGetStartTime = "getStartTime"
)

var (
	ErrPoolPaused             = errors.New("pool is paused")
	ErrBondNotStarted         = errors.New("bond not started")
	ErrBondEnded              = errors.New("bond ended")
	ErrorInvalidTokenIn       = errors.New("invalid tokenIn")
	ErrorInvalidTokenInAmount = errors.New("invalid tokenIn amount")
)
