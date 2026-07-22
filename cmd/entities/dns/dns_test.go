package dns

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "dns")
	skeletonTemplatePath = filepath.Join("..", "..", "..", "internal", "output", "skeletons", "templates", "dns")
	fixedTime            = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	testDomainID = "abc123"
	testRecordID = "rec1"

	testDomain = serverscom.DNSDomain{
		ID:               testDomainID,
		Name:             "example.com",
		Email:            "admin@example.com",
		TTL:              3600,
		DelegationStatus: serverscom.DNSDomainVerified,
		Labels:           map[string]string{"foo": "bar"},
		Created:          fixedTime,
		Updated:          fixedTime,
	}

	testRecord = serverscom.DNSRecord{
		ID:       testRecordID,
		DomainID: testDomainID,
		Name:     "www",
		Type:     serverscom.DNSRecordTypeA,
		Data:     new("1.2.3.4"),
		TTL:      new(3600),
		Created:  fixedTime,
		Updated:  fixedTime,
	}

	testDelegation = serverscom.DNSDomainDelegationData{
		Nameservers: []string{"ns1.example.com", "ns2.example.com"},
		RequiredTxt: "sc-verification=abc123",
	}
)

func TestAddDNSDomainCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockDNSService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create domain with input",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create.json")},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					CreateDomain(gomock.Any(), serverscom.DNSDomainCreateInput{
						Name:   "example.com",
						Email:  "admin@example.com",
						TTL:    3600,
						Labels: map[string]string{"foo": "bar"},
					}).
					Return(&testDomain, nil)
			},
		},
		{
			name:           "create domain with flags",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				"--name", "example.com",
				"--email", "admin@example.com",
				"--ttl", "3600",
				"--label", "foo=bar",
			},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					CreateDomain(gomock.Any(), serverscom.DNSDomainCreateInput{
						Name:   "example.com",
						Email:  "admin@example.com",
						TTL:    3600,
						Labels: map[string]string{"foo": "bar"},
					}).
					Return(&testDomain, nil)
			},
		},
		{
			name:           "skeleton for domain input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add.json")),
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().CreateDomain(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:        "create domain with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(dnsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestGetDNSDomainCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get domain in default format",
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get domain in JSON format",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get domain in YAML format",
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get domain with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			dnsServiceHandler.EXPECT().
				GetDomain(gomock.Any(), testDomainID).
				Return(&testDomain, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "get", testDomainID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestListDNSDomainsCmd(t *testing.T) {
	domain1 := testDomain
	domain2 := testDomain
	domain2.ID = "def456"
	domain2.Name = "example2.com"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.DNSDomain])
	}{
		{
			name:           "list all domains",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSDomain]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DNSDomain{domain1, domain2}, nil)
			},
		},
		{
			name:           "list domains",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSDomain]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DNSDomain{domain1}, nil)
			},
		},
		{
			name:           "list domains with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSDomain]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DNSDomain{domain1, domain2}, nil)
			},
		},
		{
			name:           "list domains with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSDomain]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DNSDomain{domain1, domain2}, nil)
			},
		},
		{
			name:        "list domains with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSDomain]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.DNSDomain](mockCtrl)

	dnsServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestUpdateDNSDomainCmd(t *testing.T) {
	newDomain := testDomain
	newDomain.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockDNSService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update domain",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					UpdateDomain(gomock.Any(), testDomainID, serverscom.DNSDomainUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newDomain, nil)
			},
		},
		{
			name: "update domain with error",
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					UpdateDomain(gomock.Any(), testDomainID, serverscom.DNSDomainUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(dnsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "update", testDomainID}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestDeleteDNSDomainCmd(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
	}{
		{name: "delete domain"},
		{name: "delete domain with error", expectError: true},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			dnsServiceHandler.EXPECT().
				DeleteDomain(gomock.Any(), testDomainID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "delete", testDomainID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
				WithArgs(args)

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestDelegationDataCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "delegation data in default format",
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "delegation_data.txt")),
		},
		{
			name:           "delegation data in JSON format",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "delegation_data.json")),
		},
		{
			name:           "delegation data in YAML format",
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "delegation_data.yaml")),
		},
		{
			name:        "delegation data with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			dnsServiceHandler.EXPECT().
				GetDelegationData(gomock.Any()).
				Return(&testDelegation, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "delegation-data"}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestListDNSRecordsCmd(t *testing.T) {
	record1 := testRecord
	record2 := testRecord
	record2.ID = "rec2"
	record2.Name = "mail"
	record2.Type = serverscom.DNSRecordTypeMX
	record2.Data = new("mail.example.com")
	record2.Priority = new(10)

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.DNSRecord])
	}{
		{
			name:           "list all records",
			output:         "json",
			args:           []string{testDomainID, "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_records_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSRecord]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DNSRecord{record1, record2}, nil)
			},
		},
		{
			name:           "list records",
			output:         "json",
			args:           []string{testDomainID},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_records.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DNSRecord{record1}, nil)
			},
		},
		{
			name:           "list records with template",
			args:           []string{testDomainID, "--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_records_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DNSRecord{record1, record2}, nil)
			},
		},
		{
			name:           "list records with pageView",
			args:           []string{testDomainID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_records_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DNSRecord{record1, record2}, nil)
			},
		},
		{
			name:        "list records with error",
			args:        []string{testDomainID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DNSRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.DNSRecord](mockCtrl)

	dnsServiceHandler.EXPECT().
		Records(testDomainID).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "list-records"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestGetDNSRecordCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get record in default format",
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_record.txt")),
		},
		{
			name:           "get record in JSON format",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_record.json")),
		},
		{
			name:           "get record in YAML format",
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_record.yaml")),
		},
		{
			name:        "get record with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			dnsServiceHandler.EXPECT().
				GetRecord(gomock.Any(), testDomainID, testRecordID).
				Return(&testRecord, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "get-record", testDomainID, "--record-id", testRecordID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestAddDNSRecordCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockDNSService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create record with input",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_record.json")),
			args:           []string{testDomainID, "--input", filepath.Join(fixtureBasePath, "create_record.json")},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					CreateRecord(gomock.Any(), testDomainID, serverscom.DNSRecordCreateInput{
						Type: serverscom.DNSRecordTypeA,
						Name: "www",
						Data: "1.2.3.4",
						TTL:  new(3600),
					}).
					Return(&testRecord, nil)
			},
		},
		{
			name:           "create record with flags",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_record.json")),
			args: []string{
				testDomainID,
				"--type", "A",
				"--name", "www",
				"--data", "1.2.3.4",
				"--ttl", "3600",
			},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					CreateRecord(gomock.Any(), testDomainID, serverscom.DNSRecordCreateInput{
						Type: serverscom.DNSRecordTypeA,
						Name: "www",
						Data: "1.2.3.4",
						TTL:  new(3600),
					}).
					Return(&testRecord, nil)
			},
		},
		{
			name:           "skeleton for record input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add-record.json")),
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().CreateRecord(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:        "create record with error",
			args:        []string{testDomainID},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(dnsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "add-record"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestUpdateDNSRecordCmd(t *testing.T) {
	updatedRecord := testRecord
	updatedRecord.Data = new("5.6.7.8")
	updatedRecord.TTL = new(7200)

	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockDNSService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update record",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_record.json")),
			args:           []string{testDomainID, "--record-id", testRecordID, "--data", "5.6.7.8", "--ttl", "7200"},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					UpdateRecord(gomock.Any(), testDomainID, testRecordID, serverscom.DNSRecordUpdateInput{
						Data: "5.6.7.8",
						TTL:  new(7200),
					}).
					Return(&updatedRecord, nil)
			},
		},
		{
			name: "update record with error",
			args: []string{testDomainID, "--record-id", testRecordID},
			configureMock: func(mock *mocks.MockDNSService) {
				mock.EXPECT().
					UpdateRecord(gomock.Any(), testDomainID, testRecordID, serverscom.DNSRecordUpdateInput{}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(dnsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "update-record"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
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

func TestDeleteDNSRecordCmd(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
	}{
		{name: "delete record"},
		{name: "delete record with error", expectError: true},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dnsServiceHandler := mocks.NewMockDNSService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.DNS = dnsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			dnsServiceHandler.EXPECT().
				DeleteRecord(gomock.Any(), testDomainID, testRecordID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			dnsCmd := NewCmd(testCmdContext)

			args := []string{"dns", "delete-record", testDomainID, "--record-id", testRecordID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(dnsCmd).
				WithArgs(args)

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
