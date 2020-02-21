package group

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

// NewHandler creates a new message handler.
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

		// todo: @aaronc is this message type supposed to be handled or for extensions only?
		//case MsgCreateGroupAccountBase:
		//case MsgProposeBase:

		case MsgCreateGroupAccountStd:
			return handleMsgCreateGroupAccountStd(ctx, k, msg)
		case MsgVote:
			return handleMsgVote(ctx, k, msg)

		//case MsgExec:

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized group message type: %T", msg)
		}
	}
}

func handleMsgVote(ctx sdk.Context, k Keeper, msg MsgVote) (*sdk.Result, error) {
	return nil, nil
}

// TODO: Do we want to introduce any new events?

func handleMsgCreateGroupAccountStd(ctx sdk.Context, k Keeper, msg MsgCreateGroupAccountStd) (*sdk.Result, error) {
	acc, err := k.CreateGroupAccount(ctx, msg.Base.Admin, msg.Base.Group, *msg.DecisionPolicy.GetThreshold(), msg.Base.Comment)
	if err != nil {
		return nil, errors.Wrap(err, "create group account")
	}
	return buildGroupAccountResult(ctx, msg.Base.Admin, acc, "created")
}

func buildGroupAccountResult(ctx sdk.Context, admin sdk.AccAddress, acc sdk.AccAddress, note string) (*sdk.Result, error) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, admin.String()),
		),
	)
	return &sdk.Result{
		Data:   acc.Bytes(),
		Log:    fmt.Sprintf("Group account %s %s", acc.String(), note),
		Events: ctx.EventManager().Events(),
	}, nil
}
