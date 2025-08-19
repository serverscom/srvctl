package cmd

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/cmd/config"
	"github.com/serverscom/srvctl/cmd/context"
	"github.com/serverscom/srvctl/cmd/entities/account"
	"github.com/serverscom/srvctl/cmd/entities/drivemodels"
	"github.com/serverscom/srvctl/cmd/entities/hosts"
	"github.com/serverscom/srvctl/cmd/entities/invoices"
	"github.com/serverscom/srvctl/cmd/entities/k8s"
	l2segments "github.com/serverscom/srvctl/cmd/entities/l2_segments"
	loadbalancers "github.com/serverscom/srvctl/cmd/entities/load_balancers"
	"github.com/serverscom/srvctl/cmd/entities/locations"
	"github.com/serverscom/srvctl/cmd/entities/racks"
	sbmmodels "github.com/serverscom/srvctl/cmd/entities/sbm_models"
	sbmosoptions "github.com/serverscom/srvctl/cmd/entities/sbm_os_options"
	serverosoptions "github.com/serverscom/srvctl/cmd/entities/server_os_options"
	serverramoptions "github.com/serverscom/srvctl/cmd/entities/server_ram_options"
	"github.com/serverscom/srvctl/cmd/entities/servermodels"
	sshkeys "github.com/serverscom/srvctl/cmd/entities/ssh-keys"
	"github.com/serverscom/srvctl/cmd/entities/ssl"
	"github.com/serverscom/srvctl/cmd/entities/uplinkbandwidths"
	"github.com/serverscom/srvctl/cmd/entities/uplinkmodels"
	"github.com/serverscom/srvctl/cmd/login"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	cobra.EnableTraverseRunHooks = true

	cmdContext := &base.CmdContext{}

	cmd := &cobra.Command{
		Use:               "srvctl",
		Short:             "CLI tool for servers.com API",
		Long:              `A command line interface for managing servers.com resources`,
		Version:           version,
		PersistentPreRunE: base.InitCmdContext(cmdContext),
		SilenceUsage:      true,
	}

	// Global flags
	base.AddGlobalFlags(cmd)

	clientFactory := &client.DefaultClientFactory{}

	// Add commands
	cmd.AddCommand(login.NewCmd(cmdContext, clientFactory))
	cmd.AddCommand(context.NewCmd(cmdContext))
	cmd.AddCommand(config.NewCmd(cmdContext))

	// resources comands
	cmd.AddCommand(sshkeys.NewCmd(cmdContext))
	cmd.AddCommand(hosts.NewCmd(cmdContext))
	cmd.AddCommand(ssl.NewCmd(cmdContext))
	cmd.AddCommand(loadbalancers.NewCmd(cmdContext))
	cmd.AddCommand(racks.NewCmd(cmdContext))
	cmd.AddCommand(invoices.NewCmd(cmdContext))
	cmd.AddCommand(account.NewCmd(cmdContext))
	cmd.AddCommand(locations.NewCmd(cmdContext))
	cmd.AddCommand(k8s.NewCmd(cmdContext))
	cmd.AddCommand(uplinkmodels.NewCmd(cmdContext))
	cmd.AddCommand(uplinkbandwidths.NewCmd(cmdContext))
	cmd.AddCommand(servermodels.NewCmd(cmdContext))
	cmd.AddCommand(drivemodels.NewCmd(cmdContext))
	cmd.AddCommand(serverosoptions.NewCmd(cmdContext))
	cmd.AddCommand(serverramoptions.NewCmd(cmdContext))
	cmd.AddCommand(sbmosoptions.NewCmd(cmdContext))
	cmd.AddCommand(sbmmodels.NewCmd(cmdContext))
	cmd.AddCommand(l2segments.NewCmd(cmdContext))

	return cmd
}
