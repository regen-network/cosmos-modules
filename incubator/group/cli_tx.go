package group

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewTxCmd(m codec.Marshaler, txg tx.Generator, ar tx.AccountRetriever) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "Group transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(NewAddGroupTxCmd(m, txg, ar))

	return txCmd
}

// NewSendTxCmd returns a CLI command handler for creating a MsgSend transaction.
func NewAddGroupTxCmd(m codec.Marshaler, txg tx.Generator, ar tx.AccountRetriever) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [admin_key_or_address] [comment]",
		Short: "Create and/or sign and broadcast a MsgCreateGroup transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithMarshaler(m)

			adminAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := &MsgCreateGroup{
				Admin: adminAddr,
				//Members: []group.Member{{
				//	Address: myAddr,
				//	Power:   sdk.NewDec(1),
				//	Comment: "foo",
				//}},
				Comment: args[1],
			}
			txf := tx.NewFactoryFromCLI(inBuf).
				WithTxGenerator(txg).
				WithAccountRetriever(ar)

			return tx.GenerateOrBroadcastTx(cliCtx, txf, msg)
		},
	}
	return flags.PostCommands(cmd)[0]
}
