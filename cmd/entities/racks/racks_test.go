package racks

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
	testId          = "testId"
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "racks")
	testRack        = serverscom.Rack{
		ID:           testId,
		Name:         "test-rack",
		LocationID:   1,
		LocationCode: "test",
		Labels:       map[string]string{"foo": "bar"},
	}
)

func TestGetRackCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get ssh key in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get ssh key in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get ssh key in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get ssh key with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	racksServiceHandler := mocks.NewMockRacksService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Racks = racksServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			racksServiceHandler.EXPECT().
				Get(gomock.Any(), testId).
				Return(&testRack, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rackCmd := NewCmd(testCmdContext)

			args := []string{"racks", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(rackCmd).
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

func TestListRacksCmd(t *testing.T) {
	testRack1 := testRack
	testRack2 := testRack
	testRack1.ID += "1"
	testRack2.Name = "test-rack 2"
	testRack1.ID += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Rack])
	}{
		{
			name:           "list all racks",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Rack]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Rack{
						testRack1,
						testRack2,
					}, nil)
			},
		},
		{
			name:           "list racks",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Rack]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Rack{
						testRack1,
					}, nil)
			},
		},
		{
			name:           "list racks with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Rack]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Rack{
						testRack1,
						testRack2,
					}, nil)
			},
		},
		{
			name:           "list racks with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Rack]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Rack{
						testRack1,
						testRack2,
					}, nil)
			},
		},
		{
			name:        "list racks with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Rack]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	racksServiceHandler := mocks.NewMockRacksService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Rack](mockCtrl)

	racksServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Racks = racksServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rackCmd := NewCmd(testCmdContext)

			args := []string{"racks", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(rackCmd).
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

func TestUpdateRackCmd(t *testing.T) {
	newRack := testRack
	newRack.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockRacksService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update rack",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockRacksService) {
				mock.EXPECT().
					Update(gomock.Any(), testId, serverscom.RackUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newRack, nil)
			},
		},
		{
			name: "update rack with error",
			id:   testId,
			configureMock: func(mock *mocks.MockRacksService) {
				mock.EXPECT().
					Update(gomock.Any(), testId, serverscom.RackUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	racksServiceHandler := mocks.NewMockRacksService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Racks = racksServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(racksServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rackCmd := NewCmd(testCmdContext)

			args := []string{"racks", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(rackCmd).
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
