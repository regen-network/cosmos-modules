package group_test

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/modules/incubator/group"
	"github.com/cosmos/modules/incubator/group/testdata"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	k := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
	ctx := group.NewContext(pKey, pTKey, groupKey)
	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(ctx, &defaultParams)

	members := []group.Member{{
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
		srcMembers []group.Member
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
			srcMembers: []group.Member{{
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
				require.False(t, k.HasGroup(ctx, group.GroupID(seq+1).Bytes()))
				return
			}
			require.NoError(t, err)

			seq++
			assert.Equal(t, group.GroupID(seq), id)

			// then all data persisted
			loadedGroup, err := k.GetGroup(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, sdk.AccAddress([]byte(spec.srcAdmin)), loadedGroup.Admin)
			assert.Equal(t, spec.srcComment, loadedGroup.Comment)
			assert.Equal(t, id, loadedGroup.Group)
			assert.Equal(t, uint64(1), loadedGroup.Version)

			// and members are stored as well
			it, err := k.GetGroupMembersByGroup(ctx, id)
			require.NoError(t, err)
			var loadedMembers []group.GroupMember
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

func TestCreateGroupAccount(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	k := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})
	ctx := group.NewContext(pKey, pTKey, groupKey)
	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(ctx, &defaultParams)

	myGroupID, err := k.CreateGroup(ctx, []byte("valid--admin-address"), nil, "test")
	require.NoError(t, err)

	specs := map[string]struct {
		srcAdmin   sdk.AccAddress
		srcGroupID group.GroupID
		srcPolicy  group.ThresholdDecisionPolicy
		srcComment string
		expErr     bool
	}{
		"all good": {
			srcAdmin:   []byte("valid--admin-address"),
			srcComment: "test",
			srcGroupID: myGroupID,
			srcPolicy: group.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Timout:    types.Duration{Seconds: 1},
			},
		},
		"decision policy threshold > total group weight": {
			srcAdmin:   []byte("valid--admin-address"),
			srcComment: "test",
			srcGroupID: myGroupID,
			srcPolicy: group.ThresholdDecisionPolicy{
				Threshold: sdk.NewDec(math.MaxInt64),
				Timout:    types.Duration{Seconds: 1},
			},
		},
		"group id does not exists": {
			srcAdmin:   []byte("valid--admin-address"),
			srcComment: "test",
			srcGroupID: 9999,
			srcPolicy: group.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Timout:    types.Duration{Seconds: 1},
			},
			expErr: true,
		},
		"admin not group admin": {
			srcAdmin:   []byte("other--admin-address"),
			srcComment: "test",
			srcGroupID: myGroupID,
			srcPolicy: group.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Timout:    types.Duration{Seconds: 1},
			},
			expErr: true,
		},
		"comment too long": {
			srcAdmin:   []byte("valid--admin-address"),
			srcComment: strings.Repeat("a", 256),
			srcGroupID: myGroupID,
			srcPolicy: group.ThresholdDecisionPolicy{
				Threshold: sdk.OneDec(),
				Timout:    types.Duration{Seconds: 1},
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			addr, err := k.CreateGroupAccount(ctx, spec.srcAdmin, spec.srcGroupID, spec.srcPolicy, spec.srcComment)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// then all data persisted
			groupAccount, err := k.GetGroupAccount(ctx, addr)
			require.NoError(t, err)
			assert.Equal(t, addr, groupAccount.Base.GroupAccount)
			assert.Equal(t, myGroupID, groupAccount.Base.Group)
			assert.Equal(t, sdk.AccAddress([]byte(spec.srcAdmin)), groupAccount.Base.Admin)
			assert.Equal(t, spec.srcComment, groupAccount.Base.Comment)
			assert.Equal(t, uint64(1), groupAccount.Base.Version)
			assert.Equal(t, &spec.srcPolicy, groupAccount.DecisionPolicy.GetDecisionPolicy())
		})
	}
}

