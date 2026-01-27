package cloudinstances

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	InputPath string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: "Add cloud instance",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.CloudComputingInstanceCreateInput{}
			if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), &input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			out, err := scClient.CloudComputingInstances.Create(ctx, input)
			if err != nil {
				return err
			}

			if out != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(out)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}
