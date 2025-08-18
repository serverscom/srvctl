package serverramoptions

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
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "server-ram-options")
	testRAMOption1  = serverscom.RAMOption{
		RAM:  16384,
		Type: "DDR4",
	}
	testRAMOption2 = serverscom.RAMOption{
		RAM:  32768,
		Type: "DDR4",
	}
	testLocationID    = int64(1)
	testServerModelID = int64(100)
)

func TestListRAMOptionsCmd(t *testing.T) {
	r1 := testRAMOption1
	r2 := testRAMOption2
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.RAMOption])
	}{
		{
			name:           "list all ram options",
			output:         "json",
			args:           []string{"-A", "--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.RAMOption]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.RAMOption{r1, r2}, nil)
			},
		},
		{
			name:           "list ram options",
			output:         "json",
			args:           []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.RAMOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.RAMOption{r1}, nil)
			},
		},
		{
			name:           "list ram options with template",
			args:           []string{"--template", "{{range .}}RAM: {{.RAM}} Type: {{.Type}}\\n{{end}}", "--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.RAMOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.RAMOption{r1, r2}, nil)
			},
		},
		{
			name:           "list ram options with pageView",
			args:           []string{"--page-view", "--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.RAMOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.RAMOption{r1, r2}, nil)
			},
		},
		{
			name:        "list ram options with error",
			args:        []string{"--location-id", fmt.Sprint(testLocationID), "--server-model-id", fmt.Sprint(testServerModelID)},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.RAMOption]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:        "list ram options missing required flags",
			expectError: true,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	locationsServiceHandler := mocks.NewMockLocationsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.RAMOption](mockCtrl)
	locationsServiceHandler.EXPECT().
		RAMOptions(gomock.Any(), gomock.Any()).
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
			ramCmd := NewCmd(testCmdContext)
			args := []string{"server-ram-options", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(ramCmd).
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
