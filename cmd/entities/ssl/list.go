package ssl

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext, sslType *SSLTypeCmd) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.SSLCertificate] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.SSLCertificates.Collection()

		if sslType != nil && sslType.typeFlag != "" {
			collection = collection.SetParam("type", sslType.typeFlag)
		}

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.SSLCertificate]{},
		&base.LabelSelectorOption[serverscom.SSLCertificate]{},
		&base.SearchPatternOption[serverscom.SSLCertificate]{},
	)

	entityName := "SSL certificates"
	if sslType != nil {
		entityName = sslType.entityName
	}

	return base.NewListCmd("list", entityName, factory, cmdContext, opts...)
}
