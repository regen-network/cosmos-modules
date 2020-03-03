package testdata

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/jsonpb"
)

const (
	msgTypeCreateProposal = "create_proposal"
	msgTypeMyMsgA         = "my_msg_a"
	msgTypeMyMsgB         = "my_msg_b"
)

var _ sdk.Msg = &MsgPropose{}

func (m MsgPropose) Route() string { return ModuleName }

func (m MsgPropose) Type() string { return msgTypeMyMsgA }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgPropose) GetSigners() []sdk.AccAddress {
	return m.Base.Proposers
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgPropose) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgPropose) ValidateBasic() error {
	if err := m.Base.ValidateBasic(); err != nil {
		return err
	}
	for i, v := range m.Msgs {
		if err := v.GetMsg().ValidateBasic(); err != nil {
			return errors.Wrapf(err, "msg %d", i)
		}
	}
	return nil
}

var _ sdk.Msg = &MyAppProposalPayloadMsgA{}

func (m MyAppProposalPayloadMsgA) Route() string { return ModuleName }

func (m MyAppProposalPayloadMsgA) Type() string { return msgTypeMyMsgA }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MyAppProposalPayloadMsgA) GetSigners() []sdk.AccAddress {
	return nil // nothing to do
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MyAppProposalPayloadMsgA) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MyAppProposalPayloadMsgA) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MyAppProposalPayloadMsgB{}

func (m MyAppProposalPayloadMsgB) Route() string { return ModuleName }

func (m MyAppProposalPayloadMsgB) Type() string { return msgTypeMyMsgB }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MyAppProposalPayloadMsgB) GetSigners() []sdk.AccAddress {
	return nil
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MyAppProposalPayloadMsgB) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MyAppProposalPayloadMsgB) ValidateBasic() error {
	return nil
}
