package group

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/pkg/errors"
)

const (
	msgTypeCreateGroup           = "create_group"
	msgTypeUpdateGroupAdmin      = "update_group_admin"
	msgTypeUpdateGroupComment    = "update_group_comment"
	msgTypeUpdateGroupMembers    = "update_group_members"
	msgTypeCreateGroupAccountStd = "create_group_account"
	msgTypeVote                  = "vote"
	msgTypeExecProposal          = "exec_proposal"
)

type MsgCreateGroupAccountI interface {
	GetBase() MsgCreateGroupAccountBase
	GetDecisionPolicy() StdDecisionPolicy
}

var _ sdk.Msg = &MsgCreateGroup{}

func (m MsgCreateGroup) Route() string { return ModuleName }
func (m MsgCreateGroup) Type() string  { return msgTypeCreateGroup }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgCreateGroup) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Admin}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgCreateGroup) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgCreateGroup) ValidateBasic() error {
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(m.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
	}
	if err := Members(m.Members).ValidateBasic(); err != nil {
		return errors.Wrap(err, "members")
	}
	for i := range m.Members {
		member := m.Members[i]
		if member.Power.Equal(sdk.ZeroDec()) {
			return sdkerrors.Wrap(ErrEmpty, "member power")
		}
	}
	return nil
}

type Members []Member

func (ms Members) ValidateBasic() error {
	index := make(map[string]struct{}, len(ms))
	for i := range ms {
		member := ms[i]
		if err := member.ValidateBasic(); err != nil {
			return err
		}
		addr := string(member.Address)
		if _, exists := index[addr]; exists {
			return errors.Wrapf(ErrDuplicate, "address: %s", member.Address)
		}
		index[addr] = struct{}{}
	}
	return nil
}

func (m Member) ValidateBasic() error {
	if m.Address.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "address")
	}
	if m.Power.IsNil() || m.Power.LT(sdk.ZeroDec()) {
		return sdkerrors.Wrap(ErrInvalid, "power")
	}
	if err := sdk.VerifyAddressFormat(m.Address); err != nil {
		return sdkerrors.Wrap(err, "address")
	}
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
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgUpdateGroupAdmin) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}

	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(m.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
	}

	if m.NewAdmin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "new admin")
	}
	if err := sdk.VerifyAddressFormat(m.NewAdmin); err != nil {
		return sdkerrors.Wrap(err, "new admin")
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
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgUpdateGroupComment) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(m.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
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
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgUpdateGroupMembers) ValidateBasic() error {
	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")

	}
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(m.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
	}

	if len(m.MemberUpdates) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "member updates")
	}
	if err := Members(m.MemberUpdates).ValidateBasic(); err != nil {
		return errors.Wrap(err, "members")
	}
	return nil
}

func (m *MsgCreateGroupAccountBase) ValidateBasic() error {
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(m.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
	}

	if m.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")
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
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgCreateGroupAccountStd) ValidateBasic() error {
	if err := m.Base.ValidateBasic(); err != nil {
		return errors.Wrap(err, "base")
	}
	if m.DecisionPolicy.GetDecisionPolicy() == nil {
		return errors.Wrap(ErrEmpty, "decision policy")
	}
	if err := m.DecisionPolicy.GetDecisionPolicy().ValidateBasic(); err != nil {
		return errors.Wrap(err, "decision policy")
	}
	return nil
}

var _ MsgCreateGroupAccountI = MsgCreateGroupAccountStd{}

func (m MsgCreateGroupAccountStd) GetBase() MsgCreateGroupAccountBase {
	return m.Base
}

func (m MsgCreateGroupAccountStd) GetDecisionPolicy() StdDecisionPolicy {
	return m.DecisionPolicy
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
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
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
	// todo: prevent duplicates in votes, ignore or normalize later?

	if m.Proposal == 0 {
		return errors.Wrap(ErrEmpty, "proposal")
	}
	if m.Choice == Choice_UNKNOWN {
		return errors.Wrap(ErrEmpty, "choice")
	}
	return nil
}

var _ sdk.Msg = &MsgExec{}

func (m MsgExec) Route() string { return ModuleName }
func (m MsgExec) Type() string  { return msgTypeExecProposal }

// GetSigners returns the addresses that must sign over msg.GetSignBytes()
func (m MsgExec) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (m MsgExec) GetSignBytes() []byte {
	var buf bytes.Buffer
	enc := jsonpb.Marshaler{}
	if err := enc.Marshal(&buf, &m); err != nil {
		panic(errors.Wrap(err, "get sign bytes"))
	}
	return sdk.MustSortJSON(buf.Bytes())
}

// ValidateBasic does a sanity check on the provided data
func (m MsgExec) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Signer); err != nil {
		return errors.Wrap(ErrInvalid, "voter")
	}
	if m.Proposal == 0 {
		return errors.Wrap(ErrEmpty, "proposal")
	}
	return nil
}
