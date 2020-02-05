package orm

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexerOnCreate(t *testing.T) {
	const myRowID uint64 = 1

	specs := map[string]struct {
		srcFunc            IndexerFunc
		expIndexKeys       [][]byte
		expRowIDs          []uint64
		expAddPolicyCalled bool
		expErr             error
	}{
		"single key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}}, nil
			},
			expAddPolicyCalled: true,
			expIndexKeys:       [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}},
			expRowIDs:          []uint64{myRowID},
		},
		"multi key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}, {1, 0, 0, 0, 0, 0, 0, 0}}, nil
			},
			expAddPolicyCalled: true,
			expIndexKeys:       [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}, {1, 0, 0, 0, 0, 0, 0, 0}},
			expRowIDs:          []uint64{myRowID, myRowID},
		},
		"empty key in slice": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{nil}, nil
			},
			expAddPolicyCalled: false,
		},
		"empty key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{}, nil
			},
			expAddPolicyCalled: false,
		},
		"nil key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, nil
			},
			expAddPolicyCalled: false,
		},
		"error case": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, errors.New("test")
			},
			expErr:             errors.New("test"),
			expAddPolicyCalled: false,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			mockPolicy := &addPolicyRecorder{}
			idx := NewIndexer(spec.srcFunc)
			idx.addPolicy = mockPolicy.add

			err := idx.OnCreate(nil, myRowID, nil)
			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, spec.expIndexKeys, mockPolicy.secondaryIndexKeys)
			assert.Equal(t, spec.expRowIDs, mockPolicy.rowIDs)
			assert.Equal(t, spec.expAddPolicyCalled, mockPolicy.called)
		})
	}
}

func TestIndexerOnDelete(t *testing.T) {
	const myRowID uint64 = 1

	specs := map[string]struct {
		srcFunc      IndexerFunc
		expIndexKeys [][]byte
		expErr       error
	}{
		"single key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}}, nil
			},
			expIndexKeys: [][]byte{makeIndexPrefixScanKey([]byte{0, 0, 0, 0, 0, 0, 0, 1}, myRowID)},
		},
		"multi key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}, {1, 0, 0, 0, 0, 0, 0, 0}}, nil
			},
			expIndexKeys: [][]byte{
				makeIndexPrefixScanKey([]byte{0, 0, 0, 0, 0, 0, 0, 1}, myRowID),
				makeIndexPrefixScanKey([]byte{1, 0, 0, 0, 0, 0, 0, 0}, myRowID),
			},
		},
		"empty key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{}, nil
			},
		},
		"nil key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, nil
			},
		},
		"error case": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, errors.New("test")
			},
			expErr: errors.New("test"),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			mockStore := &deleteKVStoreRecorder{}
			idx := NewIndexer(spec.srcFunc)
			err := idx.OnDelete(mockStore, myRowID, nil)
			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, spec.expIndexKeys, mockStore.deletes)
		})
	}
}

func TestIndexerOnUpdate(t *testing.T) {
	const myRowID uint64 = 1

	specs := map[string]struct {
		srcFunc        IndexerFunc
		expAddedKeys   [][]byte
		expDeletedKeys [][]byte
		expErr         error
	}{
		"single key - same key, no update": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}}, nil
			},
		},
		"single key - different key, replaced": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				keys := [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 2}}
				i := value.(int)
				return [][]byte{keys[i]}, nil
			},
			expAddedKeys: [][]byte{
				makeIndexPrefixScanKey([]byte{0, 0, 0, 0, 0, 0, 0, 2}, myRowID),
			},
			expDeletedKeys: [][]byte{
				makeIndexPrefixScanKey([]byte{0, 0, 0, 0, 0, 0, 0, 1}, myRowID),
			},
		},
		//"multi key": {
		//	srcFunc: func(value interface{}) ([][]byte, error) {
		//		return [][]byte{{0, 0, 0, 0, 0, 0, 0, 1}, {1, 0, 0, 0, 0, 0, 0, 0}}, nil
		//	},
		//	expIndexKeys: [][]byte{
		//		makeIndexPrefixScanKey([]byte{0, 0, 0, 0, 0, 0, 0, 1}, myRowID),
		//		makeIndexPrefixScanKey([]byte{1, 0, 0, 0, 0, 0, 0, 0}, myRowID),
		//	},
		//},
		"empty key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{}, nil
			},
		},
		"nil key": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, nil
			},
		},
		"error case": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, errors.New("test")
			},
			expErr: errors.New("test"),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			mockStore := &updateKVStoreRecorder{}
			idx := NewIndexer(spec.srcFunc)
			err := idx.OnUpdate(mockStore, myRowID, 1, 0)
			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, spec.expDeletedKeys, mockStore.deletes)
			assert.Equal(t, spec.expAddedKeys, mockStore.stored.Keys())
		})
	}
}

type addPolicyRecorder struct {
	secondaryIndexKeys [][]byte
	rowIDs             []uint64
	called             bool
}

func (c *addPolicyRecorder) add(store sdk.KVStore, key []byte, rowID uint64) error {
	c.secondaryIndexKeys = append(c.secondaryIndexKeys, key)
	c.rowIDs = append(c.rowIDs, rowID)
	c.called = true
	return nil
}

type updateKVStoreRecorder struct {
	deleteKVStoreRecorder
	stored tuples
}

func (u *updateKVStoreRecorder) Set(key, value []byte) {
	u.stored = append(u.stored, tuple{key, value})
}

func (u updateKVStoreRecorder) Has(key []byte) bool {
	for _, v := range u.stored {
		if bytes.Equal(key, v.key) {
			return true
		}
	}
	return false
}

type deleteKVStoreRecorder struct {
	alwaysPanicKVStore
	deletes [][]byte
}

func (m *deleteKVStoreRecorder) Delete(key []byte) {
	m.deletes = append(m.deletes, key)
}

type alwaysPanicKVStore struct{}

func (a alwaysPanicKVStore) GetStoreType() types.StoreType {
	panic("implement me")
}

func (a alwaysPanicKVStore) CacheWrap() types.CacheWrap {
	panic("implement me")
}

func (a alwaysPanicKVStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	panic("implement me")
}

func (a alwaysPanicKVStore) Get(key []byte) []byte {
	panic("implement me")
}

func (a alwaysPanicKVStore) Has(key []byte) bool {
	panic("implement me")
}

func (a alwaysPanicKVStore) Set(key, value []byte) {
	panic("implement me")
}

func (a alwaysPanicKVStore) Delete(key []byte) {
	panic("implement me")
}

func (a alwaysPanicKVStore) Iterator(start, end []byte) types.Iterator {
	panic("implement me")
}

func (a alwaysPanicKVStore) ReverseIterator(start, end []byte) types.Iterator {
	panic("implement me")
}

type tuple struct {
	key, val []byte
}

type tuples []tuple

func (t tuples) Keys() [][]byte {
	if t == nil {
		return nil
	}
	r := make([][]byte, len(t))
	for i, v := range t {
		r[i] = v.key
	}
	return r
}
