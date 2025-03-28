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

type AllPager interface {
	AllPages() bool
}

func NewListOptions[T any](opts ...ListOptions[T]) []ListOptions[T] {
	return opts
}

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

// label selector option
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

// search pattern option
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

// location id option
type LocationIDOption[T any] struct {
	locationID string
}

func (o *LocationIDOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.locationID, "location-id", "", "Filter results by location ID")
}

func (o *LocationIDOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.locationID != "" {
		collection.SetParam("location_id", o.locationID)
	}
}

// cluster id option
type ClusterIDOption[T any] struct {
	clusterID string
}

func (o *ClusterIDOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.clusterID, "cluster-id", "", "Filter results by cluster ID")
}

func (o *ClusterIDOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.clusterID != "" {
		collection.SetParam("cluster_id", o.clusterID)
	}
}
