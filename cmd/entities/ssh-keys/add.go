package sshkeys

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	Skeleton  bool
	InputPath string
	Name      string
	PublicKey string
	Labels    []string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an ssh key",
		Long:  "Add a new SSH key to account",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.SSHKeyCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else if flags.Skeleton {
				formatter.SetOutput("json")
				return formatter.FormatSkeleton("ssh-keys/add.json")
			} else {
				required := []string{"name", "public-key"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			sshKey, err := scClient.SSHKeys.Create(ctx, *input)
			if err != nil {
				return err
			}

			if sshKey != nil {
				return formatter.Format(sshKey)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "A name of a SSH key")
	cmd.Flags().StringVarP(&flags.PublicKey, "public-key", "", "", "A public-key of a SSH key")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func (f *AddedFlags) FillInput(cmd *cobra.Command, input *serverscom.SSHKeyCreateInput) error {
	if cmd.Flags().Changed("name") {
		input.Name = f.Name
	}
	if cmd.Flags().Changed("public-key") {
		input.PublicKey = f.PublicKey
	}
	if cmd.Flags().Changed("label") {
		labelsMap, err := base.ParseLabels(f.Labels)
		if err != nil {
			return err
		}

		input.Labels = labelsMap
	}

	return nil
}
