package orm

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNaturalKeyTablePrefixScan(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()
	const (
		testTablePrefix = iota
		testTableSeqPrefix
		testTableIndexPrefix
	)

	tb := NewNaturalKeyTableBuilder(testTablePrefix, testTableSeqPrefix, testTableIndexPrefix, storeKey, cdc, &GroupMember{}).
		Build()

	ctx := NewMockContext()

	anyWeight := sdk.NewInt(1)
	m1 := GroupMember{
		Group:  []byte("group-a"),
		Member: []byte("member-one"),
		Weight: anyWeight,
	}
	m2 := GroupMember{
		Group:  []byte("group-a"),
		Member: []byte("member-two"),
		Weight: anyWeight,
	}
	m3 := GroupMember{
		Group:  []byte("group-b"),
		Member: []byte("member-two"),
		Weight: anyWeight,
	}
	for _, g := range []GroupMember{m1, m2, m3} {
		require.NoError(t, tb.Create(ctx, &g))
	}

	specs := map[string]struct {
		start, end []byte
		expResult  []GroupMember
		expRowIDs  [][]byte
		expError   *errors.Error
		method     func(ctx HasKVStore, start, end []byte) (Iterator, error)
	}{
		"exact match with a single result": {
			start:     []byte("group-amember-one"), // == m1.NaturalKey()
			end:       []byte("group-amember-two"), // == m2.NaturalKey()
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1},
			expRowIDs: [][]byte{EncodeSequence(1)},
		},
		"one result by prefix": {
			start:     []byte("group-a"),
			end:       []byte("group-amember-two"), // == m2.NaturalKey()
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1},
			expRowIDs: [][]byte{EncodeSequence(1)},
		},
		"multi key elements by group prefix": {
			start:     []byte("group-a"),
			end:       []byte("group-b"),
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1, m2},
			expRowIDs: [][]byte{EncodeSequence(1), EncodeSequence(2)},
		},
		"open end query with second group": {
			start:     []byte("group-b"),
			end:       nil,
			method:    tb.PrefixScan,
			expResult: []GroupMember{m3},
			expRowIDs: [][]byte{EncodeSequence(3)},
		},
		"open end query with all": {
			start:     []byte("group-a"),
			end:       nil,
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1, m2, m3},
			expRowIDs: [][]byte{EncodeSequence(1), EncodeSequence(2), EncodeSequence(3)},
		},
		"open start query": {
			start:     nil,
			end:       []byte("group-b"),
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1, m2},
			expRowIDs: [][]byte{EncodeSequence(1), EncodeSequence(2)},
		},
		"open start and end query": {
			start:     nil,
			end:       nil,
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1, m2, m3},
			expRowIDs: [][]byte{EncodeSequence(1), EncodeSequence(2), EncodeSequence(3)},
		},
		"all matching prefix": {
			start:     []byte("group"),
			end:       nil,
			method:    tb.PrefixScan,
			expResult: []GroupMember{m1, m2, m3},
			expRowIDs: [][]byte{EncodeSequence(1), EncodeSequence(2), EncodeSequence(3)},
		},
		"non matching prefix": {
			start:     []byte("nobody"),
			end:       nil,
			method:    tb.PrefixScan,
			expResult: []GroupMember{},
		},
		"start equals end": {
			start:    []byte("any"),
			end:      []byte("any"),
			method:   tb.PrefixScan,
			expError: ErrArgument,
		},
		"start after end": {
			start:    []byte("b"),
			end:      []byte("a"),
			method:   tb.PrefixScan,
			expError: ErrArgument,
		},
		"reverse: exact match with a single result": {
			start:     []byte("group-amember-one"), // == m1.NaturalKey()
			end:       []byte("group-amember-two"), // == m2.NaturalKey()
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m1},
			expRowIDs: [][]byte{EncodeSequence(1)},
		},
		"reverse: one result by prefix": {
			start:     []byte("group-a"),
			end:       []byte("group-amember-two"), // == m2.NaturalKey()
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m1},
			expRowIDs: [][]byte{EncodeSequence(1)},
		},
		"reverse: multi key elements by group prefix": {
			start:     []byte("group-a"),
			end:       []byte("group-b"),
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m2, m1},
			expRowIDs: [][]byte{EncodeSequence(2), EncodeSequence(1)},
		},
		"reverse: open end query with second group": {
			start:     []byte("group-b"),
			end:       nil,
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m3},
			expRowIDs: [][]byte{EncodeSequence(3)},
		},
		"reverse: open end query with all": {
			start:     []byte("group-a"),
			end:       nil,
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m3, m2, m1},
			expRowIDs: [][]byte{EncodeSequence(3), EncodeSequence(2), EncodeSequence(1)},
		},
		"reverse: open start query": {
			start:     nil,
			end:       []byte("group-b"),
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m2, m1},
			expRowIDs: [][]byte{EncodeSequence(2), EncodeSequence(1)},
		},
		"reverse: open start and end query": {
			start:     nil,
			end:       nil,
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m3, m2, m1},
			expRowIDs: [][]byte{EncodeSequence(3), EncodeSequence(2), EncodeSequence(1)},
		},
		"reverse: all matching prefix": {
			start:     []byte("group"),
			end:       nil,
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{m3, m2, m1},
			expRowIDs: [][]byte{EncodeSequence(3), EncodeSequence(2), EncodeSequence(1)},
		},
		"reverse: non matching prefix": {
			start:     []byte("nobody"),
			end:       nil,
			method:    tb.ReversePrefixScan,
			expResult: []GroupMember{},
		},
		"reverse: start equals end": {
			start:    []byte("any"),
			end:      []byte("any"),
			method:   tb.ReversePrefixScan,
			expError: ErrArgument,
		},
		"reverse: start after end": {
			start:    []byte("b"),
			end:      []byte("a"),
			method:   tb.ReversePrefixScan,
			expError: ErrArgument,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			it, err := spec.method(ctx, spec.start, spec.end)
			require.True(t, spec.expError.Is(err), "expected #+v but got #+v", spec.expError, err)
			if spec.expError != nil {
				return
			}
			var loaded []GroupMember
			rowIDs, err := ReadAll(it, &loaded)
			require.NoError(t, err)
			assert.Equal(t, spec.expResult, loaded)
			assert.Equal(t, spec.expRowIDs, rowIDs)
		})
	}
}
