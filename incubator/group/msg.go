package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateGroup{}

const (
	msgTypeGroup = "create_group"
)

func (m MsgCreateGroup) Route() string { return ModuleName }
func (m MsgCreateGroup) Type() string  { return msgTypeGroup }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgCreateGroup) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgCreateGroup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does a sanity check on the provided data
func (m MsgCreateGroup) ValidateBasic() error {
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "members")
	}
	return nil
}
