package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/group"
)

var _ group.ProposalI = &MyAppProposal{}

func (m *MyAppProposal) GetBase() group.ProposalBase {
	return m.Base
}


func (m *MyAppProposal) SetBase(b group.ProposalBase) {
	m.Base = b
}

func (m *MyAppProposal) GetMsgs() []sdk.Msg {
	r := make([]sdk.Msg, len(m.Msgs))
	for i := range m.Msgs {
		r[i] = m.Msgs[i].GetMsg()
	}
	return r
}

func (m *MyAppProposal) SetMsgs(msgs []sdk.Msg) error {
	r := make([]MyAppMsg, len(m.Msgs))
	for i, msg := range msgs {
		m := MyAppMsg{}
		err := m.SetMsg(msg)
		if err != nil {
			return err
		}
		r[i] = m
	}
	m.Msgs = r
	return nil
}
