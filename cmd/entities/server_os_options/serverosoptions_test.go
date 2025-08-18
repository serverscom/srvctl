package serverosoptions

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
	fixtureBasePath   = filepath.Join("..", "..", "..", "testdata", "entities", "os-options")
	testOSOptionID    = int64(10)
	testLocationID    = int64(1)
	testServerModelID = int64(100)
	testOSOption      = serverscom.OperatingSystemOption{
		ID:          testOSOptionID,
		FullName:    "Ubuntu 18.04-server x86_64",
		Name:        "Ubuntu",
		Version:     "18.04-server",
		Arch:        "x86_64",
		Filesystems: []string{"ext2", "ext4", "swap", "xfs", "reiser"},
	}
)

func TestGetServerOperatingSystemOptionCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		flags          []string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get os option in default format",
			id:             testOSOptionID,
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get os option in JSON format",
			id:             testOSOptionID,
			output:         "json",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get os option in YAML format",
			id:             testOSOptionID,
			output:         "yaml",
			flags:          []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get os option with service error",
			id:          testOSOptionID,
			flags:       []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectError: true,
		},
		{
			name:        "get os option missing required flags",
			id:          testOSOptionID,
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
					GetOperatingSystemOption(gomock.Any(), testLocationID, testServerModelID, tc.id).
					Return(&testOSOption, err)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			osCmd := NewCmd(testCmdContext)
			args := []string{"server-os-options", "get", fmt.Sprint(tc.id)}
			if len(tc.flags) > 0 {
				args = append(args, tc.flags...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(osCmd).
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

func TestListServerOperatingSystemOptionsCmd(t *testing.T) {
	o1 := testOSOption
	o2 := testOSOption
	o2.ID = testOSOptionID + 1
	o2.FullName = "CentOS 7 x86_64"
	o2.Name = "CentOS"
	o2.Version = "7"
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.OperatingSystemOption])
	}{
		{
			name:           "list all os options",
			output:         "json",
			args:           []string{"-A", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.OperatingSystemOption]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.OperatingSystemOption{o1, o2}, nil)
			},
		},
		{
			name:           "list os options",
			output:         "json",
			args:           []string{"--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.OperatingSystemOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.OperatingSystemOption{o1}, nil)
			},
		},
		{
			name:           "list os options with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Arch: {{.Arch}}\\n{{end}}", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.OperatingSystemOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.OperatingSystemOption{o1, o2}, nil)
			},
		},
		{
			name:           "list os options with pageView",
			args:           []string{"--page-view", "--location-id", "1", "--server-model-id", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.OperatingSystemOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.OperatingSystemOption{o1, o2}, nil)
			},
		},
		{
			name:        "list os options with error",
			args:        []string{"--location-id", "1", "--server-model-id", "100"},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.OperatingSystemOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list os options missing required flags",
			expectError: true,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.OperatingSystemOption](mockCtrl)
	locationsServiceHandler.EXPECT().
		OperatingSystemOptions(gomock.Any(), gomock.Any()).
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
			osCmd := NewCmd(testCmdContext)
			args := []string{"server-os-options", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(osCmd).
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
