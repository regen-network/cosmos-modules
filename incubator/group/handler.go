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
		//case MsgUpdateGroupMembers:
		//	return handleMsgUpdateGroupMembers(ctx, k, msg)
		//case MsgUpdateGroupAdmin:
		//	return handleMsgUpdateGroupAdmin(ctx, k, msg)
		//case MsgUpdateGroupComment:
		//	return handleMsgUpdateGroupComment(ctx, k, msg)
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Admin.String()),
		),
	)

	//TODO: the data/ log response is not specified
	return &sdk.Result{
		Data:   id.Byte(),
		Log:    fmt.Sprintf("New group created with id %d", id),
		Events: ctx.EventManager().Events(),
	}, nil
}
