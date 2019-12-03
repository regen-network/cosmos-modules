package orm

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tm-db"
)

type MockContext struct{
	db *dbm.MemDB
}
func NewMockContext()*MockContext {
	db := dbm.NewMemDB()
	return &MockContext{	db}
}
func (m MockContext) KVStore(key sdk.StoreKey) sdk.KVStore{
	store := store.NewCommitMultiStore(m.db)
	store.MountStoreWithDB(key, sdk.StoreTypeMulti, m.db)
	return store.GetCommitKVStore(key)
}

func TestKeeperEndToEnd(t *testing.T) {
	const isCheckTx = false
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()
	ctx := NewMockContext()

	k := NewGroupKeeper(storeKey, cdc)

	//g := GroupMember{
	//	Group:  sdk.AccAddress([]byte("group-address")),
	//	Member: sdk.AccAddress([]byte("alice-address")),
	//	Weight: sdk.NewInt(100),
	//}
	g := GroupMetadata{
		Description: "my test",
		Admin:  sdk.AccAddress([]byte("admin-address")),
	}
	// when stored
	groupKey, err:= k.groupTable.Create(ctx, &g)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then we should find it
	it, err:= k.groupTable.Get(ctx, groupKey)
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
