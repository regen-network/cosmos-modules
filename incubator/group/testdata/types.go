package testdata

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/group"
)

func (m *AMyAppProposal) SetBase(new group.ProposalBase) {
	m.Base = new
}

func (m AMyAppProposal) SetMsgs([]types.Msg) error {
	panic("implement me")
}
