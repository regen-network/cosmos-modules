package orm

import (
	"encoding/binary"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// IndexerFunc creates one or multiple MultiKeyIndex keys for the source object.
type IndexerFunc func(value interface{}) ([][]byte, error)

// IndexerFunc creates exactly one index key for the source object.
type UniqueIndexerFunc func(value interface{}) ([]byte, error)

// Indexer manages the persistence for an MultiKeyIndex based on searchable keys and operations.
type Indexer struct {
	indexerFunc IndexerFunc
	addPolicy   func(store sdk.KVStore, secondaryIndexKey []byte, rowID uint64) error
}

// NewMultiKeyIndexer returns an indexer that supports multiple reference keys for an entity.
func NewMultiKeyIndexer(indexerFunc IndexerFunc) *Indexer {
	if indexerFunc == nil {
		panic("indexer func must not be nil")
	}
	return &Indexer{
		indexerFunc: indexerFunc,
		addPolicy:   multiKeyAddPolicy,
	}
}

// NewUniqueIndexer returns an indexer that requires exactly one reference keys for an entity.
func NewUniqueIndexer(f UniqueIndexerFunc) *Indexer {
	if f == nil {
		panic("indexer func must not be nil")
	}
	adaptor := func(indexerFunc UniqueIndexerFunc) IndexerFunc {
		return func(v interface{}) ([][]byte, error) {
			k, err := indexerFunc(v)
			return [][]byte{k}, err
		}
	}
	return &Indexer{
		indexerFunc: adaptor(f),
		addPolicy:   uniqueKeysAddPolicy,
	}
}

func (u Indexer) OnCreate(store sdk.KVStore, rowID uint64, value interface{}) error {
	secondaryIndexKeys, err := u.indexerFunc(value)
	if err != nil {
		return err
	}

	for _, secondaryIndexKey := range secondaryIndexKeys {
		if err := u.addPolicy(store, secondaryIndexKey, rowID); err != nil {
			return err
		}
	}
	return nil
}

func (u Indexer) OnDelete(store sdk.KVStore, rowID uint64, value interface{}) error {
	secondaryIndexKeys, err := u.indexerFunc(value)
	if err != nil {
		return err
	}
	for _, secondaryIndexKey := range secondaryIndexKeys {
		indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowID)
		store.Delete(indexKey)
	}
	return nil
}

func (u Indexer) OnUpdate(store sdk.KVStore, rowID uint64, newValue, oldValue interface{}) error {
	oldSecIdxKeys, err := u.indexerFunc(oldValue)
	if err != nil {
		return err
	}
	newSecIdxKeys, err := u.indexerFunc(newValue)
	if err != nil {
		return err
	}
	for _, oldIdxKey := range difference(oldSecIdxKeys, newSecIdxKeys) {
		store.Delete(makeIndexPrefixScanKey(oldIdxKey, rowID))
	}
	for _, newIdxKey := range difference(newSecIdxKeys, oldSecIdxKeys) {
		if err := u.addPolicy(store, newIdxKey, rowID); err != nil {
			return err
		}
	}
	return nil
}

// uniqueKeysAddPolicy enforces keys to be unique
func uniqueKeysAddPolicy(store sdk.KVStore, secondaryIndexKey []byte, rowID uint64) error {
	if len(secondaryIndexKey) == 0 {
		return errors.Wrap(ErrArgument, "empty index key")
	}
	it := store.Iterator(makeIndexPrefixScanKey(secondaryIndexKey, 0), makeIndexPrefixScanKey(secondaryIndexKey, math.MaxUint64))
	defer it.Close()
	if it.Valid() {
		return ErrUniqueConstraint
	}

	indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowID)
	store.Set(indexKey, []byte{0})
	return nil
}

// multiKeyAddPolicy allows multiple entries for a key
func multiKeyAddPolicy(store sdk.KVStore, secondaryIndexKey []byte, rowID uint64) error {
	if len(secondaryIndexKey) == 0 {
		return errors.Wrap(ErrArgument, "empty index key")
	}

	indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowID)
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

const encodedSeqLength = 8

// makeIndexPrefixScanKey combines the indexKey with the rowID
func makeIndexPrefixScanKey(indexKey []byte, rowID uint64) []byte {
	n := len(indexKey)
	res := make([]byte, n+encodedSeqLength)
	copy(res, indexKey)
	binary.BigEndian.PutUint64(res[n:], rowID)
	return res
}
