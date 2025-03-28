package invoices

import (
	"errors"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	testId          = "testId"
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "invoices")
	testInvoice     = serverscom.Invoice{
		ID:       testId,
		Number:   123,
		Status:   "paid",
		Date:     "2025-01-01",
		Type:     "1",
		TotalDue: 1.23,
		Currency: "USD",
		CsvUrl:   "http://test.csv",
		PdfUrl:   "http://test.pdf",
	}
	testInvoiceList = serverscom.InvoiceList{
		ID:       testId,
		Number:   123,
		Status:   "paid",
		Date:     "2025-01-01",
		Type:     "1",
		TotalDue: 1.23,
		Currency: "USD",
	}
)

func TestGetInvoiceCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get invoice in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get invoice in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get invoice in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get invoice with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	invoicesServiceHandler := mocks.NewMockInvoiceService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Invoices = invoicesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			invoicesServiceHandler.EXPECT().
				GetBillingInvoice(gomock.Any(), testId).
				Return(&testInvoice, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			invoiceCmd := NewCmd(testCmdContext)

			args := []string{"invoices", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(invoiceCmd).
				WithArgs(args)

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestListInvoicesCmd(t *testing.T) {
	testInvoice1 := testInvoiceList
	testInvoice2 := testInvoiceList
	testInvoice1.ID += "1"
	testInvoice2.Number = 456
	testInvoice2.ID += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.InvoiceList])
	}{
		{
			name:           "list all invoices",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.InvoiceList]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.InvoiceList{
						testInvoice1,
						testInvoice2,
					}, nil)
			},
		},
		{
			name:           "list invoices",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.InvoiceList]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.InvoiceList{
						testInvoice1,
					}, nil)
			},
		},
		{
			name:           "list invoices with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.InvoiceList]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.InvoiceList{
						testInvoice1,
						testInvoice2,
					}, nil)
			},
		},
		{
			name:           "list invoices with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.InvoiceList]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.InvoiceList{
						testInvoice1,
						testInvoice2,
					}, nil)
			},
		},
		{
			name:        "list invoices with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.InvoiceList]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	invoicesServiceHandler := mocks.NewMockInvoiceService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.InvoiceList](mockCtrl)

	invoicesServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Invoices = invoicesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			invoiceCmd := NewCmd(testCmdContext)

			args := []string{"invoices", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(invoiceCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}
