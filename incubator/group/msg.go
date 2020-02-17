package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateGroup{}

const (
	msgTypeGroup = "create_group"
)

func (m MsgCreateGroup) Route() string { return ModuleName }
func (m MsgCreateGroup) Type() string  { return msgTypeGroup }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgCreateGroup) GetSigners() []sdk.AccAddress {
	// TODO: @aaronc ok with this constraint? We can enforce the signature on creation. Also see `MsgUpdateGroupAdmin.GetSigners()`.
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
	// TODO: @aaronc can we ensure that the group is never empty? That would simplify authZ logic later as we do not have
	// to protect against this case
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "members")
	}
	index := make(map[string]struct{}, len(m.Members))
	for i := range m.Members {
		member := m.Members[i]
		if err := member.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "member")
		}
		if member.Power.LTE(sdk.ZeroDec()) {
			return sdkerrors.Wrap(ErrEmpty, "member power")
		}
		addr := member.Address.String()
		if _, exists := index[addr]; exists {
			return errors.Wrapf(ErrDuplicate, "address: %s", addr)
		}
	}
	// todo: test
	// duplicate member address
	// member address empty
	// Power -1, 0
	// comment >max
	return nil
}
