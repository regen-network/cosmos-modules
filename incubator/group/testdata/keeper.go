package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/group"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
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

func (k Keeper) CreateProposalA(ctx sdk.Context, accountAddress sdk.AccAddress, proposers []sdk.AccAddress, comment string, msgs []MyAppMsg) (uint64, error) {
	return CreateProposal(k, ctx, accountAddress, comment, proposers, msgs)
}

func CreateProposal(k Keeper, ctx sdk.Context, accountAddress sdk.AccAddress, comment string, proposers []sdk.AccAddress, msgs []MyAppMsg) (uint64, error) {
	// todo: validate
	account, err := k.groupKeeper.GetGroupAccount(ctx, accountAddress.Bytes())
	if err != nil {
		return 0, errors.Wrap(err, "load group account")
	}

	g, err := k.groupKeeper.GetGroupByGroupAccount(ctx, accountAddress)
	if err != nil {
		return 0, errors.Wrap(err, "get group by account")
	}
	blockTime, err := types.TimestampProto(ctx.BlockTime())
	if err != nil {
		return 0, errors.Wrap(err, "block time conversion")
	}
	policy := account.GetDecisionPolicy()
	window, err := types.DurationFromProto(&policy.GetThreshold().Timout)
	if err != nil {
		return 0, errors.Wrap(err, "maxVotingWindow time conversion")
	}
	endTime, err := types.TimestampProto(ctx.BlockTime().Add(window))
	if err != nil {
		return 0, errors.Wrap(err, "end time conversion")
	}
	block := ctx.BlockHeight()
	_ = block
	m := &MyAppProposal{
		Base: group.ProposalBase{
			GroupAccount:        accountAddress,
			Comment:             comment,
			Proposers:           proposers,
			SubmittedAt:         *blockTime,
			GroupVersion:        g.Version,
			GroupAccountVersion: account.Base.Version,
			Result:              group.ProposalBase_Undefined,
			Status:              group.ProposalBase_Submitted,
			Timeout:             *endTime,
		},
		Msgs:msgs,
	}
	id, err := k.groupKeeper.CreateProposal(ctx, m)
	if err != nil {
		return 0, errors.Wrap(err, "create proposal")
	}
	return id.Uint64(), nil
}