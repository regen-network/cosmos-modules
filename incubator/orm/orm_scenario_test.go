package orm

import (
	"encoding/binary"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeperEndToEndWithAutoUInt64Table(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	ctx := NewMockContext()

	k := NewGroupKeeper(storeKey)

	g := GroupMetadata{
		Description: "my test",
		Admin:       sdk.AccAddress([]byte("admin-address")),
	}
	// when stored
	rowID, err := k.groupTable.Create(ctx, &g)
	require.NoError(t, err)
	// then we should find it
	exists := k.groupTable.Has(ctx, rowID)
	require.True(t, exists)

	// and load it
	var loaded GroupMetadata

	binKey, err := k.groupTable.GetOne(ctx, rowID, &loaded)
	require.NoError(t, err)

	assert.Equal(t, rowID, binary.BigEndian.Uint64(binKey))
	assert.Equal(t, "my test", loaded.Description)
	assert.Equal(t, sdk.AccAddress([]byte("admin-address")), loaded.Admin)

	// and exists in MultiKeyIndex
	exists = k.groupByAdminIndex.Has(ctx, []byte("admin-address"))
	require.True(t, exists)

	// and when loaded
	it, err := k.groupByAdminIndex.Get(ctx, []byte("admin-address"))
	require.NoError(t, err)

	// then
	binKey, loaded = first(t, it)
	assert.Equal(t, rowID, binary.BigEndian.Uint64(binKey))
	assert.Equal(t, g, loaded)

	// when updated
	g.Admin = []byte("new-admin-address")
	err = k.groupTable.Save(ctx, rowID, &g)
	require.NoError(t, err)

	// then indexes are updated, too
	exists = k.groupByAdminIndex.Has(ctx, []byte("new-admin-address"))
	require.True(t, exists)

	exists = k.groupByAdminIndex.Has(ctx, []byte("admin-address"))
	require.False(t, exists)

	// when deleted
	err = k.groupTable.Delete(ctx, rowID)
	require.NoError(t, err)

	// then removed from primary MultiKeyIndex
	exists = k.groupTable.Has(ctx, rowID)
	require.False(t, exists)

	// and also removed from secondary MultiKeyIndex
	exists = k.groupByAdminIndex.Has(ctx, []byte("new-admin-address"))
	require.False(t, exists)
}

func TestKeeperEndToEndWithNaturalKeyTable(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	ctx := NewMockContext()

	k := NewGroupKeeper(storeKey)

	g := GroupMetadata{
		Description: "my test",
		Admin:       sdk.AccAddress([]byte("admin-address")),
	}

	m := GroupMember{
		Group:  sdk.AccAddress(EncodeSequence(1)),
		Member: sdk.AccAddress([]byte("member-address")),
		Weight: 10,
	}
	groupRowID, err := k.groupTable.Create(ctx, &g)
	require.NoError(t, err)
	require.Equal(t, uint64(1), groupRowID)
	// when stored
	err = k.groupMemberTable.Create(ctx, &m)
	require.NoError(t, err)

	// then we should find it by natural key
	naturalKey := m.NaturalKey()
	exists := k.groupMemberTable.Has(ctx, naturalKey)
	require.True(t, exists)
	// and load it by natural key
	var loaded GroupMember
	rowID, err := k.groupMemberTable.GetOne(ctx, naturalKey, &loaded)
	require.NoError(t, err)

	// then values should match expectations
	require.Equal(t, EncodeSequence(1), rowID)
	require.Equal(t, m, loaded)

	// and then the data should exists in MultiKeyIndex
	exists = k.groupMemberByGroupIndex.Has(ctx, EncodeSequence(groupRowID))
	require.True(t, exists)

	// and when loaded from MultiKeyIndex
	it, err := k.groupMemberByGroupIndex.Get(ctx, EncodeSequence(groupRowID))
	require.NoError(t, err)

	// then values should match as before
	rowID, err = First(it, &loaded)
	require.NoError(t, err)

	assert.Equal(t, EncodeSequence(groupRowID), rowID)
	assert.Equal(t, m, loaded)
	// and when we create another entry with the same natural key
	err = k.groupMemberTable.Create(ctx, &m)
	// then it should fail as the natural key must be unique
	require.True(t, ErrUniqueConstraint.Is(err))

	// and when entity updated with new natural key
	updatedMember := &GroupMember{
		Group:  m.Group,
		Member: []byte("new-member-address"),
		Weight: m.Weight,
	}
	// then it should fail as the natural key is immutable
	err = k.groupMemberTable.Save(ctx, updatedMember)
	require.Error(t, err)

	// and when entity updated with non natural key attribute modified
	updatedMember = &GroupMember{
		Group:  m.Group,
		Member: m.Member,
		Weight: 99,
	}
	// then it should not fail
	err = k.groupMemberTable.Save(ctx, updatedMember)
	require.NoError(t, err)

	// and when entity deleted
	err = k.groupMemberTable.Delete(ctx, &m)
	require.NoError(t, err)

	// then it is removed from natural key MultiKeyIndex
	exists = k.groupMemberTable.Has(ctx, naturalKey)
	require.False(t, exists)

	// and removed from secondary MultiKeyIndex
	exists = k.groupMemberByGroupIndex.Has(ctx, EncodeSequence(groupRowID))
	require.False(t, exists)
}

func TestGasCostsNaturalKeyTable(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()
	ctx := NewMockContext()

	k := NewGroupKeeper(storeKey, cdc)

	g := GroupMetadata{
		Description: "my test",
		Admin:       sdk.AccAddress([]byte("admin-address")),
	}

	m := GroupMember{
		Group:  sdk.AccAddress(EncodeSequence(1)),
		Member: sdk.AccAddress([]byte("member-address")),
		Weight: sdk.NewInt(10),
	}
	groupRowID, err := k.groupTable.Create(ctx, &g)
	require.NoError(t, err)
	require.Equal(t, uint64(1), groupRowID)
	gCtx := NewGasCountingMockContext(ctx)
	err = k.groupMemberTable.Create(gCtx, &m)
	require.NoError(t, err)
	t.Logf("gas consumed on create: %d", gCtx.GasConsumed())

	// get by natural key
	gCtx.ResetGasMeter()
	var loaded GroupMember
	_, err = k.groupMemberTable.GetOne(gCtx, m.NaturalKey(), &loaded)
	require.NoError(t, err)
	t.Logf("gas consumed on get by natural key: %d", gCtx.GasConsumed())

	// get by rowID
	gCtx.ResetGasMeter()
	_, err = k.groupMemberTable.autoTable.GetOne(gCtx, 1, &loaded)
	require.NoError(t, err)
	t.Logf("gas consumed on get by rowID: %d", gCtx.GasConsumed())

	// get by secondary index
	gCtx.ResetGasMeter()
	// and when loaded from MultiKeyIndex
	it, err := k.groupMemberByGroupIndex.Get(gCtx, EncodeSequence(groupRowID))
	require.NoError(t, err)
	var loadedSlice []GroupMember
	_, err = ReadAll(it, &loadedSlice)
	require.NoError(t, err)

	t.Logf("gas consumed on get by multi index key: %d", gCtx.GasConsumed())

	// delete
	gCtx.ResetGasMeter()
	err = k.groupMemberTable.Delete(gCtx, m)
	require.NoError(t, err)
	t.Logf("gas consumed on delete by natural key: %d", gCtx.GasConsumed())

	// with 3 elements
	for i := 1; i < 4; i++ {
		gCtx.ResetGasMeter()
		m := GroupMember{
			Group:  sdk.AccAddress(EncodeSequence(1)),
			Member: sdk.AccAddress([]byte(fmt.Sprintf("member-addres%d", i))),
			Weight: sdk.NewInt(10),
		}
		err = k.groupMemberTable.Create(gCtx, &m)
		require.NoError(t, err)
		t.Logf("%d: gas consumed on create: %d", i, gCtx.GasConsumed())
	}

	for i := 1; i < 4; i++ {
		gCtx.ResetGasMeter()
		m := GroupMember{
			Group:  sdk.AccAddress(EncodeSequence(1)),
			Member: sdk.AccAddress([]byte(fmt.Sprintf("member-addres%d", i))),
			Weight: sdk.NewInt(10),
		}
		_, err = k.groupMemberTable.GetOne(gCtx, m.NaturalKey(), &loaded)
		require.NoError(t, err)
		t.Logf("%d: gas consumed on get by natural key: %d", i, gCtx.GasConsumed())
	}

	// get by secondary index
	gCtx.ResetGasMeter()
	// and when loaded from MultiKeyIndex
	it, err = k.groupMemberByGroupIndex.Get(gCtx, EncodeSequence(groupRowID))
	require.NoError(t, err)
	_, err = ReadAll(it, &loadedSlice)
	require.NoError(t, err)
	require.Len(t, loadedSlice, 3)
	t.Logf("gas consumed on get by multi index key: %d", gCtx.GasConsumed())

	// delete
	for i, m := range loadedSlice {
		gCtx.ResetGasMeter()

		err = k.groupMemberTable.Delete(gCtx, &m)
		require.NoError(t, err)
		t.Logf("%d: gas consumed on delete: %d", i, gCtx.GasConsumed())
	}
}

func first(t *testing.T, it Iterator) ([]byte, GroupMetadata) {
	var loaded GroupMetadata
	key, err := First(it, &loaded)
	require.NoError(t, err)
	return key, loaded
}
