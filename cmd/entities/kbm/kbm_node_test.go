package kbm

import (
	"errors"
	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	testId            = "testId"
	testNetworkId     = "testNetId"
	fixtureBasePath   = filepath.Join("..", "..", "..", "testdata", "entities", "kbm")
	fixedTime         = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testPublicIP      = "1.2.3.4"
	testLocationCode  = "test"
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
	}
)

func TestNewUpdateCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_node.json")),
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

			args := []string{"kbm", "update-node", tc.id}
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

func TestNewListCmd(t *testing.T) {
	testServer1 := testKBM
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_nodes_all.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_nodes.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_nodes_template.txt")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_nodes_pageview.txt")),
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

			args := []string{"kbm", "list-nodes"}
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

func TestNewPowerCmd(t *testing.T) {
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
			args:           []string{"kbm", "set-node-power", testId, "--power=on"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_node_power.json")),
		},
		{
			name:           "power off kbm node",
			id:             testId,
			args:           []string{"kbm", "set-node-power", testId, "--power=off"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_node_power.json")),
		},
		{
			name:           "power cycle kbm node",
			id:             testId,
			args:           []string{"kbm", "set-node-power", testId, "--power=cycle"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_node_power.json")),
		},
		{
			name:        "power on kbm node with error",
			id:          testId,
			args:        []string{"kbm", "set-node-power", testId, "--power=on"},
			expectError: true,
		},
		{
			name:        "power on kbm node without flag error",
			id:          testId,
			args:        []string{"kbm", "set-node-power", testId},
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
				powerVal = strings.TrimPrefix(tc.args[len(tc.args)-1], "--power=")
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

func TestNewGetCmd(t *testing.T) {
	testServer := testKBM
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node.txt")),
		},
		{
			name:           "get KBM node in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node.json")),
		},
		{
			name:           "get KBM node in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node.yaml")),
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
				Return(&testServer, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			hostsCmd := NewCmd(testCmdContext)

			args := []string{"kbm", "get-node", tc.id}
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

func TestNewListNetworksCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_networks_all.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_networks.json")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_networks_template.txt")),
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_networks_pageview.txt")),
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

			args := []string{"kbm", "list-node-networks"}
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

func TestNewListSlotsCmd(t *testing.T) {
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
			name:           "list KBM node all slots",
			output:         "json",
			args:           []string{"testServerId", "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_slots_all.json")),
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
			name:           "list KBM node slots",
			output:         "json",
			args:           []string{"testServerId"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_slots.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.HostDriveSlot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.HostDriveSlot{
						testDriveSlot1,
					}, nil)
			},
		},
		{
			name:           "list KBM node slots with template",
			args:           []string{"testServerId", "--template", "{{range .}}Position: {{.Position}}  Interface: {{.Interface}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_slots_template.txt")),
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
			name:           "list KBM node slots with pageView",
			args:           []string{"testServerId", "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_node_slots_pageview.txt")),
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
			name:        "list KBM node slots with error",
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

			args := []string{"kbm", "list-node-slots"}
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

func TestNewGetPowerFeedsCmd(t *testing.T) {
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
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node_power_feeds.txt")),
		},
		{
			name:           "get KBM node power_feeds",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node_power_feeds.json")),
		},
		{
			name:           "get KBM node power_feeds in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node_power_feeds.yaml")),
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

			args := []string{"kbm", "list-node-feeds", tc.id}
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
