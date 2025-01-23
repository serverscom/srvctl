package base

import (
	"fmt"
	"os"

	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

// CombinePreRunE combines multiple pre-run functions into one
func CombinePreRunE(funcs ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range funcs {
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

// InitCmdContext inits cmd context and sets up necessary dependencies
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
		cmdContext.formatter = output.NewFormatter(cmd, m)

		return nil
	}
}

// CheckFields checks if the field list is enabled and lists entity fields if so
func CheckFields(cmdContext *CmdContext, entity entities.EntityInterface) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if entity == nil {
			return fmt.Errorf("entity is not initialized")
		}
		manager := cmdContext.GetManager()
		formatter := cmdContext.GetOrCreateFormatter(cmd)
		output := formatter.GetOutput()
		if output == "json" || output == "yaml" {
			return nil
		}

		fieldList, err := manager.GetResolvedBoolValue(cmd, "field-list")
		if err != nil {
			return err
		}

		if fieldList {
			formatter.ListEntityFields(entity.GetAvailableFields())
			os.Exit(0)
		}

		fields, err := manager.GetResolvedStringSliceValue(cmd, "field")
		if err != nil {
			return err
		}
		if len(fields) > 0 {
			if err := entity.Validate(fields); err != nil {
				return err
			}
		}

		return nil
	}
}
