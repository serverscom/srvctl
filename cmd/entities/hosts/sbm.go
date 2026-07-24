package hosts

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

// NewSBMCmd builds the top-level "sbm" command.
func NewSBMCmd(cmdContext *base.CmdContext) *cobra.Command {
	return newHostTypeCmd(cmdContext, HostTypeCmd{
		use:        "sbm",
		shortDesc:  "Manage scalable baremetal servers",
		entityName: "Scalable baremetal servers",
		typeFlag:   "sbm_server",
		managers: HostManagers{
			getMgr:       &SBMGetMgr{},
			powerMgr:     &SBMPowerMgr{},
			reinstallMgr: &SBMReinstallMgr{},
		},
		extraCmds: []func(*base.CmdContext) *cobra.Command{
			newAddSBMCmd,
			newUpdateSBMCmd,
			newSBMReleaseCmd,
			newListSBMCmd,
			newListSBMNetworksCmd,
			newGetSBMNetworkCmd,
			newAddSBMNetworkCmd,
			newDeleteSBMNetworkCmd,
			newGetSBMNetworkUsageCmd,
			newListSBMPTRCmd,
			newCreateSBMPTRCmd,
			newDeleteSBMPTRCmd,
		},
	})
}
