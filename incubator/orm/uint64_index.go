package orm

// UInt64IndexerFunc creates one or multiple multiKeyIndex keys of type uint64 for the source object.
type UInt64IndexerFunc func(value interface{}) ([]uint64, error)

// UInt64MultiKeyAdapter converts UInt64IndexerFunc to IndexerFunc
func UInt64MultiKeyAdapter(indexer UInt64IndexerFunc) IndexerFunc {
	return func(value interface{}) ([]RowID, error) {
		d, err := indexer(value)
		if err != nil {
			return nil, err
		}
		r := make([]RowID, len(d))
		for i, v := range d {
			r[i] = EncodeSequence(v)
		}
		return r, nil
	}
}

// UInt64Index is a typed index.
type UInt64Index struct {
	multiKeyIndex *MultiKeyIndex
}

// NewUInt64Index creates a typed secondary index
func NewUInt64Index(builder Indexable, prefix byte, indexer UInt64IndexerFunc) *UInt64Index {
	idx := UInt64Index{
		multiKeyIndex: &MultiKeyIndex{
			storeKey:  builder.StoreKey(),
			prefix:    prefix,
			rowGetter: builder.RowGetter(),
			indexer:   NewIndexer(UInt64MultiKeyAdapter(indexer)),
		},
	}
	builder.AddAfterSaveInterceptor(idx.multiKeyIndex.onSave)
	builder.AddAfterDeleteInterceptor(idx.multiKeyIndex.onDelete)
	return &idx
}

func (i UInt64Index) Has(ctx HasKVStore, key uint64) bool {
	return i.multiKeyIndex.Has(ctx, EncodeSequence(key))
}

func (i UInt64Index) Get(ctx HasKVStore, key uint64) (Iterator, error) {
	return i.multiKeyIndex.Get(ctx, EncodeSequence(key))
}

func (i UInt64Index) PrefixScan(ctx HasKVStore, start, end uint64) (Iterator, error) {
	return i.multiKeyIndex.PrefixScan(ctx, EncodeSequence(start), EncodeSequence(end))
}

func (i UInt64Index) ReversePrefixScan(ctx HasKVStore, start, end uint64) (Iterator, error) {
	return i.multiKeyIndex.ReversePrefixScan(ctx, EncodeSequence(start), EncodeSequence(end))
}
