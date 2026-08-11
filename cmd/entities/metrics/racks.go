package metrics

import (
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/metrics"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func newRacksCmd(cmdContext *base.CmdContext) *cobra.Command {
	rackMetricEntity, err := entities.Registry.GetEntityFromValue(metrics.RackMetric{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["racks"] = rackMetricEntity

	cmd := &cobra.Command{
		Use:   "racks",
		Short: "Get private racks metrics",
		Long: "Get private racks metrics: hosts count, monthly traffic and PDU/ATS power draw per rack.\n\n" +
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

			raw, err := scClient.Metrics.ListRacksMetrics(ctx)
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

			rows := metrics.BuildRackRows(samples)

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
