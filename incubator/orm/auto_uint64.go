package orm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Indexable = &AutoUInt64TableBuilder{}

func NewAutoUInt64TableBuilder(prefixData byte, prefixSeq byte, storeKey sdk.StoreKey, model Persistent) *AutoUInt64TableBuilder {
	if prefixData == prefixSeq {
		panic("prefixData and prefixSeq must be unique")
	}
	return &AutoUInt64TableBuilder{
		TableBuilder: NewTableBuilder(prefixData, storeKey, model),
		seq:          NewSequence(storeKey, prefixSeq),
	}
}

type AutoUInt64TableBuilder struct {
	*TableBuilder
	seq *Sequence
}

func (a AutoUInt64TableBuilder) Build() AutoUInt64Table {
	return AutoUInt64Table{
		table: a.TableBuilder.Build(),
		seq:   a.seq,
	}
}

// AutoUInt64Table is the table type which an auto incrementing ID.
type AutoUInt64Table struct {
	table Table
	seq   *Sequence
}

func (a AutoUInt64Table) Create(ctx HasKVStore, obj Persistent) (uint64, error) {
	autoIncID := a.seq.NextVal(ctx)
	err := a.table.Create(ctx, EncodeSequence(autoIncID), obj)
	if err != nil {
		return 0, err
	}
	return autoIncID, nil
}

func (a AutoUInt64Table) Save(ctx HasKVStore, rowID uint64, newValue Persistent) error {
	return a.table.Save(ctx, EncodeSequence(rowID), newValue)
}

func (a AutoUInt64Table) Delete(ctx HasKVStore, rowID uint64) error {
	return a.table.Delete(ctx, EncodeSequence(rowID))
}

func (a AutoUInt64Table) Has(ctx HasKVStore, rowID uint64) bool {
	return a.table.Has(ctx, EncodeSequence(rowID))
}

func (a AutoUInt64Table) GetOne(ctx HasKVStore, rowID uint64, dest Persistent) (RowID, error) {
	rawRowID := EncodeSequence(rowID)
	if err := a.table.GetOne(ctx, rawRowID, dest); err != nil {
		return nil, err
	}
	return rawRowID, nil
}

func (a AutoUInt64Table) PrefixScan(ctx HasKVStore, start, end uint64) (Iterator, error) {
	return a.table.PrefixScan(ctx, EncodeSequence(start), EncodeSequence(end))
}

func (a AutoUInt64Table) ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	return a.table.ReversePrefixScan(ctx, EncodeSequence(start), EncodeSequence(end))
}
