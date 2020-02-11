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

type GroupID uint64

type ProposalID uint64

func (g GroupMember) NaturalKey() []byte {
	result := make([]byte, 0, binary.MaxVarintLen64+len(g.Member))
	binary.PutUvarint(result, g.Group)
	result = append(result, g.Member...)
	return result
}

func (g GroupAccountMetadata) NaturalKey() []byte {
	return g.GroupAccount
}

func (v Vote) NaturalKey() []byte {
	result := make([]byte, 0, binary.MaxVarintLen64+len(v.Voter))
	binary.PutUvarint(result, v.Proposal)
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
		return []uint64{group}, nil
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
		return []uint64{group}, nil
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
	k.proposalTable = proposalTableBuilder.Build()

	//
	// Vote Table
	//
	voteTableBuilder := orm.NewNaturalKeyTableBuilder(VoteTablePrefix, VoteTableSeqPrefix, VoteTableIndexPrefix, storeKey, &Vote{})
	k.voteTable = voteTableBuilder.Build()

	return k
}

type Keeper interface {
	CreateGroup(ctx sdk.Context, members []Member, admin sdk.AccAddress, comment string) (GroupID, error)
	CreateGroupAccount(ctx sdk.Context, group GroupID, policy DecisionPolicy, admin sdk.AccAddress, comment string) (sdk.AccAddress, error)
}
