package testdata

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"
)

const (
	msgTypeCreateProposalA = "create_proposal_a"
	msgTypeCreateProposalB = "create_proposal_B"
)

var _ sdk.Msg = &MsgProposeA{}

func (m MsgProposeA) Route() string { return ModuleName }

func (m MsgProposeA) Type() string { return msgTypeCreateProposalA }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgProposeA) GetSigners() []sdk.AccAddress {
	return m.Base.Proposers
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgProposeA) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgProposeA) ValidateBasic() error {
	return m.Base.ValidateBasic()
}

var _ sdk.Msg = &MsgProposeB{}

func (m MsgProposeB) Route() string { return ModuleName }

func (m MsgProposeB) Type() string { return msgTypeCreateProposalB }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgProposeB) GetSigners() []sdk.AccAddress {
	return m.Base.Proposers
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgProposeB) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgProposeB) ValidateBasic() error {
	return m.Base.ValidateBasic()
}
