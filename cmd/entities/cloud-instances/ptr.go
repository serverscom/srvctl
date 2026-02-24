package cloudinstances

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newAddPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	var data string
	var ip string
	var ttl int
	var priority int

	cmd := &cobra.Command{
		Use:   "add-ptr <instance-id>",
		Short: "Add PTR record to cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.PTRRecordCreateInput{
				Domain: data,
				IP:     ip,
			}
			if cmd.Flags().Changed("ttl") {
				input.TTL = &ttl
			}
			if cmd.Flags().Changed("priority") {
				input.Priority = &priority
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.CreatePTRRecord(ctx, instanceID, input)
			if err != nil {
				return err
			}

			if out != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(out)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "PTR record data")
	_ = cmd.MarkFlagRequired("data")

	cmd.Flags().StringVar(&ip, "ip", "", "PTR record IP")
	_ = cmd.MarkFlagRequired("ip")

	cmd.Flags().IntVar(&ttl, "ttl", 0, "TTL for PTR record")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority for PTR record")

	return cmd
}

func newListPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.PTRRecord] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		instanceID := args[0]
		collection := scClient.CloudComputingInstances.PTRRecords(instanceID)

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.PTRRecord]{},
	)

	return base.NewListCmd("list-ptr", "PTR records", factory, cmdContext, opts...)
}

func newDeletePTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	var ptrID string

	cmd := &cobra.Command{
		Use:   "delete-ptr <instance-id>",
		Short: "Delete PTR record of cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			return scClient.CloudComputingInstances.DeletePTRRecord(ctx, instanceID, ptrID)
		},
	}

	cmd.Flags().StringVar(&ptrID, "ptr-id", "", "PTR record ID")
	if err := cmd.MarkFlagRequired("ptr-id"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
