package sbmosoptions

import (
	"fmt"
	"strconv"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64
	var sbmFlavorModelID int64

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an operating system",
		Long:  "Get an operating system for an SBM server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			osID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("operating system id should be integer")
			}

			model, err := scClient.Locations.GetSBMOperatingSystemOption(ctx, locationID, sbmFlavorModelID, int64(osID))
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
	cmd.Flags().Int64Var(&sbmFlavorModelID, "model-id", 0, "SBM flavor model id (int, required)")
	_ = cmd.MarkFlagRequired("location-id")
	_ = cmd.MarkFlagRequired("model-id")

	return cmd
}
