package locations

import (
	"errors"
	"path/filepath"
	"testing"

	"fmt"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	testId          = int64(1)
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "locations")
	testLocation    = serverscom.Location{
		ID:                testId,
		Name:              "test-location",
		Status:            "active",
		Code:              "test",
		SupportedFeatures: []string{"feature1", "feature2"},
	}
)

func TestGetLocationCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get location in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get location in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get location in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get location with error",
			id:          testId,
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
			if tc.expectError {
				err = errors.New("some error")
			}
			locationsServiceHandler.EXPECT().
				GetLocation(gomock.Any(), testId).
				Return(&testLocation, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			locationCmd := NewCmd(testCmdContext)

			args := []string{"locations", "get", fmt.Sprint(tc.id)}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(locationCmd).
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

func TestListLocationsCmd(t *testing.T) {
	testLocation1 := testLocation
	testLocation2 := testLocation
	testLocation1.ID = 1
	testLocation2.Name = "test-location 2"
	testLocation2.ID = 2

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Location])
	}{
		{
			name:           "list all locations",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Location]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Location{
						testLocation1,
						testLocation2,
					}, nil)
			},
		},
		{
			name:           "list locations",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Location]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Location{
						testLocation1,
					}, nil)
			},
		},
		{
			name:           "list locations with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Location]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Location{
						testLocation1,
						testLocation2,
					}, nil)
			},
		},
		{
			name:           "list locations with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Location]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Location{
						testLocation1,
						testLocation2,
					}, nil)
			},
		},
		{
			name:        "list locations with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Location]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Location](mockCtrl)

	locationsServiceHandler.EXPECT().
		Collection().
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
			locationCmd := NewCmd(testCmdContext)

			args := []string{"locations", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(locationCmd).
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
