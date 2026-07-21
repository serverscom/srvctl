package dns

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <domain-id>",
		Short: "Update a DNS domain",
		Long:  "Update a DNS domain by id (labels only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				log.Fatal(err)
			}
			input := serverscom.DNSDomainUpdateInput{
				Labels: labelsMap,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			domainID := args[0]
			domain, err := scClient.DNS.UpdateDomain(ctx, domainID, input)
			if err != nil {
				return err
			}

			if domain != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(domain)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}
