package group

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &MockProposalI{})
	ctx := NewContext(pKey, pTKey, groupKey)
	k.setParams(ctx, DefaultParams())

	members := []Member{{
		Address: sdk.AccAddress([]byte("one--member--address")),
		Power:   sdk.NewDec(1),
		Comment: "first",
	}, {
		Address: sdk.AccAddress([]byte("other-member-address")),
		Power:   sdk.NewDec(2),
		Comment: "second",
	}}
	specs := map[string]struct {
		srcAdmin   sdk.AccAddress
		srcMembers []Member
		srcComment string
		expErr     bool
	}{
		"all good": {
			srcAdmin:   []byte("valid--admin-address"),
			srcMembers: members,
			srcComment: "test",
		},
		"group comment too long": {
			srcAdmin:   []byte("valid--admin-address"),
			srcMembers: members,
			srcComment: strings.Repeat("a", 256),
			expErr:     true,
		},
		"member comment too long": {
			srcAdmin: []byte("valid--admin-address"),
			srcMembers: []Member{{
				Address: []byte("valid-member-address"),
				Power:   sdk.OneDec(),
				Comment: strings.Repeat("a", 256),
			}},
			srcComment: "test",
			expErr:     true,
		},
	}
	var seq uint32
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			id, err := k.CreateGroup(ctx, spec.srcAdmin, spec.srcMembers, spec.srcComment)
			if spec.expErr {
				require.Error(t, err)
				require.False(t, k.groupTable.Has(ctx, GroupID(seq+1).Bytes()))
				return
			}
			require.NoError(t, err)

			seq++
			assert.Equal(t, GroupID(seq), id)

			// then all data persisted
			loaded, err := k.GetGroup(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, sdk.AccAddress([]byte(spec.srcAdmin)), loaded.Admin)
			assert.Equal(t, spec.srcComment, loaded.Comment)
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
		})
	}

}

func TestLoadParam(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &MockProposalI{})

	ctx := NewContext(pKey, pTKey, groupKey)

	myParams := Params{MaxCommentLength: 1}
	paramSpace.SetParamSet(ctx, &myParams)

	assert.Equal(t, myParams, k.getParams(ctx))
}
