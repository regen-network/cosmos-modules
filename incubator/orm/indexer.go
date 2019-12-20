package orm

import (
	"encoding/binary"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IndexerFunc creates one or multiple index keys for the source object.
type IndexerFunc func(value interface{}) ([][]byte, error)

// Indexer manages the persistence for an index based on searchable keys and operations.
type Indexer struct {
	indexerFunc IndexerFunc
	addPolicy   func(store sdk.KVStore, secondaryIndexKey []byte, rowId uint64) error
}

func NewIndexer(indexerFunc IndexerFunc, unique bool) *Indexer {
	if indexerFunc == nil {
		panic("indexer func must not be nil")
	}
	i := &Indexer{indexerFunc: indexerFunc}
	if unique {
		i.addPolicy = uniqueKeysAddPolicy
	} else {
		i.addPolicy = multiKeyAddPolicy
	}
	return i
}

func (u Indexer) OnCreate(store sdk.KVStore, rowId uint64, value interface{}) error {
	secondaryIndexKeys, err := u.indexerFunc(value)
	if err != nil {
		return err
	}

	for _, secondaryIndexKey := range secondaryIndexKeys {
		if err := u.addPolicy(store, secondaryIndexKey, rowId); err != nil {
			return err
		}
	}
	return nil
}

func (u Indexer) OnDelete(store sdk.KVStore, rowId uint64, value interface{}) error {
	secondaryIndexKeys, err := u.indexerFunc(value)
	if err != nil {
		return err
	}
	for _, secondaryIndexKey := range secondaryIndexKeys {
		indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowId)
		store.Delete(indexKey)
	}
	return nil
}

func (u Indexer) OnUpdate(store sdk.KVStore, rowId uint64, newValue, oldValue interface{}) error {
	oldSecIdxKeys, err := u.indexerFunc(oldValue)
	if err != nil {
		return err
	}
	newSecIdxKeys, err := u.indexerFunc(newValue)
	if err != nil {
		return err
	}
	for _, oldIdxKey := range difference(oldSecIdxKeys, newSecIdxKeys) {
		store.Delete(makeIndexPrefixScanKey(oldIdxKey, rowId))
	}
	for _, newIdxKey := range difference(newSecIdxKeys, oldSecIdxKeys) {
		if err := u.addPolicy(store, newIdxKey, rowId); err != nil {
			return err
		}
	}
	return nil
}

// uniqueKeysAddPolicy enforces keys to be unique
func uniqueKeysAddPolicy(store sdk.KVStore, secondaryIndexKey []byte, rowId uint64) error {
	it := store.Iterator(makeIndexPrefixScanKey(secondaryIndexKey, 0), makeIndexPrefixScanKey(secondaryIndexKey, math.MaxUint64))
	defer it.Close()
	if it.Valid() {
		return ErrUniqueConstraint
	}

	indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowId)
	store.Set(indexKey, []byte{0})
	return nil
}

// multiKeyAddPolicy allows multiple entries for a key
func multiKeyAddPolicy(store sdk.KVStore, secondaryIndexKey []byte, rowId uint64) error {
	indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowId)
	if !store.Has(indexKey) {
		store.Set(indexKey, []byte{0})
	}
	return nil
}

// difference returns the list of elements that are in a but not in b.
func difference(a [][]byte, b [][]byte) [][]byte {
	set := make(map[string]struct{}, len(b))
	for _, v := range b {
		set[string(v)] = struct{}{}
	}
	var result [][]byte
	for _, v := range a {
		if _, ok := set[string(v)]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// todo: this feels quite complicated when reading the data. Why not store rowID(s) as payload instead?
func makeIndexPrefixScanKey(indexKey []byte, rowId uint64) []byte {
	n := len(indexKey)
	res := make([]byte, n+8)
	copy(res, indexKey)
	binary.BigEndian.PutUint64(res[n:], rowId)
	return res
}
