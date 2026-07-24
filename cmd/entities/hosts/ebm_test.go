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

	servicesFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "hosts", "services")
	testDSService           = serverscom.DedicatedServerService{
		ID:            testId,
		Name:          "Test service",
		Type:          "server_base",
		Currency:      "USD",
		StartedAt:     fixedTime,
		FinishedAt:    fixedTime,
		Total:         100.0,
		UsageQuantity: 2.0,
		Tax:           10.0,
		Subtotal:      100.0,
		DiscountRate:  5.0,
		DateFrom:      "2025-11-01",
		DateTo:        "2025-11-30",
	}

	oobFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "hosts", "oob")
	testOobCreds       = serverscom.DedicatedServerOOBCredentials{
		Login:  "test",
		Secret: "secret",
	}

	featuresFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "hosts", "features")
	testFeatureResult       = serverscom.DedicatedServerFeature{
		Name:   "disaggregated_public_ports",
		Status: "activation",
	}
	testPrivateIpxeBootResult = serverscom.DedicatedServerFeature{
		Name:   "private_ipxe_boot",
		Status: "activation",
	}
	testHostRescueModeResult = serverscom.DedicatedServerFeature{
		Name:   "host_rescue_mode",
		Status: "activation",
	}
)

func TestAddEBMCmd(t *testing.T) {
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
					DriveModelID: new(int64(3456)),
				},
				{
					Position:     2,
					DriveModelID: new(int64(3456)),
				},
			},
			Layout: []serverscom.DedicatedServerLayoutInput{
				{
					SlotPositions: []int{1, 2},
					Raid:          new(1),
					Partitions: []serverscom.DedicatedServerLayoutPartitionInput{
						{
							Target: "/boot",
							Size:   500,
							Fill:   false,
							Fs:     new("ext4"),
						},
					},
				},
			},
		},
		Hosts: []serverscom.DedicatedServerHostInput{
			{
				Hostname:             "example.aa",
				PublicIPv4NetworkID:  new("PublicNet123"),
				PrivateIPv4NetworkID: new("PrivateNet456"),
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
			name:           "create ebm server",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ebm_resp.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create_ebm_input.json")},
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), expectedInput).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:           "create ebm server with merge input with flags",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ebm_resp.json")),
			args: []string{
				"--input", filepath.Join(fixtureBasePath, "create_ebm_input.json"),
				"--layout", "slot=3,slot=4,raid=0",
				"--partition", "slot=3,slot=4,target=/boot,fs=ext4,size=500",
			},
			configureMock: func(mock *mocks.MockHostsService) {
				input := expectedInput
				input.Drives.Layout = append(input.Drives.Layout, serverscom.DedicatedServerLayoutInput{
					SlotPositions: []int{3, 4},
					Raid:          new(0),
					Partitions: []serverscom.DedicatedServerLayoutPartitionInput{
						{
							Target: "/boot",
							Size:   500,
							Fill:   false,
							Fs:     new("ext4"),
						},
					},
				})

				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), input).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:           "skeleton for ebm server input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add_ebm.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:           "create ebm server with public_ipxe_boot feature and ipxe-config",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ebm_resp.json")),
			args: []string{
				"--input", filepath.Join(fixtureBasePath, "create_ebm_input.json"),
				"--feature", "public_ipxe_boot",
				"--ipxe-config", "#!ipxe\nboot",
			},
			configureMock: func(mock *mocks.MockHostsService) {
				input := expectedInput
				input.Features = []string{"public_ipxe_boot"}
				input.IPXEConfig = new("#!ipxe\nboot")
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), input).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:           "create ebm server with ipxe-config flag",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ebm_resp.json")),
			args: []string{
				"--input", filepath.Join(fixtureBasePath, "create_ebm_input.json"),
				"--ipxe-config", "#!ipxe\nboot",
			},
			configureMock: func(mock *mocks.MockHostsService) {
				input := expectedInput
				input.IPXEConfig = new("#!ipxe\nboot")
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), input).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:           "create ebm server with ipxe-config in input file",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_ebm_resp.json")),
			args: []string{
				"--input", filepath.Join(fixtureBasePath, "create_ebm_ipxe_input.json"),
			},
			configureMock: func(mock *mocks.MockHostsService) {
				input := expectedInput
				input.IPXEConfig = new("#!ipxe\nboot")
				mock.EXPECT().
					CreateDedicatedServers(gomock.Any(), input).
					Return([]serverscom.DedicatedServer{testDS}, nil)
			},
		},
		{
			name:        "create ebm server with error",
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestGetEBMCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get ebm server in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm.txt")),
		},
		{
			name:           "get ebm server in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm.json")),
		},
		{
			name:           "get ebm server in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm.yaml")),
		},
		{
			name:        "get ebm server with error",
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestListEBMCmd(t *testing.T) {
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
			name:           "list all ebm servers",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_all.json")),
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
			name:           "list ebm servers",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServer{
						testServer1,
					}, nil)
			},
		},
		{
			name:           "list ebm servers with template",
			args:           []string{"--template", "{{range .}}Title: {{.Title}}  Type: {{.Type}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_template.txt")),
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
			name:           "list ebm servers with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_pageview.txt")),
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
			name:        "list ebm servers with error",
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestUpdateEBMCmd(t *testing.T) {
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
			name:           "update ebm server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_ebm_resp.json")),
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
			name: "update ebm server with error",
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
			sshCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "update", tc.id}
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

