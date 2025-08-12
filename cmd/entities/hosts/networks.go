package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSNetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Network] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing dedicated server ID")
		}
		id := args[0]
		return scClient.Hosts.DedicatedServerNetworks(id)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Network]{},
		&base.SearchPatternOption[serverscom.Network]{},
		&base.FamilyOption[serverscom.Network]{},
		&base.InterfaceTypeOption[serverscom.Network]{},
		&base.DistributionMethodOption[serverscom.Network]{},
		&base.AdditionalOption[serverscom.Network]{},
	)

	return base.NewListCmd("list-networks <id>", "Dedicated server networks", factory, cmdContext, opts...)
}

func newGetDSNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "get-network <id>",
		Short: ("Get a dedicated server network"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.GetDedicatedServerNetwork(ctx, id, networkID)
			if err != nil {
				return err
			}

			if entity != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(entity)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&networkID, "network-id", "", "Network id (required)")
	_ = cmd.MarkFlagRequired("network-id")

	return cmd
}

func newAddDSNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var (
		networkType        string
		mask               int
		distributionMethod string
	)

	cmd := &cobra.Command{
		Use:   "add-network <id>",
		Short: "Add private/public IPv4 network to dedicated server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			if err := validateNetworkArgs(networkType, distributionMethod, mask); err != nil {
				return err
			}

			input := serverscom.NetworkInput{
				DistributionMethod: distributionMethod,
				Mask:               mask,
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()
			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			var entity *serverscom.Network
			var err error
			switch networkType {
			case "public":
				entity, err = scClient.Hosts.AddDedicatedServerPublicIPv4Network(ctx, id, input)
			case "private":
				entity, err = scClient.Hosts.AddDedicatedServerPrivateIPv4Network(ctx, id, input)
			}
			if err != nil {
				return err
			}
			if entity != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(entity)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&networkType, "type", "", "Network type: 'public' or 'private'")
	cmd.Flags().IntVar(&mask, "mask", 0, "Network mask (required)")
	cmd.Flags().StringVar(&distributionMethod, "distribution-method", "gateway", "Distribution method ('gateway' or 'route')")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("mask")

	return cmd
}

func newDeleteDSNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "delete-network <id>",
		Short: ("Delete a dedicated server network"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.DeleteDedicatedServerNetwork(ctx, id, networkID)
			if err != nil {
				return err
			}

			if entity != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(entity)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&networkID, "network-id", "", "Network id (required)")
	_ = cmd.MarkFlagRequired("network-id")

	return cmd
}
