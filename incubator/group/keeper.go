package group

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/types"
)

const (
	// Group Table
	GroupTablePrefix        byte = 0x0
	GroupTableSeqPrefix     byte = 0x1
	GroupByAdminIndexPrefix byte = 0x2

	// Group Member Table
	GroupMemberTablePrefix         byte = 0x10
	GroupMemberByGroupIndexPrefix  byte = 0x11
	GroupMemberByMemberIndexPrefix byte = 0x12

	// Group Account Table
	GroupAccountTablePrefix        byte = 0x20
	GroupAccountByGroupIndexPrefix byte = 0x21
	GroupAccountByAdminIndexPrefix byte = 0x22

	// ProposalBase Table
	ProposalBaseTablePrefix               byte = 0x30
	ProposalBaseTableSeqPrefix            byte = 0x31
	ProposalBaseByGroupAccountIndexPrefix byte = 0x32
	ProposalBaseByProposerIndexPrefix     byte = 0x33

	// Vote Table
	VoteTablePrefix               byte = 0x40
	VoteByProposalBaseIndexPrefix byte = 0x41
	VoteByVoterIndexPrefix        byte = 0x42
)

type ProposalI interface {
	orm.Persistent
	GetBase() ProposalBase
	SetBase(ProposalBase)
	GetMsgs() []sdk.Msg
	SetMsgs([]sdk.Msg) error
}

type Keeper struct {
	key               sdk.StoreKey
	proposalModelType reflect.Type

	// Group Table
	groupTable        orm.Table
	groupByAdminIndex orm.Index

	// Group Member Table
	groupMemberTable         orm.NaturalKeyTable
	groupMemberByGroupIndex  orm.UInt64Index
	groupMemberByMemberIndex orm.Index

	// Group Account Table
	groupAccountTable        orm.NaturalKeyTable
	groupAccountByGroupIndex orm.UInt64Index
	groupAccountByAdminIndex orm.Index

	// ProposalBase Table
	proposalTable             orm.AutoUInt64Table
	ProposalGroupAccountIndex orm.Index
	ProposalByProposerIndex   orm.Index

	// Vote Table
	voteTable               orm.NaturalKeyTable
	voteByProposalBaseIndex orm.UInt64Index
	voteByVoterIndex        orm.Index
	groupSeq                orm.Sequence

	paramSpace params.Subspace
	router     sdk.Router
}

