package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	DNSDomainType              = reflect.TypeFor[serverscom.DNSDomain]()
	DNSRecordType              = reflect.TypeFor[serverscom.DNSRecord]()
	DNSDelegationDataType      = reflect.TypeFor[serverscom.DNSDomainDelegationData]()
	DNSDomainListDefaultFields = []string{"ID", "Name", "DelegationStatus", "TTL"}
	DNSRecordListDefaultFields = []string{"ID", "Type", "Name", "Data", "TTL", "Priority"}
)

func RegisterDNSDefinitions() {
	dnsDomainEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Email", Name: "Email", Path: "Email", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TTL", Name: "TTL", Path: "TTL", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DelegationStatus", Name: "DelegationStatus", Path: "DelegationStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "UnpublishDate", Name: "UnpublishDate", Path: "UnpublishDate", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": DNSDomainListDefaultFields,
		},
		eType: DNSDomainType,
	}
	if err := Registry.Register(dnsDomainEntity); err != nil {
		log.Fatal(err)
	}

	dnsRecordEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DomainID", Name: "DomainID", Path: "DomainID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Data", Name: "Data", Path: "Data", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TTL", Name: "TTL", Path: "TTL", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Priority", Name: "Priority", Path: "Priority", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list-records": DNSRecordListDefaultFields,
		},
		eType: DNSRecordType,
	}
	if err := Registry.Register(dnsRecordEntity); err != nil {
		log.Fatal(err)
	}

	dnsDelegationDataEntity := &Entity{
		fields: []Field{
			{ID: "Nameservers", Name: "Nameservers", Path: "Nameservers", ListHandlerFunc: stringHandler, PageViewHandlerFunc: slicePvHandler, Default: true},
			{ID: "RequiredTxt", Name: "RequiredTxt", Path: "RequiredTxt", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: DNSDelegationDataType,
	}
	if err := Registry.Register(dnsDelegationDataEntity); err != nil {
		log.Fatal(err)
	}
}
