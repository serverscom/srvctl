package hosts

import (
	"fmt"
	"log"
	"os"

	"maps"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type AddDSFlags struct {
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

func applyFlagsToInput(
	input *serverscom.DedicatedServerCreateInput,
	flags *AddDSFlags,
	pflags *pflag.FlagSet,
) error {
	if pflags.Changed("location-id") {
		input.LocationID = int64(flags.LocationID)
	}
	if pflags.Changed("server-model-id") {
		input.ServerModelID = int64(flags.ServerModelID)
	}
	if pflags.Changed("operating-system-id") {
		val := int64(flags.OperatingSystemID)
		input.OperatingSystemID = &val
	}
	if pflags.Changed("feature") {
		input.Features = flags.Features
	}
	if pflags.Changed("ram-size") {
		input.RAMSize = flags.RAMSize
	}
	if pflags.Changed("public-uplink-id") {
		if input.UplinkModels.Public == nil {
			input.UplinkModels.Public = &serverscom.DedicatedServerPublicUplinkInput{}
		}
		input.UplinkModels.Public.ID = int64(flags.PublicUplinkID)
	}
	if pflags.Changed("public-bandwidth-id") {
		if input.UplinkModels.Public == nil {
			input.UplinkModels.Public = &serverscom.DedicatedServerPublicUplinkInput{}
		}
		input.UplinkModels.Public.BandwidthModelID = int64(flags.PublicBandwidthID)
	}
	if pflags.Changed("private-uplink-id") {
		input.UplinkModels.Private.ID = int64(flags.PrivateUplinkID)
	}

	if pflags.Changed("drive-slots") {
		slots, err := parseDriveSlots(flags.DriveSlots)
		if err != nil {
			return err
		}
		input.Drives.Slots = slots
	}

	if pflags.Changed("layout") {
		layouts, err := parseLayout(flags.Layout)
		if err != nil {
			return err
		}
		input.Drives.Layout = mergeLayouts(input.Drives.Layout, layouts)
	}

	if pflags.Changed("partition") {
		if len(input.Drives.Layout) == 0 {
			return fmt.Errorf("partition given but layout is empty")
		}
		partitions, err := parsePartitions(flags.Partitions)
		if err != nil {
			return err
		}
		err = applyPartitions(input.Drives.Layout, partitions)
		if err != nil {
			return err
		}
	}
	if pflags.Changed("ipv6") {
		input.IPv6 = flags.IPv6
	}
	if pflags.Changed("user-data") && pflags.Changed("user-data-file") {
		return fmt.Errorf("'user-data' and 'user-data-file' can't be used together")
	}
	if pflags.Changed("user-data-file") {
		data, err := os.ReadFile(flags.UserDataFile)
		if err != nil {
			return fmt.Errorf("can't read user-data-file: %v", err)
		}
		dataStr := string(data)
		input.UserData = &dataStr
	}
	if pflags.Changed("user-data") {
		input.UserData = &flags.UserData
	}

	if pflags.Changed("labels") {
		for i := range input.Hosts {
			maps.Copy(input.Hosts[i].Labels, flags.Labels)
		}
	}
	return nil
}

func newAddDSCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddDSFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create a dedicated server",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.DedicatedServerCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), &input); err != nil {
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

			err := applyFlagsToInput(&input, flags, cmd.Flags())
			if err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := scClient.Hosts.CreateDedicatedServers(ctx, input)
			if err != nil {
				return err
			}

			if server != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(server)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")

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
	var path string
	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: "Create an SBM server",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.SBMServerCreateInput{}

			if err := base.ReadInputJSON(path, cmd.InOrStdin(), &input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := scClient.Hosts.CreateSBMServers(ctx, input)
			if err != nil {
				return err
			}

			if server != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(server)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
