package config

import (
	"fmt"
	"slices"

	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/output"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newFinalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "final",
		Short: "Show final configuration",
		Args:  cobra.ExactArgs(0),
		Long: `Show final configuration after merging all configuration levels:
- Global configurations
- Context-level configurations
- CLI-level arguments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create config manager: %w", err)
			}

			currentContext := manager.GetDefaultContextName()
			if cmd.Flags().Changed("context") {
				currentContext, err = cmd.Flags().GetString("context")
				if err != nil {
					return fmt.Errorf("can't get context from flag")
				}
			}
			ctx, err := manager.GetContext(currentContext)
			if err != nil {
				return fmt.Errorf("failed to get context %q: %w", currentContext, err)
			}

			finalConfig := buildFinalConfig(cmd, manager)

			result := output.ContextInfo{
				Context:  currentContext,
				Endpoint: ctx.Endpoint,
				Config:   finalConfig,
			}

			outputFormat, _ := manager.GetResolvedStringValue(cmd, "output")
			return output.Format(result, outputFormat)
		},
	}

	return cmd
}

func buildFinalConfig(cmd *cobra.Command, manager *config.Manager) map[string]interface{} {
	finalConfig := make(map[string]interface{})

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !slices.Contains(KnownConfigFlags, f.Name) {
			return
		}
		switch f.Value.Type() {
		case "bool":
			finalConfig[f.Name], _ = manager.GetResolvedBoolValue(cmd, f.Name)
		case "int":
			finalConfig[f.Name], _ = manager.GetResolvedIntValue(cmd, f.Name)
		default:
			finalConfig[f.Name], _ = manager.GetResolvedStringValue(cmd, f.Name)
		}
	})

	return finalConfig
}
