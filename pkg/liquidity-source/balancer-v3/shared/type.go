package shared

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type StaticExtra struct {
	HookType     HookType `json:"hookT,omitempty"`
	BufferTokens []string `json:"buffs,omitempty"`
}

type ExtraBuffer struct {
	TotalAssets *uint256.Int `json:"tA"`
	TotalSupply *uint256.Int `json:"tS"`
}

type SwapInfo struct {
	Buffers      []*ExtraBuffer
	AggregateFee *big.Int
}

type PoolMetaInfo struct {
	BufferTokenIn  string `json:"buffIn"`
	BufferTokenOut string `json:"buffOut"`
}

type AggregateFeePercentage struct {
	AggregateSwapFeePercentage  *big.Int
	AggregateYieldFeePercentage *big.Int
}

type VaultSwapParams struct {
	Kind           SwapKind
	IndexIn        int
	IndexOut       int
	AmountGivenRaw *uint256.Int
}

type PoolSwapParams struct {
	Kind                 SwapKind
	SwapFeePercentage    *uint256.Int
	AmountGivenScaled18  *uint256.Int
	BalancesLiveScaled18 []*uint256.Int
	IndexIn              int
	IndexOut             int
}

type AfterSwapParams struct {
	Kind                     SwapKind
	IndexIn                  int
	IndexOut                 int
	AmountInScaled18         *uint256.Int
	AmountOutScaled18        *uint256.Int
	TokenInBalanceScaled18   *uint256.Int
	TokenOutBalanceScaled18  *uint256.Int
	AmountCalculatedScaled18 *uint256.Int
	AmountCalculatedRaw      *uint256.Int
}

type TokenInfo struct {
	TokenType     uint8
	RateProvider  common.Address
	PaysYieldFees bool
}

type PoolDataRPC struct {
	Data struct {
		PoolConfigBits        [32]byte
		Tokens                []common.Address
		TokenInfo             []TokenInfo
		BalancesRaw           []*big.Int
		BalancesLiveScaled18  []*big.Int
		TokenRates            []*big.Int
		DecimalScalingFactors []*big.Int
	}
}

type ExtraBufferRPC struct {
	TotalAssets *big.Int
	TotalSupply *big.Int
}

type HooksConfig struct {
	EnableHookAdjustedAmounts       bool `json:"enableHookAdjustedAmounts"`
	ShouldCallComputeDynamicSwapFee bool `json:"shouldCallComputeDynamicSwapFee"`
	ShouldCallBeforeSwap            bool `json:"shouldCallBeforeSwap"`
	ShouldCallAfterSwap             bool `json:"shouldCallAfterSwap"`
}

type HooksConfigRPC struct {
	Data struct {
		EnableHookAdjustedAmounts       bool
		ShouldCallBeforeInitialize      bool
		ShouldCallAfterInitialize       bool
		ShouldCallComputeDynamicSwapFee bool
		ShouldCallBeforeSwap            bool
		ShouldCallAfterSwap             bool
		ShouldCallBeforeAddLiquidity    bool
		ShouldCallAfterAddLiquidity     bool
		ShouldCallBeforeRemoveLiquidity bool
		ShouldCallAfterRemoveLiquidity  bool
		HooksContract                   common.Address
	}
}
