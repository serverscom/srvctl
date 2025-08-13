package uplinkbandwidths

import (
	"fmt"
	"strconv"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64
	var serverModelID int64
	var uplinkModelID int64

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a bandwidth option for an uplink model",
		Long:  "Get a bandwidth option for an uplink model by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			bandwidthID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("bandwidth id should be integer")
			}

			model, err := scClient.Locations.GetBandwidthOption(ctx, locationID, serverModelID, uplinkModelID, int64(bandwidthID))
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
	cmd.Flags().Int64Var(&serverModelID, "server-model-id", 0, "Server model id (int, required)")
	cmd.Flags().Int64Var(&uplinkModelID, "uplink-model-id", 0, "Uplink model id (int, required)")
	_ = cmd.MarkFlagRequired("location-id")
	_ = cmd.MarkFlagRequired("server-model-id")
	_ = cmd.MarkFlagRequired("uplink-model-id")

	return cmd
}