func TestCreateProposal(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	k := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &testdata.MyAppProposal{})
	blockTime := time.Now()
	ctx := group.NewContext(pKey, pTKey, groupKey).WithBlockTime(blockTime)
	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(ctx, &defaultParams)

	members := []group.Member{{
		Address: []byte("valid-member-address"),
		Power:   sdk.OneDec(),
	}}
	myGroupID, err := k.CreateGroup(ctx, []byte("valid--admin-address"), members, "test")
	require.NoError(t, err)

	policy := group.ThresholdDecisionPolicy{
		Threshold: sdk.OneDec(),
		Timout:    types.Duration{Seconds: 1},
	}
	accountAddr, err := k.CreateGroupAccount(ctx, []byte("valid--admin-address"), myGroupID, policy, "test")
	require.NoError(t, err)

	policy = group.ThresholdDecisionPolicy{
		Threshold: sdk.NewDec(math.MaxInt64),
		Timout:    types.Duration{Seconds: 1},
	}
	bigThresholdAddr, err := k.CreateGroupAccount(ctx, []byte("valid--admin-address"), myGroupID, policy, "test")
	require.NoError(t, err)

	specs := map[string]struct {
		srcAccount   sdk.AccAddress
		srcProposers []sdk.AccAddress
		srcMsgs      []sdk.Msg
		srcComment   string
		expErr       bool
	}{
		"all good with minimal fields set": {
			srcAccount:   accountAddr,
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
		},
		"all good with good msg payload": {
			srcAccount:   accountAddr,
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
			srcMsgs:      []sdk.Msg{&testdata.MyAppProposalPayloadMsgA{}, &testdata.MyAppProposalPayloadMsgB{}},
		},
		"invalid payload should be rejected": {
			srcAccount:   accountAddr,
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
			srcMsgs:      []sdk.Msg{testdata.MyAppProposalPayloadMsgA{}},
			srcComment:   "payload not a pointer",
			expErr:       true,
		},
		"comment too long": {
			srcAccount:   accountAddr,
			srcComment:   strings.Repeat("a", 256),
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
			expErr:       true,
		},
		"group account required": {
			srcComment:   "test",
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
			expErr:       true,
		},
		"existing group account required": {
			srcAccount:   []byte("non-existing-account"),
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
			expErr:       true,
		},
		"impossible case: decision policy threshold > total group weight": {
			srcAccount:   bigThresholdAddr,
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address")},
			expErr:       true,
		},
		"only group members can create a proposal": {
			srcAccount:   accountAddr,
			srcProposers: []sdk.AccAddress{[]byte("non--member-address")},
			expErr:       true,
		},
		"all proposers must be in group": {
			srcAccount:   accountAddr,
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address"), []byte("non--member-address")},
			expErr:       true,
		},
		"proposers must not be nil": {
			srcAccount:   accountAddr,
			srcProposers: []sdk.AccAddress{[]byte("valid-member-address"), nil},
			expErr:       true,
		},
		"admin that is not a group member can not create proposal": {
			srcAccount:   accountAddr,
			srcComment:   "test",
			srcProposers: []sdk.AccAddress{[]byte("valid--admin-address")},
			expErr:       true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			id, err := k.CreateProposal(ctx, spec.srcAccount, spec.srcComment, spec.srcProposers, spec.srcMsgs)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// then all data persisted
			proposal, err := k.GetProposal(ctx, id)
			require.NoError(t, err)

			base := proposal.GetBase()
			assert.Equal(t, accountAddr, base.GroupAccount)
			assert.Equal(t, spec.srcComment, base.Comment)
			assert.Equal(t, spec.srcProposers, base.Proposers)

			submittedAt, err := types.TimestampFromProto(&base.SubmittedAt)
			require.NoError(t, err)
			assert.Equal(t, blockTime.UTC(), submittedAt)

			assert.Equal(t, uint64(1), base.GroupVersion)
			assert.Equal(t, uint64(1), base.GroupAccountVersion)
			assert.Equal(t, group.ProposalBase_Submitted, base.Status)
			assert.Equal(t, group.ProposalBase_Undefined, base.Result)
			assert.Equal(t, group.Tally{
				YesCount:     sdk.ZeroDec(),
				NoCount:      sdk.ZeroDec(),
				AbstainCount: sdk.ZeroDec(),
				VetoCount:    sdk.ZeroDec(),
			}, base.VoteState)

			timout, err := types.TimestampFromProto(&base.Timeout)
			require.NoError(t, err)
			assert.Equal(t, blockTime.Add(time.Second).UTC(), timout)

			if spec.srcMsgs == nil { // then empty list is ok
				assert.Len(t, proposal.GetMsgs(), 0)
			} else {
				assert.Equal(t, spec.srcMsgs, proposal.GetMsgs())
			}
		})
	}
}

