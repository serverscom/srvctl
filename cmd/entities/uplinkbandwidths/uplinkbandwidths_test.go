package uplinkbandwidths

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	bandwidthFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "uplink-bandwidths")
	testBandwidthID          = int64(10)
	testLocationID           = int64(1)
	testServerModelID        = int64(100)
	testUplinkModelID        = int64(10)
	testCommit               = int64(1000)
	testBandwidthOption      = serverscom.BandwidthOption{
		ID:     testBandwidthID,
		Name:   "20002 GB",
		Type:   "bytes",
		Commit: &testCommit,
	}
)

func TestGetUplinkBandwidthCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		flags          []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get uplink bandwidth in default format",
			id:             testBandwidthID,
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID), "--uplink-model-id", fmt.Sprint(testUplinkModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "get.txt")),
		},
		{
			name:           "get uplink bandwidth in JSON format",
			id:             testBandwidthID,
			output:         "json",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID), "--uplink-model-id", fmt.Sprint(testUplinkModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "get.json")),
		},
		{
			name:           "get uplink bandwidth in YAML format",
			id:             testBandwidthID,
			output:         "yaml",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID), "--uplink-model-id", fmt.Sprint(testUplinkModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "get.yaml")),
		},
		{
			name:        "get uplink bandwidth with service error",
			id:          testBandwidthID,
			flags:       []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID), "--uplink-model-id", fmt.Sprint(testUplinkModelID)},
			expectError: true,
		},
		{
			name:        "get uplink bandwidth missing required flags",
			id:          testBandwidthID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Locations = locationsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError && len(tc.flags) > 0 {
				err = errors.New("some error")
			}

			if len(tc.flags) > 0 {
				locationsServiceHandler.EXPECT().
					GetBandwidthOption(gomock.Any(), testLocationID, testServerModelID, testUplinkModelID, tc.id).
					Return(&testBandwidthOption, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			uplinkCmd := NewCmd(testCmdContext)

			args := []string{"uplink-bandwidths", "get", fmt.Sprint(tc.id)}
			if len(tc.flags) > 0 {
				args = append(args, tc.flags...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(uplinkCmd).
				WithArgs(args)

			cmd := builder.Build()
			execErr := cmd.Execute()

			if tc.expectError {
				g.Expect(execErr).To(HaveOccurred())
			} else {
				g.Expect(execErr).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestListUplinkBandwidthsCmd(t *testing.T) {
	b1 := testBandwidthOption
	b2 := testBandwidthOption
	b2.ID = testBandwidthID + 1
	b2.Name = "Unmetered"
	b2.Type = "unmetered"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.BandwidthOption])
	}{
		{
			name:           "list uplink bandwidths all (-A)",
			output:         "json",
			args:           []string{"-A", "--location-id", "1", "--server-model-id", "100", "--uplink-model-id", "10"},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "list_all.json")),
			configureMock: func(mc *mocks.MockCollection[serverscom.BandwidthOption]) {
				mc.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.BandwidthOption{b1, b2}, nil)
			},
		},
		{
			name:           "list uplink bandwidths",
			output:         "json",
			args:           []string{"--location-id", "1", "--server-model-id", "100", "--uplink-model-id", "10"},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "list.json")),
			configureMock: func(mc *mocks.MockCollection[serverscom.BandwidthOption]) {
				mc.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.BandwidthOption{b1}, nil)
			},
		},
		{
			name:           "list uplink bandwidths with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}", "--location-id", "1", "--server-model-id", "100", "--uplink-model-id", "10"},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "list_template.txt")),
			configureMock: func(mc *mocks.MockCollection[serverscom.BandwidthOption]) {
				mc.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.BandwidthOption{b1, b2}, nil)
			},
		},
		{
			name:           "list uplink bandwidths with pageView",
			args:           []string{"--page-view", "--location-id", "1", "--server-model-id", "100", "--uplink-model-id", "10"},
			expectedOutput: testutils.ReadFixture(filepath.Join(bandwidthFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mc *mocks.MockCollection[serverscom.BandwidthOption]) {
				mc.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.BandwidthOption{b1, b2}, nil)
			},
		},
		{
			name:        "list uplink bandwidths with error",
			args:        []string{"--location-id", "1", "--server-model-id", "100", "--uplink-model-id", "10"},
			expectError: true,
			configureMock: func(mc *mocks.MockCollection[serverscom.BandwidthOption]) {
				mc.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list uplink bandwidths missing required flags",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.BandwidthOption](mockCtrl)

	locationsServiceHandler.EXPECT().
		BandwidthOptions(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Locations = locationsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			uplinkCmd := NewCmd(testCmdContext)

			args := []string{"uplink-bandwidths", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(uplinkCmd).
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
