package weeth

import (
	"context"
	"math/big"
	"time"

	"github.com/KyberNetwork/ethrpc"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/goccy/go-json"

	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/liquidity-source/etherfi/common"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool"
	pooltrack "github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/pool/tracker"
)

type PoolTracker struct {
	ethrpcClient *ethrpc.Client
}

var _ = pooltrack.RegisterFactoryE0(DexType, NewPoolTracker)

func NewPoolTracker(ethrpcClient *ethrpc.Client) *PoolTracker {
	return &PoolTracker{
		ethrpcClient: ethrpcClient,
	}
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
	overrides map[gethcommon.Address]gethclient.OverrideAccount,
) (entity.Pool, error) {
	extra, blockNumber, err := t.getExtra(ctx, overrides)
	if err != nil {
		return p, err
	}

	extraBytes, err := json.Marshal(extra)
	if err != nil {
		return p, err
	}

	p.Reserves = []string{reserves, reserves}
	p.Extra = string(extraBytes)
	p.BlockNumber = blockNumber
	p.Timestamp = time.Now().Unix()

	return p, nil
}

func (t *PoolTracker) getExtra(
	ctx context.Context,
	overrides map[gethcommon.Address]gethclient.OverrideAccount,
) (PoolExtra, uint64, error) {
	var (
		totalShares      *big.Int
		totalPooledEther *big.Int
	)

	getPoolStateRequest := t.ethrpcClient.NewRequest().SetContext(ctx)
	if overrides != nil {
		getPoolStateRequest.SetOverrides(overrides)
	}

	getPoolStateRequest.AddCall(&ethrpc.Call{
		ABI:    common.LiquidityPoolABI,
		Target: common.LiquidityPool,
		Method: common.LiquidityPoolMethodGetTotalPooledEther,
		Params: []interface{}{},
	}, []interface{}{&totalPooledEther})

	getPoolStateRequest.AddCall(&ethrpc.Call{
		ABI:    common.EETHABI,
		Target: common.EETH,
		Method: common.EETHMethodTotalShares,
		Params: []interface{}{},
	}, []interface{}{&totalShares})

	resp, err := getPoolStateRequest.TryAggregate()
	if err != nil {
		return PoolExtra{}, 0, err
	}
	if resp.BlockNumber == nil {
		resp.BlockNumber = big.NewInt(0)
	}

	return PoolExtra{
		TotalPooledEther: totalPooledEther,
		TotalShares:      totalShares,
	}, resp.BlockNumber.Uint64(), nil
}
