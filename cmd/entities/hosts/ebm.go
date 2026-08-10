package hosts

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

// NewEBMCmd builds the top-level "ebm" command.
func NewEBMCmd(cmdContext *base.CmdContext) *cobra.Command {
	return newHostTypeCmd(cmdContext, HostTypeCmd{
		use:        "ebm",
		shortDesc:  "Manage enterprise bare metal servers",
		entityName: "Enterprise Bare Metal",
		typeFlag:   "dedicated_server",
		managers: HostManagers{
			getMgr:       &EBMGetMgr{},
			powerMgr:     &EBMPowerMgr{},
			reinstallMgr: &EBMReinstallMgr{},
		},
		extraCmds: []func(*base.CmdContext) *cobra.Command{
			newAddEBMCmd,
			newUpdateEBMCmd,
			newListEBMDriveSlotsCmd,
			newListEBMConnectionsCmd,
			newListEBMPTRCmd,
			newCreateEBMPTRCmd,
			newDeleteEBMPTRCmd,
			newEBMAbortReleaseCmd,
			newEBMScheduleReleaseCmd,
			newListEBMNetworksCmd,
			newGetEBMNetworkCmd,
			newAddEBMNetworkCmd,
			newDeleteEBMNetworkCmd,
			newActivateEBMIPv6NetworkCmd,
			newGetEBMNetworkUsageCmd,
			newListEBMCmd,
			newListEBMServicesCmd,
			newListEBMFeaturesCmd,
			newEBMFeatureSetCmd,
			newGetEBMOOBCredsCmd,
		},
	})
}
