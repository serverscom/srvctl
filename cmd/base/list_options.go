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

// status option
type StatusOption[T any] struct {
	status string
}

func (o *StatusOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.status, "status", "", "Filter results by status")
}

func (o *StatusOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.status != "" {
		collection.SetParam("status", o.status)
	}
}

// invoice type option
type InvoiceTypeOption[T any] struct {
	typeVal string
}

// use itype instead of type to avoid conflict with baseList 'type'  hidden flag which we use for subcommands
func (o *InvoiceTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.typeVal, "itype", "", "Filter results by type")
}

func (o *InvoiceTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.typeVal != "" {
		collection.SetParam("type", o.typeVal)
	}
}

// parent id option
type ParentIDOption[T any] struct {
	parentID string
}

func (o *ParentIDOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.parentID, "parent-id", "", "Filter results by parent ID")
}

func (o *ParentIDOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.parentID != "" {
		collection.SetParam("parent_id", o.parentID)
	}
}

// currency option
type CurrencyOption[T any] struct {
	currency string
}

func (o *CurrencyOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.currency, "currency", "", "Filter results by currency")
}

func (o *CurrencyOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.currency != "" {
		collection.SetParam("currency", o.currency)
	}
}

// start date option
type StartDateOption[T any] struct {
	startDate string
}

func (o *StartDateOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.startDate, "start-date", "", "Filter results by start date")
}

func (o *StartDateOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.startDate != "" {
		collection.SetParam("start_date", o.startDate)
	}
}

// end date option
type EndDateOption[T any] struct {
	endDate string
}

func (o *EndDateOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.endDate, "end-date", "", "Filter results by end date")
}

func (o *EndDateOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.endDate != "" {
		collection.SetParam("end_date", o.endDate)
	}
}

type FamilyOption[T any] struct {
	family string
}

func (o *FamilyOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.family, "family", "", "Set to 'ipv4' or 'ipv6'")
}

func (o *FamilyOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.family != "" {
		collection.SetParam("family", o.family)
	}
}

type InterfaceTypeOption[T any] struct {
	interfaceType string
}

func (o *InterfaceTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.interfaceType, "interface-type", "", "Type of network interface: public, private, oob")
}

func (o *InterfaceTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.interfaceType != "" {
		collection.SetParam("interface_type", o.interfaceType)
	}
}

type DistributionMethodOption[T any] struct {
	distributionMethod string
}

func (o *DistributionMethodOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.distributionMethod, "distribution-method", "", "Distribution method: route or gateway")
}

func (o *DistributionMethodOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.distributionMethod != "" {
		collection.SetParam("distribution_method", o.distributionMethod)
	}
}

type AdditionalOption[T any] struct {
	additional bool
}

func (o *AdditionalOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.additional, "additional", false, "Filter additional networks only")
}

func (o *AdditionalOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.additional {
		collection.SetParam("additional", "true")
	}
}
