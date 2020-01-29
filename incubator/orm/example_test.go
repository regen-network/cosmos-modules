package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GroupKeeper interface {
	CreateGroup(ctx sdk.Context, admin sdk.AccAddress, members []Member, description string) (GroupID, error)
}

type GroupID uint64

type groupKeeper struct {
	key                      sdk.StoreKey
	cdc                      *codec.Codec
	groupTable               AutoUInt64Table
	groupByAdminIndex        Index
	groupMemberTable         NaturalKeyTable
	groupMemberByGroupIndex  Index
	groupMemberByMemberIndex Index
}

var _ GroupKeeper = groupKeeper{}

type GroupTable struct {
	Description string
	Admin       sdk.AccAddress
}

type GroupMemberTable struct {
	Group       GroupID
	Member      sdk.AccAddress
	Weight      sdk.Int
	Description string
}

func (g GroupMemberTable) NaturalKey() []byte {
	result := make([]byte, 0, len(g.Group)+len(g.Member))
	result = append(result, g.Group...)
	result = append(result, g.Member...)
	return result
}

var (
	GroupTablePrefix               byte = 0x0
	GroupTableSeqPrefix            byte = 0x1
	GroupByAdminIndexPrefix        byte = 0x2
	GroupMemberTableTablePrefix         byte = 0x3
	GroupMemberTableTableSeqPrefix      byte = 0x4
	GroupMemberTableTableIndexPrefix    byte = 0x5
	GroupMemberTableByGroupIndexPrefix  byte = 0x6
	GroupMemberTableByMemberIndexPrefix byte = 0x7
)

func NewGroupKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) groupKeeper {
	k := groupKeeper{key: storeKey, cdc: cdc}

	groupTableBuilder := NewAutoUInt64TableBuilder(GroupTablePrefix, GroupTableSeqPrefix, storeKey, cdc, &GroupTable{})
	// note: quite easy to mess with Index prefixes when managed outside. no fail fast on duplicates
	k.groupByAdminIndex = NewIndex(groupTableBuilder, GroupByAdminIndexPrefix, func(val interface{}) ([][]byte, error) {
		return [][]byte{val.(*GroupTable).Admin}, nil
	})
	k.groupTable = groupTableBuilder.Build()

	groupMemberTableBuilder := NewNaturalKeyTableBuilder(GroupMemberTableTablePrefix, GroupMemberTableTableSeqPrefix, GroupMemberTableTableIndexPrefix, storeKey, cdc, &GroupMemberTable{})

	k.groupMemberByGroupIndex = NewIndex(groupMemberTableBuilder, GroupMemberTableByGroupIndexPrefix, func(val interface{}) ([][]byte, error) {
		group := val.(*GroupMemberTable).Group
		return [][]byte{group}, nil
	})
	k.groupMemberByMemberIndex = NewIndex(groupMemberTableBuilder, GroupMemberTableByMemberIndexPrefix, func(val interface{}) ([][]byte, error) {
		return [][]byte{val.(*GroupMemberTable).Member}, nil
	})
	k.groupMemberTable = groupMemberTableBuilder.Build()

	return k
}

type Member struct {
	// The address of a group member. Can be another group or a contract
	Address sdk.AccAddress `json:"address"`
	// The integral weight of this member with respect to other members
	Weight      sdk.Int `json:"weight"`
	Description string  `json:"description"`
}

func (g groupKeeper) CreateGroup(ctx sdk.Context, admin sdk.AccAddress, members []Member, description string) (groupId GroupID, err error) {
	id, err := g.groupTable.Create(ctx, GroupTable{
		Description: description,
		Admin:       admin,
	})
	if err != nil {
		return groupId, err
	}
	groupId = GroupID(id)
	for _, member := range members {
		g.groupMemberTable.Create(ctx, GroupMemberTable{
			Group:       groupId,
			Member:      member.Address,
			Weight:      member.Weight,
			Description: member.Description,
		})
	}
}
