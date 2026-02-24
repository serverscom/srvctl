package ssl

import (
	"context"
	"fmt"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	Skeleton   bool
	InputPath  string
	Name       string
	PublicKey  string
	PrivateKey string
	ChainKey   string
	Labels     []string
}

type SSLCreator interface {
	Create(ctx context.Context, client *serverscom.Client, input any) (any, error)
	NewCreateInput() any
}

type SSLCustomCreateMgr struct{}

func (c *SSLCustomCreateMgr) Create(ctx context.Context, client *serverscom.Client, input any) (any, error) {
	sslInput, ok := input.(*serverscom.SSLCertificateCreateCustomInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for custom SSL")
	}
	return client.SSLCertificates.CreateCustom(ctx, *sslInput)
}

func (c *SSLCustomCreateMgr) NewCreateInput() any {
	return &serverscom.SSLCertificateCreateCustomInput{}
}

func newAddCmd(cmdContext *base.CmdContext, sslType *SSLTypeCmd) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: fmt.Sprintf("Create a %s", sslType.entityName),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("ssl/add.json")
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := sslType.managers.createMgr.NewCreateInput()

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				required := []string{"name", "public-key", "private-key"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			sslCert, err := sslType.managers.createMgr.Create(ctx, scClient, input)
			if err != nil {
				return err
			}

			if sslCert != nil {
				return formatter.Format(sslCert)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "A name of a SSL certificate")
	cmd.Flags().StringVarP(&flags.PublicKey, "public-key", "", "", "A public-key of a SSL certificate")
	cmd.Flags().StringVarP(&flags.PrivateKey, "private-key", "", "", "A private-key of a SSL certificate")
	cmd.Flags().StringVarP(&flags.ChainKey, "chain-key", "", "", "A chain-key of a SSL certificate")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func (f *AddedFlags) FillInput(cmd *cobra.Command, input any) error {
	sslInput, ok := input.(*serverscom.SSLCertificateCreateCustomInput)
	if !ok {
		return fmt.Errorf("invalid input type for custom SSL")
	}

	if cmd.Flags().Changed("name") {
		sslInput.Name = f.Name
	}
	if cmd.Flags().Changed("public-key") {
		sslInput.PublicKey = f.PublicKey
	}
	if cmd.Flags().Changed("private-key") {
		sslInput.PrivateKey = f.PrivateKey
	}
	if cmd.Flags().Changed("chain-key") {
		sslInput.ChainKey = f.ChainKey
	}
	if cmd.Flags().Changed("label") {
		labelsMap, err := base.ParseLabels(f.Labels)
		if err != nil {
			return err
		}

		sslInput.Labels = labelsMap
	}

	return nil
}
