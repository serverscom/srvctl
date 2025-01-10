package cmd

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/cmd/config"
	"github.com/serverscom/srvctl/cmd/context"
	sshkeys "github.com/serverscom/srvctl/cmd/entities/ssh-keys"
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

	return cmd
}
