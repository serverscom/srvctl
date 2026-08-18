package entities

import (
	"log"
	"reflect"

	"github.com/serverscom/srvctl/internal/metrics"
)

var (
	HostMetricType = reflect.TypeFor[metrics.HostMetric]()
	RackMetricType = reflect.TypeFor[metrics.RackMetric]()
)

// RegisterHostMetricDefinition registers hosts metrics entity
func RegisterHostMetricDefinition() {
	hostMetricEntity := &Entity{
		fields: []Field{
			{ID: "HostID", Name: "Host ID", Path: "HostID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationCode", Name: "Location", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "Location ID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "HostType", Name: "Type", Path: "HostType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ChassisName", Name: "Chassis", Path: "ChassisName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "RackID", Name: "Rack ID", Path: "RackID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "RackType", Name: "Rack Type", Path: "RackType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicSent", Name: "Public Sent", Path: "PublicSent", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler, Default: true},
			{ID: "PublicReceived", Name: "Public Recv", Path: "PublicReceived", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler, Default: true},
			{ID: "PrivateSent", Name: "Private Sent", Path: "PrivateSent", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler, Default: true},
			{ID: "PrivateReceived", Name: "Private Recv", Path: "PrivateReceived", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler, Default: true},
			{ID: "TotalSent", Name: "Total Sent", Path: "TotalSent", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler},
			{ID: "TotalReceived", Name: "Total Recv", Path: "TotalReceived", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler},
		},
		eType: HostMetricType,
	}

	if err := Registry.Register(hostMetricEntity); err != nil {
		log.Fatal(err)
	}
}

// RegisterRackMetricDefinition registers racks metrics entity
func RegisterRackMetricDefinition() {
	rackMetricEntity := &Entity{
		fields: []Field{
			{ID: "RackID", Name: "Rack ID", Path: "RackID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationCode", Name: "Location", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "Location ID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Hosts", Name: "Hosts", Path: "Hosts", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "PublicSent", Name: "Public Sent", Path: "PublicSent", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler, Default: true},
			{ID: "PublicReceived", Name: "Public Recv", Path: "PublicReceived", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler, Default: true},
			{ID: "PrivateSent", Name: "Private Sent", Path: "PrivateSent", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler},
			{ID: "PrivateReceived", Name: "Private Recv", Path: "PrivateReceived", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler},
			{ID: "TotalSent", Name: "Total Sent", Path: "TotalSent", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler},
			{ID: "TotalReceived", Name: "Total Recv", Path: "TotalReceived", ListHandlerFunc: bytesHandler, PageViewHandlerFunc: bytesHandler},
			{ID: "PduWatts", Name: "PDU Watts", Path: "PduWatts", ListHandlerFunc: floatHandler, PageViewHandlerFunc: floatHandler, Default: true},
			{ID: "PduAmperes", Name: "PDU Amperes", Path: "PduAmperes", ListHandlerFunc: floatHandler, PageViewHandlerFunc: floatHandler, Default: true},
			{ID: "PduCount", Name: "PDUs", Path: "PduCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "AtsWatts", Name: "ATS Watts", Path: "AtsWatts", ListHandlerFunc: floatHandler, PageViewHandlerFunc: floatHandler},
			{ID: "AtsAmperes", Name: "ATS Amperes", Path: "AtsAmperes", ListHandlerFunc: floatHandler, PageViewHandlerFunc: floatHandler},
			{ID: "AtsCount", Name: "ATSs", Path: "AtsCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
		},
		eType: RackMetricType,
	}

	if err := Registry.Register(rackMetricEntity); err != nil {
		log.Fatal(err)
	}
}
