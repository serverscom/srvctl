package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	SSLCertType              = reflect.TypeOf(serverscom.SSLCertificate{})
	SSLCertCustomType        = reflect.TypeOf(serverscom.SSLCertificateCustom{})
	SSLCertLeType            = reflect.TypeOf(serverscom.SSLCertificateLE{})
	SSLCertListDefaultFields = []string{"ID", "Name", "Type", "Fingerprint"}
)

func RegisterSSLCertDefinition() {
	e := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Fingerprint", Name: "Fingerprint", Path: "Sha1Fingerprint", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Expires", Name: "Expires", Path: "Expires", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": SSLCertListDefaultFields,
		},
		eType: SSLCertType,
	}
	if err := Registry.Register(e); err != nil {
		log.Fatal(err)
	}
}

func RegisterSSLCertCustomDefinition() {
	e := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Issuer", Name: "Issuer", Path: "Issuer", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Subject", Name: "Subject", Path: "Subject", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DomainNames", Name: "DomainNames", Path: "DomainNames", PageViewHandlerFunc: slicePvHandler},
			{ID: "Fingerprint", Name: "Fingerprint", Path: "Sha1Fingerprint", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Expires", Name: "Expires", Path: "Expires", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		eType: SSLCertCustomType,
	}
	if err := Registry.Register(e); err != nil {
		log.Fatal(err)
	}
}

func RegisterSSLCertLeDefinition() {
	e := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Issuer", Name: "Issuer", Path: "Issuer", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Subject", Name: "Subject", Path: "Subject", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DomainNames", Name: "DomainNames", Path: "DomainNames", PageViewHandlerFunc: slicePvHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Expires", Name: "Expires", Path: "Expires", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		eType: SSLCertLeType,
	}
	if err := Registry.Register(e); err != nil {
		log.Fatal(err)
	}
}
