package orm

import (
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HasKVStore interface {
	KVStore(key sdk.StoreKey) sdk.KVStore
}

// Index allows efficient prefix scans is stored as key = concat(indexKeyBytes, rowIDUint64) with value empty
// so that the row ID is allows a fixed with 8 byte integer. This allows the index key bytes to be
// variable length and scanned iteratively. The
type Index interface {
	Has(ctx HasKVStore, key []byte) (bool, error)
	Get(ctx HasKVStore, key []byte) (Iterator, error)
	PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
}

type RowIDAwareIndex interface {
	RowID(ctx HasKVStore, key []byte) (uint64, error)
}

// UniqueIndex is stored as key = indexKey, and value = rowId
type UniqueIndex interface {
	Index
	GetOne(ctx HasKVStore, indexKey []byte, dest interface{}) (primaryKey []byte, error error)
}

type AfterSaveInterceptor = func(ctx HasKVStore, rowId uint64, key []byte, newValue, oldValue interface{}) error
type AfterDeleteInterceptor = func(ctx HasKVStore, rowId uint64, key []byte, value interface{}) error

type RowGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
type TableBuilder interface {
	StoreKey() sdk.StoreKey
	RowGetter() RowGetter
	AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor)
	AddAfterSaveInterceptor(interceptor AfterSaveInterceptor)
}

type UInt64Index interface {
	Has(ctx HasKVStore, key uint64) (bool, error)
	Get(ctx HasKVStore, key uint64) (Iterator, error)
	PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
}

// todo: what is the difference between this and TableBase
type Table interface {
	UniqueIndex
	// Save without the primKey may copy the object. For example when a new member is added, the primKeyer
	// would return a new key.
	Save(ctx HasKVStore, primKey []byte, value interface{}) error
}

//Iterator allows iteration through a sequence of key value pairs
type Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the key. If there
	// are no more items an error is returned
	// todo: Returns the "natural key" or the rowID in case of a naturalKeyTable?
	LoadNext(dest interface{}) (key []byte, err error)
	// Close releases the iterator and should be called at the end of iteration
	io.Closer
}

type UInt64Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the key. If there
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

// TODO: unused
type Indexer interface {
	OnCreate(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error
	BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index
}

// TableBase provides methods shared by all tables
// todo: do wee need this?
type TableBase interface {
	UniqueIndex
}

//
//// ExternalKeyTable defines a bucket where the key is stored externally to the value object
//type ExternalKeyTable interface {
//	TableBase
//	// Save saves the given key value pair
//	Save(ctx HasKVStore, key []byte, value interface{}) error
//}
//
type HasID interface {
	ID() []byte
}

//
//// NaturalKeyTable defines a bucket where all values implement HasID and the key is stored it the value and
//// returned by the HasID method
type NaturalKeyTable interface {
	TableBase
	// Save saves the value passed in
	Create(ctx HasKVStore, value HasID) error
	Save(ctx HasKVStore, value HasID) error
	Delete(ctx HasKVStore, value HasID) error
}

//
type AutoUInt64Table interface {
	Has(ctx HasKVStore, key uint64) (bool, error)
	// TODO: replace iterator by value arg, only 0..1 entitie can exist
	// TODO: Iterator does return key on load which is not uint64 type. Replace with custom iterator impl?
	Get(ctx HasKVStore, key uint64) (Iterator, error)
	PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	// Create stores the given value and returns the auto generated primary key used.
	Create(ctx HasKVStore, value interface{}) (uint64, error)
	// Save updates the entry for the the given key. The key must not be empty.
	// When no entry for the key exists, an Error is returned.
	Save(ctx HasKVStore, key uint64, value interface{}) error
	Delete(ctx HasKVStore, key uint64) error
}

//
//// AutoKeyTable specifies a bucket where keys are generated via an auto-incremented interger
//type AutoKeyTable interface {
//	ExternalKeyTable
//
//	// Create auto-generates key
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
