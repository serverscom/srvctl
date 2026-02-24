package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	CloudComputingInstanceType              = reflect.TypeOf(serverscom.CloudComputingInstance{})
	CloudComputingInstanceListDefaultFields = []string{"ID", "Name", "RegionCode", "Status", "PublicIPv4Address"}
)

func RegisterCloudComputingInstanceDefinition() {
	cloudInstanceEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RegionID", Name: "RegionID", Path: "RegionID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "RegionCode", Name: "RegionCode", Path: "RegionCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OpenstackUUID", Name: "OpenstackUUID", Path: "OpenstackUUID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "FlavorID", Name: "FlavorID", Path: "FlavorID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "FlavorName", Name: "FlavorName", Path: "FlavorName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ImageID", Name: "ImageID", Path: "ImageID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ImageName", Name: "ImageName", Path: "ImageName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocalIPv4Address", Name: "LocalIPv4Address", Path: "LocalIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv6Address", Name: "PublicIPv6Address", Path: "PublicIPv6Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "GPNEnabled", Name: "GPNEnabled", Path: "GPNEnabled", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "IPv6Enabled", Name: "IPv6Enabled", Path: "IPv6Enabled", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "BackupCopies", Name: "BackupCopies", Path: "BackupCopies", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicPortBlocked", Name: "PublicPortBlocked", Path: "PublicPortBlocked", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudComputingInstanceListDefaultFields,
		},
		eType: CloudComputingInstanceType,
	}
	if err := Registry.Register(cloudInstanceEntity); err != nil {
		log.Fatal(err)
	}
}
