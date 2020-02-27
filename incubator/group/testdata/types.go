package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/group"
)

func (m *MyAppProposal) SetBase(b group.ProposalBase) {
	m.Base = b
}

func (m *MyAppProposal) GetMsg() []sdk.Msg {
	r := make([]sdk.Msg, len(m.Msgs))
	for i := range m.Msgs {
		r[i] = m.Msgs[i].GetMsg()
	}
	return r
}