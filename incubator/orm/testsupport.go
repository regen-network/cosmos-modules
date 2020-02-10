package orm

import (
	"io"

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

type AlwaysPanicKVStore struct{}

func (a AlwaysPanicKVStore) GetStoreType() types.StoreType {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) CacheWrap() types.CacheWrap {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) Get(key []byte) []byte {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) Has(key []byte) bool {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) Set(key, value []byte) {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) Delete(key []byte) {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) Iterator(start, end []byte) types.Iterator {
	panic("Not implemented")
}

func (a AlwaysPanicKVStore) ReverseIterator(start, end []byte) types.Iterator {
	panic("Not implemented")
}
