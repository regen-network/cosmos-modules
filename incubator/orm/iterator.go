package orm

import "github.com/cosmos/cosmos-sdk/codec"

type iteratorFunc func(dest interface{}) (key []byte, err error)

func (i iteratorFunc) LoadNext(dest interface{}) (key []byte, err error) {
	return i(dest)
}

func (i iteratorFunc) Close() error {
	return nil
}

func NewSingleValueIterator(cdc *codec.Codec, key []byte, val []byte) Iterator {
	var closed bool
	return iteratorFunc(func(dest interface{}) ([]byte, error) {
		if closed || val == nil {
			return nil, ErrIteratorDone
		}
		closed = true
		return key, cdc.UnmarshalBinaryBare(val, dest)
	})
}

// First loads the first element into the given destination type and closes the iterator.
// When the iterator is closed or has no elements the according error is passed as return value.
func First(it Iterator, dest interface{}) ([]byte, error) {
	defer it.Close()
	binKey, err := it.LoadNext(dest)
	if err != nil {
		return nil, err
	}
	return binKey, nil
}
