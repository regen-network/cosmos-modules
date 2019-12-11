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

func TestKeeperEndToEnd(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()
	ctx := NewMockContext()

	k := NewGroupKeeper(storeKey, cdc)

	g := GroupMetadata{
		Description: "my test",
		Admin:       sdk.AccAddress([]byte("admin-address")),
	}
	// when stored
	groupKey, err := k.groupTable.Create(ctx, &g)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then we should find it
	exists, _ := k.groupTable.Has(ctx, groupKey)
	if exp, got := true, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	// and load it
	it, err := k.groupTable.Get(ctx, groupKey)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	binKey, loaded := first(t, it)
	if exp, got := groupKey, binary.BigEndian.Uint64(binKey); exp != got {
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
	if exp, got := groupKey, binary.BigEndian.Uint64(binKey); exp != got {
		t.Errorf("expected %v but got %v", exp, got)
	}
	if exp, got := g, loaded; !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v but got %v", exp, got)
	}

	// when deleted
	k.groupTable.Delete(ctx, groupKey)

	// then removed from primary index
	exists, _ = k.groupTable.Has(ctx, groupKey)
	if exp, got := false, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	// and removed from secondary index
	exists, _ = k.groupByAdminIndex.Has(ctx, []byte("admin-address"))
	if exp, got := false, exists; exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
}

func first(t *testing.T, it Iterator) ([]byte, GroupMetadata) {
	t.Helper()
	defer it.Close()
	var loaded GroupMetadata
	binKey, err := it.LoadNext(&loaded)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	return binKey, loaded
}
