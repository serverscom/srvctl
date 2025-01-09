package context

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/config"
)

var (
	fixtureBasePath = filepath.Join("..", "..", "testdata", "context")

	newTestConfig = func() config.Config {
		return config.Config{
			DefaultContext: "default",
			Contexts: []config.Context{
				{
					Name:     "default",
					Endpoint: "https://default.com",
					Token:    "secret",
				},
				{
					Name:     "test",
					Endpoint: "https://test.com",
					Token:    "secret",
				},
			},
		}
	}
)

func TestContextListCmd(t *testing.T) {
	testConfig := newTestConfig()
	testCases := []struct {
		name           string
		args           []string
		expectedOutput []byte
		expectError    bool
		config         config.Config
	}{
		{
			name:           "list with empty config",
			config:         config.Config{},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_empty.txt")),
		},
		{
			name:           "list",
			config:         testConfig,
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.txt")),
		},
		{
			name:           "list default",
			config:         testConfig,
			args:           []string{"--default"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_default.txt")),
		},
		{
			name:           "list no default",
			config:         testConfig,
			args:           []string{"--no-default"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_no_default.txt")),
		},
	}

	testCmdContext := testutils.NewTestCmdContext(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			testCmdContext.SetManagerConfig(&tc.config)

			contextCmd := NewCmd(testCmdContext)

			args := []string{"context", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(contextCmd).
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

func TestContextUpdateCmd(t *testing.T) {
	testConfig := newTestConfig()
	testCases := []struct {
		name                   string
		args                   []string
		contextName            string
		expectError            bool
		expectedDefaultContext string
		expectedContextName    string
		config                 config.Config
	}{
		{
			name:                   "update context",
			config:                 testConfig,
			contextName:            "test",
			expectedDefaultContext: "new-test",
			expectedContextName:    "new-test",
			args:                   []string{"test", "--name", "new-test", "--default"},
		},
	}

	testCmdContext := testutils.NewTestCmdContext(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			testCmdContext.SetManagerConfig(&tc.config)

			contextCmd := NewCmd(testCmdContext)

			args := []string{"context", "update"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(contextCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				context, err := testCmdContext.GetManager().GetContext(tc.expectedContextName)
				g.Expect(err).To(BeNil())
				defaultContext := testCmdContext.GetManager().GetDefaultContextName()

				g.Expect(defaultContext).To(BeEquivalentTo(tc.expectedDefaultContext))
				g.Expect(context.Name).To(BeEquivalentTo(tc.expectedContextName))
			}
		})
	}
}

func TestContextDeleteCmd(t *testing.T) {
	testConfig := newTestConfig()
	testCases := []struct {
		name             string
		args             []string
		expectedContexts []config.Context
		expectError      bool
		config           config.Config
	}{
		{
			name:             "delete context",
			config:           testConfig,
			args:             []string{"test"},
			expectedContexts: []config.Context{testConfig.Contexts[0]},
		},
		{
			name:        "delete default context without force",
			expectError: true,
			config:      testConfig,
			args:        []string{"default"},
		},
		{
			name:             "delete default context with force",
			config:           testConfig,
			expectedContexts: []config.Context{testConfig.Contexts[1]},
			args:             []string{"default", "--force"},
		},
	}

	testCmdContext := testutils.NewTestCmdContext(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			testCmdContext.SetManagerConfig(&tc.config)

			contextCmd := NewCmd(testCmdContext)

			args := []string{"context", "delete"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(contextCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			contexts := testCmdContext.GetManager().GetContexts()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(contexts).To(BeEquivalentTo(tc.expectedContexts))
			}
		})
	}
}
