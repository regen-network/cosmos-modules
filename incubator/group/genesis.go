package group

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	// TODO: please define requirements
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState()
}

// ValidateGenesis performs basic validation of group genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}

// InitGenesis seeds the modele from genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
}

// ExportGenesis returns a GenesisState for a given context and Keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState()
}
