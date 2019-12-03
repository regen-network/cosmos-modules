package orm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
)

type HasKVStore interface {
	KVStore(key sdk.StoreKey) sdk.KVStore
}

// Index allows efficient prefix scans is stored as storeKey = concat(indexKeyBytes, rowIDUint64) with value empty
// so that the row ID is allows a fixed with 8 byte integer. This allows the index storeKey bytes to be
// variable length and scanned iteratively. The
type Index interface {
	Has(ctx HasKVStore, key []byte) (bool, error)
	Get(ctx HasKVStore, key []byte) (Iterator, error)
	PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
}

// UniqueIndex is stored as storeKey = indexKey, and value = rowId
type UniqueIndex interface {
	Index
	GetOne(ctx HasKVStore, indexKey []byte, dest interface{}) (primaryKey []byte, error error)
}

type AfterSaveInterceptor = func(ctx HasKVStore, rowId uint64, key []byte, value interface{}) error
type AfterDeleteInterceptor = func(ctx HasKVStore, rowId uint64, key []byte) error
type RowGetter = func(ctx HasKVStore, rowId uint64) (interface{}, error)

type TableBuilder interface {
	AddAfterSaveInterceptor(interceptor AfterSaveInterceptor)
	AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor)
	RowGetter() RowGetter
}

type UInt64Index interface {
	Has(ctx HasKVStore, key uint64) (bool, error)
	Get(ctx HasKVStore, key uint64) (Iterator, error)
	PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
}

type Table interface {
	UniqueIndex
	Save(ctx HasKVStore, value interface{}) error
}

//Iterator allows iteration through a sequence of storeKey value pairs
type Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the storeKey. If there
	// are no more items an error is returned
	LoadNext(dest interface{}) (key []byte, err error)
	// Close releases the iterator and should be called at the end of iteration
	io.Closer
}

type UInt64Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the storeKey. If there
	// are no more items an error is returned
	LoadNext(dest interface{}) (key uint64, err error)
	// Close releases the iterator and should be called at the end of iteration
	io.Closer
}

type Sequence interface {
	NextVal(ctx HasKVStore) (uint64, error)
	CurVal(ctx HasKVStore) (uint64, error)
}

type SchemaManager interface {
	RegisterSchemaObject(name string, descriptor SchemaDescriptor) sdk.StoreKey
}

type SchemaDescriptor interface {
	// TODO
}
type Indexer interface {
	DoIndex(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error
	BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index
}

// TableBase provides methods shared by all tables
type TableBase interface {
	UniqueIndex
	// Delete deletes the value at the given storeKey
	Delete(ctx HasKVStore, key []byte) error
}
//
//// ExternalKeyTable defines a bucket where the storeKey is stored externally to the value object
//type ExternalKeyTable interface {
//	TableBase
//	// Save saves the given storeKey value pair
//	Save(ctx HasKVStore, storeKey []byte, value interface{}) error
//}
//
type HasID interface {
	ID() []byte
}
//
//// NaturalKeyTable defines a bucket where all values implement HasID and the storeKey is stored it the value and
//// returned by the HasID method
type NaturalKeyTable interface {
	TableBase
	// Save saves the value passed in
	Save(ctx HasKVStore, value HasID) error
}
//
type AutoUInt64Table interface {
	Has(ctx HasKVStore, key uint64) (bool, error)
	// TODO: replace iterator by value arg, only 0..1 entitie can exist
	// TODO: Iterator does return storeKey on load which is not uint64 type. Replace with custom iterator impl?
	Get(ctx HasKVStore, key uint64) (Iterator, error)
	PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	// Create stores the given value and returns the auto generated primary storeKey used.
	Create(ctx HasKVStore, value interface{}) (uint64, error)
	// Save updates the entry for the the given storeKey. The storeKey must not be empty.
	// When no entry for the storeKey exists, an Error is returned.
	Save(ctx HasKVStore, key uint64, value interface{}) error
}
//
//// AutoKeyTable specifies a bucket where keys are generated via an auto-incremented interger
//type AutoKeyTable interface {
//	ExternalKeyTable
//
//	// Create auto-generates storeKey
//	Create(ctx HasKVStore, value interface{}) ([]byte, error)
//}
//
//type TableInterceptor interface {
//	OnRead(ctx HasKVStore, value interface{}) error
//	BeforeSave(ctx HasKVStore, rowId uint64, value interface{}) error
//	AfterSave(ctx HasKVStore, rowId uint64, value interface{}) error
//	BeforeDelete(ctx HasKVStore, rowId uint64, value interface{}) error
//	AfterDelete(ctx HasKVStore, rowId uint64, value interface{}) error
//}


//type TableBuilder interface {
//	RegisterIndexer(prefix byte, indexer Indexer)
//}
//
//type NaturalKeyTableBuilder interface {
//	TableBuilder
//	Build() NaturalKeyTable
//}
//
//type AutoUInt64TableBuilder interface {
//	RegisterIndexer(prefix byte, indexer Indexer)
//	Build() AutoUInt64Table
//}

