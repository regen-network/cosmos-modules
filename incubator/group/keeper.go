package group

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/orm"
)

type keeper struct {
	key sdk.StoreKey

	// Group Table
	groupTable        orm.AutoUInt64Table
	groupByAdminIndex orm.Index

	// Group Member Table
	groupMemberTable         orm.NaturalKeyTable
	groupMemberByGroupIndex  *orm.UInt64Index
	groupMemberByMemberIndex orm.Index

	// Group Account Table
	groupAccountTable        orm.NaturalKeyTable
	groupAccountByGroupIndex *orm.UInt64Index
	groupAccountByAdminIndex orm.Index

	// Proposal Table
	proposalTable               orm.AutoUInt64Table
	proposalByGroupAccountIndex orm.Index
	proposalByProposerIndex     orm.Index

	// Vote Table
	voteTable           orm.NaturalKeyTable
	voteByProposalIndex *orm.UInt64Index
	voteByVoterIndex    orm.Index
}

func (g GroupMember) NaturalKey() []byte {
	result := make([]byte, 0, binary.MaxVarintLen64+len(g.Member))
	binary.PutUvarint(result, uint64(g.Group))
	result = append(result, g.Member...)
	return result
}

func (g GroupAccountMetadata) NaturalKey() []byte {
	return g.GroupAccount
}

func (v Vote) NaturalKey() []byte {
	result := make([]byte, 0, binary.MaxVarintLen64+len(v.Voter))
	binary.PutUvarint(result, uint64(v.Proposal))
	result = append(result, v.Voter...)
	return result
}

var (
	// Group Table
	GroupTablePrefix               byte = 0x0
	GroupTableSeqPrefix            byte = 0x1
	GroupByAdminIndexPrefix        byte = 0x2

	// Group Member Table
	GroupMemberTablePrefix         byte = 0x3
	GroupMemberTableSeqPrefix      byte = 0x4
	GroupMemberTableIndexPrefix    byte = 0x5
	GroupMemberByGroupIndexPrefix  byte = 0x6
	GroupMemberByMemberIndexPrefix byte = 0x7

	// Group Account Table
	GroupAccountTablePrefix      byte = 0x8
	GroupAccountTableSeqPrefix   byte = 0x9
	GroupAccountTableIndexPrefix byte = 0x10
	GroupAccountByGroupIndexPrefix byte = 0x11
	GroupAccountByAdminIndexPrefix byte = 0x12

	// Proposal Table
	ProposalTablePrefix byte = 0x13
	ProposalTableSeqPrefix byte = 0x14
	ProposalByGroupAccountIndexPrefix byte = 0x15
	ProposalByProposerIndexPrefix byte = 0x16

	// Vote Table
	VoteTablePrefix byte = 0x17
	VoteTableSeqPrefix byte = 0x18
	VoteTableIndexPrefix byte = 0x19
	VoteByProposalIndexPrefix byte = 0x20
	VoteByVoterIndexPrefix byte = 0x21
)

