package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			configOptions := make(map[string]any)

			cmd.Flags().Visit(func(f *pflag.Flag) {
				fillConfigOptions(cmd, f, configOptions)
			})

			// Update configuration based on scope
			scope := cmd.Parent().Name()

			if scope == "global" {
				manager.UpdateGlobalConfig(configOptions)
			} else {
				contextName := manager.GetDefaultContextName()
				ctx, err := cmd.Flags().GetString("context")
				if err != nil {
					return err
				}
				if ctx != "" {
					contextName = ctx
				}
				if err := manager.UpdateContextConfig(contextName, configOptions); err != nil {
					return fmt.Errorf("failed to update context config: %w", err)
				}
			}

			return manager.Save()
		},
	}

	addDisableFlags(cmd)

	return cmd
}

// fillConfigOptions adds cmd flag to configOptions only if flag exists in KnownConfigFlags
func fillConfigOptions(cmd *cobra.Command, f *pflag.Flag, configOptions map[string]any) {
	if strings.HasPrefix(f.Name, "no-") {
		optionName := strings.TrimPrefix(f.Name, "no-")
		configOptions[optionName] = nil
		return
	}

	if !slices.Contains(KnownConfigFlags, f.Name) {
		return
	}

	switch f.Value.Type() {
	case "bool":
		if v, err := cmd.Flags().GetBool(f.Name); err == nil {
			configOptions[f.Name] = v
		}
	case "int":
		if v, err := cmd.Flags().GetInt(f.Name); err == nil {
			configOptions[f.Name] = v
		}
	default:
		configOptions[f.Name] = f.Value.String()
	}
}

func addDisableFlags(cmd *cobra.Command) {
	for _, f := range KnownConfigFlags {
		flagName := fmt.Sprintf("no-%s", f)
		cmd.PersistentFlags().Bool(flagName, false, fmt.Sprintf("Disable %s configuration", f))
	}
}
