package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/group"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

const (
	// ProposalBase Table
	ProposalBaseTablePrefix               byte = 0x30
	ProposalBaseTableSeqPrefix            byte = 0x31
	ProposalBaseByGroupAccountIndexPrefix byte = 0x32
	ProposalBaseByProposerIndexPrefix     byte = 0x33
)

type Keeper struct {
	groupKeeper group.Keeper
	key sdk.StoreKey

	// ProposalBase Table
	myProposalTable           orm.AutoUInt64Table
	ProposalGroupAccountIndex orm.Index
	ProposalByProposerIndex   orm.Index
}

func NewTestdataKeeper(storeKey sdk.StoreKey, groupKeeper group.Keeper) Keeper {
	k := Keeper{
		groupKeeper: groupKeeper,
		key:    storeKey,
	}

	proposalTableBuilder := orm.NewAutoUInt64TableBuilder(ProposalBaseTablePrefix, ProposalBaseTableSeqPrefix, storeKey, &MyAppProposal{})
	k.ProposalGroupAccountIndex = orm.NewIndex(proposalTableBuilder, ProposalBaseByGroupAccountIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		account := value.(*MyAppProposal).Base.GroupAccount
		return []orm.RowID{account.Bytes()}, nil

	})
	k.ProposalByProposerIndex = orm.NewIndex(proposalTableBuilder, ProposalBaseByProposerIndexPrefix, func(value interface{}) ([]orm.RowID, error) {
		proposers := value.(*MyAppProposal).Base.Proposers
		r := make([]orm.RowID, len(proposers))
		for i := range proposers {
			r[i] = proposers[i].Bytes()
		}
		return r, nil
	})
	k.myProposalTable = proposalTableBuilder.Build()
	return k
}

func (k Keeper) CreateProposal(ctx sdk.Context, account sdk.AccAddress, proposers []sdk.AccAddress, comment string) (uint64, error) {
	// todo: validate

	g, err := k.groupKeeper.GetGroupByGroupAccount(ctx, account)
	if err != nil {
		return 0, errors.Wrap(err, "get group by account")
	}
	blockTime, err := types.TimestampProto(ctx.BlockTime())
	if err != nil {
		return 0, errors.Wrap(err, "block time conversion")
	}
	id, err := k.myProposalTable.Create(ctx, &MyAppProposal{
		Base: group.ProposalBase{
			GroupAccount: account,
			Comment:      comment,
			Proposers:    proposers,
			SubmittedAt:  *blockTime,
			GroupVersion: g.Version,
		},
		Msgs: nil,
	})
	if err != nil {
		return 0, errors.Wrap(err, "create proposal")
	}
	return id, nil
}
