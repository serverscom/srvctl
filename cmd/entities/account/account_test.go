package account

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
	fixtureBasePath    = filepath.Join("..", "..", "..", "testdata", "entities", "account")
	testAccountBalance = serverscom.AccountBalance{
		CurrentBalance:      100.0,
		NextInvoiceTotalDue: 0.0,
		Currency:            "EUR",
	}
)

func TestGetAccountBalanceCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get account balance in default format",
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_balance.txt")),
		},
		{
			name:           "get account balance in JSON format",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_balance.json")),
		},
		{
			name:           "get account balance in YAML format",
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_balance.yaml")),
		},
		{
			name:        "get account balance with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	accountServiceHandler := mocks.NewMockAccountService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Account = accountServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			accountServiceHandler.EXPECT().
				GetBalance(gomock.Any()).
				Return(&testAccountBalance, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			accountCmd := NewCmd(testCmdContext)

			args := []string{"account", "balance"}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(accountCmd).
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
