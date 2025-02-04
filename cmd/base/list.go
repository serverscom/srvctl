package base

import (
	"strings"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/spf13/cobra"
)

type ListOptions[T any] interface {
	AddFlags(*cobra.Command)
	ApplyToCollection(serverscom.Collection[T], *cobra.Command)
	AllPages() bool
}

// CollectionFactory is a function type that creates a typed resource collection
// with configurable verbosity level
type CollectionFactory[T any] func(verbose bool) serverscom.Collection[T]

// BaseListOptions is a base options struct for list commands
type BaseListOptions[T any] struct {
	labelSelector string
	perPage       int
	page          int
	sorting       string
	direction     string
	allPages      bool
}

// AddFlags adds common list flags to the command
func (o *BaseListOptions[T]) AddFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&o.labelSelector, "label-selector", "", "Filter by label selector")
	flags.IntVar(&o.perPage, "per-page", 0, "Number of items per page")
	flags.IntVar(&o.page, "page", 0, "Page number")
	flags.StringVar(&o.sorting, "sorting", "", "Sort field")
	flags.StringVar(&o.direction, "direction", "", "Sort direction (ASC or DESC)")
	flags.BoolVarP(&o.allPages, "all", "A", false, "Get all pages of resources")
}

// ApplyToCollection applies the options to a collection
func (o *BaseListOptions[T]) ApplyToCollection(collection serverscom.Collection[T], cmd *cobra.Command) {
	if o.labelSelector != "" {
		collection.SetParam("label_selector", o.labelSelector)
	}
	if o.sorting != "" {
		collection.SetParam("sort", o.sorting)
	}
	if o.direction != "" {
		collection.SetParam("direction", strings.ToUpper(o.direction))
	}
	if o.perPage > 0 {
		collection.SetPerPage(o.perPage)
	}
	if o.page > 0 {
		collection.SetPage(o.page)
	}
}

// AllPages returns true if all pages should be retrieved
func (o *BaseListOptions[T]) AllPages() bool {
	return o.allPages
}

// NewListCmd base list command for different collections
func NewListCmd[T any](entityName string, colFactory CollectionFactory[T], cmdContext *CmdContext, opts ListOptions[T]) *cobra.Command {
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
			opts.ApplyToCollection(collection, cmd)

			var items []T
			var err error
			if opts.AllPages() {
				items, err = collection.Collect(ctx)
			} else {
				items, err = collection.List(ctx)
			}

			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(items)
		},
	}

	opts.AddFlags(cmd)

	return cmd
}
