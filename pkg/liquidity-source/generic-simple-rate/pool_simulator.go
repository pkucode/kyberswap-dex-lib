package generic_simple_rate

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/goccy/go-json"
	"github.com/holiman/uint256"

	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/util/bignumber"
)

type PoolSimulator struct {
	pool.Pool
	paused          bool
	gas             int64
	rate            *uint256.Int
	rateUnit        *uint256.Int
	isRateInversed  bool
	isBidirectional bool
}

var _ = pool.RegisterFactory0(DexType, NewPoolSimulator)

func NewPoolSimulator(entityPool entity.Pool) (*PoolSimulator, error) {
	numTokens := len(entityPool.Tokens)
	if numTokens != 2 {
		return nil, fmt.Errorf("invalid pool tokens %v, %v", entityPool, numTokens)
	}
	if numTokens != len(entityPool.Reserves) {
		return nil, fmt.Errorf("invalid pool reserves %v, %v", entityPool, numTokens)
	}

	var tokens = make([]string, numTokens)
	var reserves = make([]*big.Int, numTokens)
	for i := 0; i < numTokens; i += 1 {
		tokens[i] = entityPool.Tokens[i].Address
		reserves[i] = bignumber.NewBig10(entityPool.Reserves[i])
	}

	var poolExtra PoolExtra
	if err := json.Unmarshal([]byte(entityPool.Extra), &poolExtra); err != nil {
		return nil, err
	}

	return &PoolSimulator{
		Pool: pool.Pool{
			Info: pool.PoolInfo{
				Address:  strings.ToLower(entityPool.Address),
				Exchange: entityPool.Exchange,
				Type:     entityPool.Type,
				Tokens:   tokens,
				Reserves: reserves,
			},
		},
		paused:          poolExtra.Paused,
		rate:            poolExtra.Rate,
		rateUnit:        poolExtra.RateUnit,
		isRateInversed:  poolExtra.IsRateInversed,
		isBidirectional: poolExtra.IsBidirectional,
		gas:             poolExtra.DefaultGas,
	}, nil
}

func (p *PoolSimulator) CalcAmountOut(param pool.CalcAmountOutParams) (*pool.CalcAmountOutResult, error) {
	if p.paused {
		return nil, ErrPoolPaused
	}

	tokenAmountIn := param.TokenAmountIn
	tokenOut := param.TokenOut

	var tokenInIndex = p.GetTokenIndex(tokenAmountIn.Token)
	var tokenOutIndex = p.GetTokenIndex(tokenOut)

	if tokenInIndex < 0 || tokenOutIndex < 0 {
		return &pool.CalcAmountOutResult{}, fmt.Errorf("tokenInIndex: %v or tokenOutIndex: %v is not correct", tokenInIndex, tokenOutIndex)
	}

	amountOut, err := p.calcAmountOut(tokenInIndex, tokenAmountIn.Amount)
	if err != nil {
		return nil, err
	}

	return &pool.CalcAmountOutResult{
		TokenAmountOut: &pool.TokenAmount{
			Token:  tokenOut,
			Amount: amountOut.ToBig(),
		},
		Fee: &pool.TokenAmount{
			Token:  tokenAmountIn.Token,
			Amount: bignumber.ZeroBI,
		},
		Gas: p.gas,
	}, nil
}

func (p *PoolSimulator) UpdateBalance(_ pool.UpdateBalanceParams) {
}

func (p *PoolSimulator) CanSwapTo(address string) []string {
	if strings.EqualFold(p.Info.Tokens[1], address) {
		return []string{p.Info.Tokens[0]}
	}
	if p.isBidirectional && strings.EqualFold(p.Info.Tokens[0], address) {
		return []string{p.Info.Tokens[1]}
	}
	return []string{}
}

func (p *PoolSimulator) CanSwapFrom(address string) []string {
	if strings.EqualFold(p.Info.Tokens[0], address) {
		return []string{p.Info.Tokens[1]}
	}
	if p.isBidirectional && strings.EqualFold(p.Info.Tokens[1], address) {
		return []string{p.Info.Tokens[0]}
	}
	return []string{}
}

func (p *PoolSimulator) GetMetaInfo(_ string, _ string) interface{} {
	return nil
}

func (p *PoolSimulator) calcAmountOut(tokenInIndex int, amountIn *big.Int) (*uint256.Int, error) {
	amountOut := new(uint256.Int)
	if amountOut.SetFromBig(amountIn) {
		return nil, ErrOverflow
	}
	if p.isRateInversed == (tokenInIndex == 0) {
		amountOut = amountOut.Mul(amountOut, p.rateUnit).Div(amountOut, p.rate)
	} else {
		amountOut = amountOut.Div(amountOut, p.rateUnit).Mul(amountOut, p.rate)
	}
	return amountOut, nil
}
