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
		ServerModelID:           new(int64(1)),
		ServerModelName:         new("server-model-123"),
		PublicUplinkID:          new(int64(2)),
		PublicUplinkName:        new("Public 1 Gbps without redundancy"),
		PrivateUplinkID:         new(int64(3)),
		PrivateUplinkName:       new("Private 1 Gbps without redundancy"),
		BandwidthID:             new(int64(4)),
		BandwidthName:           new("20000 GB"),
		OperatingSystemID:       new(int64(5)),
		OperatingSystemFullName: new("CentOS 7 x86_64"),
	}

	testNetwork = serverscom.Network{
		ID:                 testNetworkId,
		Title:              new("Some Net"),
		Status:             "active",
		Cidr:               new("100.0.8.0/29"),
		Family:             "ipv4",
		InterfaceType:      "public",
		DistributionMethod: "gateway",
		Additional:         false,
		FirstIP:            new("100.0.8.2"),
		Gateway:            new("100.0.8.1"),
		Created:            fixedTime,
		Updated:            fixedTime,
	}
	testNetworkUsage = serverscom.NetworkUsage{
		Type: "public",
		Utilization: &serverscom.Utilization{
			Value:  100,
			Commit: 50,
			Unit:   "GB",
		},
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
