package base

import (
	"log"
	"strings"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/spf13/cobra"
)

type ListOptions[T any] interface {
	AddFlags(*cobra.Command)
	ApplyToCollection(serverscom.Collection[T])
}

func NewListOptions[T any](opts ...ListOptions[T]) []ListOptions[T] {
	return opts
}

// CollectionFactory is a function type that creates a typed resource collection
// with configurable verbosity level
// type CollectionFactory[T any] func(verbose bool) serverscom.Collection[T]
type CollectionFactory[T any] func(verbose bool, args ...string) serverscom.Collection[T]

// BaseListOptions is a base options struct for list commands
type BaseListOptions[T any] struct {
	perPage   int
	page      int
	sorting   string
	direction string
	allPages  bool
}

// AddFlags adds common list flags to the command
func (o *BaseListOptions[T]) AddFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVar(&o.perPage, "per-page", 0, "Number of items per page")
	flags.IntVar(&o.page, "page", 0, "Page number")
	flags.StringVar(&o.sorting, "sorting", "", "Sort field")
	flags.StringVar(&o.direction, "direction", "", "Sort direction (ASC or DESC)")
	flags.BoolVarP(&o.allPages, "all", "A", false, "Get all pages of resources")

	flags.String("type", "", "")
	if err := flags.MarkHidden("type"); err != nil {
		log.Fatal(err)
	}
}

// ApplyToCollection applies the options to a collection
func (o *BaseListOptions[T]) ApplyToCollection(collection serverscom.Collection[T]) {

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

type LabelSelectorOption[T any] struct {
	labelSelector string
}

func (o *LabelSelectorOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.labelSelector, "label-selector", "", "Filter results by labels")
}

func (o *LabelSelectorOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.labelSelector != "" {
		collection.SetParam("label_selector", o.labelSelector)
	}
}

type SearchPatternOption[T any] struct {
	searchPattern string
}

func (o *SearchPatternOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.searchPattern, "search-pattern", "", "Return resources containing the parameter value in its name")
}

func (o *SearchPatternOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.searchPattern != "" {
		collection.SetParam("search_pattern", o.searchPattern)
	}
}

// NewListCmd base list command for different collections
func NewListCmd[T any](use string, entityName string, colFactory CollectionFactory[T], cmdContext *CmdContext, opts ...ListOptions[T]) *cobra.Command {
	aliases := []string{}
	if use == "list" {
		aliases = append(aliases, "ls")
	}
	cmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   "List " + entityName,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := SetupContext(cmd, manager)
			defer cancel()

			SetupProxy(cmd, manager)

			collection := colFactory(manager.GetVerbose(cmd), args...)
			for _, opt := range opts {
				opt.ApplyToCollection(collection)
			}

			var items []T
			var err error
			for _, opt := range opts {
				if baseOpts, ok := opt.(*BaseListOptions[T]); ok && baseOpts.AllPages() {
					items, err = collection.Collect(ctx)
					break
				}
			}

			if items == nil && err == nil {
				items, err = collection.List(ctx)
			}

			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(items)
		},
	}

	for _, opt := range opts {
		opt.AddFlags(cmd)
	}

	return cmd
}
