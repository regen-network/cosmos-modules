package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateGroup{}

const (
	msgTypeCreateGroup           = "create_group"
	msgTypeUpdateGroupAdmin      = "update_group_admin"
	msgTypeUpdateGroupComment    = "update_group_comment"
	msgTypeUpdateGroupMembers    = "update_group_members"
	msgTypeCreateGroupAccountStd = "create_group_account"
	//msgTypeProposeBase           = "create_proposal"
	msgTypeVote         = "vote"
	msgTypeExecProposal = "exec_proposal"
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
func (m MsgUpdateGroupAdmin) Type() string  { return msgTypeUpdateGroupAdmin }

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
func (m MsgUpdateGroupComment) Type() string  { return msgTypeUpdateGroupComment }

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
func (m MsgUpdateGroupMembers) Type() string  { return msgTypeUpdateGroupMembers }

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

func (m *MsgCreateGroupAccountBase) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}
	if len(m.Admin) != sdk.AddrLen {
		return sdkerrors.Wrap(ErrInvalid, "admin")
	}
	return nil
}

func (m *MsgProposeBase) ValidateBasic() error {
	if len(m.GroupAccount) != sdk.AddrLen {
		return sdkerrors.Wrap(ErrInvalid, "group account")
	}
	if len(m.Proposers) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "proposers")
	}
	for i := range m.Proposers {
		if len(m.Proposers[i]) != sdk.AddrLen {
			return sdkerrors.Wrap(ErrInvalid, "proposer account")
		}
	}
	return nil
}

var _ sdk.Msg = &MsgCreateGroupAccountStd{}

func (m MsgCreateGroupAccountStd) Route() string { return ModuleName }
func (m MsgCreateGroupAccountStd) Type() string  { return msgTypeCreateGroupAccountStd }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgCreateGroupAccountStd) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Base.Admin}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgCreateGroupAccountStd) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic does a sanity check on the provided data
func (m MsgCreateGroupAccountStd) ValidateBasic() error {
	if err := m.Base.ValidateBasic(); err != nil {
		return nil
	}
	return nil
}

var _ sdk.Msg = &MsgVote{}

func (m MsgVote) Route() string { return ModuleName }
func (m MsgVote) Type() string  { return msgTypeVote }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgVote) GetSigners() []sdk.AccAddress {
	return m.Voters
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// ValidateBasic does a sanity check on the provided data
func (m MsgVote) ValidateBasic() error {
	if len(m.Voters) == 0 {
		return errors.Wrap(ErrEmpty, "voters")
	}
	for i := range m.Voters {
		if err := sdk.VerifyAddressFormat(m.Voters[i]); err != nil {
			return errors.Wrap(ErrInvalid, "voter")
		}
	}
	if m.Proposal == 0 {
		return errors.Wrap(ErrEmpty, "proposal")
	}
	if m.Choice == Choice_UNKNOWN {
		return errors.Wrap(ErrEmpty, "choice")
	}
	return nil
}
