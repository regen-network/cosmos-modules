package testdata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	fuzz "github.com/google/gofuzz"
)

func FuzzPayloadMsg(m *sdk.Msg, c fuzz.Continue) {
	switch c.Intn(5) {
	case 0:
		*m = &MsgAlwaysSucceed{}
	case 1:
		*m = &MsgAlwaysFail{}
	case 2:
		*m = &MsgSetValue{}
	case 3:
		*m = &MsgIncCounter{}
	case 4:
		*m = &MsgIncCounter{}
		c.Fuzz(m)
	case 5:
		*m = &MsgAuthenticate{}
		c.Fuzz(m)
	}
}
