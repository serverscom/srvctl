package cloudregions

import (
	"strconv"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCredentialsCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-credentials <region-id>",
		Short: "Get credentials for cloud region by Region ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			regionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			credentials, err := scClient.CloudComputingRegions.Credentials(ctx, regionID)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(credentials)
		},
	}

	return cmd
}
