package testdata

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/modules/incubator/orm"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgPropose:
			return handleMsgPropose(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type: %T", msg)
		}
	}
}

func handleMsgPropose(ctx sdk.Context, k Keeper, msg MsgPropose) (*sdk.Result, error) {
	id, err := k.CreateProposal(ctx, msg.Base.GroupAccount, msg.Base.Proposers, msg.Base.Comment)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data:   orm.EncodeSequence(id),
		Log:    fmt.Sprintf("Proposal created :%d", id),
		Events: ctx.EventManager().Events(),
	}, nil
}