func TestVote(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	k := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &testdata.MyAppProposal{})
	blockTime := time.Now().UTC()
	parentCtx := group.NewContext(pKey, pTKey, groupKey).WithBlockTime(blockTime)
	defaultParams := group.DefaultParams()
	paramSpace.SetParamSet(parentCtx, &defaultParams)

	members := []group.Member{
		{Address: []byte("valid-member-address"), Power: sdk.OneDec()},
		{Address: []byte("power-member-address"), Power: sdk.NewDec(2)},
	}
	myGroupID, err := k.CreateGroup(parentCtx, []byte("valid--admin-address"), members, "test")
	require.NoError(t, err)

	policy := group.ThresholdDecisionPolicy{
		Threshold: sdk.NewDec(2),
		Timout:    types.Duration{Seconds: 1},
	}
	accountAddr, err := k.CreateGroupAccount(parentCtx, []byte("valid--admin-address"), myGroupID, policy, "test")
	require.NoError(t, err)
	myProposalID, err := k.CreateProposal(parentCtx, accountAddr, "integration test", []sdk.AccAddress{[]byte("valid-member-address")}, nil)
	require.NoError(t, err)

	specs := map[string]struct {
		srcProposalID     group.ProposalID
		srcVoters         []sdk.AccAddress
		srcChoice         group.Choice
		srcComment        string
		srcCtx            sdk.Context
		doBefore          func(ctx sdk.Context)
		expErr            bool
		expVoteState      group.Tally
		expProposalStatus group.ProposalBase_Status
		expResult         group.ProposalBase_Result
	}{
		"vote yes": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_YES,
			expVoteState: group.Tally{
				YesCount:     sdk.OneDec(),
				NoCount:      sdk.ZeroDec(),
				AbstainCount: sdk.ZeroDec(),
				VetoCount:    sdk.ZeroDec(),
			},
			expProposalStatus: group.ProposalBase_Submitted,
			expResult:         group.ProposalBase_Undefined,
		},
		"vote no": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			expVoteState: group.Tally{
				YesCount:     sdk.ZeroDec(),
				NoCount:      sdk.OneDec(),
				AbstainCount: sdk.ZeroDec(),
				VetoCount:    sdk.ZeroDec(),
			},
			expProposalStatus: group.ProposalBase_Submitted,
			expResult:         group.ProposalBase_Undefined,
		},
		"vote abstain": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_ABSTAIN,
			expVoteState: group.Tally{
				YesCount:     sdk.ZeroDec(),
				NoCount:      sdk.ZeroDec(),
				AbstainCount: sdk.OneDec(),
				VetoCount:    sdk.ZeroDec(),
			},
			expProposalStatus: group.ProposalBase_Submitted,
			expResult:         group.ProposalBase_Undefined,
		},
		"vote veto": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_VETO,
			expVoteState: group.Tally{
				YesCount:     sdk.ZeroDec(),
				NoCount:      sdk.ZeroDec(),
				AbstainCount: sdk.ZeroDec(),
				VetoCount:    sdk.OneDec(),
			},
			expProposalStatus: group.ProposalBase_Submitted,
			expResult:         group.ProposalBase_Undefined,
		},
		"apply decision policy early": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("power-member-address")},
			srcChoice:     group.Choice_YES,
			expVoteState: group.Tally{
				YesCount:     sdk.NewDec(2),
				NoCount:      sdk.ZeroDec(),
				AbstainCount: sdk.ZeroDec(),
				VetoCount:    sdk.ZeroDec(),
			},
			expProposalStatus: group.ProposalBase_Closed,
			expResult:         group.ProposalBase_Accepted,
		},
		"comment too long": {
			srcProposalID: myProposalID,
			srcComment:    strings.Repeat("a", 256),
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			expErr:        true,
		},
		"existing proposal required": {
			srcProposalID: 9999,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			expErr:        true,
		},
		"empty choice": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			expErr:        true,
		},
		"invalid choice": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     5,
			expErr:        true,
		},
		"all voters must be in group": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address"), []byte("non--member-address")},
			srcChoice:     group.Choice_NO,
			expErr:        true,
		},
		"voters must not include nil": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address"), nil},
			srcChoice:     group.Choice_NO,
			expErr:        true,
		},
		"voters must not be nil": {
			srcProposalID: myProposalID,
			srcChoice:     group.Choice_NO,
			expErr:        true,
		},
		"voters must not be empty": {
			srcProposalID: myProposalID,
			srcChoice:     group.Choice_NO,
			srcVoters:     []sdk.AccAddress{},
			expErr:        true,
		},
		"admin that is not a group member can not vote": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid--admin-address")},
			srcChoice:     group.Choice_NO,
			expErr:        true,
		},
		"after timeout": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			srcCtx:        parentCtx.WithBlockTime(blockTime.Add(time.Second)),
			expErr:        true,
		},
		"closed already": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			doBefore: func(ctx sdk.Context) {
				err := k.Vote(ctx, myProposalID, []sdk.AccAddress{[]byte("power-member-address")}, group.Choice_YES, "")
				require.NoError(t, err)
			},
			expErr: true,
		},
		"voted already": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			doBefore: func(ctx sdk.Context) {
				err := k.Vote(ctx, myProposalID, []sdk.AccAddress{[]byte("valid-member-address")}, group.Choice_YES, "")
				require.NoError(t, err)
			},
			expErr: true,
		},
		"with group modified": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			doBefore: func(ctx sdk.Context) {
				g, err := k.GetGroup(ctx, myGroupID)
				require.NoError(t, err)
				g.Comment = "modifed"
				require.NoError(t, k.UpdateGroup(ctx, &g))
			},
			expErr: true,
		},
		"with policy modified": {
			srcProposalID: myProposalID,
			srcVoters:     []sdk.AccAddress{[]byte("valid-member-address")},
			srcChoice:     group.Choice_NO,
			doBefore: func(ctx sdk.Context) {
				a, err := k.GetGroupAccount(ctx, accountAddr)
				require.NoError(t, err)
				a.Base.Comment = "modifed"
				require.NoError(t, k.UpdateGroupAccount(ctx, &a))
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			ctx := parentCtx
			if !spec.srcCtx.IsZero() {
				ctx = spec.srcCtx
			}
			ctx, _ = ctx.CacheContext()

			if spec.doBefore != nil {
				spec.doBefore(ctx)
			}
			err := k.Vote(ctx, spec.srcProposalID, spec.srcVoters, spec.srcChoice, spec.srcComment)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// and all votes are stored
			for _, voter := range spec.srcVoters {
				// then all data persisted
				loaded, err := k.GetVote(ctx, spec.srcProposalID, voter)
				require.NoError(t, err)
				assert.Equal(t, spec.srcProposalID, loaded.Proposal)
				assert.Equal(t, voter, loaded.Voter)
				assert.Equal(t, spec.srcChoice, loaded.Choice)
				assert.Equal(t, spec.srcComment, loaded.Comment)
				submittedAt, err := types.TimestampFromProto(&loaded.SubmittedAt)
				require.NoError(t, err)
				assert.Equal(t, blockTime, submittedAt)
			}

			// and proposal is updated
			proposal, err := k.GetProposal(ctx, spec.srcProposalID)
			require.NoError(t, err)
			assert.Equal(t, spec.expVoteState, proposal.GetBase().VoteState)
			assert.Equal(t, spec.expResult, proposal.GetBase().Result)
			assert.Equal(t, spec.expProposalStatus, proposal.GetBase().Status)
		})
	}
}

func TestLoadParam(t *testing.T) {
	amino := codec.New()
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(amino, pKey, pTKey, group.DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(group.StoreKeyName)
	k := group.NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &group.MockProposalI{})

	ctx := group.NewContext(pKey, pTKey, groupKey)

	myParams := group.Params{MaxCommentLength: 1}
	paramSpace.SetParamSet(ctx, &myParams)

	assert.Equal(t, myParams, k.GetParams(ctx))
}
