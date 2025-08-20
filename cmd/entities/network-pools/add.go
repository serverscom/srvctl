package networkpools

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newCreateSubnetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var title string
	var cidr string
	var mask int

	cmd := &cobra.Command{
		Use:   "add-subnet <network_pool_id>",
		Short: "Create a subnetwork",
		Long:  "Create a subnetwork for a network pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.SubnetworkCreateInput{}
			if cmd.Flags().Changed("title") {
				input.Title = &title
			}
			if cmd.Flags().Changed("cidr") {
				input.CIDR = &cidr
			}
			if cmd.Flags().Changed("mask") {
				input.Mask = &mask
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			networkPoolID := args[0]
			sshKey, err := scClient.NetworkPools.CreateSubnetwork(ctx, networkPoolID, input)
			if err != nil {
				return err
			}

			if sshKey != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(sshKey)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Subnetwork name")
	cmd.Flags().StringVar(&cidr, "cidr", "", "If this parameter is filled in, a subnetwork will be created according to the specified Classless Inter-Domain Routing pattern: x.x.x.x/x")
	cmd.Flags().IntVar(&mask, "mask", 0, "Alternative parameter to CIDR. If it's specified, a subnetwork will be allocated from a pool by the mask")

	return cmd
}
