package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAutoUInt64TableBuilder(prefix byte, key sdk.StoreKey, cdc *codec.Codec) AutoUInt64TableBuilder {
	return autoUInt64TableBuilder{prefix: prefix, key: key, cdc: cdc}
}

type autoUInt64TableBuilder struct {
	prefix      byte
	key         sdk.StoreKey
	cdc         *codec.Codec
	indexerRefs []indexRef
}

func (a autoUInt64TableBuilder) RegisterIndexer(prefix byte, indexer Indexer) {
	a.indexerRefs = append(a.indexerRefs, indexRef{prefix: prefix, indexer: indexer})
}

func (a autoUInt64TableBuilder) Build() AutoUInt64Table {
	return autoUInt64Table{}
}

type autoUInt64Table struct {
}

func (a autoUInt64Table) Has(ctx HasKVStore, key uint64) (bool, error) {
	panic("implement me")
}

func (a autoUInt64Table) Get(ctx HasKVStore, key uint64) (Iterator, error) {
	panic("implement me")
}

func (a autoUInt64Table) PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	panic("implement me")
}

func (a autoUInt64Table) ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	panic("implement me")
}

func (a autoUInt64Table) Save(ctx HasKVStore, key []byte, value interface{}) error {
	panic("implement me")
}

