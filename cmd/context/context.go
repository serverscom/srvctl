package context

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage contexts",
		Long:  `Manage authentication contexts for different API accounts`,
	}

	cmd.AddCommand(
		newListCmd(),
		newUpdateCmd(),
		newDeleteCmd(),
	)

	return cmd
}
