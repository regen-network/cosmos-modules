package group

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
	cdc.RegisterConcrete(MsgCreateGroup{}, "cosmos-sdk/MsgCreateGroup", nil)
	cdc.RegisterConcrete(MsgUpdateGroupMembers{}, "cosmos-sdk/MsgUpdateGroupMembers", nil)
	cdc.RegisterConcrete(MsgUpdateGroupAdmin{}, "cosmos-sdk/MsgUpdateGroupAdmin", nil)
	cdc.RegisterConcrete(MsgUpdateGroupComment{}, "cosmos-sdk/MsgUpdateGroupComment", nil)
	cdc.RegisterConcrete(MsgCreateGroupAccountStd{}, "cosmos-sdk/MsgCreateGroupAccountStd", nil)
	cdc.RegisterConcrete(MsgVote{}, "cosmos-sdk/group/MsgVote", nil)
	cdc.RegisterConcrete(MsgExec{}, "cosmos-sdk/group/MsgExec", nil)

	// oh man... amino
	cdc.RegisterConcrete(StdDecisionPolicy{}, "cosmos-sdk/StdDecisionPolicy", nil)
	cdc.RegisterConcrete(&StdDecisionPolicy_Threshold{}, "cosmos-sdk/StdDecisionPolicy_Threshold", nil)
	cdc.RegisterConcrete(ThresholdDecisionPolicy{}, "cosmos-sdk/ThresholdDecisionPolicy", nil)
	cdc.RegisterInterface((*isStdDecisionPolicy_Sum)(nil), nil)
}

// generic sealed codec to be used throughout module
var ModuleCdc *Codec

func init() {
	ModuleCdc = NewCodec(codec.New())
	RegisterCodec(ModuleCdc.amino)
	codec.RegisterCrypto(ModuleCdc.amino)
	ModuleCdc.amino.Seal()
}
