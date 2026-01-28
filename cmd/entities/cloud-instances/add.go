package cloudinstances

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	Skeleton  bool
	InputPath string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add cloud instance",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("cloud-instances/add.json")
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.CloudComputingInstanceCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), &input); err != nil {
					return err
				}
			} else {
				required := []string{"input"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			out, err := scClient.CloudComputingInstances.Create(ctx, input)
			if err != nil {
				return err
			}

			if out != nil {
				return formatter.Format(out)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	return cmd
}
