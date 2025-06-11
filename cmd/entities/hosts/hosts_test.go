package hosts

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
	testId          = "testId"
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "hosts")
	fixedTime       = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testHost        = serverscom.Host{
		ID:      testId,
		Title:   "example.aa",
		Status:  "active",
		Created: fixedTime,
		Updated: fixedTime,
	}
	testDS = serverscom.DedicatedServer{
		ID:      testId,
		RackID:  testId,
		Type:    "dedicated_server",
		Title:   "example.aa",
		Status:  "active",
		Created: fixedTime,
		Updated: fixedTime,
	}
	testKBM = serverscom.KubernetesBaremetalNode{
		ID:                          testId,
		RackID:                      testId,
		KubernetesClusterID:         testId,
		KubernetesClusterNodeID:     testId,
		KubernetesClusterNodeNumber: 1,
		Type:                        "kubernetes_baremetal_node",
		Title:                       "example.aa",
		Status:                      "active",
		Created:                     fixedTime,
		Updated:                     fixedTime,
	}
	testSBM = serverscom.SBMServer{
		ID:      testId,
		RackID:  testId,
		Type:    "sbm_server",
		Title:   "example.aa",
		Status:  "active",
		Created: fixedTime,
		Updated: fixedTime,
	}
)

func TestAddDSCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockHostsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create dedicated server",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ds_resp.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create_ds_input.json")},
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), serverscom.DedicatedServerCreateInput{
						ServerModelID: 1234,
						LocationID:    5678,
						RAMSize:       16,
						UplinkModels: serverscom.DedicatedServerUplinkModelsInput{
							Public: &serverscom.DedicatedServerPublicUplinkInput{
								ID:               4321,
								BandwidthModelID: 8765,
							},
							Private: serverscom.DedicatedServerPrivateUplinkInput{
								ID: 7890,
							},
						},
						Drives: serverscom.DedicatedServerDrivesInput{
							Slots: []serverscom.DedicatedServerSlotInput{
								{
									Position:     1,
									DriveModelID: testutils.PtrInt64(3456),
								},
								{
									Position:     2,
									DriveModelID: testutils.PtrInt64(3456),
								},
							},
							Layout: []serverscom.DedicatedServerLayoutInput{
								{
									SlotPositions: []int{1, 2},
									Raid:          testutils.PtrInt(1),
									Partitions: []serverscom.DedicatedServerLayoutPartitionInput{
										{
											Target: "/boot",
											Size:   500,
											Fill:   false,
											Fs:     testutils.PtrString("ext4"),
										},
									},
								},
							},
						},
						Hosts: []serverscom.DedicatedServerHostInput{
							{
								Hostname:             "example.aa",
								PublicIPv4NetworkID:  testutils.PtrString("PublicNet123"),
								PrivateIPv4NetworkID: testutils.PtrString("PrivateNet456"),
								Labels: map[string]string{
									"environment": "testing",
								},
							},
						},
					}).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:        "create dedicated server with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(hostsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestAddSBMCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockHostsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create SBM server",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_sbm_resp.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create_sbm_input.json")},
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					CreateSBMServers(gomock.Any(), serverscom.SBMServerCreateInput{
						FlavorModelID: 1234,
						LocationID:    5678,
						Hosts: []serverscom.SBMServerHostInput{
							{
								Hostname:             "example.aa",
								PublicIPv4NetworkID:  testutils.PtrString("PublicNetTest123"),
								PrivateIPv4NetworkID: testutils.PtrString("PrivateNetTest456"),
								Labels: map[string]string{
									"environment": "testing",
								},
							},
						},
					}).
					Return([]serverscom.SBMServer{testSBM}, nil)
			},
		},
		{
			name:        "create SBM server with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(hostsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestGetDSCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get dedicated server in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds.txt")),
		},
		{
			name:           "get dedicated server in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds.json")),
		},
		{
			name:           "get dedicated server in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds.yaml")),
		},
		{
			name:        "get dedicated server with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostsServiceHandler.EXPECT().
				GetDedicatedServer(gomock.Any(), testId).
				Return(&testDS, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestGetKBMCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get KBM node in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_kbm.txt")),
		},
		{
			name:           "get KBM node in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_kbm.json")),
		},
		{
			name:           "get KBM node in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_kbm.yaml")),
		},
		{
			name:        "get KBM node with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostsServiceHandler.EXPECT().
				GetKubernetesBaremetalNode(gomock.Any(), testId).
				Return(&testKBM, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "kbm", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestGetSBMCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get SBM server in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm.txt")),
		},
		{
			name:           "get SBM server in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm.json")),
		},
		{
			name:           "get SBM server in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm.yaml")),
		},
		{
			name:        "get SBM server with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostsServiceHandler.EXPECT().
				GetSBMServer(gomock.Any(), testId).
				Return(&testSBM, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestListDSCmd(t *testing.T) {
	testServer1 := testHost
	testServer1.Type = "dedicated_server"
	testServer2 := testServer1
	testServer2.ID = "testId2"
	testServer2.Title = "example.bb"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Host])
	}{
		{
			name:           "list all dedicated servers",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list dedicated servers",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list dedicated servers with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list dedicated servers with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:        "list dedicated servers with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Host](mockCtrl)

	hostsServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestListKBMCmd(t *testing.T) {
	testServer1 := testHost
	testServer1.Type = "kubernetes_baremetal_node"
	testServer2 := testServer1
	testServer2.ID = "testId2"
	testServer2.Title = "example.bb"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Host])
	}{
		{
			name:           "list all KMB nodes",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list KBM nodes",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list KBM nodes with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list KBM nodes with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:        "list KBM nodes with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Host](mockCtrl)

	hostsServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "kbm", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestListSBMCmd(t *testing.T) {
	testServer1 := testHost
	testServer1.Type = "sbm_server"
	testServer2 := testServer1
	testServer2.ID = "testId2"
	testServer2.Title = "example.bb"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Host])
	}{
		{
			name:           "list all SBM servers",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list SBM servers",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list SBM servers with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list SBM servers with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:        "list SBM servers with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Host](mockCtrl)

	hostsServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestUpdateDSCmd(t *testing.T) {
	newServer := testDS
	newServer.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockHostsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update dedicated server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_ds_resp.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					UpdateDedicatedServer(gomock.Any(), testId, serverscom.DedicatedServerUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newServer, nil)
			},
		},
		{
			name: "update dedicated server with error",
			id:   testId,
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					UpdateDedicatedServer(gomock.Any(), testId, serverscom.DedicatedServerUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(hostsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestUpdateKBMCmd(t *testing.T) {
	newServer := testKBM
	newServer.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockHostsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update KBM node",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_kbm_resp.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					UpdateKubernetesBaremetalNode(gomock.Any(), testId, serverscom.KubernetesBaremetalNodeUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newServer, nil)
			},
		},
		{
			name: "update KBM node with error",
			id:   testId,
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					UpdateKubernetesBaremetalNode(gomock.Any(), testId, serverscom.KubernetesBaremetalNodeUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(hostsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "kbm", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestUpdateSBMCmd(t *testing.T) {
	newServer := testSBM
	newServer.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockHostsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update SBM server node",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_sbm_resp.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					UpdateSBMServer(gomock.Any(), testId, serverscom.SBMServerUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newServer, nil)
			},
		},
		{
			name: "update SBM server with error",
			id:   testId,
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					UpdateSBMServer(gomock.Any(), testId, serverscom.SBMServerUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(hostsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestScheduleReleaseDSCmd(t *testing.T) {
	releasedServer := testDS
	releasedServer.ScheduledRelease = &fixedTime
	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "release dedicated server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "release_ds_resp.json")),
		},
		// TODO add after implementing release-after in client
		// {
		// 	name:           "release dedicated server with --release-after",
		// 	id:             testId,
		// 	output:         "json",
		// 	expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "release_ds_resp.json")),
		// 	args:           []string{"--release-after", "2025-01-01T12:34:56+03:00"},
		// },
		{
			name:        "release dedicated server with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostsServiceHandler.EXPECT().
				ScheduleReleaseForDedicatedServer(gomock.Any(), testId).
				Return(&releasedServer, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "schedule-release", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestAbortReleaseDSCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "abort release dedicated server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds.json")),
		},
		{
			name:        "abort release dedicated server with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostsServiceHandler.EXPECT().
				AbortReleaseForDedicatedServer(gomock.Any(), testId).
				Return(&testDS, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "abort-release", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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

func TestReleaseSBMCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "release SBM server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm.json")),
		},
		{
			name:        "release SBM server with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			hostsServiceHandler.EXPECT().
				ReleaseSBMServer(gomock.Any(), testId).
				Return(&testSBM, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "release", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
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
