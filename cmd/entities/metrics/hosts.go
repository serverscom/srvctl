package metrics

import (
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/metrics"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func newHostsCmd(cmdContext *base.CmdContext) *cobra.Command {
	hostMetricEntity, err := entities.Registry.GetEntityFromValue(metrics.HostMetric{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["hosts"] = hostMetricEntity

	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Get hosts metrics",
		Long: "Get hosts metrics: monthly traffic per host.\n\n" +
			"Use --output raw to get metrics in the Prometheus text exposition format as returned by the API.",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlagsWithOutputs(cmdContext, entitiesMap, []string{rawOutput}),
			checkPaginationFlags(cmdContext),
		),
		Args: base.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			raw, err := scClient.Metrics.ListHostsMetrics(ctx)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			if formatter.GetOutput() == rawOutput {
				return printRaw(cmd, raw)
			}

			samples, err := metrics.Parse(raw)
			if err != nil {
				return err
			}

			rows := metrics.BuildHostRows(samples)

			rows, err = paginate(cmd, rows)
			if err != nil {
				return err
			}

			return formatter.Format(rows)
		},
	}

	addFlags(cmd)

	return cmd
}
