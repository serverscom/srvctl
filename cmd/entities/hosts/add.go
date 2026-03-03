package hosts

import (
	"fmt"
	"maps"
	"os"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddEBMFlags struct {
	Skeleton          bool
	InputPath         string
	LocationID        int
	ServerModelID     int
	OperatingSystemID int
	Features          []string
	RAMSize           int
	PublicUplinkID    int
	PublicBandwidthID int
	PrivateUplinkID   int
	DriveSlots        map[string]int
	Layout            []string
	Partitions        []string
	IPv6              bool
	UserDataFile      string
	UserData          string
	Labels            map[string]string
}

type AddSBMFlags struct {
	Skeleton          bool
	InputPath         string
	LocationID        int
	SBMFlavorModelID  int
	OperatingSystemID int
	SSHKeyFingerprint []string
	UserDataFile      string
	UserData          string
}

func (f *AddEBMFlags) FillInput(cmd *cobra.Command, input *serverscom.DedicatedServerCreateInput) error {
	pflags := cmd.Flags()
	if pflags.Changed("location-id") {
		input.LocationID = int64(f.LocationID)
	}
	if pflags.Changed("server-model-id") {
		input.ServerModelID = int64(f.ServerModelID)
	}
	if pflags.Changed("operating-system-id") {
		val := int64(f.OperatingSystemID)
		input.OperatingSystemID = &val
	}
	if pflags.Changed("feature") {
		input.Features = f.Features
	}
	if pflags.Changed("ram-size") {
		input.RAMSize = f.RAMSize
	}
	if pflags.Changed("public-uplink-id") {
		if input.UplinkModels.Public == nil {
			input.UplinkModels.Public = &serverscom.DedicatedServerPublicUplinkInput{}
		}
		input.UplinkModels.Public.ID = int64(f.PublicUplinkID)
	}
	if pflags.Changed("public-bandwidth-id") {
		if input.UplinkModels.Public == nil {
			input.UplinkModels.Public = &serverscom.DedicatedServerPublicUplinkInput{}
		}
		input.UplinkModels.Public.BandwidthModelID = int64(f.PublicBandwidthID)
	}
	if pflags.Changed("private-uplink-id") {
		input.UplinkModels.Private.ID = int64(f.PrivateUplinkID)
	}

	if pflags.Changed("drive-slots") {
		slots, err := parseDriveSlots(f.DriveSlots)
		if err != nil {
			return err
		}
		input.Drives.Slots = slots
	}

	if pflags.Changed("layout") {
		layouts, err := parseLayout(f.Layout)
		if err != nil {
			return err
		}
		input.Drives.Layout = mergeLayouts(input.Drives.Layout, layouts)
	}

	if pflags.Changed("partition") {
		if len(input.Drives.Layout) == 0 {
			return fmt.Errorf("partition given but layout is empty")
		}
		partitions, err := parsePartitions(f.Partitions)
		if err != nil {
			return err
		}
		err = applyPartitions(input.Drives.Layout, partitions)
		if err != nil {
			return err
		}
	}
	if pflags.Changed("ipv6") {
		input.IPv6 = f.IPv6
	}
	if pflags.Changed("user-data") && pflags.Changed("user-data-file") {
		return fmt.Errorf("'user-data' and 'user-data-file' can't be used together")
	}
	if pflags.Changed("user-data-file") {
		data, err := os.ReadFile(f.UserDataFile)
		if err != nil {
			return fmt.Errorf("can't read user-data-file: %v", err)
		}
		dataStr := string(data)
		input.UserData = &dataStr
	}
	if pflags.Changed("user-data") {
		input.UserData = &f.UserData
	}

	if pflags.Changed("labels") {
		for i := range input.Hosts {
			maps.Copy(input.Hosts[i].Labels, f.Labels)
		}
	}
	return nil
}

func (f *AddSBMFlags) FillInput(cmd *cobra.Command, input *serverscom.SBMServerCreateInput) error {
	pflags := cmd.Flags()
	if pflags.Changed("location-id") {
		input.LocationID = int64(f.LocationID)
	}
	if pflags.Changed("sbm-flavor-model-id") {
		input.FlavorModelID = int64(f.SBMFlavorModelID)
	}
	if pflags.Changed("operating-system-id") {
		val := int64(f.OperatingSystemID)
		input.OperatingSystemID = &val
	}
	if pflags.Changed("ssh-key-fingerprint") {
		input.SSHKeyFingerprints = f.SSHKeyFingerprint
	}
	if pflags.Changed("user-data") && pflags.Changed("user-data-file") {
		return fmt.Errorf("'user-data' and 'user-data-file' can't be used together")
	}
	if pflags.Changed("user-data-file") {
		data, err := os.ReadFile(f.UserDataFile)
		if err != nil {
			return fmt.Errorf("can't read user-data-file: %v", err)
		}
		dataStr := string(data)
		input.UserData = &dataStr
	}
	if pflags.Changed("user-data") {
		input.UserData = &f.UserData
	}
	return nil
}

func newAddEBMCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddEBMFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create an enterprise bare metal server",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("hosts/add_ebm.json")
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.DedicatedServerCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), &input); err != nil {
					return err
				}
			} else {
				required := []string{"location-id", "server-model-id", "private-uplink-id", "ram-size", "drive-slots", "layout"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if len(input.Hosts) == 0 && len(args) == 0 {
				return fmt.Errorf("no hosts found from positional args and no hosts found from input, can't continue")
			}

			for _, hostname := range args {
				input.Hosts = append(input.Hosts, serverscom.DedicatedServerHostInput{
					Hostname: hostname,
					Labels:   make(map[string]string),
				})
			}

			err := flags.FillInput(cmd, &input)
			if err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := scClient.Hosts.CreateDedicatedServers(ctx, input)
			if err != nil {
				return err
			}

			if server != nil {
				return formatter.Format(server)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	cmd.Flags().IntVar(&flags.LocationID, "location-id", 0, "Create the server(s) in the specific location ID")
	cmd.Flags().IntVar(&flags.ServerModelID, "server-model-id", 0, "Use specific server model ID to create the server")
	cmd.Flags().IntVar(&flags.OperatingSystemID, "operating-system-id", 0, "Install the specific operating system")
	cmd.Flags().StringSliceVar(&flags.Features, "feature", nil, "Set of features")
	cmd.Flags().IntVar(&flags.RAMSize, "ram-size", 0, "Desired amount of RAM in GB")
	cmd.Flags().IntVar(&flags.PublicUplinkID, "public-uplink-id", 0, "The public uplink ID, can be omitted if do not want public uplink")
	cmd.Flags().IntVar(&flags.PrivateUplinkID, "private-uplink-id", 0, "The private uplink ID")
	cmd.Flags().IntVar(&flags.PublicBandwidthID, "public-bandwidth-id", 0, "The public bandwidth ID, MUST be omitted if public uplink id is not passed")
	cmd.Flags().StringToIntVar(&flags.DriveSlots, "drive-slots", nil, "mapping of the specific slot to the specific drive model")
	cmd.Flags().StringArrayVar(&flags.Layout, "layout", nil, "Configuration of drives layout")
	cmd.Flags().StringArrayVar(&flags.Partitions, "partition", nil, "Configuration of the specific partitions")
	cmd.Flags().BoolVar(&flags.IPv6, "ipv6", false, "Enable IPv6")
	cmd.Flags().StringVar(&flags.UserDataFile, "user-data-file", "", "Path to user data which should be readed")
	cmd.Flags().StringVar(&flags.UserData, "user-data", "", "Content of user data")
	cmd.Flags().StringToStringVar(&flags.Labels, "labels", nil, "The set of labels which will be applied to the all hosts of this operation")

	return cmd
}

func newAddSBMCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddSBMFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create an SBM server",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("hosts/add_sbm.json")
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.SBMServerCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), &input); err != nil {
					return err
				}
			} else {
				required := []string{"location-id", "sbm-flavor-model-id", "operating-system-id"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if len(input.Hosts) == 0 && len(args) == 0 {
				return fmt.Errorf("no hosts found from positional args and no hosts found from input, can't continue")
			}

			for _, hostname := range args {
				input.Hosts = append(input.Hosts, serverscom.SBMServerHostInput{
					Hostname: hostname,
				})
			}

			if err := flags.FillInput(cmd, &input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := scClient.Hosts.CreateSBMServers(ctx, input)
			if err != nil {
				return err
			}

			if server != nil {
				return formatter.Format(server)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	cmd.Flags().IntVar(&flags.LocationID, "location-id", 0, "A unique identifier of a location")
	cmd.Flags().IntVar(&flags.SBMFlavorModelID, "sbm-flavor-model-id", 0, "A unique identifier of an SBM flavor")
	cmd.Flags().IntVar(&flags.OperatingSystemID, "operating-system-id", 0, "A unique identifier of an operating system")
	cmd.Flags().StringArrayVar(&flags.SSHKeyFingerprint, "ssh-key-fingerprint", nil, "Fingerprint of an SSH key to access the server")
	cmd.Flags().StringVar(&flags.UserDataFile, "user-data-file", "", "Path to user data which should be read")
	cmd.Flags().StringVar(&flags.UserData, "user-data", "", "Content of user data")

	return cmd
}
