package servermodels

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
	fixtureBasePath   = filepath.Join("..", "..", "..", "testdata", "entities", "server-models")
	testLocationID    = int64(1)
	testServerModelID = int64(100)

	testServerModelOption = serverscom.ServerModelOption{
		ID:                 testServerModelID,
		Name:               "server-model-123",
		CPUName:            "Intel Xeon Silver 4214",
		CPUCount:           2,
		CPUCoresCount:      24,
		CPUFrequency:       2200,
		RAM:                64,
		RAMType:            "DDR4 ECC",
		MaxRAM:             2048,
		HasRAIDController:  true,
		RAIDControllerName: "PERC H740P",
		DriveSlotsCount:    16,
	}

	testServerModelOptionDetail = serverscom.ServerModelOptionDetail{
		ID:                 testServerModelID,
		Name:               "server-model-123",
		CPUName:            "Intel Xeon Silver 4214",
		CPUCount:           2,
		CPUCoresCount:      24,
		CPUFrequency:       2200,
		RAM:                64,
		RAMType:            "DDR4 ECC",
		MaxRAM:             2048,
		HasRAIDController:  true,
		RAIDControllerName: "PERC H740P",
		DriveSlotsCount:    16,
		DriveSlots: []serverscom.ServerModelDriveSlot{
			{Position: 1, Interface: "SAS", FormFactor: "2.5\"", DriveModelID: 101, HotSwappable: true},
			{Position: 2, Interface: "SAS", FormFactor: "2.5\"", DriveModelID: 102, HotSwappable: true},
		},
	}
)

func TestGetServerModelOptionCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		flags          []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get server model in default format",
			id:             testServerModelID,
			flags:          []string{"--location-id", fmt.Sprint(testLocationID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get server model in JSON format",
			id:             testServerModelID,
			output:         "json",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get server model in YAML format",
			id:             testServerModelID,
			output:         "yaml",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get server model with service error",
			id:          testServerModelID,
			flags:       []string{"--location-id", fmt.Sprint(testLocationID)},
			expectError: true,
		},
		{
			name:        "get server model missing required flags",
			id:          testServerModelID,
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
					GetServerModelOption(gomock.Any(), testLocationID, tc.id).
					Return(&testServerModelOptionDetail, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			serverCmd := NewCmd(testCmdContext)

			args := []string{"server-models", "get", fmt.Sprint(tc.id)}
			if len(tc.flags) > 0 {
				args = append(args, tc.flags...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(serverCmd).
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

func TestListServerModelOptionsCmd(t *testing.T) {
	s1 := testServerModelOption
	s2 := testServerModelOption
	s2.ID = testServerModelID + 1
	s2.Name = "server-model-456"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.ServerModelOption])
	}{
		{
			name:           "list all server models",
			output:         "json",
			args:           []string{"-A", "--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.ServerModelOption]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.ServerModelOption{s1, s2}, nil)
			},
		},
		{
			name:           "list server models",
			output:         "json",
			args:           []string{"--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.ServerModelOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.ServerModelOption{s1}, nil)
			},
		},
		{
			name:           "list server models with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}", "--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.ServerModelOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.ServerModelOption{s1, s2}, nil)
			},
		},
		{
			name:           "list server models with pageView",
			args:           []string{"--page-view", "--location-id", "1"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.ServerModelOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.ServerModelOption{s1, s2}, nil)
			},
		},
		{
			name:        "list server models with error",
			args:        []string{"--location-id", "1"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.ServerModelOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list server models missing required flags",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.ServerModelOption](mockCtrl)

	locationsServiceHandler.EXPECT().
		ServerModelOptions(gomock.Any()).
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
			serverCmd := NewCmd(testCmdContext)

			args := []string{"server-models", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(serverCmd).
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
