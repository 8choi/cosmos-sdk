package lib

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Mapping interface {
	Get(interface{}, interface{}) error
	Has(interface{}) bool
	Set(interface{}, interface{})
	Delete(interface{})
}

type mapping struct {
	cdc   *codec.Codec
	store sdk.KVStore
}

func NewMapping(cdc *codec.Codec, store sdk.KVStore) Mapping {
	return mapping{
		cdc:   cdc,
		store: store,
	}
}

func NewPrimitiveMapping(store sdk.KVStore) Mapping {
	return mapping{
		cdc:   codec.New(),
		store: store,
	}
}

func (m mapping) Get(key interface{}, ptr interface{}) error {
	bz := m.store.Get(m.cdc.MustMarshalBinary(key))
	return m.cdc.UnmarshalBinary(bz, ptr)
}

func (m mapping) Has(key interface{}) bool {
	return m.store.Has(m.cdc.MustMarshalBinary(key))
}

func (m mapping) Set(key interface{}, val interface{}) {
	m.store.Set(m.cdc.MustMarshalBinary(key), m.cdc.MustMarshalBinary(val))
}

func (m mapping) Delete(key interface{}) {
	m.store.Delete(m.cdc.MustMarshalBinary(key))
}
