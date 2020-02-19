package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateGroup{}

const (
	msgTypeCreateGroup         = "create_group"
	msgTypeCUpdateGroupAdmin   = "update_group_admin"
	msgTypeCUpdateGroupComment = "update_group_comment"
	msgTypeCUpdateGroupMembers = "update_group_members"
)

func (m MsgCreateGroup) Route() string { return ModuleName }
func (m MsgCreateGroup) Type() string  { return msgTypeCreateGroup }

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

var _ sdk.Msg = &MsgUpdateGroupAdmin{}

func (m MsgUpdateGroupAdmin) Route() string { return ModuleName }
func (m MsgUpdateGroupAdmin) Type() string  { return msgTypeCUpdateGroupAdmin }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgUpdateGroupAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgUpdateGroupAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic does a sanity check on the provided data
func (m MsgUpdateGroupAdmin) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}

	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if m.NewAdmin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "new admin")
	}
	if m.Admin.Equals(m.NewAdmin) {
		return sdkerrors.Wrap(ErrInvalid, "new and old admin are the same")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateGroupComment{}

func (m MsgUpdateGroupComment) Route() string { return ModuleName }
func (m MsgUpdateGroupComment) Type() string  { return msgTypeCUpdateGroupComment }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgUpdateGroupComment) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgUpdateGroupComment) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic does a sanity check on the provided data
func (m MsgUpdateGroupComment) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateGroupMembers{}

func (m MsgUpdateGroupMembers) Route() string { return ModuleName }
func (m MsgUpdateGroupMembers) Type() string  { return msgTypeCUpdateGroupMembers }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgUpdateGroupMembers) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgUpdateGroupMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic does a sanity check on the provided data
func (m MsgUpdateGroupMembers) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if len(m.MemberUpdates) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "member updates")
	}
	index := make(map[string]struct{}, len(m.MemberUpdates))
	for i := range m.MemberUpdates {
		member := m.MemberUpdates[i]
		if err := member.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "member")
		}
		if member.Power.LT(sdk.ZeroDec()) {
			return sdkerrors.Wrap(ErrInvalid, "member power")
		}
		addr := member.Address.String()
		if _, exists := index[addr]; exists {
			return errors.Wrapf(ErrDuplicate, "address: %s", addr)
		}
	}
	return nil
}
