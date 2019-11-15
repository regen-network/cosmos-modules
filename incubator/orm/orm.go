/* Package orm (object-relational mapping) provides a set of tools on top of the KV store interface to handle
things like secondary indexes and auto-generated ID's that would otherwise need to be hand-generated on a case by
case basis.
*/
package orm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
)

type HasKVStore interface {
	KVStore(key sdk.StoreKey) sdk.KVStore
}

type Indexer interface {
	DoIndex(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error
	BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index
}

func makeIndexPrefixScanKey(indexKey []byte, rowId uint64) []byte {
	n := len(indexKey)
	res := make([]byte, n+8)
	copy(res, indexKey)
	binary.LittleEndian.PutUint64(res[n:], rowId)
	return res
}

type uniqueIndexer struct {
	indexFn func(key []byte, value interface{}) ([]byte, error)
}

type uniqueIndex struct {
	uniqueIndexer
	storeKey    sdk.StoreKey
	prefix      []byte
	modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
}

func (u uniqueIndexer) DoIndex(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error {
	indexKey, err := u.indexFn(key, value)
	if err != nil {
		return err
	}
	bz := store.Get(indexKey)
	if bz != nil {
		existingRowId := binary.LittleEndian.Uint64(bz)
		if existingRowId != rowId {
			return fmt.Errorf("unique index constraint violated")
		}
		return nil
	}
	bz = make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, rowId)
	store.Set(indexKey, bz)
	return nil
}

func (u uniqueIndexer) BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index {
	return uniqueIndex{u, storeKey, prefix, modelGetter}
}

func (u uniqueIndex) Has(ctx HasKVStore, key []byte) (bool, error) {
	panic("implement me")
}

func (u uniqueIndex) getRowId(ctx HasKVStore, key []byte) uint64 {
	store := prefix.NewStore(ctx.KVStore(u.storeKey), u.prefix)
	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(bz)
}

func (u uniqueIndex) GetOne(ctx HasKVStore, indexKey []byte, dest interface{}) ([]byte, error) {
	rowId := u.getRowId(ctx, indexKey)
	return u.modelGetter(ctx, rowId, dest)
}

func (u uniqueIndex) Get(ctx HasKVStore, key []byte) (Iterator, error) {
	panic("implement me")
}

func (u uniqueIndex) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func (u uniqueIndex) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func newPrimaryKeyIndexer() Indexer {
	return uniqueIndexer{indexFn: func(key []byte, value interface{}) (i []byte, e error) {
		return key, nil
	}}
}

type IndexFunc = func(value interface{}) ([]byte, error)

func NewUniqueIndexer(fn IndexFunc) Indexer {
	return uniqueIndexer{indexFn: func(key []byte, value interface{}) (i []byte, e error) {
		return fn(value)
	}}
}

var _ Indexer = uniqueIndexer{}

func NewIndexer(fn IndexFunc) Indexer {
	return indexer{fn}
}

type indexer struct {
	indexFn IndexFunc
}

func (i indexer) DoIndex(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error {
	key, err := i.indexFn(value)
	if err != nil {
		return err
	}
	indexKey := makeIndexPrefixScanKey(key, rowId)
	if !store.Has(indexKey) {
		store.Set(indexKey, []byte{0})
	}
	return nil
}

func (i indexer) BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index {
	return index{storeKey: storeKey, prefix: prefix, modelGetter: modelGetter}
}

var _ Indexer = indexer{}

type index struct {
	storeKey    sdk.StoreKey
	prefix      []byte
	modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
}

func (i index) Has(ctx HasKVStore, key []byte) (bool, error) {
	panic("implement me")
}

func (i index) Get(ctx HasKVStore, key []byte) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(key, nil)
	return indexIterator{ctx: ctx, it: it, end: key, modelGetter: i.modelGetter}, nil
}

func (i index) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func (i index) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

type Index interface {
	Has(ctx HasKVStore, key []byte) (bool, error)
	Get(ctx HasKVStore, key []byte) (Iterator, error)
	PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error)
}

type UniqueIndex interface {
	Index
	GetOne(ctx HasKVStore, indexKey []byte, dest interface{}) (primaryKey []byte, error error)
}

type UInt64Index interface {
	Has(ctx HasKVStore, key uint64) (bool, error)
	Get(ctx HasKVStore, key uint64) (Iterator, error)
	PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
}

type indexIterator struct {
	ctx         HasKVStore
	modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
	it          types.Iterator
	end         []byte
	reverse     bool
}

