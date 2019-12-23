package orm

import (
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type iteratorFunc func(dest interface{}) (key []byte, err error)

func (i iteratorFunc) LoadNext(dest interface{}) (key []byte, err error) {
	return i(dest)
}

func (i iteratorFunc) Close() error {
	return nil
}

func NewSingleValueIterator(cdc *codec.Codec, rowID []byte, val []byte) Iterator {
	var closed bool
	return iteratorFunc(func(dest interface{}) ([]byte, error) {
		if closed || val == nil {
			return nil, ErrIteratorDone
		}
		closed = true
		return rowID, cdc.UnmarshalBinaryBare(val, dest)
	})
}

// Iterator that return ErrIteratorInvalid only.
func NewInvalidIterator() Iterator {
	return iteratorFunc(func(dest interface{}) ([]byte, error) {
		return nil, ErrIteratorInvalid
	})
}

// First loads the first element into the given destination type and closes the iterator.
// When the iterator is closed or has no elements the according error is passed as return value.
func First(it Iterator, dest interface{}) ([]byte, error) {
	if it == nil {
		return nil, errors.Wrap(ErrArgument, "iterator must not be nil")
	}
	defer it.Close()
	binKey, err := it.LoadNext(dest)
	if err != nil {
		return nil, err
	}
	return binKey, nil
}

// ModelSlicePtr represents a pointer to a slice of models. Think of it as
// *[]Model Because of Go's type system, using []Model type would not work for us.
// Instead we use a placeholder type and the validation is done during the
// runtime.
type ModelSlicePtr interface{}

// ReadAll consumes all values for the iterator and stores them in a new slice at the passed ModelSlicePtr.
// The slice can be empty when the iterator does not return any values but not nil. The iterator
// is closed afterwards.
// Example:
// 			var loaded []GroupMetadata
//			rowIDs, err := ReadAll(it, &loaded)
//			require.NoError(t, err)
//
func ReadAll(it Iterator, dest ModelSlicePtr) ([][]byte, error) {
	if it == nil {
		return nil, errors.Wrap(ErrArgument, "iterator must not be nil")
	}
	defer it.Close()
	if dest == nil {
		return nil, errors.Wrap(ErrArgument, "destination must not be nil")
	}
	tp := reflect.ValueOf(dest)
	if tp.Kind() != reflect.Ptr {
		return nil, errors.Wrap(ErrArgument, "destination must be a pointer to a slice")
	}
	if tp.Elem().Kind() != reflect.Slice {
		return nil, errors.Wrap(ErrArgument, "destination must point to a slice")
	}
	slice := tp.Elem()
	if !slice.CanSet() {
		return nil, errors.Wrap(ErrArgument, "destination not assignable")
	}

	typ := reflect.TypeOf(dest).Elem().Elem()
	t := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	var rowIDs [][]byte
	for {
		obj := reflect.New(typ)
		binKey, err := it.LoadNext(obj.Interface())
		switch {
		case err == nil:
			t = reflect.Append(t, obj.Elem())
		case ErrIteratorDone.Is(err):
			slice.Set(t)
			return rowIDs, nil
		default:
			return nil, err
		}
		rowIDs = append(rowIDs, binKey)
	}
	return rowIDs, nil
}
