package testdata

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", fmt.Sprintf("x/%s", ModuleName))
		switch msg := msg.(type) {
		case MsgPropose:
			return handleMsgPropose(ctx, k, msg)
		case *MyAppProposalPayloadMsgA:
			logger.Info("executed MyAppProposalPayloadMsgA msg")
			return &sdk.Result{
				Data:   nil,
				Log:    "MyAppProposalPayloadMsgA executed",
				Events: ctx.EventManager().Events(),
			}, nil
		case *MyAppProposalPayloadMsgB:
			logger.Info("executed MyAppProposalPayloadMsgB msg")
			return nil, errors.New("execution of MyAppProposalPayloadMsgB testdata always fails")
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type: %T", msg)
		}
	}
}

func handleMsgPropose(ctx sdk.Context, k Keeper, msg MsgPropose) (*sdk.Result, error) {
	// todo: vaidate
	// check execNow
	id, err := k.CreateProposal(ctx, msg.Base.GroupAccount, msg.Base.Proposers, msg.Base.Comment, msg.Msgs)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data:   id.Byte(),
		Log:    fmt.Sprintf("Proposal created :%d", id),
		Events: ctx.EventManager().Events(),
	}, nil
}
