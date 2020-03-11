package group

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/jsonpb"
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

// AccountCondition returns a condition to build a group account address.
func AccountCondition(id uint64) Condition {
	return NewCondition("group", "account", orm.EncodeSequence(id))
}

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
	RegisterCodec(cdc) // can not be removed until sdk.StdTx support protobuf
}

func (a AppModule) DefaultGenesis() json.RawMessage {
	var buf bytes.Buffer
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(&buf, NewGenesisState())
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal default genesis"))
	}
	return buf.Bytes()
}

func (a AppModule) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := jsonpb.Unmarshal(bytes.NewReader(bz), &data); err != nil {
		return errors.Wrapf(err, "validate genesis")
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

func (a AppModule) InitGenesis(ctx sdk.Context, bz json.RawMessage) []abci.ValidatorUpdate {
	var data GenesisState
	if err := jsonpb.Unmarshal(bytes.NewReader(bz), &data); err != nil {
		panic(errors.Wrapf(err, "init genesis"))
	}

	if err := data.Validate(); err != nil {
		panic(fmt.Sprintf("failed to validate %s genesis state: %s", ModuleName, err))
	}
	a.keeper.setParams(ctx, data.Params)
	return []abci.ValidatorUpdate{}

}

func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	var buf bytes.Buffer
	marshaller := jsonpb.Marshaler{}
	if err := marshaller.Marshal(&buf, ExportGenesis(ctx, a.keeper)); err != nil {
		panic(errors.Wrap(err, "export genesis"))
	}
	return buf.Bytes()
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
