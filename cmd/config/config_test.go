package config

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/config"
)

var (
	fixtureBasePath = filepath.Join("..", "..", "testdata", "config")

	newTestConfig = func() config.Config {
		return config.Config{
			DefaultContext: "test",
			Contexts: []config.Context{
				{
					Name:     "test",
					Endpoint: "https://test.com",
					Token:    "secret",
					Config: config.ConfigOptions{
						"http-timeout": 99,
						"proxy":        "https://test-proxy.com",
					},
				},
			},
		}
	}
)

func TestConfigFinalCmd(t *testing.T) {
	testConfig := newTestConfig()
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		config         config.Config
	}{
		{
			name:        "empty config",
			expectError: true,
			config:      config.Config{},
		},
		{
			name:           "config final",
			config:         testConfig,
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "final.txt")),
		},
		{
			name:           "config final with flag override",
			config:         testConfig,
			args:           []string{"--http-timeout", "100"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "final_override.txt")),
		},
	}

	testCmdContext := testutils.NewTestCmdContext(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			testCmdContext.SetManagerConfig(&tc.config)

			configCmd := NewCmd(testCmdContext)

			args := []string{"config", "final"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(configCmd).
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

func TestConfigGlobalUpdateCmd(t *testing.T) {
	testConfig := newTestConfig()
	testCases := []struct {
		name                string
		args                []string
		expectError         bool
		expectedConfigOpts  config.ConfigOptions
		expectedContextName string
		config              config.Config
	}{
		{
			name:                "config global update",
			args:                []string{"--http-timeout", "100", "--no-proxy"},
			config:              testConfig,
			expectedContextName: "test",
			expectedConfigOpts: config.ConfigOptions{
				"http-timeout": 100,
			},
		},
	}

	testCmdContext := testutils.NewTestCmdContext(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			testCmdContext.SetManagerConfig(&tc.config)

			configCmd := NewCmd(testCmdContext)

			args := []string{"config", "global", "update"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(configCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				globalConfig := testCmdContext.GetManager().GetGlobalConfig()
				g.Expect(globalConfig).To(BeEquivalentTo(tc.expectedConfigOpts))
			}
		})
	}
}

func TestConfigContextUpdateCmd(t *testing.T) {
	testConfig := newTestConfig()
	testCases := []struct {
		name                string
		args                []string
		expectError         bool
		expectedConfigOpts  config.ConfigOptions
		expectedContextName string
		config              config.Config
	}{
		{
			name:                "config context update",
			args:                []string{"--http-timeout", "100", "--no-proxy"},
			config:              testConfig,
			expectedContextName: "test",
			expectedConfigOpts: config.ConfigOptions{
				"http-timeout": 100,
			},
		},
	}

	testCmdContext := testutils.NewTestCmdContext(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			testCmdContext.SetManagerConfig(&tc.config)

			configCmd := NewCmd(testCmdContext)

			args := []string{"config", "context", "update"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(configCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				context, err := testCmdContext.GetManager().GetContext(tc.expectedContextName)
				g.Expect(err).To(BeNil())

				g.Expect(context.Config).To(BeEquivalentTo(tc.expectedConfigOpts))
			}
		})
	}
}
