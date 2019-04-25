package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/types"
)

var oneUatom = sdk.NewCoins(sdk.NewCoin("uatom", sdk.OneInt()))

func TestSupplier(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	require.Panics(t, func() { keeper.GetSupplier(ctx) }, "should panic when supplier is not set")

	expectedSupplier := types.DefaultSupplier()
	keeper.SetSupplier(ctx, expectedSupplier)

	// test inflation
	expectedSupplier.Inflate(types.TypeCirculating, oneUatom)
	keeper.InflateSupply(ctx, types.TypeCirculating, oneUatom)

	supplier := keeper.GetSupplier(ctx)
	require.Equal(t, expectedSupplier.CirculatingSupply, supplier.CirculatingSupply)

	// test deflation
	expectedSupplier.Deflate(types.TypeCirculating, oneUatom)
	keeper.DeflateSupply(ctx, types.TypeCirculating, oneUatom)

	supplier = keeper.GetSupplier(ctx)
	require.Equal(t, expectedSupplier.CirculatingSupply, supplier.CirculatingSupply)
}

func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {

	keySupply := sdk.NewKVStoreKey(StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "supply-chain"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)

	cdc := codec.New()

	keeper := NewKeeper(cdc, keySupply)

	return ctx, keeper
}
