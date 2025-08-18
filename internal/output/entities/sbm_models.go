package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	SBMFlavorType = reflect.TypeOf(serverscom.SBMFlavor{})
)

func RegisterSBMModelOptionDefinition() {
	sbmFlavorEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUName", Name: "CPU Name", Path: "CPUName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "CPUCount", Name: "CPU Count", Path: "CPUCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUCoresCount", Name: "CPU Cores Count", Path: "CPUCoresCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUFrequency", Name: "CPU Frequency", Path: "CPUFrequency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "RAMSize", Name: "RAM Size", Path: "RAMSize", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DrivesConfiguration", Name: "Drives Configuration", Path: "DrivesConfiguration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicUplinkModelID", Name: "Public Uplink Model ID", Path: "PublicUplinkModelID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicUplinkModelName", Name: "Public Uplink Model Name", Path: "PublicUplinkModelName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateUplinkModelID", Name: "Private Uplink Model ID", Path: "PrivateUplinkModelID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateUplinkModelName", Name: "Private Uplink Model Name", Path: "PrivateUplinkModelName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "BandwidthID", Name: "Bandwidth ID", Path: "BandwidthID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "BandwidthName", Name: "Bandwidth Name", Path: "BandwidthName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},

		eType: SBMFlavorType,
	}
	if err := Registry.Register(sbmFlavorEntity); err != nil {
		log.Fatal(err)
	}
}
