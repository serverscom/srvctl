package cmd

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/cmd/config"
	"github.com/serverscom/srvctl/cmd/context"
	"github.com/serverscom/srvctl/cmd/entities/account"
	"github.com/serverscom/srvctl/cmd/entities/hosts"
	"github.com/serverscom/srvctl/cmd/entities/invoices"
	"github.com/serverscom/srvctl/cmd/entities/k8s"
	loadbalancers "github.com/serverscom/srvctl/cmd/entities/load_balancers"
	"github.com/serverscom/srvctl/cmd/entities/locations"
	"github.com/serverscom/srvctl/cmd/entities/racks"
	sshkeys "github.com/serverscom/srvctl/cmd/entities/ssh-keys"
	"github.com/serverscom/srvctl/cmd/entities/ssl"
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

	return cmd
}
