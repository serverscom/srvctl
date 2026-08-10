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
)

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
								PublicIPv4NetworkID:  new("PublicNetTest123"),
								PrivateIPv4NetworkID: new("PrivateNetTest456"),
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
			name:           "create SBM server with params",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "create_sbm_resp.json")),
			args: []string{
				"--location-id", "5678",
				"--sbm-flavor-model-id", "1234",
				"--operating-system-id", "9999",
				"--ssh-key-fingerprint", "48:81:0c:43:99:12:71:5e:ba:fd:e7:2f:20:d7:95:e8",
				"example.aa",
			},
			configureMock: func(mock *mocks.MockHostsService) {
				osID := int64(9999)
				mock.EXPECT().
					CreateSBMServers(gomock.Any(), serverscom.SBMServerCreateInput{
						FlavorModelID:      1234,
						LocationID:         5678,
						OperatingSystemID:  &osID,
						SSHKeyFingerprints: []string{"48:81:0c:43:99:12:71:5e:ba:fd:e7:2f:20:d7:95:e8"},
						Hosts: []serverscom.SBMServerHostInput{
							{
								Hostname: "example.aa",
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sshCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "update", tc.id}
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "release", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "list-power-feeds", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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

func TestGetSBMNetworkCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		networkID      string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get SBM network in default format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network.txt")),
		},
		{
			name:           "get SBM network in JSON format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network.json")),
		},
		{
			name:           "get SBM network in YAML format",
			id:             testId,
			networkID:      testNetworkId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network.yaml")),
		},
		{
			name:        "get SBM network with error",
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
				GetSBMServerNetwork(gomock.Any(), tc.id, tc.networkID).
				Return(&testNetwork, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{
				"sbm", "get-network", tc.id,
				"--network-id", tc.networkID,
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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

func TestListSBMNetworksCmd(t *testing.T) {
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
			name:           "list all SBM networks",
			output:         "json",
			args:           []string{"testServerId", "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_networks_all.json")),
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
			name:           "list SBM networks",
			output:         "json",
			args:           []string{"testServerId"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_networks.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{
						testNetwork1,
					}, nil)
			},
		},
		{
			name:           "list SBM networks with template",
			args:           []string{"testServerId", "--template", "{{range .}}Network: {{.ID}}  Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_networks_template.txt")),
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
			name:           "list SBM networks with pageView",
			args:           []string{"testServerId", "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_sbm_networks_pageview.txt")),
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
			name:        "list SBM networks with error",
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
		SBMServerNetworks(gomock.Any()).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "list-networks"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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

func TestAddSBMNetworkCmd(t *testing.T) {
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
			name:   "add private SBM network",
			output: "json",
			args: []string{
				testId,
				"--type", "private",
				"--mask", "29",
				"--distribution-method", "gateway",
			},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					AddSBMServerPrivateIPv4Network(gomock.Any(), testId, expectedInputPrivate).
					Return(&testNetwork, nil)
			},
		},
		{
			name: "add public SBM network is not supported",
			args: []string{
				testId,
				"--type", "public",
				"--mask", "29",
			},
			expectError: true,
		},
		{
			name: "add SBM network with unsupported mask",
			args: []string{
				testId,
				"--type", "private",
				"--mask", "24",
			},
			expectError: true,
		},
		{
			name:   "add SBM network with error",
			output: "json",
			args: []string{
				testId,
				"--type", "private",
				"--mask", "29",
			},
			expectError: true,
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					AddSBMServerPrivateIPv4Network(gomock.Any(), testId, expectedInputPrivate).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := append([]string{"sbm", "add-network"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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

func TestDeleteSBMNetworkCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		networkID      string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "delete SBM network",
			id:             testId,
			networkID:      testNetworkId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network.json")),
		},
		{
			name:        "delete SBM network with error",
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
				DeleteSBMServerNetwork(gomock.Any(), tc.id, tc.networkID).
				Return(&testNetwork, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "delete-network", tc.id, "--network-id", tc.networkID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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

func TestGetSBMNetworkUsageCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get SBM network usage in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network_usage.txt")),
		},
		{
			name:           "get SBM network usage in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network_usage.json")),
		},
		{
			name:           "get SBM network usage in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_sbm_network_usage.yaml")),
		},
		{
			name:        "get SBM network usage with error",
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
				GetSBMServerNetworkUsage(gomock.Any(), tc.id).
				Return(&testNetworkUsage, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "network-usage", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := append([]string{"sbm", "list-ptr"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "add-ptr"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
			sbmCmd := NewSBMCmd(testCmdContext)

			args := []string{"sbm", "delete-ptr", tc.serverID, "--ptr-id", tc.ptrID}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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
