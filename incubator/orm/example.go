package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dvsekhvalnov/jose2go/arrays"
)

type GroupKeeper struct {
	key                      sdk.StoreKey
	cdc                      *codec.Codec
	groupTable               AutoUInt64Table
	groupByAdminIndex        Index
	groupMemberTable         NaturalKeyTable
	groupMemberByGroupIndex  Index
	groupMemberByMemberIndex Index
	groupAccountTable        AutoKeyTable
	groupAccountByGroupIndex Index
	groupAccountByAdminIndex Index
	proposalTable            AutoUInt64Table
	proposalByGroupIndex     Index
	voteTable                NaturalKeyTable
	voteByProposalIndex      UInt64Index
	voteByVoterIndex         Index
}

type GroupMetadata struct {
	Description string
	Admin       sdk.AccAddress
}

type GroupMember struct {
	Group  sdk.AccAddress
	Member sdk.AccAddress
	Weight sdk.Int
}

var	(
	GroupTablePrefix               = []byte{0x0}
	GroupByAdminIndexPrefix        = []byte{0x1}
	GroupMemberTablePrefix         = []byte{0x2}
	GroupMemberByGroupIndexPrefix  = []byte{0x3}
	GroupMemberByMemberIndexPrefix = []byte{0x3}
)

func NewGroupKeeper(key sdk.StoreKey, cdc *codec.Codec) GroupKeeper {
	k := GroupKeeper{key: key, cdc: cdc}

	groupTableBuilder := NewAutoUInt64TableBuilder(GroupTablePrefix, key, cdc)
	k.groupByAdminIndex = NewIndex(groupTableBuilder, GroupByAdminIndexPrefix, func(val interface{}) []byte {
		return val.(GroupMetadata).Admin
	})
	k.groupTable = groupTableBuilder.Build()

	groupMemberTableBuilder := NewNaturalKeyTableBuilder(GroupMemberTablePrefix, key, cdc, func(val interface{}) []byte {
		gm := val.(GroupMember)
		return arrays.Concat(gm.Group, gm.Member)
	})
	k.groupMemberByGroupIndex = NewIndex(groupMemberTableBuilder, GroupMemberByGroupIndexPrefix, func(val interface{}) []byte {
		return val.(GroupMember).Group
	})
	k.groupMemberByMemberIndex = NewIndex(groupMemberTableBuilder, GroupMemberByMemberIndexPrefix, func(val interface{}) []byte {
		return val.(GroupMember).Member
	})
	k.groupMemberTable = groupMemberTableBuilder.Build()

	return k
}

func NewGroupKeeper2(mgr SchemaManager) GroupKeeper {
	k := GroupKeeper{}

	groupTableBuilder := NewAutoUInt64TableBuilder(mgr, "group", func setId(model interface{}, id uint64) {
		model.ID = id
	})
	k.groupByAdminIndex = NewIndex(groupTableBuilder, func(val interface{}) []byte {
		return val.(GroupMetadata).Admin
	})
	k.groupTable = groupTableBuilder.Build()

	groupMemberTableBuilder := NewNaturalKeyTableBuilder(GroupMemberTablePrefix, key, cdc, func(val interface{}) []byte {
		gm := val.(GroupMember)
		return arrays.Concat(gm.Group, gm.Member)
	})
	k.groupMemberByGroupIndex = NewIndex(groupMemberTableBuilder, GroupMemberByGroupIndexPrefix, func(val interface{}) []byte {
		return val.(GroupMember).Group
	})
	k.groupMemberByMemberIndex = NewIndex(groupMemberTableBuilder, GroupMemberByMemberIndexPrefix, func(val interface{}) []byte {
		return val.(GroupMember).Member
	})
	k.groupMemberTable = groupMemberTableBuilder.Build()

	return k
}
