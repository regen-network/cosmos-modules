package group_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	subspace "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/modules/incubator/group"
	"github.com/cosmos/modules/incubator/group/testdata"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/types"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenesisImportExportParameters(t *testing.T) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	srcK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})

	f := fuzz.New()
	for i := 0; i < 100; i++ {
		var params group.Params
		f.Fuzz(&params)
		srcCtx := group.NewContext(pKey, pTKey, groupKey)
		paramSpace.SetParamSet(srcCtx, &params)

		raw := group.NewAppModule(srcK).ExportGenesis(srcCtx, cdc)

		destK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
		destCtx := group.NewContext(pKey, pTKey, groupKey)

		_ = group.NewAppModule(destK).InitGenesis(destCtx, cdc, raw)
		require.Equal(t, srcK.GetParams(srcCtx), destK.GetParams(destCtx))
	}
}

func TestGenesisImportExportGroups(t *testing.T) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	srcK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
	srcCtx := group.NewContext(pKey, pTKey, groupKey)

	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(srcCtx, &defaultParams)

	f := fuzz.New().Funcs(group.FuzzAddr, group.FuzzPositiveDec, group.FuzzComment)
	for i := 0; i < 100; i++ {
		var (
			members []group.Member
			admin   sdk.AccAddress
			descr   string
		)
		f.Fuzz(&members)
		f.NilChance(0).Fuzz(&admin)
		f.Fuzz(&descr)

		groupID, err := srcK.CreateGroup(srcCtx, admin, members, descr)
		require.NoError(t, err)

		raw := group.NewAppModule(srcK).ExportGenesis(srcCtx, cdc)

		destK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
		destCtx := group.NewContext(pKey, pTKey, groupKey)

		_ = group.NewAppModule(destK).InitGenesis(destCtx, cdc, raw)

		got, err := srcK.GetGroup(srcCtx, groupID)
		require.NoError(t, err)
		exp, err := destK.GetGroup(destCtx, groupID)
		require.NoError(t, err)
		assert.Equal(t, got, exp)

		// sequence
		assert.Equal(t, srcK.GetGroupSeqValue(srcCtx), destK.GetGroupSeqValue(srcCtx))
	}
	// sanity check indexes
	srcIT, err := srcK.GroupByAdminIndex.PrefixScan(srcCtx, nil, nil)
	require.NoError(t, err)
	destIT, err := srcK.GroupByAdminIndex.PrefixScan(srcCtx, nil, nil)
	require.NoError(t, err)

	for {
		var exp, got group.GroupMetadata
		expID, srcErr := srcIT.LoadNext(&exp)
		gotID, gotErr := destIT.LoadNext(&got)

		if orm.ErrIteratorDone.Is(srcErr) {
			assert.True(t, orm.ErrIteratorDone.Is(gotErr))
			return
		}
		require.NoError(t, srcErr)
		require.NoError(t, gotErr)
		assert.Equal(t, exp, got)
		assert.Equal(t, expID, gotID)
	}
}

func TestGenesisImportExportGroupMembers(t *testing.T) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	srcK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
	srcCtx := group.NewContext(pKey, pTKey, groupKey)

	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(srcCtx, &defaultParams)
	f := fuzz.New().Funcs(group.FuzzAddr, group.FuzzPositiveDec, group.FuzzComment)
	for i := 0; i < 100; i++ {
		var (
			admin sdk.AccAddress
			descr string
		)
		f.NilChance(0).Fuzz(&admin)
		f.Fuzz(&descr)

		var members []group.Member
		f.Funcs(group.FuzzGroupMember).Fuzz(&members)

		groupID, err := srcK.CreateGroup(srcCtx, admin, members, descr)
		require.NoError(t, err)

		raw := group.NewAppModule(srcK).ExportGenesis(srcCtx, cdc)

		destK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
		destCtx := group.NewContext(pKey, pTKey, groupKey)

		_ = group.NewAppModule(destK).InitGenesis(destCtx, cdc, raw)

		for i, m := range members {
			exp, err := srcK.GetGroupMember(srcCtx, groupID, m.Address)
			require.NoError(t, err)
			got, err := destK.GetGroupMember(destCtx, groupID, m.Address)
			require.NoError(t, err)
			assert.Equal(t, exp, got, i)
		}
	}
}

func TestGenesisImportExportGroupAccount(t *testing.T) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	srcK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
	srcCtx := group.NewContext(pKey, pTKey, groupKey)

	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(srcCtx, &defaultParams)
	f := fuzz.New().Funcs(group.FuzzAddr, group.FuzzPositiveDec, group.FuzzComment, group.FuzzPositiveDuration)
	for i := 0; i < 100; i++ {
		var (
			admin   sdk.AccAddress
			comment string
			policy  group.ThresholdDecisionPolicy
		)
		f.NilChance(0).Fuzz(&admin)
		f.Fuzz(&policy)
		f.Fuzz(&comment)

		groupID, err := srcK.CreateGroup(srcCtx, admin, nil, comment)
		require.NoError(t, err)
		acc, err := srcK.CreateGroupAccount(srcCtx, admin, groupID, policy, comment)
		require.NoError(t, err)

		raw := group.NewAppModule(srcK).ExportGenesis(srcCtx, cdc)

		destK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
		destCtx := group.NewContext(pKey, pTKey, groupKey)

		_ = group.NewAppModule(destK).InitGenesis(destCtx, cdc, raw)

		exp, err := srcK.GetGroupAccount(srcCtx, acc)
		require.NoError(t, err)
		got, err := destK.GetGroupAccount(destCtx, acc)
		require.NoError(t, err)
		assert.Equal(t, exp, got, i)
		assert.Equal(t, srcK.GetGroupAccountSeqValue(srcCtx), destK.GetGroupAccountSeqValue(destCtx))
	}
}

