package orm

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/orm/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUInt64Index(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")

	const anyPrefix = 0x10
	tableBuilder := NewNaturalKeyTableBuilder(anyPrefix, storeKey, &testdata.GroupMember{}, Max255DynamicLengthIndexKeyCodec{})
	myIndex := NewUInt64Index(tableBuilder, GroupMemberByMemberIndexPrefix, func(val interface{}) ([]uint64, error) {
		return []uint64{uint64(val.(*testdata.GroupMember).Member[0])}, nil
	})
	myTable := tableBuilder.Build()

	ctx := NewMockContext()

	m := testdata.GroupMember{
		Group:  sdk.AccAddress(EncodeSequence(1)),
		Member: sdk.AccAddress([]byte("member-address")),
		Weight: 10,
	}
	err := myTable.Create(ctx, &m)
	require.NoError(t, err)

	indexedKey := uint64('m')

	// Has
	assert.True(t, myIndex.Has(ctx, indexedKey))

	// Get
	it, err := myIndex.Get(ctx, indexedKey)
	require.NoError(t, err)
	var loaded testdata.GroupMember
	rowID, err := it.LoadNext(&loaded)
	require.NoError(t, err)
	require.Equal(t, uint64(1), DecodeSequence(rowID))
	require.Equal(t, m, loaded)

	// PrefixScan match
	it, err = myIndex.PrefixScan(ctx, 0, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.NoError(t, err)
	require.Equal(t, uint64(1), DecodeSequence(rowID))
	require.Equal(t, m, loaded)

	// PrefixScan no match
	it, err = myIndex.PrefixScan(ctx, indexedKey+1, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.Error(t, ErrIteratorDone, err)

	// ReversePrefixScan match
	it, err = myIndex.ReversePrefixScan(ctx, 0, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.NoError(t, err)
	require.Equal(t, uint64(1), DecodeSequence(rowID))
	require.Equal(t, m, loaded)

	// ReversePrefixScan no match
	it, err = myIndex.ReversePrefixScan(ctx, indexedKey+1, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.Error(t, ErrIteratorDone, err)
}

func TestUInt64MultiKeyAdapter(t *testing.T) {
	specs := map[string]struct {
		srcFunc UInt64IndexerFunc
		exp     []RowID
		expErr  error
	}{
		"single key": {
			srcFunc: func(value interface{}) ([]uint64, error) {
				return []uint64{1}, nil
			},
			exp: []RowID{{0, 0, 0, 0, 0, 0, 0, 1}},
		},
		"multi key": {
			srcFunc: func(value interface{}) ([]uint64, error) {
				return []uint64{1, 1 << 56}, nil
			},
			exp: []RowID{{0, 0, 0, 0, 0, 0, 0, 1}, {1, 0, 0, 0, 0, 0, 0, 0}},
		},
		"empty key": {
			srcFunc: func(value interface{}) ([]uint64, error) {
				return []uint64{}, nil
			},
			exp: []RowID{},
		},
		"nil key": {
			srcFunc: func(value interface{}) ([]uint64, error) {
				return nil, nil
			},
			exp: []RowID{},
		},
		"error case": {
			srcFunc: func(value interface{}) ([]uint64, error) {
				return nil, errors.New("test")
			},
			expErr: errors.New("test"),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			fn := UInt64MultiKeyAdapter(spec.srcFunc)
			r, err := fn(nil)
			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, spec.exp, r)
		})
	}
}

func TestVirtualUInt64Index(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	const testTablePrefix = iota
	builder := NewNaturalKeyTableBuilder(testTablePrefix, storeKey, &testdata.GroupMember{}, Max255DynamicLengthIndexKeyCodec{})
	table := builder.Build()

	ctx := NewMockContext()
	const anyWeight = 1
	m1 := testdata.GroupMember{
		Group:  EncodeSequence(1),
		Member: []byte("member-one"),
		Weight: anyWeight,
	}
	m2 := testdata.GroupMember{
		Group:  EncodeSequence(1),
		Member: []byte("member-two"),
		Weight: anyWeight,
	}
	m3 := testdata.GroupMember{
		Group:  EncodeSequence(2),
		Member: []byte("member-two"),
		Weight: anyWeight,
	}
	for _, g := range []testdata.GroupMember{m1, m2, m3} {
		require.NoError(t, table.Create(ctx, &g))
	}
	idx := AsUInt64Index(NewVirtualIndex(builder))

	it, err := idx.Get(ctx, 1)
	require.NoError(t, err)
	assert.True(t, idx.Has(ctx, 1))

	var loaded []testdata.GroupMember
	rowIDs, err := ReadAll(it, &loaded)
	require.NoError(t, err)
	assert.Equal(t, []testdata.GroupMember{m1, m2}, loaded)
	assert.Equal(t, []RowID{m1.NaturalKey(), m2.NaturalKey()}, rowIDs)

	// and with prefix scan
	it, err = idx.PrefixScan(ctx, 1, 9999)
	require.NoError(t, err)
	rowIDs, err = ReadAll(it, &loaded)
	require.NoError(t, err)
	assert.Equal(t, []testdata.GroupMember{m1, m2, m3}, loaded)
	assert.Equal(t, []RowID{m1.NaturalKey(), m2.NaturalKey(), m3.NaturalKey()}, rowIDs)
	// and reverse
	it, err = idx.ReversePrefixScan(ctx, 1, 9999)
	require.NoError(t, err)
	rowIDs, err = ReadAll(it, &loaded)
	require.NoError(t, err)
	assert.Equal(t, []testdata.GroupMember{m3, m2, m1}, loaded)
	assert.Equal(t, []RowID{m3.NaturalKey(), m2.NaturalKey(), m1.NaturalKey()}, rowIDs)

	// and when one entry removed
	require.NoError(t, table.Delete(ctx, &m2))
	it, err = idx.Get(ctx, 1)
	require.NoError(t, err)

	assert.True(t, idx.Has(ctx, 1))

	rowIDs, err = ReadAll(it, &loaded)
	require.NoError(t, err)
	assert.Equal(t, []testdata.GroupMember{m1}, loaded)
	assert.Equal(t, []RowID{m1.NaturalKey()}, rowIDs)

	// and when other entry removed
	require.NoError(t, table.Delete(ctx, &m1))
	it, err = idx.Get(ctx, 1)
	require.NoError(t, err)

	assert.False(t, idx.Has(ctx, 1))

	rowIDs, err = ReadAll(it, &loaded)
	require.NoError(t, err)
	assert.Equal(t, []testdata.GroupMember{}, loaded)
	assert.Nil(t, rowIDs)
}
