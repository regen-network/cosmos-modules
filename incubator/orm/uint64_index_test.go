package orm

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUInt64Index(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()

	groupMemberTableBuilder := NewNaturalKeyTableBuilder(GroupMemberTablePrefix, GroupMemberTableSeqPrefix, GroupMemberTableIndexPrefix, storeKey, cdc, &GroupMember{})
	idx := NewUInt64Index(groupMemberTableBuilder, GroupMemberByMemberIndexPrefix, func(val interface{}) ([]uint64, error) {
		return []uint64{uint64(val.(*GroupMember).Member[0])}, nil
	})
	groupMemberTable := groupMemberTableBuilder.Build()

	ctx := NewMockContext()

	m := GroupMember{
		Group:  sdk.AccAddress(EncodeSequence(1)),
		Member: sdk.AccAddress([]byte("member-address")),
		Weight: sdk.NewInt(10),
	}
	err := groupMemberTable.Create(ctx, &m)
	require.NoError(t, err)

	indexedKey := uint64('m')

	// Has
	assert.True(t, idx.Has(ctx, indexedKey))

	// Get
	it, err := idx.Get(ctx, indexedKey)
	require.NoError(t, err)
	var loaded GroupMember
	rowID, err := it.LoadNext(&loaded)
	require.NoError(t, err)
	require.Equal(t, uint64(1), DecodeSequence(rowID))
	require.Equal(t, m, loaded)

	// PrefixScan match
	it, err = idx.PrefixScan(ctx, 0, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.NoError(t, err)
	require.Equal(t, uint64(1), DecodeSequence(rowID))
	require.Equal(t, m, loaded)

	// PrefixScan no match
	it, err = idx.PrefixScan(ctx, indexedKey+1, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.Error(t, ErrIteratorDone, err)

	// ReversePrefixScan match
	it, err = idx.ReversePrefixScan(ctx, 0, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.NoError(t, err)
	require.Equal(t, uint64(1), DecodeSequence(rowID))
	require.Equal(t, m, loaded)

	// ReversePrefixScan no match
	it, err = idx.ReversePrefixScan(ctx, indexedKey+1, 255)
	require.NoError(t, err)
	rowID, err = it.LoadNext(&loaded)
	require.Error(t, ErrIteratorDone, err)
}