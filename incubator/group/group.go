package group

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "group"

	// StoreKey defines the primary module store key
	StoreKeyName = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	DefaultParamspace = ModuleName
)

type AppModule struct {
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

func (a AppModule) Name() string {
	return ModuleName
}

func (a AppModule) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (a AppModule) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(NewGenesisState())
}

func (a AppModule) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}
	return data.Validate()
}

func (a AppModule) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	//rest.RegisterRoutes(ctx, r, ModuleCdc, RouterKey)
	// todo: what client functions do we want to support?
	panic("implement me")
}

func (a AppModule) GetTxCmd(*codec.Codec) *cobra.Command {
	//return cli.GetTxCmd(StoreKey, cdc)
	panic("implement me")
}

func (a AppModule) GetQueryCmd(*codec.Codec) *cobra.Command {
	//return cli.GetQueryCmd(StoreKey, cdc)
	panic("implement me")
}

func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	if err := genesisState.Validate(); err != nil {
		panic(fmt.Sprintf("failed to validate %s genesis state: %s", ModuleName, err))
	}
	a.keeper.setParams(ctx, genesisState.Params)
	return []abci.ValidatorUpdate{}

}

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, a.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	// todo: anything to check?
	// todo: check that tally sums must never have less than block before ?
}

func (a AppModule) Route() string {
	return RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.keeper)
}

func (a AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(a.keeper)
}

func (a AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {}

func (a AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
