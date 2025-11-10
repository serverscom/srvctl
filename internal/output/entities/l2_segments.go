package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	L2SegmentType                    = reflect.TypeOf(serverscom.L2Segment{})
	L2SegmentLocationGroupType       = reflect.TypeOf(serverscom.L2LocationGroup{})
	L2SegmentMemberType              = reflect.TypeOf(serverscom.L2Member{})
	L2SegmentNetworkType             = reflect.TypeOf(serverscom.Network{})
	L2SegmentListDefaultFields       = []string{"ID", "Name", "Type", "Status", "LocationGroupID", "LocationGroupCode"}
	L2SegmentMemberListDefaultFields = []string{"ID", "Title", "Mode", "VLAN", "Status"}
)

func RegisterL2SegmentDefinitions() {
	l2SegmentEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationGroupID", Name: "LocationGroupID", Path: "LocationGroupID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationGroupCode", Name: "LocationGroupCode", Path: "LocationGroupCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},

		cmdDefaultFields: map[string][]string{
			"list": L2SegmentListDefaultFields,
		},
		eType: L2SegmentType,
	}
	if err := Registry.Register(l2SegmentEntity); err != nil {
		log.Fatal(err)
	}

	l2SegmentLocationGroupEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Code", Name: "Code", Path: "Code", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "GroupType", Name: "GroupType", Path: "GroupType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationIDs", Name: "LocationIDs", Path: "LocationIDs", PageViewHandlerFunc: slicePvHandler},
			{ID: "Hyperscalers", Name: "Hyperscalers", Path: "Hyperscalers", PageViewHandlerFunc: slicePvHandler},
		},
		eType: L2SegmentLocationGroupType,
	}
	if err := Registry.Register(l2SegmentLocationGroupEntity); err != nil {
		log.Fatal(err)
	}

	l2SegmentMemberEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Mode", Name: "Mode", Path: "Mode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "VLAN", Name: "VLAN", Path: "VLAN", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},

		cmdDefaultFields: map[string][]string{
			"list-members": L2SegmentMemberListDefaultFields,
		},
		eType: L2SegmentMemberType,
	}
	if err := Registry.Register(l2SegmentMemberEntity); err != nil {
		log.Fatal(err)
	}
}
