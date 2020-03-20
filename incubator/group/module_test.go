package group

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestInitWithDefaultGenesis(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &MockProposalI{})
	ctx := NewContext(pKey, pTKey, groupKey)

	module := NewAppModule(k)
	assert.Equal(t, []abci.ValidatorUpdate{}, module.InitGenesis(ctx, module.DefaultGenesis()))
	assert.Equal(t, DefaultParams(), k.GetParams(ctx))
}