func (i indexIterator) LoadNext(dest interface{}) (key []byte, err error) {
	if !i.it.Valid() {
		return nil, fmt.Errorf("not found")
	}
	indexPrefixKey := i.it.Key()
	n := len(indexPrefixKey)
	indexKey := indexPrefixKey[:n-8]
	cmp := bytes.Compare(indexKey, i.end)
	if i.end != nil {
		if !i.reverse && cmp > 0 {
			return nil, fmt.Errorf("not found")
		} else if i.reverse && cmp < 0 {
			return nil, fmt.Errorf("not found")
		}
	}
	rowId := binary.LittleEndian.Uint64(indexPrefixKey[n-8:])
	i.it.Next()
	return i.modelGetter(i.ctx, rowId, dest)
}

func (i indexIterator) Close() error {
	i.it.Close()
	return nil
}

//type IntIndex interface {
//	Has(ctx HasKVStore, key sdk.Int) (bool, error)
//	Get(ctx HasKVStore, key sdk.Int) (Iterator, error)
//	PrefixScan(ctx HasKVStore, start sdk.Int, end sdk.Int) (Iterator, error)
//	ReversePrefixScan(ctx HasKVStore, start sdk.Int, end sdk.Int) (Iterator, error)
//}
//
type BucketBuilder interface {
	CreateIndex(prefix []byte, indexer Indexer) Index
	Build() BucketBase
}

// BucketBase provides methods shared by all buckets
type BucketBase interface {
	UniqueIndex
	// Delete deletes the value at the given key
	Delete(ctx HasKVStore, key []byte) error
}

// ExternalKeyBucket defines a bucket where the key is stored externally to the value object
type ExternalKeyBucket interface {
	BucketBase
	// Save saves the given key value pair
	Save(ctx HasKVStore, key []byte, value interface{}) error
}

type HasID interface {
	ID() []byte
}

// NaturalKeyTable defines a bucket where all values implement HasID and the key is stored it the value and
// returned by the HasID method
type NaturalKeyTable interface {
	BucketBase
	// Save saves the value passed in
	Save(ctx HasKVStore, value HasID) error
}

type AutoUInt64Table interface {
	Has(ctx HasKVStore, key uint64) (bool, error)
	Get(ctx HasKVStore, key uint64) (Iterator, error)
	PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error)
	Save(ctx HasKVStore, key []byte, value interface{}) error
}

// AutoKeyTable specifies a bucket where keys are generated via an auto-incremented interger
type AutoKeyTable interface {
	ExternalKeyBucket

	// Create auto-generates key
	Create(ctx HasKVStore, value interface{}) ([]byte, error)
}

//Iterator allows iteration through a sequence of key value pairs
type Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the key. If there
	// are no more items an error is returned
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

type bucketBase struct {
	uniqueIndex
	key          sdk.StoreKey
	bucketPrefix []byte
	cdc          *codec.Codec
	indexers     []Indexer
}

var _ BucketBase = bucketBase{}

func (b bucketBase) rootStore(ctx HasKVStore) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
}

func (b bucketBase) indexStore(ctx HasKVStore, indexName string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(fmt.Sprintf("%s/%s", b.bucketPrefix, indexName)))
}

type externalKeyBucket struct {
	bucketBase
}

//func NewExternalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) ExternalKeyBucket {
//	return &externalKeyBucket{bucketBase{
//		key,
//		bucketPrefix,
//		cdc,
//		indexes,
//	}}
//}

func (b bucketBase) save(ctx HasKVStore, key []byte, value interface{}) error {
	rootStore := b.rootStore(ctx)
	bz, err := b.cdc.MarshalBinaryBare(value)
	if err != nil {
		return err
	}
	rootStore.Set(key, bz)
	for _, idx := range b.indexers {
		indexStore := b.indexStore(ctx, panic("TODO: index prefix"))
		i, err := idx.DoIndex(store, key, value)
		if err != nil {
			return err
		}
		indexStore.Set([]byte(fmt.Sprintf("%x/%x", i, key)), []byte{0})
	}
	return nil
}
func (b externalKeyBucket) Save(ctx HasKVStore, key []byte, value interface{}) error {
	return b.save(ctx, key, value)
}

func (b bucketBase) Delete(ctx HasKVStore, key []byte) error {
	rootStore := b.rootStore(ctx)
	rootStore.Delete(key)
	// TODO: delete indexes
	return nil
}

