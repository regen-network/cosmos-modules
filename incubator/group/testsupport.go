package group

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
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

func createGroupKeeper() (Keeper, sdk.Context) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &MockProposalI{})
	ctx := NewContext(pKey, pTKey, groupKey)
	k.setParams(ctx, DefaultParams())
	return k, ctx
}

type MockProposalI struct {
}

func (m MockProposalI) Marshal() ([]byte, error) {
	panic("implement me")
}

func (m MockProposalI) Unmarshal([]byte) error {
	panic("implement me")
}

func (m MockProposalI) GetBase() ProposalBase {
	panic("implement me")
}

func (m MockProposalI) SetBase(ProposalBase) {
	panic("implement me")
}

func (m MockProposalI) GetMsgs() []sdk.Msg {
	panic("implement me")
}

func (m MockProposalI) SetMsgs([]sdk.Msg) error {
	panic("implement me")
}

func (m MockProposalI) ValidateBasic() error {
	panic("implement me")
}
