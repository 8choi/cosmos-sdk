package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/types"
)

// Keeper defines the keeper of the supply store
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

// NewKeeper creates a new supply Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

// GetSupplier retrieves the Supplier from store
func (k Keeper) GetSupplier(ctx sdk.Context) (supplier types.Supplier) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(supplierKey)
	if b == nil {
		panic("Stored supplier should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supplier)
	return
}

// SetSupplier sets the Supplier to store
func (k Keeper) SetSupplier(ctx sdk.Context, supplier types.Supplier) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supplier)
	store.Set(supplierKey, b)
}

// InflateSupply adds tokens to the supplier
func (k Keeper) InflateSupply(ctx sdk.Context, supplyType string, amount sdk.Coins) {
	supplier := k.GetSupplier(ctx)
	supplier.Inflate(supplyType, amount)

	k.SetSupplier(ctx, supplier)
}

// DeflateSupply subtracts tokens to the suplier
func (k Keeper) DeflateSupply(ctx sdk.Context, supplyType string, amount sdk.Coins) {
	supplier := k.GetSupplier(ctx)
	supplier.Deflate(supplyType, amount)

	k.SetSupplier(ctx, supplier)
}
