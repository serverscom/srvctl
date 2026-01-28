package hosts

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	testId               = "testId"
	testNetworkId        = "testNetId"
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "hosts")
	skeletonTemplatePath = filepath.Join("..", "..", "..", "internal", "output", "skeletons", "templates", "hosts")
	fixedTime            = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testPublicIP         = "1.2.3.4"
	testLocationCode     = "test"
	testHost             = serverscom.Host{
		ID:                testId,
		Title:             "example.aa",
		Status:            "active",
		PublicIPv4Address: &testPublicIP,
		LocationCode:      testLocationCode,
		Created:           fixedTime,
		Updated:           fixedTime,
	}
	testConfigDetails = serverscom.ConfigurationDetails{
		RAMSize:                 2,
		ServerModelID:           testutils.PtrInt64(1),
		ServerModelName:         testutils.PtrString("server-model-123"),
		PublicUplinkID:          testutils.PtrInt64(2),
		PublicUplinkName:        testutils.PtrString("Public 1 Gbps without redundancy"),
		PrivateUplinkID:         testutils.PtrInt64(3),
		PrivateUplinkName:       testutils.PtrString("Private 1 Gbps without redundancy"),
		BandwidthID:             testutils.PtrInt64(4),
		BandwidthName:           testutils.PtrString("20000 GB"),
		OperatingSystemID:       testutils.PtrInt64(5),
		OperatingSystemFullName: testutils.PtrString("CentOS 7 x86_64"),
	}
	testDS = serverscom.DedicatedServer{
		ID:                   testId,
		RackID:               testId,
		Type:                 "dedicated_server",
		Title:                "example.aa",
		Status:               "active",
		LocationCode:         testLocationCode,
		PublicIPv4Address:    &testPublicIP,
		ConfigurationDetails: testConfigDetails,
		Created:              fixedTime,
		Updated:              fixedTime,
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
		LocationCode:                testLocationCode,
		PublicIPv4Address:           &testPublicIP,
		ConfigurationDetails:        testConfigDetails,
		Created:                     fixedTime,
		Updated:                     fixedTime,
	}
	testSBM = serverscom.SBMServer{
		ID:                   testId,
		RackID:               testId,
		Type:                 "sbm_server",
		Title:                "example.aa",
		Status:               "active",
		LocationCode:         testLocationCode,
		PublicIPv4Address:    &testPublicIP,
		ConfigurationDetails: testConfigDetails,
		Created:              fixedTime,
		Updated:              fixedTime,
	}
	netTitle    = "Some Net"
	cidr        = "100.0.8.0/29"
	testNetwork = serverscom.Network{
		ID:                 testNetworkId,
		Title:              &netTitle,
		Status:             "active",
		Cidr:               &cidr,
		Family:             "ipv4",
		InterfaceType:      "public",
		DistributionMethod: "gateway",
		Additional:         false,
		Created:            fixedTime,
		Updated:            fixedTime,
	}
	testDriveModel = serverscom.DriveModel{
		ID:         int64(10),
		Name:       "ssd-model-749",
		Capacity:   100,
		Interface:  "SATA3",
		FormFactor: "2.5",
		MediaType:  "SSD",
	}
	testDriveSlot = serverscom.HostDriveSlot{
		Position:   1,
		Interface:  "SAS",
		FormFactor: "2.5\"",
		DriveModel: &testDriveModel,
	}
	testPowerFeed = serverscom.HostPowerFeed{
		Name:   "testPowerFeed",
		Status: "on",
		Type:   "physical",
	}
)

func TestListHostsCmd(t *testing.T) {
	testServer1 := testHost
	testServer1.Type = "dedicated_server"
	testServer2 := testServer1
	testServer2.ID = "testId2"
	testServer2.Type = "sbm_server"
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
			name:           "list all hosts",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_hosts.json")),
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
			name:           "list hosts with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_hosts_template.txt")),
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
			name:           "list hosts with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_hosts_pageview.txt")),
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
			name:        "list hosts with error",
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

			args := []string{"hosts", "list"}
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