func NewGroupKeeper(storeKey sdk.StoreKey, paramSpace params.Subspace, router sdk.Router, proposalModel ProposalI) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(params.NewKeyTable().RegisterParamSet(&Params{}))
	}
	if storeKey == nil {
		panic("storeKey must not be nil")
	}

	if proposalModel == nil {
		panic("proposalModel must not be nil")
	}
	if router == nil {
		panic("router must not be nil")
	}
	tp := reflect.TypeOf(proposalModel)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	k := Keeper{key: storeKey, paramSpace: paramSpace, proposalModelType: tp, router: router}

	//
	// Group Table
	//
	groupTableBuilder := orm.NewTableBuilder(GroupTablePrefix, storeKey, &GroupMetadata{}, orm.FixLengthIndexKeys(orm.EncodedSeqLength))
	k.groupSeq = orm.NewSequence(storeKey, GroupTableSeqPrefix)
	k.groupByAdminIndex = orm.NewIndex(groupTableBuilder, GroupByAdminIndexPrefix, func(val interface{}) ([]orm.RowID, error) {
		return []orm.RowID{val.(*GroupMetadata).Admin.Bytes()}, nil
	})
	k.groupTable = groupTableBuilder.Build()

	//
	// Group Member Table
	//
	groupMemberTableBuilder := orm.NewNaturalKeyTableBuilder(GroupMemberTablePrefix, storeKey, &GroupMember{}, orm.Max255DynamicLengthIndexKeyCodec{})
	k.groupMemberByGroupIndex = orm.NewUInt64Index(groupMemberTableBuilder, GroupMemberByGroupIndexPrefix, func(val interface{}) ([]uint64, error) {
		group := val.(*GroupMember).Group
		return []uint64{uint64(group)}, nil
	})
	k.groupMemberByMemberIndex = orm.NewIndex(groupMemberTableBuilder, GroupMemberByMemberIndexPrefix, func(val interface{}) ([]orm.RowID, error) {
		member := val.(*GroupMember).Member
		return []orm.RowID{member.Bytes()}, nil
	})
	k.groupMemberTable = groupMemberTableBuilder.Build()

	//
	// Group Account Table
	//
	groupAccountTableBuilder := orm.NewNaturalKeyTableBuilder(GroupAccountTablePrefix, storeKey, &StdGroupAccountMetadata{}, orm.Max255DynamicLengthIndexKeyCodec{})
	k.groupAccountByGroupIndex = orm.NewUInt64Index(groupAccountTableBuilder, GroupAccountByGroupIndexPrefix, func(value interface{}) ([]uint64, error) {
		group := value.(*StdGroupAccountMetadata).Base.Group
		return []uint64{uint64(group)}, nil
	})
	k.groupAccountByAdminIndex = orm.NewIndex(groupAccountTableBuilder, GroupAccountByAdminIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		admin := value.(*StdGroupAccountMetadata).Base.Admin
		return []orm.RowID{admin.Bytes()}, nil
	})
	k.groupAccountTable = groupAccountTableBuilder.Build()

	// Proposal Table
	proposalTableBuilder := orm.NewAutoUInt64TableBuilder(ProposalBaseTablePrefix, ProposalBaseTableSeqPrefix, storeKey, proposalModel)
	k.ProposalGroupAccountIndex = orm.NewIndex(proposalTableBuilder, ProposalBaseByGroupAccountIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		account := value.(ProposalI).GetBase().GroupAccount
		return []orm.RowID{account.Bytes()}, nil

	})
	k.ProposalByProposerIndex = orm.NewIndex(proposalTableBuilder, ProposalBaseByProposerIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		proposers := value.(ProposalI).GetBase().Proposers
		r := make([]orm.RowID, len(proposers))
		for i := range proposers {
			r[i] = proposers[i].Bytes()
		}
		return r, nil
	})
	k.proposalTable = proposalTableBuilder.Build()

	//
	// Vote Table
	//
	voteTableBuilder := orm.NewNaturalKeyTableBuilder(VoteTablePrefix, storeKey, &Vote{}, orm.Max255DynamicLengthIndexKeyCodec{})
	k.voteByProposalBaseIndex = orm.NewUInt64Index(voteTableBuilder, VoteByProposalBaseIndexPrefix, func(value interface{}) ([]uint64, error) {
		return []uint64{uint64(value.(*Vote).Proposal)}, nil
	})
	k.voteByVoterIndex = orm.NewIndex(voteTableBuilder, VoteByVoterIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		return []orm.RowID{value.(*Vote).Voter.Bytes()}, nil
	})
	k.voteTable = voteTableBuilder.Build()

	return k
}

// MaxCommentSize returns the maximum length of a comment
func (k Keeper) MaxCommentSize(ctx sdk.Context) int {
	var result uint32
	k.paramSpace.Get(ctx, ParamMaxCommentLength, &result)
	return int(result)
}

func (k Keeper) CreateGroup(ctx sdk.Context, admin sdk.AccAddress, members []Member, comment string) (GroupID, error) {
	// todo: validate
	// deduplicate
	maxCommentSize := k.MaxCommentSize(ctx)
	if len(comment) > maxCommentSize {
		return 0, errors.Wrap(ErrMaxLimit, "group comment")
	}

	totalWeight := sdk.ZeroDec()
	for i := range members {
		m := members[i]
		if len(m.Comment) > maxCommentSize {
			return 0, errors.Wrap(ErrMaxLimit, "group comment")
		}
		totalWeight = totalWeight.Add(m.Power)
	}

	id := k.groupSeq.NextVal(ctx)
	var groupID = GroupID(id)
	err := k.groupTable.Create(ctx, orm.EncodeSequence(id), &GroupMetadata{
		Group:       groupID,
		Admin:       admin,
		Comment:     comment,
		Version:     1,
		TotalWeight: totalWeight,
	})
	if err != nil {
		return 0, errors.Wrap(err, "could not create group")
	}

	for i := range members {
		m := members[i]
		err := k.groupMemberTable.Create(ctx, &GroupMember{
			Group:   groupID,
			Member:  m.Address,
			Weight:  m.Power,
			Comment: m.Comment,
		})
		if err != nil {
			return 0, errors.Wrapf(err, "could not store member %d", i)
		}
	}
	return groupID, nil
}

