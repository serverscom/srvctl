package dns

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	Skeleton  bool
	InputPath string
	Name      string
	Email     string
	TTL       int
	Labels    []string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a DNS domain",
		Long:  "Create/Register a new DNS domain/zone",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("dns/add.json")
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.DNSDomainCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				required := []string{"name", "email"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			domain, err := scClient.DNS.CreateDomain(ctx, *input)
			if err != nil {
				return err
			}

			if domain != nil {
				return formatter.Format(domain)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "A name of a DNS domain")
	cmd.Flags().StringVarP(&flags.Email, "email", "", "", "An email used for the zone's SOA record")
	cmd.Flags().IntVarP(&flags.TTL, "ttl", "", 0, "A default TTL in seconds (0-604800)")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func (f *AddedFlags) FillInput(cmd *cobra.Command, input *serverscom.DNSDomainCreateInput) error {
	if cmd.Flags().Changed("name") {
		input.Name = f.Name
	}
	if cmd.Flags().Changed("email") {
		input.Email = f.Email
	}
	if cmd.Flags().Changed("ttl") {
		input.TTL = f.TTL
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
