/*
Package orm is a convenient object to data store mapper.
*/
package orm

import (
	"io"
	"reflect"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const ormCodespace = "orm"

var (
	ErrNotFound         = errors.Register(ormCodespace, 100, "not found")
	ErrIteratorDone     = errors.Register(ormCodespace, 101, "iterator done")
	ErrIteratorInvalid  = errors.Register(ormCodespace, 102, "iterator invalid")
	ErrType             = errors.Register(ormCodespace, 110, "invalid type")
	ErrUniqueConstraint = errors.Register(ormCodespace, 111, "unique constraint violation")
	ErrArgument         = errors.Register(ormCodespace, 112, "invalid argument")
)

// HasKVStore is a subset of the cosmos-sdk context defined for loose coupling and simpler test setups.
type HasKVStore interface {
	KVStore(key sdk.StoreKey) sdk.KVStore
}

// Persistent supports Marshal and Unmarshal
//
// This is separated from Marshal, as this almost always requires
// a pointer, and functions that only need to marshal bytes can
// use the Marshaller interface to access non-pointers.
//
// As with Marshaller, this may do internal validation on the data
// and errors should be expected.
type Persistent interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// Index allows efficient prefix scans is stored as key = concat(indexKeyBytes, rowIDUint64) with value empty
// so that the row NaturalKey is allows a fixed with 8 byte integer. This allows the MultiKeyIndex key bytes to be
// variable length and scanned iteratively. The
type Index interface {
	Has(ctx HasKVStore, key []byte) bool
	Get(ctx HasKVStore, key []byte) (Iterator, error)
	PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
}

// Iterator allows iteration through a sequence of key value pairs
type Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the key. If there
	// are no more items the ErrIteratorDone error is returned
	// The key is the rowID and not any MultiKeyIndex key.
	LoadNext(dest Persistent) (key []byte, err error)
	// Close releases the iterator and should be called at the end of iteration
	io.Closer
}

// Indexable types are used to setup new tables.
// This interface provides a set of functions that can be called by indexes to register and interact with the tables.
type Indexable interface {
	StoreKey() sdk.StoreKey
	AddAfterSaveInterceptor(interceptor AfterSaveInterceptor)
	AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor)
	RowGetter() RowGetter
}

// AfterSaveInterceptor defines a callback function to be called on Create + Update.
type AfterSaveInterceptor func(ctx HasKVStore, rowID uint64, newValue, oldValue Persistent) error

// AfterDeleteInterceptor defines a callback function to be called on Delete operations.
type AfterDeleteInterceptor func(ctx HasKVStore, rowID uint64, value Persistent) error

// RowGetter loads a persistent object by row ID into the destination object. The dest parameter must therefore be a pointer.
// The key returned is the serialized row ID.
// Any implementation must return `ErrNotFound` when no object for the rowID exists
type RowGetter func(ctx HasKVStore, rowID uint64, dest Persistent) (key []byte, err error)

// NewTypeSafeRowGetter returns a `RowGetter` with type check on the dest parameter.
func NewTypeSafeRowGetter(storeKey sdk.StoreKey, prefixKey byte, model reflect.Type) RowGetter {
	return func(ctx HasKVStore, rowID uint64, dest Persistent) ([]byte, error) {
		if err := assertCorrectType(model, dest); err != nil {
			return nil, err
		}
		store := prefix.NewStore(ctx.KVStore(storeKey), []byte{prefixKey})
		key := EncodeSequence(rowID)
		it := store.Iterator(key, EncodeSequence(rowID+1))
		if !it.Valid() {
			return nil, ErrNotFound
		}
		defer it.Close()
		return key, dest.Unmarshal(it.Value())
	}
}

func assertCorrectType(model reflect.Type, obj Persistent) error {
	tp := reflect.TypeOf(obj)
	if tp.Kind() != reflect.Ptr {
		return errors.Wrap(ErrType, "model destination must be a pointer")
	}
	if model != tp.Elem() {
		return errors.Wrapf(ErrType, "can not use %T with this bucket", obj)
	}
	return nil
}
