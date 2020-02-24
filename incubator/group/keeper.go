package group

import (
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
	GetBase() ProposalBase
	SetBase(ProposalBase)
}

type ProposalModel interface {
	orm.Persistent
	GetProposalI() ProposalI
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
}

func NewGroupKeeper(storeKey sdk.StoreKey, paramSpace params.Subspace, proposalModel ProposalModel) Keeper {
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
	tp := reflect.TypeOf(proposalModel)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	k := Keeper{key: storeKey, paramSpace: paramSpace, proposalModelType: tp}

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
		account := value.(ProposalModel).GetProposalI().GetBase().GroupAccount
		return []orm.RowID{account.Bytes()}, nil

	})
	k.ProposalByProposerIndex = orm.NewIndex(proposalTableBuilder, ProposalBaseByProposerIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		proposers := value.(ProposalModel).GetProposalI().GetBase().Proposers
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
	proposalModel, err := k.getProposalModel(ctx, id)
	if err != nil {
		return err
	}
	proposal := proposalModel.GetProposalI().GetBase()
	if proposal.Status != ProposalBase_Submitted {
		return errors.Wrap(ErrInvalid, "proposal not open")
	}
	votingPeriodEnd, err := types.TimestampFromProto(&proposal.VotingEndTime)
	if err != nil {
		return err
	}
	end := votingPeriodEnd.UTC()
	if end.Before(ctx.BlockTime()) || end.Equal(ctx.BlockTime()) {
		return errors.Wrap(ErrExpired, "voting period has ended already")
	}
	electorate, err := k.GetGroupByGroupAccount(ctx, proposal.GroupAccount)
	if err != nil {
		return err
	}
	if electorate.Version != proposal.GroupVersion {
		// todo: this is not the voters fault.
		return errors.Wrap(ErrModified, "group was modified")
	}
	for _, v := range voters {
		newVote := Vote{
			Proposal:    id,
			Voter:       v,
			Choice:      choice,
			Comment:     comment,
			SubmittedAt: *blockTime,
		}
		member := GroupMember{
			Group:  electorate.Group,
			Member: v,
		}

		if err := k.groupMemberTable.GetOne(ctx, member.NaturalKey(), &member); err != nil {
			return errors.Wrapf(err, "get member by group and address")
		}

		var oldVote *Vote
		err = k.voteTable.GetOne(ctx, newVote.NaturalKey(), oldVote)
		switch {
		case orm.ErrNotFound.Is(err):
		case err != nil:
			return errors.Wrap(err, "load old vote")
		default:
			if err := proposal.VoteState.Sub(*oldVote, member.Weight); err != nil {
				return errors.Wrap(err, "previous vote")
			}
		}
		if err := proposal.VoteState.Add(newVote, member.Weight); err != nil {
			return errors.Wrap(err, "new vote")
		}

		// todo: here a put would be nicer again
		if oldVote == nil {
			if err := k.voteTable.Create(ctx, &newVote); err != nil {
				return errors.Wrap(err, "create vote")
			}
		} else {
			if err := k.voteTable.Save(ctx, &newVote); err != nil {
				return errors.Wrap(err, "update vote")
			}
		}
	}
	proposalModel.GetProposalI().SetBase(proposal)
	return k.proposalTable.Save(ctx, id.Uint64(), proposalModel)
}

func (k Keeper) ExecProposal(ctx sdk.Context, id ProposalID) error {
	proposalModel, err := k.getProposalModel(ctx, id)
	if err != nil {
		return err
	}
	proposal := proposalModel.GetProposalI().GetBase()

	if proposal.Status != ProposalBase_Submitted {
		return errors.Wrap(ErrInvalid, "proposal not open")
	}

	var accountMetadata StdGroupAccountMetadata
	if err := k.groupAccountTable.GetOne(ctx, proposal.GroupAccount.Bytes(), &accountMetadata); err != nil {
		return errors.Wrap(err, "load group account")
	}

	electorate, err := k.GetGroupByGroupAccount(ctx, proposal.GroupAccount)
	if err != nil {
		return err
	}

	if electorate.Version != proposal.GroupVersion {
		proposal.Status = ProposalBase_Aborted
		return nil
		// todo: or error?
		//return errors.Wrap(ErrModified, "group was modified")
	}
	// todo: validate
	height := ctx.BlockHeight()
	_ = height

	votingPeriodEnd, err := types.TimestampFromProto(&proposal.VotingEndTime)
	if err != nil {
		return err
	}
	end := votingPeriodEnd.UTC()
	block := ctx.BlockTime()
	if block.Before(end) {
		return errors.Wrap(ErrExpired, "voting period not ended yet")
	}

	proposal.Status = ProposalBase_Closed

	// run decision policy
	policy := accountMetadata.DecisionPolicy.GetThreshold()
	if policy == nil {
		return errors.Wrap(ErrInvalid, "unknown decision policy")
	}

	submittedAt, err := types.TimestampFromProto(&proposal.SubmittedAt)
	if err != nil {
		return errors.Wrap(err, "from proto time")
	}
	switch accepted, err := policy.Allow(proposal.VoteState, electorate.TotalWeight, ctx.BlockTime().Sub(submittedAt)); {
	case err != nil:
		return errors.Wrap(err, "policy execution")
	case accepted:
		proposal.Result = ProposalBase_Accepted
	default:
		proposal.Result = ProposalBase_Rejected
	}
	proposalModel.GetProposalI().SetBase(proposal)
	return k.proposalTable.Save(ctx, id.Uint64(), proposalModel)
}

func (k Keeper) getProposalModel(ctx sdk.Context, id ProposalID) (ProposalModel, error) {
	loaded := reflect.New(k.proposalModelType).Interface().(ProposalModel)
	if _, err := k.proposalTable.GetOne(ctx, id.Uint64(), loaded); err != nil {
		return nil, errors.Wrap(err, "load proposal source")
	}
	return loaded, nil
}
func (k Keeper) GetProposal(ctx sdk.Context, id ProposalID) (ProposalI, error) {
	s, err := k.getProposalModel(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.GetProposalI(), nil
}

func (k Keeper) CreateProposal(ctx sdk.Context, p ProposalModel) (ProposalID, error) {
	id, err := k.proposalTable.Create(ctx, p)
	if err != nil {
		return 0, errors.Wrap(err, "create proposal")
	}
	return ProposalID(id), nil
}

//
//func (k Keeper) UpdateGroupAccountAdmin(ctx orm.HasKVStore, groupAcc sdk.AccAddress, newAdmin sdk.AccAddress) error {
//	panic("implement me")
//}
//
//func (k Keeper) UpdateGroupAccountDecisionPolicy(ctx orm.HasKVStore, groupAcc sdk.AccAddress, newPolicy DecisionPolicy) error {
//	panic("implement me")
//}
//
//func (k Keeper) UpdateGroupAccountComment(ctx orm.HasKVStore, groupAcc sdk.AccAddress, newComment string) error {
//	panic("implement me")
//}
//
//func (k Keeper) Propose(ctx orm.HasKVStore, groupAcc sdk.AccAddress, approvers []sdk.AccAddress, msgs []sdk.Msg, comment string, execNow bool) (id ProposalID, execResult sdk.Result) {
//	panic("implement me")
//}
//
//func (k Keeper) Vote(ctx orm.HasKVStore, id ProposalID, voters []sdk.AccAddress, choice Choice) error {
//	panic("implement me")
//}
//
//func (k Keeper) Exec(ctx orm.HasKVStore, id ProposalID) sdk.Result {
//	panic("implement me")
//}

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
