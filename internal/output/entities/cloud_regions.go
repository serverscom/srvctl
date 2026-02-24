package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	CloudComputingRegionType              = reflect.TypeOf(serverscom.CloudComputingRegion{})
	CloudComputingRegionListDefaultFields = []string{"ID", "Name", "Code"}
	CloudComputingImageType               = reflect.TypeOf(serverscom.CloudComputingImage{})
	CloudComputingImageListDefaultFields  = []string{"ID", "Name"}
	CloudComputingFlavorType              = reflect.TypeOf(serverscom.CloudComputingFlavor{})
	CloudComputingFlavorListDefaultFields = []string{"ID", "Name"}
	CloudSnapshotType                     = reflect.TypeOf(serverscom.CloudSnapshot{})
	CloudSnapshotListDefaultFields        = []string{"ID", "Name", "ImageSize", "MinDisk", "Status", "IsBackup", "FileURL"}
	CloudComputingRegionCredentialsType   = reflect.TypeOf(serverscom.CloudComputingRegionCredentials{})
)

func RegisterCloudComputingRegionDefinitions() {
	// CloudComputingRegion
	cloudRegionEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Code", Name: "Code", Path: "Code", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudComputingRegionListDefaultFields,
		},
		eType: CloudComputingRegionType,
	}
	if err := Registry.Register(cloudRegionEntity); err != nil {
		log.Fatal(err)
	}

	// CloudComputingImage
	cloudImageEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudComputingImageListDefaultFields,
		},
		eType: CloudComputingImageType,
	}
	if err := Registry.Register(cloudImageEntity); err != nil {
		log.Fatal(err)
	}

	// CloudComputingFlavor
	cloudFlavorEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudComputingFlavorListDefaultFields,
		},
		eType: CloudComputingFlavorType,
	}
	if err := Registry.Register(cloudFlavorEntity); err != nil {
		log.Fatal(err)
	}

	// CloudSnapshot
	cloudSnapshotEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ImageSize", Name: "Image Size", Path: "ImageSize", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "MinDisk", Name: "Min Disk", Path: "MinDisk", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "IsBackup", Name: "Is Backup", Path: "IsBackup", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "FileURL", Name: "File URL", Path: "FileURL", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudSnapshotListDefaultFields,
		},
		eType: CloudSnapshotType,
	}
	if err := Registry.Register(cloudSnapshotEntity); err != nil {
		log.Fatal(err)
	}

	// CloudComputingRegionCredentials
	cloudCredentialsEntity := &Entity{
		fields: []Field{
			{ID: "Password", Name: "Password", Path: "Password", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TenantName", Name: "Tenant Name", Path: "TenantName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "URL", Name: "URL", Path: "URL", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Username", Name: "Username", Path: "Username", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: CloudComputingRegionCredentialsType,
	}
	if err := Registry.Register(cloudCredentialsEntity); err != nil {
		log.Fatal(err)
	}
}
