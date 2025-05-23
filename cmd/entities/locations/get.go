package locations

import (
	"fmt"
	"strconv"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a location",
		Long:  "Get a location by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("location id should be integer")
			}
			location, err := scClient.Locations.GetLocation(ctx, int64(id))
			if err != nil {
				return err
			}

			if location != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(location)
			}
			return nil
		},
	}

	return cmd
}
