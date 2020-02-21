package group_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/modules/incubator/group"
	"github.com/cosmos/modules/incubator/group/testdata"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func createTestApp(isCheckTx bool) (*testdata.SimApp, sdk.Context) {
	db := dbm.NewMemDB()
	app := testdata.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
	genesisState := testdata.ModuleBasics.DefaultGenesis()
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
	myKey, _, myAddr := types.KeyTestPubAddr()
	myAccount := app.AccountKeeper.NewAccountWithAddress(ctx, myAddr)
	app.AccountKeeper.SetAccount(ctx, myAccount)

	balances := sdk.NewCoins(sdk.NewInt64Coin("atom", 1000))
	require.NoError(t, app.BankKeeper.SetBalances(ctx, myAddr, balances))

	fee := types.NewTestStdFee()

	msgs := []sdk.Msg{group.MsgCreateGroup{
		Admin: myAddr,
		Members: []group.Member{{
			Address: myAddr,
			Power:   sdk.NewDec(1),
			Comment: "foo",
		}},
		Comment: "integration test",
	}}

	privs, accNums, seqs := []crypto.PrivKey{myKey}, myAccount.GetAccountNumber(), myAccount.GetSequence()
	tx := types.NewTestTx(ctx, msgs, privs, []uint64{accNums}, []uint64{seqs}, fee)

	resp := app.DeliverTx(abci.RequestDeliverTx{Tx: app.Codec().MustMarshalBinaryLengthPrefixed(tx)})

	require.Equal(t, uint32(0), resp.Code, resp.Log)
	assert.Equal(t, orm.EncodeSequence(1), resp.Data)
	assert.True(t, app.GroupKeeper.HasGroup(ctx, resp.Data))
}

func TestCreateProposal(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup account
	myKey, _, myAddr := types.KeyTestPubAddr()
	myAccount := app.AccountKeeper.NewAccountWithAddress(ctx, myAddr)
	app.AccountKeeper.SetAccount(ctx, myAccount)

	balances := sdk.NewCoins(sdk.NewInt64Coin("atom", 100000))
	require.NoError(t, app.BankKeeper.SetBalances(ctx, myAddr, balances))

	// setup group
	msgs := []sdk.Msg{
		group.MsgCreateGroup{
			Admin: myAddr,
			Members: []group.Member{{
				Address: myAddr,
				Power:   sdk.NewDec(1),
				Comment: "me",
			}},
			Comment: "integration test",
		},
		// setup group account
		group.MsgCreateGroupAccountStd{
			Base: group.MsgCreateGroupAccountBase{
				Admin:   myAddr,
				Group:   1,
				Comment: "first account",
			},
			DecisionPolicy: group.StdDecisionPolicy{
				Sum: &group.StdDecisionPolicy_Threshold{
					Threshold: &group.ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
					},
				},
			},
		},
		// submit proposal
		testdata.MsgPropose{
			Base: group.MsgProposeBase{
				GroupAccount: make([]byte, sdk.AddrLen), // todo: see comment in keeper.CreateGroupAccount
				Proposers:    []sdk.AccAddress{myAddr},
				Comment:      "ok",
			},
		},
		// vote
		group.MsgVote{
			Proposal: 0,
			Voters:   []sdk.AccAddress{myAddr},
			Choice:   group.Choice_YES,
			Comment:  "all in",
		},
		// execute
	}
	fee := types.NewTestStdFee()
	privs, accNums, seqs := []crypto.PrivKey{myKey}, myAccount.GetAccountNumber(), myAccount.GetSequence()
	tx := types.NewTestTx(ctx, msgs, privs, []uint64{accNums}, []uint64{seqs}, fee)

	resp := app.DeliverTx(abci.RequestDeliverTx{Tx: app.Codec().MustMarshalBinaryLengthPrefixed(tx)})
	require.Equal(t, uint32(0), resp.Code, resp.Log)

	// and then register a group account

	// and then create proposal

	//	msgs = []sdk.Msg{P{
	//		Admin: myAddr,
	//		Members: []group.Member{{
	//			Address: myAddr,
	//			Power:   sdk.NewDec(1),
	//			Comment: "me",
	//		}},
	//		Comment: "integration test",
	//	}}
	//	fee := types.NewTestStdFee()
	//	privs, accNums, seqs := []crypto.PrivKey{myKey}, myAccount.GetAccountNumber(), myAccount.GetSequence()
	//	tx := types.NewTestTx(ctx, msgs, privs, []uint64{accNums}, []uint64{seqs}, fee)
	//
	//	resp := app.DeliverTx(abci.RequestDeliverTx{Tx: app.Codec().MustMarshalBinaryLengthPrefixed(tx)})
	//	require.Equal(t, uint32(0), resp.Code, resp.Log)
	//
	//
}