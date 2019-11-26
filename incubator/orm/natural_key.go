package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewNaturalKeyTableBuilder(prefix byte, key sdk.StoreKey, cdc *codec.Codec, getPrimaryKey func(val interface{}) []byte) NaturalKeyTableBuilder {
	return naturalKeyTableBuilder{prefix: prefix, key: key, cdc: cdc, getPrimaryKey: getPrimaryKey}

}

type naturalKeyTableBuilder struct {
	prefix        byte
	key           sdk.StoreKey
	cdc           *codec.Codec
	getPrimaryKey func(val interface{}) []byte
}

func (n naturalKeyTableBuilder) RegisterIndexer(prefix byte, indexer Indexer) {
	panic("implement me")
}

func (n naturalKeyTableBuilder) Build() NaturalKeyTable {
	panic("implement me")
}



