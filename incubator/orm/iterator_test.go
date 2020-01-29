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
			srcIT:     IteratorFunc(nil),
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

func TestLimitedIterator(t *testing.T) {
	sliceIter := func(s ...string) Iterator {
		var pos int
		return IteratorFunc(func(dest interface{}) (key []byte, err error) {
			if pos == len(s) {
				return nil, ErrIteratorDone
			}
			v := s[pos]

			*dest.(*string) = v // dest is a pointer so we set the value here
			pos++
			return []byte(v), nil
		})
	}
	specs := map[string]struct {
		src Iterator
		exp []string
	}{
		"all from range with max > length": {
			src: LimitIterator(sliceIter("a", "b", "c"), 4),
			exp: []string{"a", "b", "c"},
		},
		"up to max": {
			src: LimitIterator(sliceIter("a", "b", "c"), 2),
			exp: []string{"a", "b"},
		},
		"none when max = 0": {
			src: LimitIterator(sliceIter("a", "b", "c"), 0),
			exp: []string{},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			var loaded []string
			_, err := ReadAll(spec.src, &loaded)
			require.NoError(t, err)
			assert.Equal(t, spec.exp, loaded)
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
	return IteratorFunc(func(dest interface{}) (key []byte, err error) {
		return nil, nil
	})
}
