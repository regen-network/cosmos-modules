package group

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func createTestApp(isCheckTx bool) (*SimApp, sdk.Context) {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		var genesisState map[string]json.RawMessage
		stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}
	ctx := app.NewContext(isCheckTx, abci.Header{})
	// oh man.... :-/
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	return app, ctx
}

func TestCreateGroupScenario(t *testing.T) {
	app, ctx := createTestApp(false)
	priv1, _, addr1 := types.KeyTestPubAddr()
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)

	balances := sdk.NewCoins(sdk.NewInt64Coin("atom", 1000))
	require.NoError(t, app.BankKeeper.SetBalances(ctx, addr1, balances))

	ctx = ctx.WithBlockHeight(1)

	fee := types.NewTestStdFee()

	msgs := []sdk.Msg{MsgCreateGroup{
		Admin: addr1,
		Members: []Member{{
			Address: addr1,
			Power:   sdk.NewDec(1),
			Comment: "foo",
		}},
		Comment: "integration test",
	}}

	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	resp := app.DeliverTx(abci.RequestDeliverTx{Tx: app.Codec().MustMarshalBinaryLengthPrefixed(tx)})
	require.Equal(t, uint32(0), resp.Code, resp.Log)
	//require.NoError(t, resp.)
	_ = resp
	//assert.Equal(t, orm.EncodeSequence(1), result.Data)
	//assert.NotEmpty(t, gas.GasUsed)
}
