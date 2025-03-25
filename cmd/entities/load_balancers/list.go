package loadbalancers

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext, lbType *LBTypeCmd) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.LoadBalancer] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.LoadBalancers.Collection()

		if lbType != nil && lbType.typeFlag != "" {
			collection = collection.SetParam("type", lbType.typeFlag)
		}

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.LoadBalancer]{},
		&base.LabelSelectorOption[serverscom.LoadBalancer]{},
		&base.SearchPatternOption[serverscom.LoadBalancer]{},
		// TODO location_id, cluster_id
	)

	entityName := "load balancers"
	if lbType != nil {
		entityName = lbType.entityName
	}

	return base.NewListCmd("list", entityName, factory, cmdContext, opts...)
}
