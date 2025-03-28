package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	InvoiceType     = reflect.TypeOf(serverscom.Invoice{})
	InvoiceListType = reflect.TypeOf(serverscom.InvoiceList{})
)

func RegisterInvoiceDefinition() {
	invoiceEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Number", Name: "Number", Path: "Number", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ParentID", Name: "ParentID", Path: "ParentID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Date", Name: "Date", Path: "Date", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TotalDue", Name: "TotalDue", Path: "TotalDue", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Currency", Name: "Currency", Path: "Currency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CsvUrl", Name: "CsvUrl", Path: "CsvUrl", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PdfUrl", Name: "PdfUrl", Path: "PdfUrl", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
		},
		eType: InvoiceType,
	}
	if err := Registry.Register(invoiceEntity); err != nil {
		log.Fatal(err)
	}

	invoiceListEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Number", Name: "Number", Path: "Number", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ParentID", Name: "ParentID", Path: "ParentID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Date", Name: "Date", Path: "Date", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TotalDue", Name: "TotalDue", Path: "TotalDue", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Currency", Name: "Currency", Path: "Currency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: InvoiceListType,
	}
	if err := Registry.Register(invoiceListEntity); err != nil {
		log.Fatal(err)
	}
}