func TestScheduleReleaseEBMCmd(t *testing.T) {
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
			name:           "release ebm server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "release_ebm_resp.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ScheduleReleaseForDedicatedServer(gomock.Any(), testId, serverscom.ScheduleReleaseInput{}).
					Return(&releasedServer, nil)
			},
		},
		{
			name:           "release ebm server with --release-after",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "release_ebm_scheduled_resp.json")),
			args:           []string{"--release-after", "2025-01-01T12:34:56+03:00"},
			configureMock: func(mock *mocks.MockHostsService) {
				releasedServer.ScheduledRelease = &fixedTime
				mock.EXPECT().
					ScheduleReleaseForDedicatedServer(gomock.Any(), testId, serverscom.ScheduleReleaseInput{ReleaseAfter: "2025-01-01T12:34:56+03:00"}).
					Return(&releasedServer, nil)
			},
		},
		{
			name:        "release ebm server with error",
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "schedule-release", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestAbortReleaseEBMCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "abort release ebm server",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm.json")),
		},
		{
			name:        "abort release ebm server with error",
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "abort-release", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestGetEBMNetworkCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.txt")),
		},
		{
			name:           "get DS network in JSON format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.json")),
		},
		{
			name:           "get DS network in YAML format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.yaml")),
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{
				"ebm", "get-network", tc.id,
				"--network-id", tc.networkID,
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestListEBMNetworksCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_networks_all.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_networks.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_networks_template.txt")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_ebm_networks_pageview.txt")),
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "list-networks"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestAddEBMNetworkCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.json")),
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := append([]string{"ebm", "add-network"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestDeleteEBMNetworkCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.json")),
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "delete-network", tc.id, "--network-id", tc.networkID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestListEBMPowerFeedsCmd(t *testing.T) {
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "list-power-feeds", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestActivateEBMIPv6NetworkCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "activate DS IPv6 network",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network.json")),
		},
		{
			name:        "activate DS IPv6 network with error",
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
				ActivateDedicatedServerPubliIPv6Network(gomock.Any(), tc.id).
				Return(&testNetwork, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "activate-ipv6-network", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestGetEBMNetworkUsageCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get DS network usage in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network_usage.txt")),
		},
		{
			name:           "get DS network usage in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network_usage.json")),
		},
		{
			name:           "get DS network usage in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_ebm_network_usage.yaml")),
		},
		{
			name:        "get DS network usage with error",
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
				GetDedicatedServerNetworkUsage(gomock.Any(), tc.id).
				Return(&testNetworkUsage, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "network-usage", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestListEBMDriveSlotsCmd(t *testing.T) {
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "list-drive-slots"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestListEBMServicesCmd(t *testing.T) {
	testService1 := testDSService
	testService2 := testDSService
	testService2.ID += "2"
	testService2.Name = "Test service 2"
	testService2.Type = "uplink"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.DedicatedServerService])
	}{
		{
			name:           "list all ds services",
			output:         "json",
			args:           []string{"-A", testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
						testService2,
					}, nil)
			},
		},
		{
			name:           "list ds services",
			output:         "json",
			args:           []string{testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
					}, nil)
			},
		},
		{
			name:           "list ds services with template",
			args:           []string{testServerID, "--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
						testService2,
					}, nil)
			},
		},
		{
			name:           "list ds services with page-view",
			args:           []string{testServerID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
						testService2,
					}, nil)
			},
		},
		{
			name:        "list ds services with error",
			args:        []string{testServerID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.DedicatedServerService](mockCtrl)

	hostService.EXPECT().DedicatedServerServices(gomock.Any()).Return(collection).AnyTimes()
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
			ebmCmd := NewEBMCmd(testCmdContext)

			args := append([]string{"ebm", "list-services"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestGetEBMOOBCredsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get ebm server oob creds in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(oobFixtureBasePath, "get.txt")),
		},
		{
			name:           "get ebm server oob creds in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(oobFixtureBasePath, "get.json")),
		},
		{
			name:           "get ebm server oob creds in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(oobFixtureBasePath, "get.yaml")),
		},
		{
			name:        "get ebm server oob creds with error",
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
				GetDedicatedServerOOBCredentials(gomock.Any(), testId, map[string]string{"fingerprint": "test"}).
				Return(&testOobCreds, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			ebmCmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "get-oob-credentials", tc.id, "--fingerprint", "test"}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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

func TestListEBMPTRCmd(t *testing.T) {
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
			cmd := NewEBMCmd(testCmdContext)

			args := append([]string{"ebm", "list-ptr"}, tc.args...)
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

func TestCreateEBMPTRCmd(t *testing.T) {
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
			cmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "add-ptr"}
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

func TestDeleteEBMPTRCmd(t *testing.T) {
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
			cmd := NewEBMCmd(testCmdContext)

			args := []string{"ebm", "delete-ptr", tc.serverID, "--ptr-id", tc.ptrID}
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

func TestListEBMFeaturesCmd(t *testing.T) {
	testFeature1 := serverscom.DedicatedServerFeature{
		Name:   "disaggregated_public_ports",
		Status: "deactivated",
	}
	testFeature2 := serverscom.DedicatedServerFeature{
		Name:   "no_public_network",
		Status: "unavailable",
	}

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.DedicatedServerFeature])
	}{
		{
			name:           "list all ds features",
			output:         "json",
			args:           []string{"-A", testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerFeature]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DedicatedServerFeature{
						testFeature1,
						testFeature2,
					}, nil)
			},
		},
		{
			name:           "list ds features",
			output:         "json",
			args:           []string{testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerFeature]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerFeature{
						testFeature1,
					}, nil)
			},
		},
		{
			name:           "list ds features with template",
			args:           []string{testServerID, "--template", "{{range .}}Name: {{.Name}} Status: {{.Status}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerFeature]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerFeature{
						testFeature1,
						testFeature2,
					}, nil)
			},
		},
		{
			name:           "list ds features with page-view",
			args:           []string{testServerID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerFeature]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerFeature{
						testFeature1,
						testFeature2,
					}, nil)
			},
		},
		{
			name:        "list ds features with error",
			args:        []string{testServerID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerFeature]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.DedicatedServerFeature](mockCtrl)

	hostService.EXPECT().DedicatedServerFeatures(gomock.Any()).Return(collection).AnyTimes()
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
			cmd := NewEBMCmd(testCmdContext)

			args := append([]string{"ebm", "list-features"}, tc.args...)
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

func TestEBMFeatureSetCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockHostsService)
	}{
		{
			name:           "activate feature",
			args:           []string{testServerID, "--feature", "disaggregated_public_ports", "--command", "activate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivateDisaggregatedPublicPortsFeature(gomock.Any(), testServerID).
					Return(&testFeatureResult, nil)
			},
		},
		{
			name:           "deactivate feature",
			args:           []string{testServerID, "--feature", "disaggregated_public_ports", "--command", "deactivate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					DeactivateDisaggregatedPublicPortsFeature(gomock.Any(), testServerID).
					Return(&testFeatureResult, nil)
			},
		},
		{
			name:           "activate private_ipxe_boot with ipxe-config",
			args:           []string{testServerID, "--feature", "private_ipxe_boot", "--command", "activate", "--ipxe-config", "#!ipxe\nchain http://boot.example.com", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_private_ipxe_boot.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivatePrivateIpxeBootFeature(gomock.Any(), testServerID, serverscom.PrivateIpxeBootFeatureInput{
						IPXEConfig: "#!ipxe\nchain http://boot.example.com",
					}).
					Return(&testPrivateIpxeBootResult, nil)
			},
		},
		{
			name:           "deactivate private_ipxe_boot",
			args:           []string{testServerID, "--feature", "private_ipxe_boot", "--command", "deactivate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_private_ipxe_boot.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					DeactivatePrivateIpxeBootFeature(gomock.Any(), testServerID).
					Return(&testPrivateIpxeBootResult, nil)
			},
		},
		{
			name:           "activate host_rescue_mode with password auth",
			args:           []string{testServerID, "--feature", "host_rescue_mode", "--command", "activate", "--auth-method", "password", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_host_rescue_mode.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivateHostRescueModeFeature(gomock.Any(), testServerID, serverscom.HostRescueModeFeatureInput{
						AuthMethods: []string{"password"},
					}).
					Return(&testHostRescueModeResult, nil)
			},
		},
		{
			name:           "activate host_rescue_mode with ssh_key auth",
			args:           []string{testServerID, "--feature", "host_rescue_mode", "--command", "activate", "--auth-method", "ssh_key", "--ssh-key-fingerprint", "aa:bb:cc", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_host_rescue_mode.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivateHostRescueModeFeature(gomock.Any(), testServerID, serverscom.HostRescueModeFeatureInput{
						AuthMethods:        []string{"ssh_key"},
						SSHKeyFingerprints: []string{"aa:bb:cc"},
					}).
					Return(&testHostRescueModeResult, nil)
			},
		},
		{
			name:           "deactivate host_rescue_mode",
			args:           []string{testServerID, "--feature", "host_rescue_mode", "--command", "deactivate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_host_rescue_mode.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					DeactivateHostRescueModeFeature(gomock.Any(), testServerID).
					Return(&testHostRescueModeResult, nil)
			},
		},
		{
			name:        "api error",
			args:        []string{testServerID, "--feature", "disaggregated_public_ports", "--command", "activate"},
			expectError: true,
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivateDisaggregatedPublicPortsFeature(gomock.Any(), testServerID).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "invalid command",
			args:        []string{testServerID, "--feature", "disaggregated_public_ports", "--command", "invalid"},
			expectError: true,
		},
		{
			name:        "unsupported feature",
			args:        []string{testServerID, "--feature", "unknown_feature", "--command", "activate"},
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
			if tc.configureMock != nil {
				tc.configureMock(hostService)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			ebmCmd := NewEBMCmd(testCmdContext)

			args := append([]string{"ebm", "feature-set"}, tc.args...)
			builder := testutils.NewTestCommandBuilder().
				WithCommand(ebmCmd).
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
