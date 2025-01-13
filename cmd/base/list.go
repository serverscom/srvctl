package base

import (
	"strings"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/spf13/cobra"
)

// CollectionFactory is a function type that creates a typed resource collection
// with configurable verbosity level
type CollectionFactory[T any] func(verbose bool) serverscom.Collection[T]

type listOptions struct {
	labelSelector string
	perPage       int
	page          int
	sorting       string
	direction     string
	allPages      bool
}

// NewListCmd base list command for different collections
func NewListCmd[T any](entityName string, colFactory CollectionFactory[T], cmdContext *CmdContext) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List " + entityName,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := SetupContext(cmd, manager)
			defer cancel()

			SetupProxy(cmd, manager)

			collection := colFactory(manager.GetVerbose(cmd))
			if opts.labelSelector != "" {
				collection.SetParam("label_selector", opts.labelSelector)
			}
			if opts.sorting != "" {
				collection.SetParam("sort", opts.sorting)
			}
			if opts.direction != "" {
				collection.SetParam("direction", strings.ToUpper(opts.direction))
			}

			var items []T
			var err error
			if opts.allPages {
				items, err = collection.Collect(ctx)
			} else {
				if opts.perPage > 0 {
					collection.SetPerPage(opts.perPage)
				}
				if opts.page > 0 {
					collection.SetPage(opts.page)
				}
				items, err = collection.List(ctx)
			}

			if err != nil {
				return err
			}

			outputFormat, _ := manager.GetResolvedStringValue(cmd, "output")
			formatter := output.NewFormatter(cmd.OutOrStdout())
			return formatter.FormatList(items, outputFormat)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.labelSelector, "label-selector", "", "Filter by label selector")
	flags.IntVar(&opts.perPage, "per-page", 0, "Number of items per page")
	flags.IntVar(&opts.page, "page", 0, "Page number")
	flags.StringVar(&opts.sorting, "sorting", "", "Sort field")
	flags.StringVar(&opts.direction, "direction", "", "Sort direction (ASC or DESC)")
	flags.BoolVarP(&opts.allPages, "all", "A", false, "Get all pages of resources")

	return cmd
}