func NewIndex(builder TableBuilder, prefix byte, indexer func(val interface{}) []byte) Index {
	return index{}
}

type TableBuilder interface {
	RegisterIndexer(prefix byte, indexer Indexer)
}

func NewAutoUInt64TableBuilder(prefix byte, key sdk.StoreKey, cdc *codec.Codec) AutoUInt64TableBuilder {
	return AutoUInt64TableBuilder{prefix: prefix, key: key, cdc: cdc}
}

type indexRef struct {
	prefix  byte
	indexer Indexer
}

type AutoUInt64TableBuilder struct {
	prefix      byte
	key         sdk.StoreKey
	cdc         *codec.Codec
	indexerRefs []indexRef
}

func (a AutoUInt64TableBuilder) RegisterIndexer(prefix byte, indexer Indexer) {
	a.indexerRefs = append(a.indexerRefs, indexRef{prefix: prefix, indexer: indexer})
}

func (a AutoUInt64TableBuilder) Build() AutoUInt64Table {
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

//type naturalKeyBucket struct {
//	bucketBase
//}
//
//func NewNaturalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) NaturalKeyTable {
//	return &naturalKeyBucket{bucketBase{key, bucketPrefix, cdc, indexes}}
//}
//
//func (n naturalKeyBucket) Save(ctx HasKVStore, value HasID) error {
//	return n.save(ctx, value.ID(), value)
//}
//
//func NewAutoIDBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index, idGenerator func(x uint64) []byte) AutoKeyTable {
//	return &autoIDBucket{externalKeyBucket{bucketBase{key, bucketPrefix, cdc, indexes}}, idGenerator}
//}
//
//type autoIDBucket struct {
//	externalKeyBucket
//	idGenerator func(x uint64) []byte
//}
//
//func writeUInt64(x uint64) []byte {
//	buf := make([]byte, binary.MaxVarintLen64)
//	n := binary.PutUvarint(buf, x)
//	return buf[:n]
//}
//
//func readUInt64(bz []byte) (uint64, error) {
//	x, n := binary.Uvarint(bz)
//	if n <= 0 {
//		return 0, fmt.Errorf("can't read var uint64")
//	}
//	return x, nil
//}
//
//func (a autoIDBucket) Create(ctx HasKVStore, value interface{}) ([]byte, error) {
//	st := a.indexStore(ctx, "$")
//	bz := st.Get([]byte("$"))
//	var nextID uint64 = 0
//	var err error
//	if bz != nil {
//		nextID, err = readUInt64(bz)
//		if err != nil {
//			return nil, err
//		}
//	}
//	st.Set([]byte("$"), writeUInt64(nextID))
//	return a.idGenerator(nextID), nil
//}
//
//type iterator struct {
//	cdc *codec.Codec
//	it  sdk.Iterator
//}
//
//func (i *iterator) LoadNext(dest interface{}) (key []byte, err error) {
//	if !i.it.Valid() {
//		return nil, fmt.Errorf("invalid")
//	}
//	key = i.it.Key()
//	err = i.cdc.UnmarshalBinaryBare(i.it.Value(), dest)
//	if err != nil {
//		return nil, err
//	}
//	i.it.Next()
//	return key, nil
//}
//
//func (i *iterator) Close() {
//	i.it.Close()
//}
//
//type indexIterator struct {
//	bucketBase
//	ctx HasKVStore
//	it    sdk.Iterator
//	start []byte
//	end   []byte
//}
//
//func (i indexIterator) LoadNext(dest interface{}) (key []byte, err error) {
//	if !i.it.Valid() {
//		return nil, fmt.Errorf("invalid")
//	}
//	pieces := strings.Split(string(i.it.Key()), "/")
//	if len(pieces) != 2 {
//		return nil, fmt.Errorf("unexpected index key")
//	}
//	indexPrefix, err := hex.DecodeString(pieces[0])
//	if err != nil {
//		return nil, err
//	}
//	// check out of range
//	if !((i.start == nil || bytes.Compare(i.start, indexPrefix) >= 0) && (i.end == nil || bytes.Compare(indexPrefix, i.end) <= 0)) {
//		return nil, fmt.Errorf("done")
//	}
//	key, err = hex.DecodeString(pieces[1])
//	if err != nil {
//		return nil, err
//	}
//	err = i.bucketBase.GetOne(i.ctx, key, dest)
//	if err != nil {
//		return nil, err
//	}
//	return key, nil
//}
//
//func (i indexIterator) Close() {
//	i.it.Close()
//}
//
//
