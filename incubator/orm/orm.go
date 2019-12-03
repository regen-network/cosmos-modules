/* Package orm (object-relational mapping) provides a set of tools on top of the KV store interface to handle
things like secondary indexes and auto-generated ID's that would otherwise need to be hand-generated on a case by
case basis.
*/
package orm

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//type IntIndex interface {
//	Has(ctx HasKVStore, storeKey sdk.Int) (bool, error)
//	Get(ctx HasKVStore, storeKey sdk.Int) (Iterator, error)
//	PrefixScan(ctx HasKVStore, start sdk.Int, end sdk.Int) (Iterator, error)
//	ReversePrefixScan(ctx HasKVStore, start sdk.Int, end sdk.Int) (Iterator, error)
//}
//
type bucketBase struct {
	uniqueIndex
	key          sdk.StoreKey
	bucketPrefix []byte
	cdc          *codec.Codec
	indexers     []Indexer
}

var _ TableBase = bucketBase{}

func (b bucketBase) rootStore(ctx HasKVStore) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
}

func (b bucketBase) indexStore(ctx HasKVStore, indexName string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(fmt.Sprintf("%s/%s", b.bucketPrefix, indexName)))
}

type externalKeyBucket struct {
	bucketBase
}

//func NewExternalKeyBucket(storeKey sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) ExternalKeyTable {
//	return &externalKeyBucket{bucketBase{
//		storeKey,
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

//type naturalKeyBucket struct {
//	bucketBase
//}
//
//func NewNaturalKeyBucket(storeKey sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) NaturalKeyTable {
//	return &naturalKeyBucket{bucketBase{storeKey, bucketPrefix, cdc, indexes}}
//}
//
//func (n naturalKeyBucket) Save(ctx HasKVStore, value HasID) error {
//	return n.save(ctx, value.ID(), value)
//}
//
//func NewAutoIDBucket(storeKey sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index, idGenerator func(x uint64) []byte) AutoKeyTable {
//	return &autoIDBucket{externalKeyBucket{bucketBase{storeKey, bucketPrefix, cdc, indexes}}, idGenerator}
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
//func (i *iterator) LoadNext(dest interface{}) (storeKey []byte, err error) {
//	if !i.it.Valid() {
//		return nil, fmt.Errorf("invalid")
//	}
//	storeKey = i.it.Key()
//	err = i.cdc.UnmarshalBinaryBare(i.it.Value(), dest)
//	if err != nil {
//		return nil, err
//	}
//	i.it.Next()
//	return storeKey, nil
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
//func (i indexIterator) LoadNext(dest interface{}) (storeKey []byte, err error) {
//	if !i.it.Valid() {
//		return nil, fmt.Errorf("invalid")
//	}
//	pieces := strings.Split(string(i.it.Key()), "/")
//	if len(pieces) != 2 {
//		return nil, fmt.Errorf("unexpected index storeKey")
//	}
//	indexPrefix, err := hex.DecodeString(pieces[0])
//	if err != nil {
//		return nil, err
//	}
//	// check out of range
//	if !((i.start == nil || bytes.Compare(i.start, indexPrefix) >= 0) && (i.end == nil || bytes.Compare(indexPrefix, i.end) <= 0)) {
//		return nil, fmt.Errorf("done")
//	}
//	storeKey, err = hex.DecodeString(pieces[1])
//	if err != nil {
//		return nil, err
//	}
//	err = i.bucketBase.GetOne(i.ctx, storeKey, dest)
//	if err != nil {
//		return nil, err
//	}
//	return storeKey, nil
//}
//
//func (i indexIterator) Close() {
//	i.it.Close()
//}
//
//

