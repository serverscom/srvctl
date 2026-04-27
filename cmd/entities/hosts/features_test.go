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
	featuresFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "hosts", "features")

	testFeatureResult = serverscom.DedicatedServerFeature{
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
			cmd := NewCmd(testCmdContext)

			args := append([]string{"hosts", "ebm", "list-features"}, tc.args...)
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
			args:           []string{testServerID, "--feature", "disaggregated_public_ports", "--state", "activate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivateDisaggregatedPublicPortsFeature(gomock.Any(), testServerID).
					Return(&testFeatureResult, nil)
			},
		},
		{
			name: "deactivate feature",
			args: []string{testServerID, "--feature", "disaggregated_public_ports", "--state", "deactivate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					DeactivateDisaggregatedPublicPortsFeature(gomock.Any(), testServerID).
					Return(&testFeatureResult, nil)
			},
		},
		{
			name:           "activate private_ipxe_boot with ipxe-config",
			args:           []string{testServerID, "--feature", "private_ipxe_boot", "--state", "activate", "--ipxe-config", "#!ipxe\nchain http://boot.example.com", "--output", "json"},
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
			args:           []string{testServerID, "--feature", "private_ipxe_boot", "--state", "deactivate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_private_ipxe_boot.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					DeactivatePrivateIpxeBootFeature(gomock.Any(), testServerID).
					Return(&testPrivateIpxeBootResult, nil)
			},
		},
		{
			name:           "activate host_rescue_mode with password auth",
			args:           []string{testServerID, "--feature", "host_rescue_mode", "--state", "activate", "--auth-methods", "password", "--output", "json"},
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
			args:           []string{testServerID, "--feature", "host_rescue_mode", "--state", "activate", "--auth-methods", "ssh_key", "--ssh-key-fingerprints", "aa:bb:cc", "--output", "json"},
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
			args:           []string{testServerID, "--feature", "host_rescue_mode", "--state", "deactivate", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(featuresFixtureBasePath, "feature_set_host_rescue_mode.json")),
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					DeactivateHostRescueModeFeature(gomock.Any(), testServerID).
					Return(&testHostRescueModeResult, nil)
			},
		},
		{
			name:        "api error",
			args:        []string{testServerID, "--feature", "disaggregated_public_ports", "--state", "activate"},
			expectError: true,
			configureMock: func(mock *mocks.MockHostsService) {
				mock.EXPECT().
					ActivateDisaggregatedPublicPortsFeature(gomock.Any(), testServerID).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "invalid state",
			args:        []string{testServerID, "--feature", "disaggregated_public_ports", "--state", "invalid"},
			expectError: true,
		},
		{
			name:        "unsupported feature",
			args:        []string{testServerID, "--feature", "unknown_feature", "--state", "activate"},
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
			hostsCmd := NewCmd(testCmdContext)

			args := append([]string{"hosts", "ebm", "feature-set"}, tc.args...)
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
