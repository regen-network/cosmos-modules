package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/modules/incubator/orm"
)

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (s GenesisState) Validate() error {
	return s.Params.Validate()
}

// ExportGenesis returns a GenesisState for a given context and Keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) (*GenesisState, error) {
	groups, _, err := orm.ExportTableData(ctx, k.groupTable)
	if err != nil {
		return nil, errors.Wrap(err, "groups")
	}
	groupMembers, _, err := orm.ExportTableData(ctx, k.groupMemberTable)
	if err != nil {
		return nil, errors.Wrap(err, "group members")
	}
	groupAccounts, _, err := orm.ExportTableData(ctx, k.groupAccountTable)
	if err != nil {
		return nil, errors.Wrap(err, "group accounts")
	}
	proposals, proposalSeq, err := orm.ExportTableData(ctx, k.proposalTable)
	if err != nil {
		return nil, errors.Wrap(err, "proposals")
	}
	votes, _, err := orm.ExportTableData(ctx, k.voteTable)
	if err != nil {
		return nil, errors.Wrap(err, "proposals")
	}
	return &GenesisState{
		Params:          k.GetParams(ctx),
		Groups:          groups,
		GroupSeq:        k.groupSeq.CurVal(ctx),
		GroupMembers:    groupMembers,
		GroupAccounts:   groupAccounts,
		GroupAccountSeq: k.groupAccountSeq.CurVal(ctx),
		ProposalSeq:     proposalSeq,
		Proposals:       proposals,
		Votes:           votes,
	}, nil
}

// ImportGenesis sets state for a given context and Keeper.
func ImportGenesis(ctx sdk.Context, k Keeper, g GenesisState) error {
	if err := g.Validate(); err != nil {
		return err
	}
	k.setParams(ctx, g.Params)
	if err := orm.ImportTableData(ctx, k.groupTable, g.Groups, 0); err != nil {
		return errors.Wrap(err, "groups")
	}
	if err := k.groupSeq.InitVal(ctx, g.GroupSeq); err != nil {
		return errors.Wrap(err, "group seq")
	}
	if err := orm.ImportTableData(ctx, k.groupMemberTable, g.GroupMembers, 0); err != nil {
		return errors.Wrap(err, "group members")
	}
	if err := orm.ImportTableData(ctx, k.groupAccountTable, g.GroupAccounts, 0); err != nil {
		return errors.Wrap(err, "group accounts")
	}
	if err := k.groupAccountSeq.InitVal(ctx, g.GroupAccountSeq); err != nil {
		return errors.Wrap(err, "group account seq")
	}
	if err := orm.ImportTableData(ctx, k.proposalTable, g.Proposals, g.ProposalSeq); err != nil {
		return errors.Wrap(err, "proposals")
	}
	if err := orm.ImportTableData(ctx, k.voteTable, g.Votes, 0); err != nil {
		return errors.Wrap(err, "votes")
	}
	return nil
}
