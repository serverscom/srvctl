package ssl

import (
	"context"
	"fmt"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"strings"
)

var (
	path       string
	name       string
	publicKey  string
	privateKey string
	chainKey   string
	labels     []string
)

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
	cmd := &cobra.Command{
		Use:   "add",
		Short: fmt.Sprintf("Create a %s", sslType.entityName),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := sslType.managers.createMgr.NewCreateInput()

			if cmd.Flags().Changed("input") {
				if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				if err := validateCustomSSLFlags(cmd); err != nil {
					return err
				}
				if err := fillCustomSSLInput(cmd, input); err != nil {
					return err
				}
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			sslCert, err := sslType.managers.createMgr.Create(ctx, scClient, input)
			if err != nil {
				return err
			}

			if sslCert != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(sslCert)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().StringVarP(&name, "name", "n", "", "A name of a SSL certificate")
	cmd.Flags().StringVarP(&publicKey, "public-key", "", "", "A public-key of a SSL certificate")
	cmd.Flags().StringVarP(&privateKey, "private-key", "", "", "A private-key of a SSL certificate")
	cmd.Flags().StringVarP(&chainKey, "chain-key", "", "", "A chain-key of a SSL certificate")
	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func validateCustomSSLFlags(cmd *cobra.Command) error {
	required := []string{"name", "public-key", "private-key", "chain-key"}
	var missing []string

	for _, flag := range required {
		if !cmd.Flags().Changed(flag) {
			missing = append(missing, "--"+flag)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"use --input or provide all required flags (missing: %s)",
			strings.Join(missing, ", "),
		)
	}

	return nil
}

func fillCustomSSLInput(cmd *cobra.Command, input any) error {
	sslInput, ok := input.(*serverscom.SSLCertificateCreateCustomInput)
	if !ok {
		return fmt.Errorf("invalid input type for custom SSL")
	}

	if cmd.Flags().Changed("name") {
		sslInput.Name = name
	}
	if cmd.Flags().Changed("public-key") {
		sslInput.PublicKey = publicKey
	}
	if cmd.Flags().Changed("private-key") {
		sslInput.PrivateKey = privateKey
	}
	if cmd.Flags().Changed("chain-key") {
		sslInput.ChainKey = chainKey
	}
	if cmd.Flags().Changed("label") {
		labelsMap, err := base.ParseLabels(labels)
		if err != nil {
			return err
		}

		sslInput.Labels = labelsMap
	}

	return nil
}
