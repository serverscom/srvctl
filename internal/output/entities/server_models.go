package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	ServerModelOptionType       = reflect.TypeOf(serverscom.ServerModelOption{})
	ServerModelOptionDetailType = reflect.TypeOf(serverscom.ServerModelOptionDetail{})
)

func getServerModelDriveSlotsField() Field {
	f := Field{
		ID:                  "DriveSlots",
		Name:                "DriveSlots",
		Path:                "DriveSlots",
		PageViewHandlerFunc: structPVHandler,
	}
	childs := []Field{
		{ID: "Position", Name: "Position", Path: "Position", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "Interface", Name: "Interface", Path: "Interface", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "FormFactor", Name: "FormFactor", Path: "FormFactor", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "DriveModelID", Name: "DriveModelID", Path: "DriveModelID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "HotSwappable", Name: "HotSwappable", Path: "HotSwappable", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
	}
	f.ChildFields = append(f.ChildFields, childs...)
	return f
}

func RegisterServerModelOptionDefinitions() {

	// used for list cmd
	serverModelOptionEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUName", Name: "CPUName", Path: "CPUName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "CPUCount", Name: "CPUCount", Path: "CPUCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUCoresCount", Name: "CPUCoresCount", Path: "CPUCoresCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUFrequency", Name: "CPUFrequency", Path: "CPUFrequency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "RAM", Name: "RAM", Path: "RAM", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RAMType", Name: "RAMType", Path: "RAMType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "MaxRAM", Name: "MaxRAM", Path: "MaxRAM", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "HasRAIDController", Name: "HasRAIDController", Path: "HasRAIDController", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RAIDControllerName", Name: "RAIDControllerName", Path: "RAIDControllerName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "DriveSlotsCount", Name: "DriveSlotsCount", Path: "DriveSlotsCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: ServerModelOptionType,
	}
	if err := Registry.Register(serverModelOptionEntity); err != nil {
		log.Fatal(err)
	}

	// used for get cmd
	serverModelOptionDetailEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUName", Name: "CPUName", Path: "CPUName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUCount", Name: "CPUCount", Path: "CPUCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUCoresCount", Name: "CPUCoresCount", Path: "CPUCoresCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CPUFrequency", Name: "CPUFrequency", Path: "CPUFrequency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RAM", Name: "RAM", Path: "RAM", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RAMType", Name: "RAMType", Path: "RAMType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "MaxRAM", Name: "MaxRAM", Path: "MaxRAM", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "HasRAIDController", Name: "HasRAIDController", Path: "HasRAIDController", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RAIDControllerName", Name: "RAIDControllerName", Path: "RAIDControllerName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DriveSlotsCount", Name: "DriveSlotsCount", Path: "DriveSlotsCount", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			getServerModelDriveSlotsField(),
		},
		eType: ServerModelOptionDetailType,
	}
	if err := Registry.Register(serverModelOptionDetailEntity); err != nil {
		log.Fatal(err)
	}
}
