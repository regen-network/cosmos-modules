package group

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrEmpty = sdkerrors.Register(ModuleName, 1, "value is empty")
)
