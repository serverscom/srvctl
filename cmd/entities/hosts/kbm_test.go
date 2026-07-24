package hosts

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
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
)

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
			kbmCmd := NewKBMCmd(testCmdContext)

			args := []string{"kbm", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(kbmCmd).
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
			name:           "list all KBM nodes",
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
			kbmCmd := NewKBMCmd(testCmdContext)

			args := []string{"kbm", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(kbmCmd).
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
			sshCmd := NewKBMCmd(testCmdContext)

			args := []string{"kbm", "update", tc.id}
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
			kbmCmd := NewKBMCmd(testCmdContext)

			args := []string{"kbm", "list-power-feeds", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(kbmCmd).
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
			kbmCmd := NewKBMCmd(testCmdContext)

			args := []string{"kbm", "list-networks"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(kbmCmd).
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
			kbmCmd := NewKBMCmd(testCmdContext)

			args := []string{"kbm", "list-drive-slots"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(kbmCmd).
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
			args:           []string{"kbm", "power", testId, "--command=on"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_kbm_power.json")),
		},
		{
			name:           "power off kbm node",
			id:             testId,
			args:           []string{"kbm", "power", testId, "--command=off"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_kbm_power.json")),
		},
		{
			name:           "power cycle kbm node",
			id:             testId,
			args:           []string{"kbm", "power", testId, "--command=cycle"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "set_kbm_power.json")),
		},
		{
			name:        "power on kbm node with error",
			id:          testId,
			args:        []string{"kbm", "power", testId, "--command=on"},
			expectError: true,
		},
		{
			name:        "power on kbm node without flag error",
			id:          testId,
			args:        []string{"kbm", "power", testId},
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
			kbmCmd := NewKBMCmd(testCmdContext)

			if tc.output != "" {
				tc.args = append(tc.args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(kbmCmd).
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
