package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateGroup:
			return handleMsgCreateGroup(ctx, k, msg)
		//case MsgUpdateGroupMembers:
		//	return handleMsgUpdateGroupMembers(ctx, k, msg)
		//case MsgUpdateGroupAdmin:
		//	return handleMsgUpdateGroupAdmin(ctx, k, msg)
		//case MsgUpdateGroupComment:
		//	return handleMsgUpdateGroupComment(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized bank message type: %T", msg)
		}
	}

}

func handleMsgCreateGroup(ctx sdk.Context, k Keeper, msg MsgCreateGroup) (*sdk.Result, error) {

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Admin.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
