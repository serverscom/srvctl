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
	servicesFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "hosts", "services")
	testDSService           = serverscom.DedicatedServerService{
		ID:            testId,
		Name:          "Test service",
		Type:          "server_base",
		Currency:      "USD",
		StartedAt:     fixedTime,
		FinishedAt:    fixedTime,
		Total:         100.0,
		UsageQuantity: 2.0,
		Tax:           10.0,
		Subtotal:      100.0,
		DiscountRate:  5.0,
		DateFrom:      "2025-11-01",
		DateTo:        "2025-11-30",
	}
)

func TestListDSServicesCmd(t *testing.T) {
	testService1 := testDSService
	testService2 := testDSService
	testService2.ID += "2"
	testService2.Name = "Test service 2"
	testService2.Type = "uplink"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.DedicatedServerService])
	}{
		{
			name:           "list all ds services",
			output:         "json",
			args:           []string{"-A", testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
						testService2,
					}, nil)
			},
		},
		{
			name:           "list ds services",
			output:         "json",
			args:           []string{testServerID},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
					}, nil)
			},
		},
		{
			name:           "list ds services with template",
			args:           []string{testServerID, "--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
						testService2,
					}, nil)
			},
		},
		{
			name:           "list ds services with page-view",
			args:           []string{testServerID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(servicesFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.DedicatedServerService{
						testService1,
						testService2,
					}, nil)
			},
		},
		{
			name:        "list ds services with error",
			args:        []string{testServerID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.DedicatedServerService]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	hostService := mocks.NewMockHostsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.DedicatedServerService](mockCtrl)

	hostService.EXPECT().DedicatedServerServices(gomock.Any()).Return(collection).AnyTimes()
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

			args := append([]string{"hosts", "ds", "list-services"}, tc.args...)
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
