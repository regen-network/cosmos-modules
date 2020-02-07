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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic does a sanity check on the provided data
func (m MsgCreateGroup) ValidateBasic() error {
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
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
	// empty members list allowed
	// max members list allowed???
	// duplicate member address
	// member address empty
	// Power -1, 0
	// comment >max
	return nil
}
