package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

const (
	GroupTablePrefix        = 0x0
	GroupByAdminIndexPrefix = 0x1
	GroupMemberTablePrefix  = 0x2
)

func NewGroupKeeper(key sdk.StoreKey, cdc *codec.Codec) GroupKeeper {
	k := GroupKeeper{key: key, cdc: cdc}
	groupTableBuilder := NewAutoUInt64TableBuilder(GroupTablePrefix, key, cdc)
	k.groupByAdminIndex = NewIndex(groupTableBuilder, GroupByAdminIndexPrefix, func(val interface{}) []byte {
		return val.(GroupMetadata).Admin
	})
	k.groupTable = groupTableBuilder.Build()
	return k
}