func NewGroupKeeper(storeKey sdk.StoreKey) keeper {
	k := keeper{key: storeKey}

	//
	// Group Table
	//
	groupTableBuilder := orm.NewAutoUInt64TableBuilder(GroupTablePrefix, GroupTableSeqPrefix, storeKey, &GroupMetadata{})
	k.groupByAdminIndex = orm.NewIndex(groupTableBuilder, GroupByAdminIndexPrefix, func(val interface{}) ([][]byte, error) {
		return [][]byte{val.(*GroupMetadata).Admin}, nil
	})
	k.groupTable = groupTableBuilder.Build()

	//
	// Group Member Table
	//
	groupMemberTableBuilder := orm.NewNaturalKeyTableBuilder(GroupMemberTablePrefix, GroupMemberTableSeqPrefix, GroupMemberTableIndexPrefix, storeKey, &GroupMember{})
	k.groupMemberByGroupIndex = orm.NewUInt64Index(groupMemberTableBuilder, GroupMemberByGroupIndexPrefix, func(val interface{}) ([]uint64, error) {
		group := val.(*GroupMember).Group
		return []uint64{uint64(group)}, nil
	})
	k.groupMemberByMemberIndex = orm.NewIndex(groupMemberTableBuilder, GroupMemberByMemberIndexPrefix, func(val interface{}) ([][]byte, error) {
		return [][]byte{val.(*GroupMember).Member}, nil
	})
	k.groupMemberTable = groupMemberTableBuilder.Build()

	//
	// Group Account Table
	//
	groupAccountTableBuilder := orm.NewNaturalKeyTableBuilder(GroupAccountTablePrefix, GroupAccountTableSeqPrefix, GroupAccountTableIndexPrefix, storeKey, &GroupAccountMetadata{})
	k.groupAccountByGroupIndex = orm.NewUInt64Index(groupAccountTableBuilder, GroupAccountByGroupIndexPrefix, func(value interface{}) ([]uint64, error) {
		group := value.(*GroupAccountMetadata).Group
		return []uint64{uint64(group)}, nil
	})
	k.groupAccountByAdminIndex = orm.NewIndex(groupAccountTableBuilder, GroupAccountByAdminIndexPrefix, func(value interface{}) ([][]byte, error) {
		admin := value.(*GroupAccountMetadata).Admin
		return [][]byte{admin}, nil
	})
	k.groupAccountTable = groupAccountTableBuilder.Build()

	//
	// Proposal Table
	//
	proposalTableBuilder := orm.NewAutoUInt64TableBuilder(ProposalTablePrefix, ProposalTableSeqPrefix, storeKey, &Proposal{})
	k.proposalByGroupAccountIndex = orm.NewIndex(proposalTableBuilder, ProposalByGroupAccountIndexPrefix, func(value interface{}) ([][]byte, error) {
		return [][]byte{value.(*Proposal).GroupAccount}, nil

	})
	k.proposalByProposerIndex = orm.NewIndex(proposalTableBuilder, ProposalByProposerIndexPrefix, func(value interface{}) ([][]byte, error) {
		return value.(*Proposal).Proposers, nil
	})
	k.proposalTable = proposalTableBuilder.Build()

	//
	// Vote Table
	//
	voteTableBuilder := orm.NewNaturalKeyTableBuilder(VoteTablePrefix, VoteTableSeqPrefix, VoteTableIndexPrefix, storeKey, &Vote{})
	k.voteByProposalIndex = orm.NewUInt64Index(voteTableBuilder, VoteByProposalIndexPrefix, func(value interface{}) ([]uint64, error) {
		return []uint64{uint64(value.(*Vote).Proposal)}, nil
	})
	k.voteByVoterIndex = orm.NewIndex(voteTableBuilder, VoteByVoterIndexPrefix, func(value interface{}) ([][]byte, error) {
		return [][]byte{value.(*Vote).Voter}, nil
	})
	k.voteTable = voteTableBuilder.Build()

	return k
}

type Keeper interface {
	// Groups
	CreateGroup(ctx sdk.Context, admin sdk.AccAddress, members []Member, comment string) (GroupID, error)
	UpdateGroupMembers(ctx sdk.Context, group GroupID, membersUpdates []Member) error
	UpdateGroupAdmin(ctx sdk.Context, group GroupID, newAdmin sdk.AccAddress) error
	UpdateGroupComment(ctx sdk.Context, group GroupID, newComment string) error

	// Group Accounts
	CreateGroupAccount(ctx sdk.Context, admin sdk.AccAddress, group GroupID, policy DecisionPolicy, comment string) (sdk.AccAddress, error)
	UpdateGroupAccountAdmin(ctx sdk.Context, groupAcc sdk.AccAddress, newAdmin sdk.AccAddress) error
	UpdateGroupAccountDecisionPolicy(ctx sdk.Context, groupAcc sdk.AccAddress, newPolicy DecisionPolicy) error
	UpdateGroupAccountComment(ctx sdk.Context, groupAcc sdk.AccAddress, newComment string) error

	// Proposals

	// Propose returns a new proposal ID and a populated sdk.Result which could return an error
	// or the result of execution if execNow was set to true
	Propose(ctx sdk.Context, groupAcc sdk.AccAddress, approvers []sdk.AccAddress, msgs []sdk.Msg, comment string, execNow bool) (id ProposalID, execResult sdk.Result)

	Vote(ctx sdk.Context, id ProposalID, voters []sdk.AccAddress, choice Choice) error

	// Exec attempts to execute the specified proposal. If the proposal is in a valid
	// state and has enough approvals, then it will be executed and its result will be
	// returned, otherwise the result will contain an error
	Exec(ctx sdk.Context, id ProposalID) sdk.Result
}
