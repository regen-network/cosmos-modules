package orm

import (
	stdErrors "errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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
				return [][]byte{{}}, nil
			},
			expAddPolicyCalled: false,
		},
		"nil key in slice": {
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
				return nil, stdErrors.New("test")
			},
			expErr:             stdErrors.New("test"),
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
		"empty key in slice": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{}}, nil
			},
		},
		"nil key in slice": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{nil}, nil
			},
		},
		"error case": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, stdErrors.New("test")
			},
			expErr: stdErrors.New("test"),
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
				return [][]byte{EncodeSequence(1)}, nil
			},
		},
		"single key - different key, replaced": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				keys := [][]byte{EncodeSequence(1), EncodeSequence(2)}
				return [][]byte{keys[value.(int)]}, nil
			},
			expAddedKeys: [][]byte{
				makeIndexPrefixScanKey(EncodeSequence(2), myRowID),
			},
			expDeletedKeys: [][]byte{
				makeIndexPrefixScanKey(EncodeSequence(1), myRowID),
			},
		},
		"multi key - same key, no update": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{EncodeSequence(1), EncodeSequence(2)}, nil
			},
		},
		"multi key - replaced": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				keys := [][]byte{EncodeSequence(1), EncodeSequence(2), EncodeSequence(3), EncodeSequence(4)}
				return [][]byte{keys[value.(int)], keys[value.(int)+2]}, nil
			},
			expAddedKeys: [][]byte{
				makeIndexPrefixScanKey(EncodeSequence(2), myRowID),
				makeIndexPrefixScanKey(EncodeSequence(4), myRowID),
			},
			expDeletedKeys: [][]byte{
				makeIndexPrefixScanKey(EncodeSequence(1), myRowID),
				makeIndexPrefixScanKey(EncodeSequence(3), myRowID),
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
		"empty key in slice": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{{}}, nil
			},
		},
		"nil key in slice": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return [][]byte{nil}, nil
			},
		},
		"error case": {
			srcFunc: func(value interface{}) ([][]byte, error) {
				return nil, stdErrors.New("test")
			},
			expErr: stdErrors.New("test"),
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

func TestUniqueKeyAddPolicy(t *testing.T) {
	const myRowID = 1
	myPresetKey := makeIndexPrefixScanKey([]byte("my-preset-key"), myRowID)
	specs := map[string]struct {
		srcKey           []byte
		expErr           *errors.Error
		expExistingEntry []byte
	}{

		"create when not exists": {
			srcKey:           []byte("my-index-key"),
			expExistingEntry: append([]byte("my-index-key"), EncodeSequence(myRowID)...),
		},
		"error when exists already": {
			srcKey: []byte("my-preset-key"),
			expErr: ErrUniqueConstraint,
		},
		"nil key not allowed": {
			srcKey: nil,
			expErr: ErrArgument,
		},
		"empty key not allowed": {
			srcKey: []byte{},
			expErr: ErrArgument,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			storeKey := sdk.NewKVStoreKey("test")
			store := NewMockContext().KVStore(storeKey)
			store.Set(myPresetKey, []byte{})
			err := uniqueKeysAddPolicy(store, spec.srcKey, myRowID)
			require.True(t, spec.expErr.Is(err))
			if spec.expErr != nil {
				return
			}
			assert.True(t, store.Has(spec.expExistingEntry), "not found")
		})
	}
}

func TestMultiKeyAddPolicy(t *testing.T) {
	const myRowID = 1
	myPresetKey := makeIndexPrefixScanKey([]byte("my-preset-key"), myRowID)
	specs := map[string]struct {
		srcKey           []byte
		expErr           *errors.Error
		expExistingEntry []byte
	}{

		"create when not exists": {
			srcKey:           []byte("my-index-key"),
			expExistingEntry: append([]byte("my-index-key"), EncodeSequence(myRowID)...),
		},
		"noop when exists already": {
			srcKey:           []byte("my-preset-key"),
			expExistingEntry: myPresetKey,
		},
		"nil key not allowed": {
			srcKey: nil,
			expErr: ErrArgument,
		},
		"empty key not allowed": {
			srcKey: []byte{},
			expErr: ErrArgument,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			storeKey := sdk.NewKVStoreKey("test")
			store := NewMockContext().KVStore(storeKey)
			store.Set(myPresetKey, []byte{})
			err := multiKeyAddPolicy(store, spec.srcKey, myRowID)
			require.True(t, spec.expErr.Is(err))
			if spec.expErr != nil {
				return
			}
			assert.True(t, store.Has(spec.expExistingEntry))
		})
	}
}

func TestDifference(t *testing.T) {
	asByte := func(s []string) [][]byte {
		r := make([][]byte, len(s))
		for i := 0; i < len(s); i++ {
			r[i] = []byte(s[i])
		}
		return r
	}

	specs := map[string]struct {
		srcA      []string
		srcB      []string
		expResult [][]byte
	}{
		"all of A": {
			srcA:      []string{"a", "b"},
			srcB:      []string{"c"},
			expResult: [][]byte{[]byte("a"), []byte("b")},
		},
		"A - B": {
			srcA:      []string{"a", "b"},
			srcB:      []string{"b", "c", "d"},
			expResult: [][]byte{[]byte("a")},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got := difference(asByte(spec.srcA), asByte(spec.srcB))
			assert.Equal(t, spec.expResult, got)
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

type deleteKVStoreRecorder struct {
	AlwaysPanicKVStore
	deletes [][]byte
}

func (m *deleteKVStoreRecorder) Delete(key []byte) {
	m.deletes = append(m.deletes, key)
}

type updateKVStoreRecorder struct {
	deleteKVStoreRecorder
	stored    tuples
	hasResult bool
}

func (u *updateKVStoreRecorder) Set(key, value []byte) {
	u.stored = append(u.stored, tuple{key, value})
}

func (u updateKVStoreRecorder) Has(key []byte) bool {
	return u.hasResult
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
