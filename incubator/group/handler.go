package group

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateGroup:
			return handleMsgCreateGroup(ctx, k, msg)
		case MsgUpdateGroupAdmin:
			return handleMsgUpdateGroupAdmin(ctx, k, msg)
		case MsgUpdateGroupComment:
			return handleMsgUpdateGroupComment(ctx, k, msg)
		case MsgUpdateGroupMembers:
			return handleMsgUpdateGroupMembers(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized group message type: %T", msg)
		}
	}
}

func handleMsgCreateGroup(ctx sdk.Context, k Keeper, msg MsgCreateGroup) (*sdk.Result, error) {
	id, err := k.CreateGroup(ctx, msg.Admin, msg.Members, msg.Comment)
	if err != nil {
		return nil, errors.Wrap(err, "create group")
	}
	return buildGroupResult(ctx, msg.Admin, id, "created")
}

func handleMsgUpdateGroupAdmin(ctx sdk.Context, k Keeper, msg MsgUpdateGroupAdmin) (*sdk.Result, error) {
	action := func(m *GroupMetadata) error {
		m.Admin = msg.NewAdmin
		return k.UpdateGroup(ctx, m)
	}
	return doAuthenticated(k, ctx, &msg, action, "admin updated")
}

func handleMsgUpdateGroupComment(ctx sdk.Context, k Keeper, msg MsgUpdateGroupComment) (*sdk.Result, error) {
	action := func(m *GroupMetadata) error {
		m.Comment = msg.Comment
		return k.UpdateGroup(ctx, m)
	}
	return doAuthenticated(k, ctx, &msg, action, "comment updated")
}

func handleMsgUpdateGroupMembers(ctx sdk.Context, k Keeper, msg MsgUpdateGroupMembers) (*sdk.Result, error) {
	action := func(m *GroupMetadata) error {
		for i := range msg.MemberUpdates {
			member := GroupMember{Group: msg.Group,
				Member:  msg.MemberUpdates[i].Address,
				Weight:  msg.MemberUpdates[i].Power,
				Comment: msg.MemberUpdates[i].Comment,
			}
			if member.Weight.Equal(sdk.ZeroDec()) {
				if err := k.groupMemberTable.Delete(ctx, &member); err != nil {
					return errors.Wrap(err, "delete member")
				}
				continue
			}

			// todo: a PUT would be nicer in this scenario to save the extra cost for the additional Has operation
			// todo: revisit `Has` syntax: groupMemberTable.Has(ctx, member) ?
			if k.groupMemberTable.Has(ctx, member.NaturalKey()) {
				if err := k.groupMemberTable.Save(ctx, &member); err != nil {
					return errors.Wrap(err, "add member")
				}
			} else {
				if err := k.groupMemberTable.Create(ctx, &member); err != nil {
					return errors.Wrap(err, "add member")
				}
			}
		}
		return k.UpdateGroup(ctx, m)
	}
	return doAuthenticated(k, ctx, &msg, action, "members updated")

}

type authNGroupMsg interface {
	GetGroup() GroupID
	GetAdmin() sdk.AccAddress // equal GetSigners()
}

func doAuthenticated(k Keeper, ctx sdk.Context, msg authNGroupMsg, action func(*GroupMetadata) error, note string) (*sdk.Result, error) {
	group, err := k.GetGroup(ctx, msg.GetGroup())
	if err != nil {
		return nil, err
	}
	if !group.Admin.Equals(msg.GetAdmin()) {
		return nil, errors.Wrap(ErrUnauthorized, "not group admin")
	}
	if err := action(&group); err != nil {
		return nil, errors.Wrap(err, note)
	}
	return buildGroupResult(ctx, msg.GetAdmin(), msg.GetGroup(), note)
}

func buildGroupResult(ctx sdk.Context, admin sdk.AccAddress, group GroupID, note string) (*sdk.Result, error) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, admin.String()),
		),
	)
	return &sdk.Result{
		Data:   group.Byte(),
		Log:    fmt.Sprintf("Group %d %s", group, note),
		Events: ctx.EventManager().Events(),
	}, nil
}
