package group

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func createTestApp(isCheckTx bool) (*SimApp, sdk.Context) {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
	genesisState := ModuleBasics.DefaultGenesis()
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
	app.Commit()
	header := abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := app.NewContext(isCheckTx, header)
	return app, ctx
}

func TestCreateGroupScenario(t *testing.T) {
	app, ctx := createTestApp(false)
	priv1, _, addr1 := types.KeyTestPubAddr()
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)

	balances := sdk.NewCoins(sdk.NewInt64Coin("atom", 1000))
	require.NoError(t, app.BankKeeper.SetBalances(ctx, addr1, balances))

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

	privs, accNums, seqs := []crypto.PrivKey{priv1}, acc1.GetAccountNumber(), acc1.GetSequence()
	tx := types.NewTestTx(ctx, msgs, privs, []uint64{accNums}, []uint64{seqs}, fee)

	resp := app.DeliverTx(abci.RequestDeliverTx{Tx: app.Codec().MustMarshalBinaryLengthPrefixed(tx)})

	require.Equal(t, uint32(0), resp.Code, resp.Log)
	assert.Equal(t, orm.EncodeSequence(1), resp.Data)
	assert.True(t, app.GroupKeeper.groupTable.Has(ctx, resp.Data))
}
