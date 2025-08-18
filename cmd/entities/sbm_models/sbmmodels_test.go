package sbmmodels

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
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "sbm-models")
	testSBMFlavorID = int64(10)
	testLocationID  = int64(1)
	testSBMFlavor   = serverscom.SBMFlavor{
		ID:                     testSBMFlavorID,
		Name:                   "P-101",
		CPUName:                "cpu_name",
		CPUCount:               1,
		CPUCoresCount:          2,
		CPUFrequency:           "3.8",
		RAMSize:                4096,
		DrivesConfiguration:    "",
		PublicUplinkModelID:    843,
		PublicUplinkModelName:  "uplink-model-name-30",
		PrivateUplinkModelID:   842,
		PrivateUplinkModelName: "uplink-model-name-29",
		BandwidthID:            844,
		BandwidthName:          "public-bandwidth-model-35",
	}
)

func TestGetSBMModelCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		flags          []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get sbm model in default format",
			id:             testSBMFlavorID,
			flags:          []string{"--location-id", fmt.Sprint(testLocationID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get sbm model in JSON format",
			id:             testSBMFlavorID,
			output:         "json",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get sbm model in YAML format",
			id:             testSBMFlavorID,
			output:         "yaml",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get sbm model with service error",
			id:          testSBMFlavorID,
			flags:       []string{"--location-id", fmt.Sprint(testLocationID)},
			expectError: true,
		},
		{
			name:        "get sbm model missing required flags",
			id:          testSBMFlavorID,
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
					GetSBMFlavorOption(gomock.Any(), testLocationID, tc.id).
					Return(&testSBMFlavor, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sbmCmd := NewCmd(testCmdContext)

			args := []string{"sbm-models", "get", fmt.Sprint(tc.id)}
			if len(tc.flags) > 0 {
				args = append(args, tc.flags...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sbmCmd).
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

func TestListUplinkModelsCmd(t *testing.T) {
	f1 := testSBMFlavor
	f2 := testSBMFlavor
	f2.ID = testSBMFlavorID + 1
	f2.Name = "P-102"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.SBMFlavor])
	}{
		{
			name:           "list all sbm models",
			output:         "json",
			args:           []string{"-A", "--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMFlavor]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.SBMFlavor{f1, f2}, nil)
			},
		},
		{
			name:           "list sbm models",
			output:         "json",
			args:           []string{"--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SBMFlavor{f1}, nil)
			},
		},
		{
			name:           "list sbm models with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}", "--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SBMFlavor{f1, f2}, nil)
			},
		},
		{
			name:           "list sbm models with pageView",
			args:           []string{"--page-view", "--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SBMFlavor{f1, f2}, nil)
			},
		},
		{
			name:        "list sbm models with error",
			args:        []string{"--location-id", "1"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.SBMFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list sbm models missing required flags",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.SBMFlavor](mockCtrl)

	locationsServiceHandler.EXPECT().
		SBMFlavorOptions(gomock.Any()).
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

			args := []string{"sbm-models", "list"}
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
