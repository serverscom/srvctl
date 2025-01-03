package cmd

import (
	"github.com/serverscom/srvctl/cmd/config"
	"github.com/serverscom/srvctl/cmd/context"
	"github.com/serverscom/srvctl/cmd/ssh"
	"github.com/spf13/cobra"
)

var verbose bool

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "srvctl",
		Short: "CLI tool for servers.com API",
		Long:  `A command line interface for managing servers.com resources`,
	}

	// Global flags
	cmd.PersistentFlags().String("config", "", "config file path")
	cmd.PersistentFlags().String("context", "", "context name")
	cmd.PersistentFlags().String("proxy", "", "proxy url")
	cmd.PersistentFlags().Int("http-timeout", 30, "HTTP timeout ( seconds )")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringP("output", "o", "text", "output format (text/json/yaml)")

	// Add commands
	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(context.NewCmd())
	cmd.AddCommand(config.NewCmd())
	cmd.AddCommand(ssh.NewCmd())

	return cmd
}
