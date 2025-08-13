package servermodels

import (
	"fmt"
	"strconv"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a server model",
		Long:  "Get a server model by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			serverModelID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("server model id should be integer")
			}

			model, err := scClient.Locations.GetServerModelOption(ctx, locationID, int64(serverModelID))
			if err != nil {
				return err
			}

			if model != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(model)
			}
			return nil
		},
	}

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location id (int, required)")
	_ = cmd.MarkFlagRequired("location-id")

	return cmd
}
