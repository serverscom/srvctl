package networkpools

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var title string
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <network_pool_id>",
		Short: "Update a network pool",
		Long:  "Update a network pool by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				log.Fatal(err)
			}
			input := serverscom.NetworkPoolInput{
				Labels: labelsMap,
			}
			if cmd.Flags().Changed("title") {
				input.Title = &title
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			networkPoolID := args[0]
			sshKey, err := scClient.NetworkPools.Update(ctx, networkPoolID, input)
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

	cmd.Flags().StringVar(&title, "title", "", "string")
	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func newUpdateSubnetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var title string
	var subnetID string

	cmd := &cobra.Command{
		Use:   "update-subnet <network_pool_id>",
		Short: "Update a subnetwork for a network pool",
		Long:  "Update a subnetwork for a network pool by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.SubnetworkUpdateInput{}
			if cmd.Flags().Changed("title") {
				input.Title = &title
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			networkPoolID := args[0]
			sshKey, err := scClient.NetworkPools.UpdateSubnetwork(ctx, networkPoolID, subnetID, input)
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
	cmd.Flags().StringVar(&subnetID, "network-id", "", "A unique identifier of a pool")
	_ = cmd.MarkFlagRequired("network-id")

	return cmd
}
