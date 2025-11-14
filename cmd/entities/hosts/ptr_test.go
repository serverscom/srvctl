package hosts

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
	ptrFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "ptr")
	testServerID       = "serverId"
	testPTRID          = "ptrId"

	testPTR = serverscom.PTRRecord{
		ID:       testPTRID,
		IP:       "192.0.2.5",
		Domain:   "ptr-test.example.com",
		Priority: 10,
		TTL:      300,
	}
)

func TestListDSPTRCmd(t *testing.T) {
	testPTR1 := testPTR
	testPTR2 := testPTR
	testPTR1.ID += "1"
	testPTR2.ID += "2"
	testPTR2.Domain = "another.example.com"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.PTRRecord])
	}{
		{
			name:           "list all ptr records",
			output:         "json",
			args:           []string{"-A", testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:           "list ptr records",
			output:         "json",
			args:           []string{testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
					}, nil)
			},
		},
		{
			name:           "list ptr records with template",
			args:           []string{testServerID, "--template", "{{range .}}ID: {{.ID}} PTR: {{.Domain}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:           "list ptr records with page-view",
			args:           []string{testServerID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:        "list ptr records with error",
			args:        []string{testServerID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.PTRRecord](mockCtrl)

	hostService.EXPECT().DedicatedServerPTRRecords(gomock.Any()).Return(collection).AnyTimes()
	collection.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collection).AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collection)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := append([]string{"hosts", "ds", "list-ptr"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestCreateDSPTRCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		output         string
		expectedOutput []byte
		configureMock  func(*mocks.MockHostsService)
		expectError    bool
	}{
		{
			name:           "create ds ptr",
			args:           []string{testServerID, "--ip", testPTR.IP, "--domain", testPTR.Domain, "--ttl", "300", "--priority", "10"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				in := serverscom.PTRRecordCreateInput{
					IP:       testPTR.IP,
					Domain:   testPTR.Domain,
					TTL:      &testPTR.TTL,
					Priority: &testPTR.Priority,
				}
				mock.EXPECT().
					CreatePTRRecordForDedicatedServer(gomock.Any(), testServerID, in).
					Return(&testPTR, nil)
			},
		},
		{
			name:        "create ds ptr error",
			args:        []string{testServerID, "--ip", testPTR.IP, "--domain", testPTR.Domain},
			expectError: true,
			configureMock: func(mock *mocks.MockHostsService) {
				in := serverscom.PTRRecordCreateInput{
					IP:     testPTR.IP,
					Domain: testPTR.Domain,
				}
				mock.EXPECT().
					CreatePTRRecordForDedicatedServer(gomock.Any(), testServerID, in).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(hostService)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "add-ptr"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestDeleteDSPTRCmd(t *testing.T) {
	testCases := []struct {
		name        string
		serverID    string
		ptrID       string
		expectError bool
	}{
		{
			name:     "delete ds ptr",
			serverID: testServerID,
			ptrID:    testPTRID,
		},
		{
			name:        "delete ds ptr error",
			serverID:    testServerID,
			ptrID:       testPTRID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostService.EXPECT().
				DeletePTRRecordForDedicatedServer(gomock.Any(), tc.serverID, tc.ptrID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "delete-ptr", tc.serverID, "--ptr-id", tc.ptrID}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err = c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestListSBMPTRCmd(t *testing.T) {
	testPTR1 := testPTR
	testPTR2 := testPTR
	testPTR1.ID += "1"
	testPTR2.ID += "2"
	testPTR2.Domain = "another.example.com"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.PTRRecord])
	}{
		{
			name:           "list all ptr records",
			output:         "json",
			args:           []string{"-A", testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:           "list ptr records",
			output:         "json",
			args:           []string{testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
					}, nil)
			},
		},
		{
			name:           "list ptr records with template",
			args:           []string{testServerID, "--template", "{{range .}}ID: {{.ID}} PTR: {{.Domain}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:           "list ptr records with page-view",
			args:           []string{testServerID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:        "list ptr records with error",
			args:        []string{testServerID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.PTRRecord](mockCtrl)

	hostService.EXPECT().SBMServerPTRRecords(gomock.Any()).Return(collection).AnyTimes()
	collection.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collection).AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collection)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := append([]string{"hosts", "sbm", "list-ptr"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestCreateSBMPTRCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		output         string
		expectedOutput []byte
		configureMock  func(*mocks.MockHostsService)
		expectError    bool
	}{
		{
			name:           "create sbm ptr",
			args:           []string{testServerID, "--ip", testPTR.IP, "--domain", testPTR.Domain, "--ttl", "300", "--priority", "10"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				in := serverscom.PTRRecordCreateInput{
					IP:       testPTR.IP,
					Domain:   testPTR.Domain,
					TTL:      &testPTR.TTL,
					Priority: &testPTR.Priority,
				}
				mock.EXPECT().
					CreatePTRRecordForSBMServer(gomock.Any(), testServerID, in).
					Return(&testPTR, nil)
			},
		},
		{
			name:        "create sbm ptr error",
			args:        []string{testServerID, "--ip", testPTR.IP, "--domain", testPTR.Domain},
			expectError: true,
			configureMock: func(mock *mocks.MockHostsService) {
				in := serverscom.PTRRecordCreateInput{
					IP:     testPTR.IP,
					Domain: testPTR.Domain,
				}
				mock.EXPECT().
					CreatePTRRecordForSBMServer(gomock.Any(), testServerID, in).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(hostService)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "add-ptr"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestDeleteSBMPTRCmd(t *testing.T) {
	testCases := []struct {
		name        string
		serverID    string
		ptrID       string
		expectError bool
	}{
		{
			name:     "delete ds ptr",
			serverID: testServerID,
			ptrID:    testPTRID,
		},
		{
			name:        "delete ds ptr error",
			serverID:    testServerID,
			ptrID:       testPTRID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostService.EXPECT().
				DeletePTRRecordForSBMServer(gomock.Any(), tc.serverID, tc.ptrID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "delete-ptr", tc.serverID, "--ptr-id", tc.ptrID}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err = c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
