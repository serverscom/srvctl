package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type hostListOptions struct {
	base.BaseListOptions[serverscom.Host]
	rackID     string
	locationID string
}

func (o *hostListOptions) AddFlags(cmd *cobra.Command) {
	o.BaseListOptions.AddFlags(cmd)

	flags := cmd.Flags()
	flags.StringVar(&o.rackID, "rack-id", "", "Filter by rack ID")
	flags.StringVar(&o.locationID, "location-id", "", "Filter by location ID")
}

func (o *hostListOptions) ApplyToCollection(collection serverscom.Collection[serverscom.Host]) {
	o.BaseListOptions.ApplyToCollection(collection)

	if o.rackID != "" {
		collection.SetParam("rack_id", o.rackID)
	}
	if o.locationID != "" {
		collection.SetParam("location_id", o.locationID)
	}
}

func newListCmd(cmdContext *base.CmdContext, hostType *HostTypeCmd) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Host] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.Hosts.Collection()

		if hostType != nil && hostType.typeFlag != "" {
			collection = collection.SetParam("type", hostType.typeFlag)
		}

		return collection
	}

	opts := base.NewListOptions(
		&hostListOptions{},
		&base.LabelSelectorOption[serverscom.Host]{},
		&base.SearchPatternOption[serverscom.Host]{},
	)

	entityName := "Hosts"
	if hostType != nil {
		entityName = hostType.entityName
	}

	return base.NewListCmd("list", entityName, factory, cmdContext, opts...)
}