func TestGenesisImportExportProposals(t *testing.T) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	srcK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &testdata.MyAppProposal{})
	srcCtx := group.NewContext(pKey, pTKey, groupKey)

	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(srcCtx, &defaultParams)
	f := fuzz.New().Funcs(group.FuzzAddr, group.FuzzPositiveDec, group.FuzzComment, group.FuzzPositiveDuration)
	var (
		admin   sdk.AccAddress
		policy  group.ThresholdDecisionPolicy
		members []group.Member
	)
	f.Fuzz(&admin)
	f.Fuzz(&policy)
	f.NilChance(0).Fuzz(&members)
	policy.Threshold = sdk.OneDec()

	groupID, err := srcK.CreateGroup(srcCtx, admin, members, "test")
	require.NoError(t, err)
	acc, err := srcK.CreateGroupAccount(srcCtx, admin, groupID, policy, "test")
	require.NoError(t, err)

	f = fuzz.New()
	for i := 0; i < 100; i++ {
		var (
			comment   string
			msgs      []sdk.Msg
			proposers []sdk.AccAddress
			blockTime time.Time
		)
		f.Fuzz(&comment)
		f.Funcs(group.FuzzAddr, testdata.FuzzPayloadMsg).Fuzz(&msgs)
		f.NilChance(0).Funcs(fuzzAddressesFrom(members)).Fuzz(&proposers)
		f.NilChance(0).Fuzz(&blockTime)

		srcCtx = srcCtx.WithBlockTime(blockTime)
		propID, err := srcK.CreateProposal(srcCtx, acc, comment, proposers, msgs)
		require.NoError(t, err)
		raw := group.NewAppModule(srcK).ExportGenesis(srcCtx, cdc)

		destK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &testdata.MyAppProposal{})
		destCtx := group.NewContext(pKey, pTKey, groupKey)

		_ = group.NewAppModule(destK).InitGenesis(destCtx, cdc, raw)

		exp, err := srcK.GetProposal(srcCtx, propID)
		require.NoError(t, err)
		got, err := destK.GetProposal(destCtx, propID)
		require.NoError(t, err)
		assert.Equal(t, exp, got, i)
		assert.Equal(t, srcK.GetGroupAccountSeqValue(srcCtx), destK.GetGroupAccountSeqValue(destCtx))
	}
}

func TestGenesisImportExportVotes(t *testing.T) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	srcK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &testdata.MyAppProposal{})
	srcCtx := group.NewContext(pKey, pTKey, groupKey)

	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(srcCtx, &defaultParams)
	f := fuzz.New().Funcs(group.FuzzAddr, group.FuzzPositiveDec, group.FuzzComment, group.FuzzPositiveDuration)
	var (
		admin   sdk.AccAddress
		comment string
		members []group.Member
	)
	f.Fuzz(&admin)
	f.NilChance(0).NumElements(1, 10).Fuzz(&members)
	f.Fuzz(&comment)

	groupID, err := srcK.CreateGroup(srcCtx, admin, members, comment)
	require.NoError(t, err)
	loadedGroup, _ := srcK.GetGroup(srcCtx, groupID)
	allAgreePolicy := group.ThresholdDecisionPolicy{Timout: types.Duration{Seconds: 999}, Threshold: loadedGroup.TotalWeight.Sub(sdk.OneDec())}
	acc, err := srcK.CreateGroupAccount(srcCtx, admin, groupID, allAgreePolicy, comment)
	require.NoError(t, err)

	f = fuzz.New()
	for i := 0; i < 100; i++ {
		var (
			voters    []sdk.AccAddress
			choice    group.Choice
			comment   string
			blockTime time.Time
		)
		f.NilChance(0).Funcs(fuzzAddressesFrom(members)).Fuzz(&voters)
		f.Fuzz(&comment)
		f.NilChance(0).Funcs(group.FuzzChoice).Fuzz(&choice)
		f.NilChance(0).Fuzz(&blockTime)

		srcCtx = srcCtx.WithBlockTime(blockTime)
		propID, err := srcK.CreateProposal(srcCtx, acc, comment, []sdk.AccAddress{members[0].Address}, nil)
		require.NoError(t, err)

		err = srcK.Vote(srcCtx, propID, voters, choice, comment)
		require.NoError(t, err)

		raw := group.NewAppModule(srcK).ExportGenesis(srcCtx, cdc)

		destK := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &testdata.MyAppProposal{})
		destCtx := group.NewContext(pKey, pTKey, groupKey)

		_ = group.NewAppModule(destK).InitGenesis(destCtx, cdc, raw)

		for _, voter := range voters {
			exp, err := srcK.GetVote(srcCtx, propID, voter)
			require.NoError(t, err)
			got, err := destK.GetVote(destCtx, propID, voter)
			require.NoError(t, err)
			assert.Equal(t, exp, got)
		}
	}
}

func fuzzAddressesFrom(members []group.Member) func(m *[]sdk.AccAddress, c fuzz.Continue) {
	return func(m *[]sdk.AccAddress, c fuzz.Continue) {
		n := 1
		if len(members) > 1 {
			n = c.Intn(len(members)-1) + 1
		}
		a := make([]sdk.AccAddress, n)
		for i := 0; i < n; i++ {
			a[i] = members[i].Address
		}
		*m = a
	}
}
