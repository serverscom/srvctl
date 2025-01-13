package base

import (
	"fmt"

	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
)

// InitCmdContext creates cmd
func InitCmdContext(cmdContext *CmdContext) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		m, err := config.NewManager(configPath)
		if err != nil {
			return fmt.Errorf("failed to initialize config manager: %w", err)
		}

		c := client.NewClient(
			m.GetToken(),
			m.GetEndpoint(),
		)
		version := cmd.Root().Version
		c.SetUserAgent(userAgent(version))

		cmdContext.manager = m
		cmdContext.client = c

		return nil
	}
}
