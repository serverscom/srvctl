package base

import (
	"fmt"
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
	flags.StringVar(&o.direction, "direction", "", "Sort direction (asc, desc)")
	flags.BoolVarP(&o.allPages, "all", "A", false, "Get all pages of resources")
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

// HiddenTypeOption adds hidden type flag.
// Used in commands that determine type from sub command rather than user input
type HiddenTypeOption[T any] struct{}

func (o *HiddenTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().String("type", "", "")
	_ = cmd.Flags().MarkHidden("type")
}

func (o *HiddenTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	// stub for compatibility with other options
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

// invoice status option
type InvoiceStatusOption[T any] struct {
	status string
}

func (o *InvoiceStatusOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.status, "status", "", "Filter results by status (pending, outstanding, overdue, paid, canceled, reissued)")
}

func (o *InvoiceStatusOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
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
	cmd.Flags().StringVar(&o.typeVal, "type", "", "Filter results by type (invoice, credit_note)")
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
	cmd.Flags().StringVar(&o.family, "family", "", "Filter results by IP family (ipv4, ipv6)")
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
	cmd.Flags().StringVar(&o.interfaceType, "interface-type", "", "Filter results by network interface type (public, private, oob)")
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
	cmd.Flags().StringVar(&o.distributionMethod, "distribution-method", "", "Filter results by distribution method (route, gateway)")
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

// redundancy option
type RedundancyOption[T any] struct {
	redundancy bool
}

func (o *RedundancyOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.redundancy, "redundancy", false, "Filter uplinks by redundancy (true, false)")
}

func (o *RedundancyOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.redundancy {
		collection.SetParam("redundancy", "true")
	}
}

// uplink type option
type UplinkTypeOption[T any] struct {
	uplinkType string
}

func (o *UplinkTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.uplinkType, "type", "", "Filter uplinks by type (public, private)")
}

func (o *UplinkTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.uplinkType != "" {
		collection.SetParam("type", o.uplinkType)
	}
}

// operating system id option
type OperatingSystemIDOption[T any] struct {
	osID int64
}

func (o *OperatingSystemIDOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().Int64Var(&o.osID, "operating-system-id", 0, "Filter uplinks by operating system ID")
}

func (o *OperatingSystemIDOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.osID != 0 {
		collection.SetParam("operating_system_id", fmt.Sprintf("%d", o.osID))
	}
}

// bandwidth type option
type BandwidthTypeOption[T any] struct {
	bandwidthType string
}

func (o *BandwidthTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.bandwidthType, "type", "", `Filter bandwidth options by type (bytes, bandwidth, unmetered)`)
}

func (o *BandwidthTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.bandwidthType != "" {
		collection.SetParam("type", o.bandwidthType)
	}
}

type HasRaidControllerOption[T any] struct {
	hasRaidController bool
}

func (o *HasRaidControllerOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.hasRaidController, "has-raid-controller", false,
		"Filter only servers with RAID controller")
}

func (o *HasRaidControllerOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.hasRaidController {
		collection.SetParam("has_raid_controller", "true")
	}
}

type DriveMediaTypeOption[T any] struct {
	mediaType string
}

func (o *DriveMediaTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.mediaType, "media-type", "",
		"Filter drives by media type (HDD, SSD")
}

func (o *DriveMediaTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.mediaType != "" {
		collection.SetParam("media_type", o.mediaType)
	}
}

type DriveInterfaceOption[T any] struct {
	iface string
}

func (o *DriveInterfaceOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.iface, "interface", "",
		"Filter drives by interface (SATA1, SATA2, SATA3, SAS, NVMe-PCIe)")
}

func (o *DriveInterfaceOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.iface != "" {
		collection.SetParam("interface", o.iface)
	}
}

type SBMFlavorsShowAllOption[T any] struct {
	all bool
}

func (o *SBMFlavorsShowAllOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.all, "show-all", false,
		"Filter to show all SBM flavors including unavailable ones")
}

func (o *SBMFlavorsShowAllOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.all {
		collection.SetParam("show_all", "true")
	}
}

type L2SegmentGroupTypeOption[T any] struct {
	group string
}

func (o *L2SegmentGroupTypeOption[T]) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.group, "group-type", "",
		"Filter l2 location groups by type (public, private)")
}

func (o *L2SegmentGroupTypeOption[T]) ApplyToCollection(collection serverscom.Collection[T]) {
	if o.group != "" {
		collection.SetParam("group_type", o.group)
	}
}
