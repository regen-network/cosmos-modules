package group

import (
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func NewContext(keys ...sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := rootmulti.NewStore(db)
	for _, v := range keys {
		storeType := sdk.StoreTypeIAVL
		if _, ok := v.(*sdk.TransientStoreKey); ok {
			storeType = sdk.StoreTypeTransient
		}
		cms.MountStoreWithDB(v, storeType, db)
		cms.LoadLatestVersion()
	}
	return sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())

}
