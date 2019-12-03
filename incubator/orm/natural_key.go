package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewNaturalKeyTableBuilder(prefix byte, key sdk.StoreKey, cdc *codec.Codec, getPrimaryKey func(val interface{}) []byte) *naturalKeyTableBuilder {
	return &naturalKeyTableBuilder{prefix: prefix, key: key, cdc: cdc, getPrimaryKey: getPrimaryKey}

}

type naturalKeyTableBuilder struct {
	prefix        byte
	key           sdk.StoreKey
	cdc           *codec.Codec
	getPrimaryKey func(val interface{}) []byte
}

func (n *naturalKeyTableBuilder) AddAfterSaveInterceptor(interceptor AfterSaveInterceptor) {
	panic("implement me")
}

func (n *naturalKeyTableBuilder) AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor) {
	panic("implement me")
}

func (n *naturalKeyTableBuilder) RowGetter() RowGetter {
	panic("implement me")
}

func (n *naturalKeyTableBuilder) RegisterIndexer(prefix byte, indexer Indexer) {
	panic("implement me")
}

func (n naturalKeyTableBuilder) Build() naturalKeyTable {
	panic("implement me")
}

type naturalKeyTable struct {
}

func (a naturalKeyTable) Has(ctx HasKVStore, key uint64) (bool, error) {
	panic("implement me")
}

func (a naturalKeyTable) Get(ctx HasKVStore, key uint64) (Iterator, error) {
	panic("implement me")
}

func (a naturalKeyTable) PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	panic("implement me")
}

func (a naturalKeyTable) ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	panic("implement me")
}

func (a naturalKeyTable) Save(ctx HasKVStore, key []byte, value interface{}) error {
	panic("implement me")
}





