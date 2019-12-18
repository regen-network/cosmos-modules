package orm

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then we should find it
	exists, _ := k.groupTable.Has(ctx, rowID)
	if exp, got := true, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	// and load it
	it, err := k.groupTable.Get(ctx, rowID)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	binKey, loaded := first(t, it)
	if exp, got := rowID, binary.BigEndian.Uint64(binKey); exp != got {
		t.Errorf("expected %v but got %v", exp, got)
	}
	if exp, got := "my test", loaded.Description; exp != got {
		t.Errorf("expected %v but got %v", exp, got)
	}
	if exp, got := sdk.AccAddress([]byte("admin-address")), loaded.Admin; !bytes.Equal(exp, got) {
		t.Errorf("expected %X but got %X", exp, got)
	}
	// and exists in index
	exists, err = k.groupByAdminIndex.Has(ctx, []byte("admin-address"))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if !exists {
		t.Fatalf("expected entry to exist")
	}
	// and when loaded
	it, err = k.groupByAdminIndex.Get(ctx, []byte("admin-address"))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then
	binKey, loaded = first(t, it)
	if exp, got := rowID, binary.BigEndian.Uint64(binKey); exp != got {
		t.Errorf("expected %v but got %v", exp, got)
	}
	if exp, got := g, loaded; !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v but got %v", exp, got)
	}
	// when updated
	g.Admin = []byte("new-admin-address")
	err = k.groupTable.Save(ctx, rowID, &g)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then indexes are updated, too
	exists, _ = k.groupByAdminIndex.Has(ctx, []byte("new-admin-address"))
	if exp, got := true, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	exists, _ = k.groupByAdminIndex.Has(ctx, []byte("admin-address"))
	if exp, got := false, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}

	// when deleted
	k.groupTable.Delete(ctx, rowID)

	// then removed from primary index
	exists, _ = k.groupTable.Has(ctx, rowID)
	if exp, got := false, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	// and also removed from secondary index
	exists, _ = k.groupByAdminIndex.Has(ctx, []byte("new-admin-address"))
	if exp, got := false, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
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
	groupPrimKey := EncodeSequence(1)
	m := GroupMember{
		Group:  sdk.AccAddress(groupPrimKey),
		Member: sdk.AccAddress([]byte("member-address")),
		Weight: sdk.NewInt(10),
	}
	if _, err := k.groupTable.Create(ctx, &g); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// when stored
	err := k.groupMemberTable.Create(ctx, &m)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	// then we should find it by natural key
	naturalKey := m.ID()
	exists, _ := k.groupMemberTable.Has(ctx, naturalKey)
	if exp, got := true, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	// and load it by natural key
	it, err := k.groupMemberTable.Get(ctx, naturalKey)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	var loaded GroupMember
	rowID, err := First(it, &loaded)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then values should match expectations
	if exp, got := EncodeSequence(1), rowID; !bytes.Equal(exp, got) {
		t.Fatalf("expected %X but got %X", exp, got)
	}
	if exp, got := m, loaded; !reflect.DeepEqual(exp, got) {
		t.Fatalf("expected %#v but got %#v", exp, got)
	}

	// and then the data should exists in index
	exists, err = k.groupMemberByGroupIndex.Has(ctx, groupPrimKey)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if !exists {
		t.Fatalf("expected entry to exist")
	}
	// and when loaded from index
	it, err = k.groupMemberByGroupIndex.Get(ctx, groupPrimKey)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then values should match as before
	rowID, err = First(it, &loaded)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if exp, got := groupPrimKey, rowID; !bytes.Equal(exp, got) {
		t.Errorf("expected %X but got %X", exp, got)
	}
	if exp, got := m, loaded; !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v but got %v", exp, got)
	}
	// and when we create another entry with the same natural key
	err = k.groupMemberTable.Create(ctx, &m)
	// then it should fail as the natural key must be unique
	if !ErrUniqueConstraint.Is(err) {
		t.Fatal("expected error but got %#V", err)
	}

	// and when entity updated with new natural key
	updatedMember := &GroupMember{
		Group:  m.Group,
		Member: []byte("new-member-address"),
		Weight: m.Weight,
	}
	err = k.groupMemberTable.Save(ctx, updatedMember)

	// then it should fail as the natural key is immutable
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	// and when entity updated with non natural key attribute modified
	updatedMember = &GroupMember{
		Group:  m.Group,
		Member: m.Member,
		Weight: sdk.NewInt(99),
	}
	err = k.groupMemberTable.Save(ctx, updatedMember)

	// then it should not fail
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	// and when entity deleted
	err = k.groupMemberTable.Delete(ctx, m)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then it is removed from natural key index
	exists, err = k.groupMemberTable.Has(ctx, naturalKey)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if exp, got := false, exists; exp != got {
		t.Errorf("expected %v but got %v", exp, got)
	}
	// and removed from secondary index
	exists, _ = k.groupMemberByGroupIndex.Has(ctx, groupPrimKey)
	if exp, got := false, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
}

func first(t *testing.T, it Iterator) ([]byte, GroupMetadata) {
	t.Helper()
	var loaded GroupMetadata
	key, err := First(it, &loaded)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	return key, loaded
}
