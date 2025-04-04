package sfrxeth

import (
	"context"
	"math/big"
	"time"

	"github.com/KyberNetwork/ethrpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/goccy/go-json"

	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	frax_common "github.com/KyberNetwork/kyberswap-dex-lib/pkg/liquidity-source/frax/common"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
	pooltrack "github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool/tracker"
)

type PoolTracker struct {
	config       *Config
	ethrpcClient *ethrpc.Client
}

var _ = pooltrack.RegisterFactoryCE(DexType, NewPoolTracker)

func NewPoolTracker(config *Config, ethrpcClient *ethrpc.Client) (*PoolTracker, error) {
	return &PoolTracker{
		config:       config,
		ethrpcClient: ethrpcClient,
	}, nil
}

func (t *PoolTracker) GetNewPoolState(
	ctx context.Context,
	p entity.Pool,
	params pool.GetNewPoolStateParams,
) (entity.Pool, error) {
	return t.getNewPoolState(ctx, p, params, nil)
}

func (t *PoolTracker) GetNewPoolStateWithOverrides(
	ctx context.Context,
	p entity.Pool,
	params pool.GetNewPoolStateWithOverridesParams,
) (entity.Pool, error) {
	return t.getNewPoolState(ctx, p, pool.GetNewPoolStateParams{Logs: params.Logs}, params.Overrides)
}

func (t *PoolTracker) getNewPoolState(
	ctx context.Context,
	p entity.Pool,
	_ pool.GetNewPoolStateParams,
	overrides map[common.Address]gethclient.OverrideAccount,
) (entity.Pool, error) {
	totalSupply, totalAssets, extra, blockNumber, err := getState(
		ctx, p.Address, p.Tokens[1].Address, t.ethrpcClient,
		overrides,
	)
	if err != nil {
		return p, err
	}

	extraBytes, err := json.Marshal(extra)
	if err != nil {
		return p, err
	}

	p.Reserves = entity.PoolReserves{totalAssets.String(), totalSupply.String()}
	p.Extra = string(extraBytes)
	p.BlockNumber = blockNumber
	p.Timestamp = time.Now().Unix()

	return p, nil
}

func getState(
	ctx context.Context,
	minterAddress string,
	sfrxETHAddress string,
	ethrpcClient *ethrpc.Client,
	overrides map[common.Address]gethclient.OverrideAccount,
) (*big.Int, *big.Int, PoolExtra, uint64, error) {
	var (
		submitPaused bool
		totalSupply  *big.Int
		totalAssets  *big.Int
	)

	calls := ethrpcClient.NewRequest().SetContext(ctx)
	if overrides != nil {
		calls.SetOverrides(overrides)
	}

	calls.AddCall(&ethrpc.Call{
		ABI:    frax_common.FrxETHMinterABI,
		Target: minterAddress,
		Method: minterMethodSubmitPaused,
	}, []interface{}{&submitPaused})
	calls.AddCall(&ethrpc.Call{
		ABI:    frax_common.SfrxETHABI,
		Target: sfrxETHAddress,
		Method: SfrxETHMethodTotalAssets,
	}, []interface{}{&totalAssets})
	calls.AddCall(&ethrpc.Call{
		ABI:    frax_common.SfrxETHABI,
		Target: sfrxETHAddress,
		Method: SfrxETHMethodTotalSupply,
	}, []interface{}{&totalSupply})

	resp, err := calls.Aggregate()
	if err != nil {
		return nil, nil, PoolExtra{}, 0, err
	}

	if resp.BlockNumber == nil {
		resp.BlockNumber = big.NewInt(0)
	}

	poolExtra := PoolExtra{
		SubmitPaused: submitPaused,
	}

	return totalSupply, totalAssets, poolExtra, resp.BlockNumber.Uint64(), nil
}
