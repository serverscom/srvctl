package login

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	fixtureBasePath = filepath.Join("..", "..", "testdata", "login")
)

func TestLoginCmd(t *testing.T) {
	testCases := []struct {
		name           string
		input          io.Reader
		args           []string
		configureMock  func(*mocks.MockCollection[serverscom.Host])
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:        "login with empty token",
			args:        []string{"test-context"},
			input:       strings.NewReader(""),
			expectError: true,
		},
		{
			name:           "successful login",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "success.txt")),
			args:           []string{"test-context"},
			input:          strings.NewReader("token\n"),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{}, nil)
			},
		},
		{
			name:           "login with force",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "success.txt")),
			args:           []string{"--force"},
			input:          strings.NewReader("token\n"),
			configureMock: func(mock *mocks.MockCollection[serverscom.Host]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Host{}, nil)
			},
		},
		{
			name:        "login with invalid context name",
			args:        []string{"_invalid"},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostsServiceHandler := mocks.NewMockHostsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.Host](mockCtrl)

	hostsServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Hosts = hostsServiceHandler
	testCmdContext := testutils.NewTestCmdContext(scClient)
	testClientFactory := &client.TestClientFactory{
		TestClient: client.NewWithClient(scClient),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			loginCmd := NewCmd(testCmdContext, testClientFactory)

			args := []string{"login"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(loginCmd).
				WithInput(tc.input).
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
