package group

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerier(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &MockProposalI{})
	ctx := NewContext(pKey, pTKey, groupKey)
	k.setParams(ctx, DefaultParams())

	_, err := k.CreateGroup(ctx, []byte("one-admin-address"), nil, "example1")
	require.NoError(t, err)
	_, err = k.CreateGroup(ctx, []byte("other-admin-address"), nil, "example2")
	require.NoError(t, err)

	q := NewQuerier(k)
	specs := map[string]struct {
		srcPath     string
		srcData     []byte
		expModelLen int
		expErr      *errors.Error
	}{
		"query table for single entry": {
			srcPath:     "xgroup",
			srcData:     orm.EncodeSequence(1),
			expModelLen: 1,
		},
		"query table to find all in range": {
			srcPath:     "xgroup?range",
			expModelLen: 2,
		},
		"query table to find all with prefix": {
			srcPath:     "xgroup?prefix",
			srcData:     []byte{0},
			expModelLen: 2,
		},
		"query index for single entity": {
			srcPath:     "xgroup/admin",
			srcData:     []byte("one-admin-address"),
			expModelLen: 1,
		},
		"query index to find all in range": {
			srcPath:     "xgroup/admin?range",
			expModelLen: 2,
		},
		"query index to find all with prefix": {
			srcPath:     "xgroup/admin?prefix",
			srcData:     []byte("o"),
			expModelLen: 2,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			data, err := q(ctx, []string{spec.srcPath}, abci.RequestQuery{Path: spec.srcPath, Data: spec.srcData})
			require.True(t, spec.expErr.Is(err), "unexpected error", err)
			t.Logf("%s", string(data))
			var res map[string]interface{}
			require.NoError(t, json.Unmarshal(data, &res))
			require.Contains(t, res, "data")
			require.Len(t, res["data"], spec.expModelLen)
		},
		)
	}
}
