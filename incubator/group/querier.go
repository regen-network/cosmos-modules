package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/orm"
)

func NewQuerier(k Keeper) sdk.Querier {
	qh := orm.NewQueryHandler()
	qh.AddTableRoute("xgroup", k.groupTable)
	qh.AddIndexRoute("xgroup/admin", k.groupByAdminIndex)
	return qh.Handle
}
