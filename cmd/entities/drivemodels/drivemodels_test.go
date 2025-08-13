package drivemodels

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
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "drive-models")
	testDriveModelID     = int64(10)
	testLocationID       = int64(1)
	testServerModelID    = int64(100)
	testDriveModelOption = serverscom.DriveModel{
		ID:         testDriveModelID,
		Name:       "ssd-model-749",
		Capacity:   100,
		Interface:  "SATA3",
		FormFactor: "2.5",
		MediaType:  "SSD",
	}
)

func TestGetDriveModelOptionCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		flags          []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get drive model in default format",
			id:             testDriveModelID,
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get drive model in JSON format",
			id:             testDriveModelID,
			output:         "json",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get drive model in YAML format",
			id:             testDriveModelID,
			output:         "yaml",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get drive model with service error",
			id:          testDriveModelID,
			flags:       []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectError: true,
		},
		{
			name:        "get drive model missing required flags",
			id:          testDriveModelID,
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
					GetDriveModelOption(gomock.Any(), testLocationID, testServerModelID, tc.id).
					Return(&testDriveModelOption, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			driveCmd := NewCmd(testCmdContext)

			args := []string{"drive-models", "get", fmt.Sprint(tc.id)}
			if len(tc.flags) > 0 {
				args = append(args, tc.flags...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(driveCmd).
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

func TestListDriveModelOptionsCmd(t *testing.T) {
	d1 := testDriveModelOption
	d2 := testDriveModelOption
	d2.ID = testDriveModelID + 1
	d2.Name = "hdd-model-741"
	d2.MediaType = "HDD"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.DriveModel])
	}{
		{
			name:           "list all drive models",
			output:         "json",
			args:           []string{"-A", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DriveModel]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DriveModel{d1, d2}, nil)
			},
		},
		{
			name:           "list drive models",
			output:         "json",
			args:           []string{"--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DriveModel]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DriveModel{d1}, nil)
			},
		},
		{
			name:           "list drive models with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} MediaType: {{.MediaType}}\n{{end}}", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DriveModel]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DriveModel{d1, d2}, nil)
			},
		},
		{
			name:           "list drive models with pageView",
			args:           []string{"--page-view", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DriveModel]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DriveModel{d1, d2}, nil)
			},
		},
		{
			name:        "list drive models with error",
			args:        []string{"--location-id", "1", "--server-model-id", "100"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DriveModel]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list drive models missing required flags",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.DriveModel](mockCtrl)

	locationsServiceHandler.EXPECT().
		DriveModelOptions(gomock.Any(), gomock.Any()).
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
			driveCmd := NewCmd(testCmdContext)

			args := []string{"drive-models", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(driveCmd).
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
