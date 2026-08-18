package metrics

import (
	"fmt"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

const (
	// rawOutput prints metrics as returned by the API
	rawOutput = "raw"

	// defaultPerPage limits the number of rows printed by default. All the metrics
	// come in a single response, so unlike the list commands there is no page size
	// coming from the API.
	defaultPerPage = 20
)

// addFlags adds flags supported by metrics commands
func addFlags(cmd *cobra.Command) {
	base.AddFormatFlags(cmd)

	// shadows the global output flag, as metrics support their own set of formats
	cmd.PersistentFlags().StringP("output", "o", "text", "output format (text/raw)")

	flags := cmd.Flags()
	flags.Int("per-page", defaultPerPage, "Number of items per page")
	flags.Int("page", 0, "Page number")
	flags.BoolP("all", "A", false, "Get all pages of resources")
}

// checkPaginationFlags rejects pagination flags with the raw output, as in that
// mode metrics are printed exactly as returned by the API
func checkPaginationFlags(cmdContext *base.CmdContext) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cmdContext.GetOrCreateFormatter(cmd).GetOutput() != rawOutput {
			return nil
		}
		for _, flag := range []string{"all", "page", "per-page"} {
			if cmd.Flags().Changed(flag) {
				return fmt.Errorf("--%s can't be used with the raw output", flag)
			}
		}
		return nil
	}
}

// printRaw prints metrics as returned by the API
func printRaw(cmd *cobra.Command, raw string) error {
	_, err := fmt.Fprint(cmd.OutOrStdout(), raw)
	return err
}

// paginate returns a page of rows according to the page and per-page flags.
// The whole set of metrics comes in a single response, so rows are paginated locally.
func paginate[T any](cmd *cobra.Command, rows []T) ([]T, error) {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return nil, err
	}
	if all {
		return rows, nil
	}

	perPage, err := cmd.Flags().GetInt("per-page")
	if err != nil {
		return nil, err
	}
	if perPage <= 0 {
		return rows, nil
	}

	page, err := cmd.Flags().GetInt("page")
	if err != nil {
		return nil, err
	}
	if page <= 0 {
		page = 1
	}

	start := (page - 1) * perPage
	if start >= len(rows) {
		return rows[:0], nil
	}

	return rows[start:min(start+perPage, len(rows))], nil
}
