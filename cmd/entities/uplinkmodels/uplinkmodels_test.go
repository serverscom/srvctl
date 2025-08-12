package uplinkmodels

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
	uplinkFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "uplink-models")
	testUplinkID          = int64(10)
	testLocationID        = int64(1)
	testServerModelID     = int64(100)
	testUplinkOption      = serverscom.UplinkOption{
		ID:         testUplinkID,
		Name:       "Public 1 Gbps with redundancy",
		Speed:      1000,
		Type:       "public",
		Redundancy: true,
	}
)

func TestGetUplinkModelCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		flags          []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get uplink model in default format",
			id:             testUplinkID,
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "get.txt")),
		},
		{
			name:           "get uplink model in JSON format",
			id:             testUplinkID,
			output:         "json",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "get.json")),
		},
		{
			name:           "get uplink model in YAML format",
			id:             testUplinkID,
			output:         "yaml",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "get.yaml")),
		},
		{
			name:        "get uplink model with service error",
			id:          testUplinkID,
			flags:       []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectError: true,
		},
		{
			name:        "get uplink missing required flags",
			id:          testUplinkID,
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
					GetUplinkOption(gomock.Any(), testLocationID, testServerModelID, tc.id).
					Return(&testUplinkOption, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			uplinkCmd := NewCmd(testCmdContext)

			args := []string{"uplink-models", "get", fmt.Sprint(tc.id)}
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

func TestListUplinkModelsCmd(t *testing.T) {
	u1 := testUplinkOption
	u2 := testUplinkOption
	u2.ID = testUplinkID + 1
	u2.Name = "Private 1 Gbps with redundancy"
	u2.Type = "private"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.UplinkOption])
	}{
		{
			name:           "list all uplink models",
			output:         "json",
			args:           []string{"-A", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.UplinkOption]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.UplinkOption{u1, u2}, nil)
			},
		},
		{
			name:           "list uplink models",
			output:         "json",
			args:           []string{"--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.UplinkOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.UplinkOption{u1}, nil)
			},
		},
		{
			name:           "list uplink models with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Type: {{.Type}}\n{{end}}", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.UplinkOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.UplinkOption{u1, u2}, nil)
			},
		},
		{
			name:           "list uplink models with pageView",
			args:           []string{"--page-view", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(uplinkFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.UplinkOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.UplinkOption{u1, u2}, nil)
			},
		},
		{
			name:        "list uplink models with error",
			args:        []string{"--location-id", "1", "--server-model-id", "100"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.UplinkOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list uplink models missing required flags",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.UplinkOption](mockCtrl)

	locationsServiceHandler.EXPECT().
		UplinkOptions(gomock.Any(), gomock.Any()).
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

			args := []string{"uplink-models", "list"}
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
