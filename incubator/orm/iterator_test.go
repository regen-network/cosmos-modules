package orm

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAll(t *testing.T) {
	specs := map[string]struct {
		srcIT     Iterator
		destSlice func() ModelSlicePtr
		expErr    *errors.Error
		expIDs    [][]byte
		expResult ModelSlicePtr
	}{
		"all good with object slice": {
			srcIT: mockIter(EncodeSequence(1), &GroupMetadata{Description: "test"}),
			destSlice: func() ModelSlicePtr {
				x := make([]GroupMetadata, 1)
				return &x
			},
			expIDs:    [][]byte{EncodeSequence(1)},
			expResult: &[]GroupMetadata{{Description: "test"}},
		},
		"all good with pointer slice": {
			srcIT: mockIter(EncodeSequence(1), &GroupMetadata{Description: "test"}),
			destSlice: func() ModelSlicePtr {
				x := make([]*GroupMetadata, 1)
				return &x
			},
			expIDs:    [][]byte{EncodeSequence(1)},
			expResult: &[]*GroupMetadata{{Description: "test"}},
		},
		"dest slice empty": {
			srcIT: mockIter(EncodeSequence(1), &GroupMetadata{}),
			destSlice: func() ModelSlicePtr {
				x := make([]GroupMetadata, 0)
				return &x
			},
			expResult: &[]GroupMetadata{},
		},
		"dest pointer with nil value": {
			srcIT: mockIter(EncodeSequence(1), &GroupMetadata{}),
			destSlice: func() ModelSlicePtr {
				return (*[]GroupMetadata)(nil)
			},
			expErr: ErrArgument,
		},
		"iterator is nil": {
			srcIT:     nil,
			destSlice: func() ModelSlicePtr { return new([]GroupMetadata) },
			expErr:    ErrArgument,
		},
		"dest slice is nil": {
			srcIT:     noopIter(),
			destSlice: func() ModelSlicePtr { return nil },
			expErr:    ErrArgument,
		},
		"dest slice is not a pointer": {
			srcIT:     iteratorFunc(nil),
			destSlice: func() ModelSlicePtr { return make([]GroupMetadata, 1) },
			expErr:    ErrArgument,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			loaded := spec.destSlice()
			ids, err := ReadAll(spec.srcIT, loaded)
			require.True(t, spec.expErr.Is(err), "expected %s but got %s", spec.expErr, err)
			assert.Equal(t, spec.expIDs, ids)
			if err == nil {
				assert.Equal(t, spec.expResult, loaded)
			}
		})
	}
}

// mockIter amino encodes + decodes value object.
func mockIter(rowID []byte, val interface{}) Iterator {
	cdc := codec.New()
	b, err := cdc.MarshalBinaryBare(val)
	if err != nil {
		panic(err)
	}
	return NewSingleValueIterator(cdc, rowID, b)
}

func noopIter() Iterator {
	return iteratorFunc(func(dest interface{}) (key []byte, err error) {
		return nil, nil
	})
}
