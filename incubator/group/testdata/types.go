package testdata

import (
	"github.com/cosmos/modules/incubator/group"
)

func (m *AMyAppProposal) SetBase(new group.ProposalBase) {
	m.Base = new
}

func (m *BMyAppProposal) SetBase(new group.ProposalBase) {
	m.Base = new
}

