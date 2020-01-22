package orm

import (
	"encoding/binary"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"
)

type MockContext struct {
	db    *dbm.MemDB
	store types.CommitMultiStore
}

func NewMockContext() *MockContext {
	db := dbm.NewMemDB()
	return &MockContext{
		db:    dbm.NewMemDB(),
		store: store.NewCommitMultiStore(db),
	}

}
func (m MockContext) KVStore(key sdk.StoreKey) sdk.KVStore {
	if s := m.store.GetCommitKVStore(key); s != nil {
		return s
	}
	m.store.MountStoreWithDB(key, sdk.StoreTypeIAVL, m.db)
	if err := m.store.LoadLatestVersion(); err != nil {
		panic(err)
	}
	return m.store.GetCommitKVStore(key)
}

func TestKeeperEndToEndWithAutoUInt64Table(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()
	ctx := NewMockContext()

	k := NewGroupKeeper(storeKey, cdc)

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
		Weight: sdk.NewInt(99),
	}
	// then it should not fail
	err = k.groupMemberTable.Save(ctx, updatedMember)
	require.NoError(t, err)

	// and when entity deleted
	err = k.groupMemberTable.Delete(ctx, m)
	require.NoError(t, err)

	// then it is removed from natural key MultiKeyIndex
	exists = k.groupMemberTable.Has(ctx, naturalKey)
	require.False(t, exists)

	// and removed from secondary MultiKeyIndex
	exists = k.groupMemberByGroupIndex.Has(ctx, EncodeSequence(groupRowID))
	require.False(t, exists)
}

func first(t *testing.T, it Iterator) ([]byte, GroupMetadata) {
	var loaded GroupMetadata
	key, err := First(it, &loaded)
	require.NoError(t, err)
	return key, loaded
}
