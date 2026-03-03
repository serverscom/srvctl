package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListEBMPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.PTRRecord] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerPTRRecords(args[0])
	}

	opts := &base.BaseListOptions[serverscom.PTRRecord]{}

	cmd := base.NewListCmd("list-ptr", "Enterprise bare metal server PTR records", factory, cmdContext, opts)
	cmd.Use = "list-ptr <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func newCreateEBMPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	var (
		ip       string
		domain   string
		ttl      int
		priority int
	)

	cmd := &cobra.Command{
		Use:   "add-ptr <server_id>",
		Short: "Create a PTR record",
		Long:  "Create a PTR record for an enterprise bare metal server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			serverID := args[0]

			input := &serverscom.PTRRecordCreateInput{
				IP:     ip,
				Domain: domain,
			}
			if cmd.Flags().Changed("ttl") {
				input.TTL = &ttl
			}
			if cmd.Flags().Changed("priority") {
				input.Priority = &priority
			}

			prtRecord, err := scClient.Hosts.CreatePTRRecordForDedicatedServer(ctx, serverID, *input)
			if err != nil {
				return err
			}

			if prtRecord != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(prtRecord)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&ip, "ip", "", "An IP address associated with a PTR record (required)")
	_ = cmd.MarkFlagRequired("ip")

	cmd.Flags().StringVar(&domain, "domain", "", "A domain name for a PTR record (required)")
	_ = cmd.MarkFlagRequired("domain")

	cmd.Flags().IntVar(&ttl, "ttl", 0, "TTL (time to live) in seconds")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority (lower value means higher priority)")

	return cmd
}

func newDeleteEBMPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	var recordID string

	cmd := &cobra.Command{
		Use:   "delete-ptr <server_id>",
		Short: "Delete a PTR record",
		Long:  "Delete a PTR record for an enterprise bare metal server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			serverID := args[0]
			return scClient.Hosts.DeletePTRRecordForDedicatedServer(ctx, serverID, recordID)
		},
	}

	cmd.Flags().StringVar(&recordID, "ptr-id", "", "Record ID (required)")
	_ = cmd.MarkFlagRequired("ptr-id")

	return cmd
}

func newListSBMPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.PTRRecord] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.SBMServerPTRRecords(args[0])
	}

	opts := &base.BaseListOptions[serverscom.PTRRecord]{}

	cmd := base.NewListCmd("list-ptr", "SBM server PTR records", factory, cmdContext, opts)
	cmd.Use = "list-ptr <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func newCreateSBMPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	var (
		ip       string
		domain   string
		ttl      int
		priority int
	)

	cmd := &cobra.Command{
		Use:   "add-ptr <server_id>",
		Short: "Create a PTR record",
		Long:  "Create a PTR record for an SBM server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			serverID := args[0]

			input := &serverscom.PTRRecordCreateInput{
				IP:     ip,
				Domain: domain,
			}
			if cmd.Flags().Changed("ttl") {
				input.TTL = &ttl
			}
			if cmd.Flags().Changed("priority") {
				input.Priority = &priority
			}

			prtRecord, err := scClient.Hosts.CreatePTRRecordForSBMServer(ctx, serverID, *input)
			if err != nil {
				return err
			}

			if prtRecord != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(prtRecord)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&ip, "ip", "", "An IP address associated with a PTR record (required)")
	_ = cmd.MarkFlagRequired("ip")

	cmd.Flags().StringVar(&domain, "domain", "", "A domain name for a PTR record (required)")
	_ = cmd.MarkFlagRequired("domain")

	cmd.Flags().IntVar(&ttl, "ttl", 0, "TTL (time to live) in seconds")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority (lower value means higher priority)")

	return cmd
}

func newDeleteSBMPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	var recordID string

	cmd := &cobra.Command{
		Use:   "delete-ptr <server_id>",
		Short: "Delete a PTR record",
		Long:  "Delete a PTR record for an SBM server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			serverID := args[0]
			return scClient.Hosts.DeletePTRRecordForSBMServer(ctx, serverID, recordID)
		},
	}

	cmd.Flags().StringVar(&recordID, "ptr-id", "", "Record ID (required)")
	_ = cmd.MarkFlagRequired("ptr-id")

	return cmd
}