func (k Keeper) GetGroup(ctx sdk.Context, id GroupID) (GroupMetadata, error) {
	var obj GroupMetadata
	return obj, k.groupTable.GetOne(ctx, id.Byte(), &obj)
}

func (k Keeper) HasGroup(ctx sdk.Context, rowID orm.RowID) bool {
	return k.groupTable.Has(ctx, rowID)
}

func (k Keeper) UpdateGroup(ctx sdk.Context, g *GroupMetadata) error {
	g.Version++
	return k.groupTable.Save(ctx, g.Group.Byte(), g)
}

func (k Keeper) getParams(ctx sdk.Context) Params {
	var p Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

func (k Keeper) setParams(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// CreateGroupAccount creates and persists a `StdGroupAccountMetadata`
//func (k Keeper) CreateGroupAccount(ctx sdk.Context, admin sdk.AccAddress, groupID GroupID, policy DecisionPolicy, comment string) (sdk.AccAddress, error) {
func (k Keeper) CreateGroupAccount(ctx sdk.Context, admin sdk.AccAddress, groupID GroupID, policy ThresholdDecisionPolicy, comment string) (sdk.AccAddress, error) {
	maxCommentSize := k.MaxCommentSize(ctx)
	if len(comment) > maxCommentSize {
		return nil, errors.Wrap(ErrMaxLimit, "group account comment")
	}

	// todo: other validations
	// todo: where to store decision policy?
	//var accountAddr sdk.AccAddress   // todo: how do we generate deterministic address??? as in weave with conditions?

	accountAddr := make([]byte, sdk.AddrLen)
	groupAccount := StdGroupAccountMetadata{
		Base: GroupAccountMetadataBase{
			GroupAccount: accountAddr,
			Group:        groupID,
			Admin:        admin,
			Comment:      comment,
			Version:      1,
		},
		DecisionPolicy: StdDecisionPolicy{Sum: &StdDecisionPolicy_Threshold{&policy}},
	}
	if err := k.groupAccountTable.Create(ctx, &groupAccount); err != nil {
		return nil, errors.Wrap(err, "could not create group account")
	}
	return accountAddr, nil
}

func (k Keeper) HasGroupAccount(ctx sdk.Context, address sdk.AccAddress) bool {
	return k.groupAccountTable.Has(ctx, address.Bytes())
}

func (k Keeper) GetGroupAccount(ctx sdk.Context, accountAddress sdk.AccAddress) (StdGroupAccountMetadata, error) {
	var obj StdGroupAccountMetadata
	return obj, k.groupAccountTable.GetOne(ctx, accountAddress.Bytes(), &obj)
}

func (k Keeper) GetGroupByGroupAccount(ctx sdk.Context, accountAddress sdk.AccAddress) (GroupMetadata, error) {
	obj, err := k.GetGroupAccount(ctx, accountAddress)
	if err != nil {
		return GroupMetadata{}, errors.Wrap(err, "load group account")
	}
	return k.GetGroup(ctx, obj.Base.Group)
}

func (k Keeper) Vote(ctx sdk.Context, id ProposalID, voters []sdk.AccAddress, choice Choice, comment string) error {
	// voters !=0
	// comment within range
	// within voting period
	blockTime, err := types.TimestampProto(ctx.BlockTime())
	if err != nil {
		return err
	}
	proposal, err := k.GetProposal(ctx, id)
	if err != nil {
		return err
	}
	base := proposal.GetBase()
	if base.Status != ProposalBase_Submitted {
		return errors.Wrap(ErrInvalid, "proposal not open")
	}
	votingPeriodEnd, err := types.TimestampFromProto(&base.Timeout)
	if err != nil {
		return err
	}
	if votingPeriodEnd.Before(ctx.BlockTime()) || votingPeriodEnd.Equal(ctx.BlockTime()) {
		return errors.Wrap(ErrExpired, "voting period has ended already")
	}
	var accountMetadata StdGroupAccountMetadata
	if err := k.groupAccountTable.GetOne(ctx, base.GroupAccount.Bytes(), &accountMetadata); err != nil {
		return errors.Wrap(err, "load group account")
	}
	if base.GroupAccountVersion != accountMetadata.Base.Version {
		// todo: this is not the voters fault so we return an error to rollback the TX
		return errors.Wrap(ErrModified, "group account was modified")
	}

	electorate, err := k.GetGroup(ctx, accountMetadata.Base.Group)
	if err != nil {
		return err
	}
	if electorate.Version != base.GroupVersion {
		// todo: this is not the voters fault so we return an error to rollback the TX
		return errors.Wrap(ErrModified, "group was modified")
	}

	// count and store votes
	for _, voterAddr := range voters {
		voter := GroupMember{Group: electorate.Group, Member: voterAddr}
		if err := k.groupMemberTable.GetOne(ctx, voter.NaturalKey(), &voter); err != nil {
			return errors.Wrapf(err, "address: %s", voterAddr)
		}
		newVote := Vote{
			Proposal:    id,
			Voter:       voterAddr,
			Choice:      choice,
			Comment:     comment,
			SubmittedAt: *blockTime,
		}
		if err := base.VoteState.Add(newVote, voter.Weight); err != nil {
			return errors.Wrap(err, "add new vote")
		}

		if err := k.voteTable.Create(ctx, &newVote); err != nil {
			return errors.Wrap(err, "store vote")
		}
	}

	// run tally with new votes to close early

	policy := accountMetadata.DecisionPolicy.GetThreshold()
	submittedAt, err := types.TimestampFromProto(&base.SubmittedAt)
	if err != nil {
		return err
	}
	switch result, err := policy.Allow(base.VoteState, electorate.TotalWeight, ctx.BlockTime().Sub(submittedAt)); {
	case err != nil:
		return errors.Wrap(err, "policy execution")
	case result == DecisionPolicyResult{Allow: true, Final: true}:
		base.Result = ProposalBase_Accepted
		base.Status = ProposalBase_Closed
	case result == DecisionPolicyResult{Allow: false, Final: true}:
		base.Result = ProposalBase_Rejected
		base.Status = ProposalBase_Closed
	}

	proposal.SetBase(base)
	return k.proposalTable.Save(ctx, id.Uint64(), proposal)
}

// ExecProposal can be executed n times before the timeout. It will update the proposal status and executes the msg payload.
// There are no separate transactions for the payload messages so that it is a full atomic operation that
// would either succeed or fail.
func (k Keeper) ExecProposal(ctx sdk.Context, id ProposalID) error {
	proposal, err := k.GetProposal(ctx, id)
	if err != nil {
		return err
	}
	// check constraints
	base := proposal.GetBase()

	if base.Status != ProposalBase_Submitted && base.Status != ProposalBase_Closed {
		return errors.Wrapf(ErrInvalid, "not possible with proposal status %s", base.Status.String())
	}
	votingPeriodEnd, err := types.TimestampFromProto(&base.Timeout)
	if err != nil {
		return err
	}
	if ctx.BlockTime().After(votingPeriodEnd) {
		return errors.Wrap(ErrExpired, "proposal has timed out already")
	}

	var accountMetadata StdGroupAccountMetadata
	if err := k.groupAccountTable.GetOne(ctx, base.GroupAccount.Bytes(), &accountMetadata); err != nil {
		return errors.Wrap(err, "load group account")
	}

	storeUpdates := func() error {
		proposal.SetBase(base)
		return k.proposalTable.Save(ctx, id.Uint64(), proposal)
	}

	if base.GroupAccountVersion != accountMetadata.Base.Version {
		base.Result = ProposalBase_Undefined
		base.Status = ProposalBase_Aborted
		return storeUpdates()
	}

	electorate, err := k.GetGroup(ctx, accountMetadata.Base.Group)
	if err != nil {
		return errors.Wrap(err, "load group")
	}

	if electorate.Version != base.GroupVersion {
		base.Result = ProposalBase_Undefined
		base.Status = ProposalBase_Aborted
		return storeUpdates()
	}

	if base.Status == ProposalBase_Submitted {
		// proposal was not closed early so run decision policy
		policy := accountMetadata.DecisionPolicy.GetThreshold()
		if policy == nil {
			return errors.Wrap(ErrInvalid, "unknown decision policy")
		}

		submittedAt, err := types.TimestampFromProto(&base.SubmittedAt)
		if err != nil {
			return errors.Wrap(err, "from proto time")
		}
		switch result, err := policy.Allow(base.VoteState, electorate.TotalWeight, ctx.BlockTime().Sub(submittedAt)); {
		case err != nil:
			return errors.Wrap(err, "policy execution")
		case result == DecisionPolicyResult{Allow: true, Final: true}:
			base.Result = ProposalBase_Accepted
			base.Status = ProposalBase_Closed
		case result == DecisionPolicyResult{Allow: false, Final: true}:
			base.Result = ProposalBase_Rejected
			base.Status = ProposalBase_Closed
		default:
			// there might be votes coming so we can not close it
		}
	}

	if base.Status == ProposalBase_Closed && base.Result == ProposalBase_Accepted {

		logger := ctx.Logger().With("module", fmt.Sprintf("x/%s", ModuleName))
		proposalType := reflect.TypeOf(proposal).String()

		msgs := proposal.GetMsgs()
		results := make([]sdk.Result, len(msgs))
		for i, msg := range msgs {
			for _, acct := range msg.GetSigners() {
				_ = acct
				if !accountMetadata.Base.GroupAccount.Equals(acct) {
					return errors.Wrap(errors.ErrUnauthorized, "proposal msg does not have permission")
				}
			}

			handler := k.router.Route(ctx, msg.Route())
			if handler == nil {
				logger.Debug("no handler found", "type", proposalType, "proposalID", id, "route", msg.Route(), "pos", i)
				return errors.Wrap(ErrInvalid, "no message handler found")
			}
			r, err := handler(ctx, msg)
			if err != nil {
				return errors.Wrapf(err, "message %q at position %d", msg.Type(), i)
			}
			results[i] = *r
		}
		_ = results // todo: merge results
	}
	return storeUpdates()
}

func (k Keeper) GetProposal(ctx sdk.Context, id ProposalID) (ProposalI, error) {
	loaded := reflect.New(k.proposalModelType).Interface().(ProposalI)
	if _, err := k.proposalTable.GetOne(ctx, id.Uint64(), loaded); err != nil {
		return nil, errors.Wrap(err, "load proposal source")
	}
	return loaded, nil
}

func (k Keeper) CreateProposal(ctx sdk.Context, p ProposalI) (ProposalID, error) {
	id, err := k.proposalTable.Create(ctx, p)
	if err != nil {
		return 0, errors.Wrap(err, "create proposal")
	}
	return ProposalID(id), nil
}

type KeeperDELME interface { // obsolete when Keeper implements all functions
	// Groups
	CreateGroup(ctx orm.HasKVStore, admin sdk.AccAddress, members []Member, comment string) (GroupID, error)
	UpdateGroupMembers(ctx orm.HasKVStore, group GroupID, membersUpdates []Member) error
	UpdateGroupAdmin(ctx orm.HasKVStore, group GroupID, newAdmin sdk.AccAddress) error
	UpdateGroupComment(ctx orm.HasKVStore, group GroupID, newComment string) error

	// Group Accounts
	CreateGroupAccount(ctx orm.HasKVStore, admin sdk.AccAddress, group GroupID, policy DecisionPolicy, comment string) (sdk.AccAddress, error)
	UpdateGroupAccountAdmin(ctx orm.HasKVStore, groupAcc sdk.AccAddress, newAdmin sdk.AccAddress) error
	UpdateGroupAccountDecisionPolicy(ctx orm.HasKVStore, groupAcc sdk.AccAddress, newPolicy DecisionPolicy) error
	UpdateGroupAccountComment(ctx orm.HasKVStore, groupAcc sdk.AccAddress, newComment string) error

	// ProposalBases

	// Propose returns a new ProposalBase ID and a populated sdk.Result which could return an error
	// or the result of execution if execNow was set to true
	Propose(ctx orm.HasKVStore, groupAcc sdk.AccAddress, approvers []sdk.AccAddress, msgs []sdk.Msg, comment string, execNow bool) (id ProposalID, execResult sdk.Result)

	Vote(ctx orm.HasKVStore, id ProposalID, voters []sdk.AccAddress, choice Choice) error

	// Exec attempts to execute the specified ProposalBase. If the ProposalBase is in a valid
	// state and has enough approvals, then it will be executed and its result will be
	// returned, otherwise the result will contain an error
	Exec(ctx orm.HasKVStore, id ProposalID) sdk.Result
}
