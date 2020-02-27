package testdata

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

type Codec struct {
	codec.Marshaler

	// Keep reference to the amino codec to allow backwards compatibility along
	// with type, and interface registration.
	amino *codec.Codec
}

func NewCodec(amino *codec.Codec) *Codec {
	return &Codec{Marshaler: codec.NewHybridCodec(amino), amino: amino}
}

// ----------------------------------------------------------------------------

// RegisterCodec registers all the necessary crisis module concrete types and
// interfaces with the provided codec reference.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgPropose{}, "testdata/MsgPropose", nil)
	// oh man... amino
	cdc.RegisterInterface((*isMyAppMsg_Sum)(nil), nil)
	cdc.RegisterConcrete(&MyAppProposalPayloadMsgA{}, "testdata/MyAppProposalPayloadMsgA", nil)
	cdc.RegisterConcrete(&MyAppProposalPayloadMsgB{}, "testdata/MyAppProposalPayloadMsgB", nil)
	cdc.RegisterConcrete(&MyAppMsg_A{}, "testdata/MyAppMsg_A", nil)
	cdc.RegisterConcrete(&MyAppMsg_B{}, "testdata/MyAppMsg_B", nil)
}

// generic sealed codec to be used throughout module
var ModuleCdc *Codec

func init() {
	ModuleCdc = NewCodec(codec.New())
	RegisterCodec(ModuleCdc.amino)
	codec.RegisterCrypto(ModuleCdc.amino)
	ModuleCdc.amino.Seal()
}
