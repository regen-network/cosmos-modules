package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/group"
)

type Keeper struct {
	groupKeeper group.Keeper
	key         sdk.StoreKey
}

func NewTestdataKeeper(storeKey sdk.StoreKey, groupKeeper group.Keeper) Keeper {
	k := Keeper{
		groupKeeper: groupKeeper,
		key:         storeKey,
	}
	return k
}

func (k Keeper) CreateProposal(ctx sdk.Context, accountAddress sdk.AccAddress, proposers []sdk.AccAddress, comment string, msgs MyAppMsgs) (group.ProposalID, error) {
	return k.groupKeeper.CreateProposal(ctx, accountAddress, comment, proposers, msgs.AsSDKMsgs())
}
