package hosts

import (
	"context"
	"fmt"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListEBMFeaturesCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DedicatedServerFeature] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerFeatures(args[0])
	}

	opts := &base.BaseListOptions[serverscom.DedicatedServerFeature]{}

	cmd := base.NewListCmd("list-features", "Enterprise bare metal server features", factory, cmdContext, opts)
	cmd.Use = "list-features <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

type featureSetFlags struct {
	Feature           string
	State             string
	IPXEConfig        string
	AuthMethod        []string
	SSHKeyFingerprint []string
}

func newEBMFeatureSetCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &featureSetFlags{}

	cmd := &cobra.Command{
		Use:   "feature-set <id>",
		Short: "Activate or deactivate a feature on an enterprise bare metal server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()
			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			id := args[0]

			var (
				result *serverscom.DedicatedServerFeature
				err    error
			)

			switch flags.State {
			case "activate":
				result, err = activateEBMFeature(ctx, scClient, id, flags)
			case "deactivate":
				result, err = deactivateEBMFeature(ctx, scClient, id, flags.Feature)
			default:
				return fmt.Errorf("invalid state %q: must be activate or deactivate", flags.State)
			}

			if err != nil {
				return err
			}

			if result != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(result)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flags.Feature, "feature", "", "feature name (required)")
	cmd.Flags().StringVar(&flags.State, "state", "", "desired state: activate or deactivate (required)")
	cmd.Flags().StringVar(&flags.IPXEConfig, "ipxe-config", "", "iPXE config script (for private_ipxe_boot)")
	cmd.Flags().StringArrayVar(&flags.AuthMethod, "auth-method", nil, "auth method: password, ssh_key (for host_rescue_mode)")
	cmd.Flags().StringArrayVar(&flags.SSHKeyFingerprint, "ssh-key-fingerprint", nil, "SSH key fingerprint (for host_rescue_mode with ssh_key auth)")

	_ = cmd.MarkFlagRequired("feature")
	_ = cmd.MarkFlagRequired("state")

	return cmd
}

func activateEBMFeature(ctx context.Context, client *serverscom.Client, id string, flags *featureSetFlags) (*serverscom.DedicatedServerFeature, error) {
	switch flags.Feature {
	case "disaggregated_public_ports":
		return client.Hosts.ActivateDisaggregatedPublicPortsFeature(ctx, id)
	case "disaggregated_private_ports":
		return client.Hosts.ActivateDisaggregatedPrivatePortsFeature(ctx, id)
	case "no_public_ip_address":
		return client.Hosts.ActivateNoPublicIpAddressFeature(ctx, id)
	case "no_private_ip":
		return client.Hosts.ActivateNoPrivateIpFeature(ctx, id)
	case "oob_public_access":
		return client.Hosts.ActivateOobPublicAccessFeature(ctx, id)
	case "no_public_network":
		return client.Hosts.ActivateNoPublicNetworkFeature(ctx, id)
	case "host_rescue_mode":
		return client.Hosts.ActivateHostRescueModeFeature(ctx, id, serverscom.HostRescueModeFeatureInput{
			AuthMethods:        flags.AuthMethod,
			SSHKeyFingerprints: flags.SSHKeyFingerprint,
		})
	case "private_ipxe_boot":
		return client.Hosts.ActivatePrivateIpxeBootFeature(ctx, id, serverscom.PrivateIpxeBootFeatureInput{
			IPXEConfig: flags.IPXEConfig,
		})
	default:
		return nil, fmt.Errorf("unsupported feature: %s", flags.Feature)
	}
}

func deactivateEBMFeature(ctx context.Context, client *serverscom.Client, id, feature string) (*serverscom.DedicatedServerFeature, error) {
	switch feature {
	case "disaggregated_public_ports":
		return client.Hosts.DeactivateDisaggregatedPublicPortsFeature(ctx, id)
	case "disaggregated_private_ports":
		return client.Hosts.DeactivateDisaggregatedPrivatePortsFeature(ctx, id)
	case "no_public_ip_address":
		return client.Hosts.DeactivateNoPublicIpAddressFeature(ctx, id)
	case "no_private_ip":
		return client.Hosts.DeactivateNoPrivateIpFeature(ctx, id)
	case "oob_public_access":
		return client.Hosts.DeactivateOobPublicAccessFeature(ctx, id)
	case "no_public_network":
		return client.Hosts.DeactivateNoPublicNetworkFeature(ctx, id)
	case "host_rescue_mode":
		return client.Hosts.DeactivateHostRescueModeFeature(ctx, id)
	case "private_ipxe_boot":
		return client.Hosts.DeactivatePrivateIpxeBootFeature(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported feature: %s", feature)
	}
}
