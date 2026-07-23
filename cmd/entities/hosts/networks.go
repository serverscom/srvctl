package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListEBMNetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Network] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerNetworks(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Network]{},
		&base.SearchPatternOption[serverscom.Network]{},
		&base.FamilyOption[serverscom.Network]{},
		&base.InterfaceTypeOption[serverscom.Network]{},
		&base.DistributionMethodOption[serverscom.Network]{},
		&base.AdditionalOption[serverscom.Network]{},
	)

	cmd := base.NewListCmd("list-networks", "Enterprise bare metal server networks", factory, cmdContext, opts...)
	cmd.Use = "list-networks <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func newGetEBMNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "get-network <id>",
		Short: ("Get an enterprise bare metal server network"),
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

func newAddEBMNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var (
		networkType        string
		mask               int
		distributionMethod string
	)

	cmd := &cobra.Command{
		Use:   "add-network <id>",
		Short: "Add private/public IPv4 network to enterprise bare metal server",
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

func newDeleteEBMNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "delete-network <id>",
		Short: ("Delete an enterprise bare metal server network"),
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

func newListKBMNetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Network] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.KubernetesBaremetalNodeNetworks(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Network]{},
		&base.SearchPatternOption[serverscom.Network]{},
		&base.FamilyOption[serverscom.Network]{},
		&base.InterfaceTypeOption[serverscom.Network]{},
		&base.DistributionMethodOption[serverscom.Network]{},
		&base.AdditionalOption[serverscom.Network]{},
	)

	cmd := base.NewListCmd("list-networks", "KBM node networks", factory, cmdContext, opts...)
	cmd.Use = "list-networks <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func newListSBMNetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Network] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.SBMServerNetworks(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Network]{},
		&base.SearchPatternOption[serverscom.Network]{},
		&base.FamilyOption[serverscom.Network]{},
		&base.InterfaceTypeOption[serverscom.Network]{},
		&base.DistributionMethodOption[serverscom.Network]{},
		&base.AdditionalOption[serverscom.Network]{},
	)

	cmd := base.NewListCmd("list-networks", "Scalable baremetal server networks", factory, cmdContext, opts...)
	cmd.Use = "list-networks <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func newGetSBMNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "get-network <id>",
		Short: ("Get a scalable baremetal server network"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.GetSBMServerNetwork(ctx, id, networkID)
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

func newAddSBMNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var (
		mask               int
		distributionMethod string
	)

	cmd := &cobra.Command{
		Use:   "add-network <id>",
		Short: "Add private IPv4 network to scalable baremetal server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			if err := validateNetworkArgs("private", distributionMethod, mask); err != nil {
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

			entity, err := scClient.Hosts.AddSBMServerPrivateIPv4Network(ctx, id, input)
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

	cmd.Flags().IntVar(&mask, "mask", 0, "Network mask (required)")
	cmd.Flags().StringVar(&distributionMethod, "distribution-method", "gateway", "Distribution method ('gateway')")
	_ = cmd.MarkFlagRequired("mask")

	return cmd
}

func newDeleteSBMNetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "delete-network <id>",
		Short: ("Delete a scalable baremetal server network"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.DeleteSBMServerNetwork(ctx, id, networkID)
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

func newActivateEBMIPv6NetworkCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-ipv6-network <id>",
		Short: ("Activate a public IPv6 network for an enterprise bare metal server"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.ActivateDedicatedServerPubliIPv6Network(ctx, id)
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

	return cmd
}

func newGetEBMNetworkUsageCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-usage <id>",
		Short: ("Get network utilization for an enterprise bare metal server"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.GetDedicatedServerNetworkUsage(ctx, id)
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

	return cmd
}

func newGetSBMNetworkUsageCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-usage <id>",
		Short: ("Get network utilization for a scalable baremetal server"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.GetSBMServerNetworkUsage(ctx, id)
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

	return cmd
}
