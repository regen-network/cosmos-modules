package testdata

import (
	"fmt"

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
	// register all proposals with an executor
	groupKeeper.ExecRouter.Add(&AMyAppProposal{}, ProposalAExecutor(k))
	groupKeeper.ExecRouter.Add(&BMyAppProposal{}, k.ProposalExecutor)
	return k
}

func (k Keeper) CreateProposalA(ctx sdk.Context, accountAddress sdk.AccAddress, proposers []sdk.AccAddress, comment string) (uint64, error) {
	return CreateProposalA(k, ctx, accountAddress, comment, proposers)
}

func CreateProposalA(k Keeper, ctx sdk.Context, accountAddress sdk.AccAddress, comment string, proposers []sdk.AccAddress) (uint64, error) {
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
	window, err := types.DurationFromProto(&policy.GetThreshold().MaxVotingWindow)
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
		Sum: &MyAppProposal_A{A: &AMyAppProposal{
			Base: group.ProposalBase{
				GroupAccount:        accountAddress,
				Comment:             comment,
				Proposers:           proposers,
				SubmittedAt:         *blockTime,
				GroupVersion:        g.Version,
				GroupAccountVersion: account.Base.Version,
				Result:              group.ProposalBase_Undefined,
				Status:              group.ProposalBase_Submitted,
				ExecutorResult:      group.ProposalBase_NotRun,
				VotingEndTime:       *endTime,
			},
		},
		},
	}
	id, err := k.groupKeeper.CreateProposal(ctx, m)
	if err != nil {
		return 0, errors.Wrap(err, "create proposal")
	}
	return id.Uint64(), nil
}

func (k Keeper) CreateProposalB(ctx sdk.Context, accountAddress sdk.AccAddress, proposers []sdk.AccAddress, comment string) (uint64, error) {
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
	window, err := types.DurationFromProto(&policy.GetThreshold().MaxVotingWindow)
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
		Sum: &MyAppProposal_B{B: &BMyAppProposal{
			Base: group.ProposalBase{
				GroupAccount:        accountAddress,
				Comment:             comment,
				Proposers:           proposers,
				SubmittedAt:         *blockTime,
				GroupVersion:        g.Version,
				GroupAccountVersion: account.Base.Version,
				Result:              group.ProposalBase_Undefined,
				Status:              group.ProposalBase_Submitted,
				ExecutorResult:      group.ProposalBase_NotRun,
				VotingEndTime:       *endTime,
			},
		},
		},
	}
	id, err := k.groupKeeper.CreateProposal(ctx, m)
	if err != nil {
		return 0, errors.Wrap(err, "create proposal")
	}
	return id.Uint64(), nil
}

// within a keeper
func (k Keeper) ProposalExecutor(ctx sdk.Context, proposalI group.ProposalI) error {
	logger := ctx.Logger().With("module", fmt.Sprintf("x/%s", ModuleName))

	switch proposalI.(type) {
	case *AMyAppProposal:
		logger.Info("executing AMyAppProposal")
		return nil
	case *BMyAppProposal:
		logger.Info("executing BMyAppProposal")
		return errors.New("exec should fail by intention")
	default:
		return errors.Wrapf(group.ErrType, "%T", proposalI)
	}
}

// or as standalone function
func ProposalAExecutor(k Keeper) func(ctx sdk.Context, proposalI group.ProposalI) error {
	return func(ctx sdk.Context, proposalI group.ProposalI) error {
		logger := ctx.Logger().With("module", fmt.Sprintf("x/%s", ModuleName))
		switch p := proposalI.(type) {
		case *AMyAppProposal:
			_ = p
		default:
			return errors.Wrapf(group.ErrType, "got %T", proposalI)
		}
		logger.Info("executing AMyAppProposal")
		return nil
	}
}
