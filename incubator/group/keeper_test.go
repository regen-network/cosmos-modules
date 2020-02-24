package group

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(ModuleCdc.amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, &MockProposalModel{})
	ctx := NewContext(pKey, pTKey, groupKey)
	k.setParams(ctx, DefaultParams())

	members := []Member{{
		Address: sdk.AccAddress([]byte("member-address")),
		Power:   sdk.NewDec(1),
		Comment: "first",
	}, {
		Address: sdk.AccAddress([]byte("other-member-address")),
		Power:   sdk.NewDec(2),
		Comment: "second",
	}}

	id, err := k.CreateGroup(ctx, []byte("admin-address"), members, "test")
	require.NoError(t, err)
	assert.Equal(t, GroupID(1), id)
	// then all data persisted
	loaded, err := k.GetGroup(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, sdk.AccAddress([]byte("admin-address")), loaded.Admin)
	assert.Equal(t, "test", loaded.Comment)
	assert.Equal(t, id, loaded.Group)
	assert.Equal(t, uint64(1), loaded.Version)

	// and members are stored as well
	it, err := k.groupMemberByGroupIndex.Get(ctx, uint64(id))
	require.NoError(t, err)
	var loadedMembers []GroupMember
	_, err = orm.ReadAll(it, &loadedMembers)
	require.NoError(t, err)
	assert.Equal(t, len(members), len(loadedMembers))
	for i := range loadedMembers {
		assert.Equal(t, members[i].Comment, loadedMembers[i].Comment)
		assert.Equal(t, members[i].Address, loadedMembers[i].Member)
		assert.Equal(t, members[i].Power, loadedMembers[i].Weight)
		assert.Equal(t, id, loadedMembers[i].Group)
	}
}

func TestLoadParam(t *testing.T) {
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(ModuleCdc.amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, &MockProposalModel{})

	ctx := NewContext(pKey, pTKey, groupKey)

	myParams := Params{MaxCommentLength: 1}
	paramSpace.SetParamSet(ctx, &myParams)

	assert.Equal(t, myParams, k.getParams(ctx))
}
