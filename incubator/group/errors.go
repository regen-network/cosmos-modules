package group

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrEmpty     = sdkerrors.Register(ModuleName, 1, "value is empty")
	ErrDuplicate = sdkerrors.Register(ModuleName, 2, "duplicate value")
	ErrMaxLimit  = sdkerrors.Register(ModuleName, 3, "limit exceeded")
	ErrType      = sdkerrors.Register(ModuleName, 4, "invalid type")
	ErrInvalid   = sdkerrors.Register(ModuleName, 5, "invalid value")
)
