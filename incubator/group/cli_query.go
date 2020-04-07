package group

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/spf13/cobra"
)

func NewQueryCmd(m codec.Marshaler) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "Querying commands for the group module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(ListGroupsCmd(m))

	return cmd
}

func ListGroupsCmd(m codec.Marshaler) *cobra.Command {
	var limit uint8
	var cursor string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Query for groups",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithMarshaler(m)

			//var params interface{}
			route := fmt.Sprintf("custom/%s/%s?range&limit=%d", QuerierRoute, "xgroup", limit)
			if cursor != "" {
				route = route + "&cursor=" + cursor
			}
			//bz, err := m.MarshalJSON(params)
			//if err != nil {
			//	return fmt.Errorf("failed to marshal params: %w", err)
			//}

			res, _, err := cliCtx.QueryWithData(route, []byte{})
			if err != nil {
				return err
			}
			var result orm.QueryResult
			if err := m.UnmarshalJSON(res, &result); err != nil {
				return err
			}
			println("response received: "+ string(res))
			// todo: fails with encoding issues :-/
			//return cliCtx.PrintOutput(result)
			return nil
		},
	}
	cmd.Flags().Uint8Var(&limit,"limit", 25, "max number of results")
	cmd.Flags().StringVar(&cursor,"cursor", "", "cursor to navigate through result set")
	return flags.GetCommands(cmd)[0]
}