func TestAddDSCmd(t *testing.T) {
	expectedInput := serverscom.DedicatedServerCreateInput{
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
	}

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
					CreateDedicatedServers(gomock.Any(), expectedInput).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:           "create dedicated server with merge input with flags",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ds_resp.json")),
			args: []string{
				"--input", filepath.Join(fixtureBasePath, "create_ds_input.json"),
				"--layout", "slot=3,slot=4,raid=0",
				"--partition", "slot=3,slot=4,target=/boot,fs=ext4,size=500",
			},
			configureMock: func(mock *mocks.MockHostsService) {
				input := expectedInput
				input.Drives.Layout = append(input.Drives.Layout, serverscom.DedicatedServerLayoutInput{
					SlotPositions: []int{3, 4},
					Raid:          testutils.PtrInt(0),
					Partitions: []serverscom.DedicatedServerLayoutPartitionInput{
						{
							Target: "/boot",
							Size:   500,
							Fill:   false,
							Fs:     testutils.PtrString("ext4"),
						},
					},
				})

				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), input).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:           "skeleton for dedicated server input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add_ds.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), gomock.Any()).
					Times(0)
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
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
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
			name:           "skeleton for SBM server input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add_sbm.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					CreateSBMServers(gomock.Any(), gomock.Any()).
					Times(0)
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
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
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
	testServer1 := testDS
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
		configureMock  func(*mocks.MockCollection[serverscom.DedicatedServer])
	}{
		{
			name:           "list all dedicated servers",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServer]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DedicatedServer{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list dedicated servers",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServer{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list dedicated servers with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServer{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list dedicated servers with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServer{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:        "list dedicated servers with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.DedicatedServer](mockCtrl)

	hostsServiceHandler.EXPECT().
		ListDedicatedServers().
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
	testServer1 := testKBM
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
		configureMock  func(*mocks.MockCollection[serverscom.KubernetesBaremetalNode])
	}{
		{
			name:           "list all KMB nodes",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesBaremetalNode]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.KubernetesBaremetalNode{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list KBM nodes",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesBaremetalNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesBaremetalNode{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list KBM nodes with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesBaremetalNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesBaremetalNode{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list KBM nodes with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesBaremetalNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesBaremetalNode{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:        "list KBM nodes with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesBaremetalNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.KubernetesBaremetalNode](mockCtrl)

	hostsServiceHandler.EXPECT().
		ListKubernetesBaremetalNodes().
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
	testServer1 := testSBM
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
		configureMock  func(*mocks.MockCollection[serverscom.SBMServer])
	}{
		{
			name:           "list all SBM servers",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMServer]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.SBMServer{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list SBM servers",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SBMServer{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list SBM servers with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SBMServer{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:           "list SBM servers with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SBMServer{
						testServer1,
						testServer2,
					}, nil)
			},
		},
		{
			name:        "list SBM servers with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.SBMServer](mockCtrl)

	hostsServiceHandler.EXPECT().
		ListSBMServers().
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
	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockHostsService)
	}{
		{
			name:           "release dedicated server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "release_ds_resp.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ScheduleReleaseForDedicatedServer(gomock.Any(), testId, serverscom.ScheduleReleaseInput{}).
					Return(&releasedServer, nil)
			},
		},
		{
			name:           "release dedicated server with --release-after",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "release_ds_scheduled_resp.json")),
			args:           []string{"--release-after", "2025-01-01T12:34:56+03:00"},
			configureMock: func(mock *mocks.MockHostsService) {
				releasedServer.ScheduledRelease = &fixedTime
				mock.EXPECT().
					ScheduleReleaseForDedicatedServer(gomock.Any(), testId, serverscom.ScheduleReleaseInput{ReleaseAfter: "2025-01-01T12:34:56+03:00"}).
					Return(&releasedServer, nil)
			},
		},
		{
			name:        "release dedicated server with error",
			id:          testId,
			expectError: true,
			configureMock: func(mock *mocks.MockHostsService) {
				releasedServer.ScheduledRelease = &fixedTime
				mock.EXPECT().
					ScheduleReleaseForDedicatedServer(gomock.Any(), testId, serverscom.ScheduleReleaseInput{}).
					Return(nil, errors.New("some error"))
			},
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

func TestGetDSNetworkCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		networkID      string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get DS network in default format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds_network.txt")),
		},
		{
			name:           "get DS network in JSON format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds_network.json")),
		},
		{
			name:           "get DS network in YAML format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds_network.yaml")),
		},
		{
			name:        "get DS network with error",
			id:          testId,
			networkID:   testNetworkId,
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
				GetDedicatedServerNetwork(gomock.Any(), tc.id, tc.networkID).
				Return(&testNetwork, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{
				"hosts", "ds", "get-network", tc.id,
				"--network-id", tc.networkID,
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

func TestListDSNetworksCmd(t *testing.T) {
	testNetwork1 := testNetwork
	testNetwork1.ID = testNetworkId
	testNetwork2 := testNetwork1
	testNetwork2.ID = "testNetId2"
	netTitle2 := "Another Net"
	testNetwork2.Title = &netTitle2

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Network])
	}{
		{
			name:           "list all DS networks",
			output:         "json",
			args:           []string{"testServerId", "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_networks_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
						testNetwork2,
					}, nil)
			},
		},
		{
			name:           "list DS networks",
			output:         "json",
			args:           []string{"testServerId"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_networks.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
					}, nil)
			},
		},
		{
			name:           "list DS networks with template",
			args:           []string{"testServerId", "--template", "{{range .}}Network: {{.ID}}  Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_networks_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
						testNetwork2,
					}, nil)
			},
		},
		{
			name:           "list DS networks with pageView",
			args:           []string{"testServerId", "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ds_networks_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
						testNetwork2,
					}, nil)
			},
		},
		{
			name:        "list DS networks with error",
			args:        []string{"testServerId"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Network](mockCtrl)

	hostsServiceHandler.EXPECT().
		DedicatedServerNetworks(gomock.Any()).
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

			args := []string{"hosts", "ds", "list-networks"}
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

func TestAddDSNetworkCmd(t *testing.T) {
	expectedInputPublic := serverscom.NetworkInput{
		DistributionMethod: "route",
		Mask:               32,
	}
	expectedInputPrivate := serverscom.NetworkInput{
		DistributionMethod: "gateway",
		Mask:               29,
	}

	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockHostsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:   "add public DS network",
			output: "json",
			args: []string{
				testId,
				"--type", "public",
				"--mask", "32",
				"--distribution-method", "route",
			},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds_network.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					AddDedicatedServerPublicIPv4Network(gomock.Any(), testId, expectedInputPublic).
					Return(&testNetwork, nil)
			},
		},
		{
			name:   "add private DS network",
			output: "json",
			args: []string{
				testId,
				"--type", "private",
				"--mask", "29",
				"--distribution-method", "gateway",
			},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds_network.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					AddDedicatedServerPrivateIPv4Network(gomock.Any(), testId, expectedInputPrivate).
					Return(&testNetwork, nil)
			},
		},
		{
			name: "add DS network with unsupported mask",
			args: []string{
				testId,
				"--type", "public",
				"--mask", "24",
				"--distribution-method", "gateway",
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
			hostsCmd := NewCmd(testCmdContext)

			args := append([]string{"hosts", "ds", "add-network"}, tc.args...)
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

func TestDeleteDSNetworkCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		networkID      string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "delete DS network",
			id:             testId,
			networkID:      testNetworkId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ds_network.json")),
		},
		{
			name:        "delete DS network with error",
			id:          testId,
			networkID:   testNetworkId,
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
				DeleteDedicatedServerNetwork(gomock.Any(), tc.id, tc.networkID).
				Return(&testNetwork, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "delete-network", tc.id, "--network-id", tc.networkID}
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

func TestListKBMPowerFeedsCmd(t *testing.T) {
	testPowerFeed1 := testPowerFeed
	testPowerFeed2 := testPowerFeed1

	testPowerFeed2.Name = "testPowerFeed2"
	testPowerFeed2.Status = "off"

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.HostPowerFeed])
	}{
		{
			name:           "get KBM node power_feeds in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.txt")),
		},
		{
			name:           "get KBM node power_feeds",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.json")),
		},
		{
			name:           "get KBM node power_feeds in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.yaml")),
		},
		{
			name:        "get KBM node power_feeds with error",
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
				KubernetesBaremetalNodePowerFeeds(gomock.Any(), testId).
				Return([]serverscom.HostPowerFeed{testPowerFeed1, testPowerFeed2}, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "kbm", "list-power-feeds", tc.id}
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

func TestListDSPowerFeedsCmd(t *testing.T) {
	testPowerFeed1 := testPowerFeed
	testPowerFeed2 := testPowerFeed1

	testPowerFeed2.Name = "testPowerFeed2"
	testPowerFeed2.Status = "off"

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.HostPowerFeed])
	}{
		{
			name:           "get ds power_feeds in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.txt")),
		},
		{
			name:           "get ds power_feeds",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.json")),
		},
		{
			name:           "get ds power_feeds in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.yaml")),
		},
		{
			name:        "get ds power_feeds with error",
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
				DedicatedServerPowerFeeds(gomock.Any(), testId).
				Return([]serverscom.HostPowerFeed{testPowerFeed1, testPowerFeed2}, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "ds", "list-power-feeds", tc.id}
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

func TestListSBMPowerFeedsCmd(t *testing.T) {
	testPowerFeed1 := testPowerFeed
	testPowerFeed2 := testPowerFeed1

	testPowerFeed2.Name = "testPowerFeed2"
	testPowerFeed2.Status = "off"

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.HostPowerFeed])
	}{
		{
			name:           "get sbm power_feeds in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.txt")),
		},
		{
			name:           "get sbm power_feeds",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.json")),
		},
		{
			name:           "get sbm power_feeds in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_power_feeds.yaml")),
		},
		{
			name:        "get sbm power_feeds with error",
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
				SBMServerPowerFeeds(gomock.Any(), testId).
				Return([]serverscom.HostPowerFeed{testPowerFeed1, testPowerFeed2}, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"hosts", "sbm", "list-power-feeds", tc.id}
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

func TestListKBMNetworksCmd(t *testing.T) {
	testNetwork1 := testNetwork
	testNetwork1.ID = testNetworkId
	testNetwork2 := testNetwork1
	testNetwork2.ID = "testNetId2"
	netTitle2 := "Another Net"
	testNetwork2.Title = &netTitle2

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Network])
	}{
		{
			name:           "list KBM node all networks",
			output:         "json",
			args:           []string{"testServerId", "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_networks_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
						testNetwork2,
					}, nil)
			},
		},
		{
			name:           "list KBM node networks",
			output:         "json",
			args:           []string{"testServerId"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_networks.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
					}, nil)
			},
		},
		{
			name:           "list DS networks with template",
			args:           []string{"testServerId", "--template", "{{range .}}Network: {{.ID}}  Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_networks_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
						testNetwork2,
					}, nil)
			},
		},
		{
			name:           "list DS networks with pageView",
			args:           []string{"testServerId", "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_kbm_networks_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
						testNetwork2,
					}, nil)
			},
		},
		{
			name:        "list DS networks with error",
			args:        []string{"testServerId"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Network](mockCtrl)

	hostsServiceHandler.EXPECT().
		KubernetesBaremetalNodeNetworks(gomock.Any()).
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

			args := []string{"hosts", "kbm", "list-networks"}
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

func TestListDSDriveSlotsCmd(t *testing.T) {
	testDriveSlot1 := testDriveSlot
	testDriveSlot2 := testDriveSlot1
	testDriveSlot2.Position = 2

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.HostDriveSlot])
	}{
		{
			name:           "list ds all drive slots",
			output:         "json",
			args:           []string{"testServerId", "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
						testDriveSlot2,
					}, nil)
			},
		},
		{
			name:           "list ds drive slots",
			output:         "json",
			args:           []string{"testServerId"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
					}, nil)
			},
		},
		{
			name:           "list ds drive slots with template",
			args:           []string{"testServerId", "--template", "{{range .}}Position: {{.Position}}  Interface: {{.Interface}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
						testDriveSlot2,
					}, nil)
			},
		},
		{
			name:           "list ds drive slots with pageView",
			args:           []string{"testServerId", "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
						testDriveSlot2,
					}, nil)
			},
		},
		{
			name:        "list ds drive slots with error",
			args:        []string{"testServerId"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.HostDriveSlot](mockCtrl)

	hostsServiceHandler.EXPECT().
		DedicatedServerDriveSlots(gomock.Any()).
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

			args := []string{"hosts", "ds", "list-drive-slots"}
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

func TestListKBMDriveSlotsCmd(t *testing.T) {
	testDriveSlot1 := testDriveSlot
	testDriveSlot2 := testDriveSlot1
	testDriveSlot2.Position = 2

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.HostDriveSlot])
	}{
		{
			name:           "list KBM node all drive slots",
			output:         "json",
			args:           []string{"testServerId", "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
						testDriveSlot2,
					}, nil)
			},
		},
		{
			name:           "list KBM node drive slots",
			output:         "json",
			args:           []string{"testServerId"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
					}, nil)
			},
		},
		{
			name:           "list KBM node drive slots with template",
			args:           []string{"testServerId", "--template", "{{range .}}Position: {{.Position}}  Interface: {{.Interface}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
						testDriveSlot2,
					}, nil)
			},
		},
		{
			name:           "list KBM node drive slots with pageView",
			args:           []string{"testServerId", "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_drive_slots_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
						testDriveSlot2,
					}, nil)
			},
		},
		{
			name:        "list KBM node drive slots with error",
			args:        []string{"testServerId"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.HostDriveSlot](mockCtrl)

	hostsServiceHandler.EXPECT().
		KubernetesBaremetalNodeDriveSlots(gomock.Any()).
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

			args := []string{"hosts", "kbm", "list-drive-slots"}
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

func TestKBMPowerCmd(t *testing.T) {
	testServer := testKBM
	testServer.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "power on kbm node",
			id:             testId,
			args:           []string{"hosts", "kbm", "power", testId, "--command=on"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_kbm_power.json")),
		},
		{
			name:           "power off kbm node",
			id:             testId,
			args:           []string{"hosts", "kbm", "power", testId, "--command=off"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_kbm_power.json")),
		},
		{
			name:           "power cycle kbm node",
			id:             testId,
			args:           []string{"hosts", "kbm", "power", testId, "--command=cycle"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_kbm_power.json")),
		},
		{
			name:        "power on kbm node with error",
			id:          testId,
			args:        []string{"hosts", "kbm", "power", testId, "--command=on"},
			expectError: true,
		},
		{
			name:        "power on kbm node without flag error",
			id:          testId,
			args:        []string{"hosts", "kbm", "power", testId},
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

			var powerVal string
			if tc.args != nil {
				powerVal = strings.TrimPrefix(tc.args[len(tc.args)-1], "--command=")
			}
			expectPowerCall(hostsServiceHandler, powerVal, testId, &testServer, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			if tc.output != "" {
				tc.args = append(tc.args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(hostsCmd).
				WithArgs(tc.args)

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

func expectPowerCall(m *mocks.MockHostsService, action string, id string, s *serverscom.KubernetesBaremetalNode, err error) {
	calls := map[string]func() *gomock.Call{
		"on":    func() *gomock.Call { return m.EXPECT().PowerOnKubernetesBaremetalNode(gomock.Any(), id) },
		"off":   func() *gomock.Call { return m.EXPECT().PowerOffKubernetesBaremetalNode(gomock.Any(), id) },
		"cycle": func() *gomock.Call { return m.EXPECT().PowerCycleKubernetesBaremetalNode(gomock.Any(), id) },
	}

	if action == "" {
		m.EXPECT().
			PowerOnKubernetesBaremetalNode(gomock.Any(), id).
			Times(0)
		m.EXPECT().
			PowerOffKubernetesBaremetalNode(gomock.Any(), id).
			Times(0)
		m.EXPECT().
			PowerCycleKubernetesBaremetalNode(gomock.Any(), id).
			Times(0)

		gomock.InOrder(
			calls["on"]().Times(0),
			calls["off"]().Times(0),
			calls["cycle"]().Times(0),
		)

		return
	}

	if h, ok := calls[action]; ok {
		h().Return(s, err).Times(1)
	}
}
