package testutils

import (
	"bytes"
	"io"
	"os"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
)

type TestCommandBuilder struct {
	rootCmd *cobra.Command
	output  *bytes.Buffer
}

func NewTestCommandBuilder() *TestCommandBuilder {
	var output bytes.Buffer

	return &TestCommandBuilder{
		rootCmd: NewTestCmd(),
		output:  &output,
	}
}

func (b *TestCommandBuilder) WithCommand(cmd *cobra.Command) *TestCommandBuilder {
	b.rootCmd.AddCommand(cmd)
	return b
}

func (b *TestCommandBuilder) WithInput(in io.Reader) *TestCommandBuilder {
	b.rootCmd.SetIn(in)
	return b
}

func (b *TestCommandBuilder) WithArgs(args []string) *TestCommandBuilder {
	b.rootCmd.SetArgs(args)
	return b
}

func (b *TestCommandBuilder) GetOutput() string {
	return b.output.String()
}

func (b *TestCommandBuilder) Build() *cobra.Command {
	b.rootCmd.SetOut(b.output)
	return b.rootCmd
}

func ReadFixture(filePath string) []byte {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}

func NewTestCmd() *cobra.Command {
	c := &cobra.Command{Use: "srvctl"}
	base.AddGlobalFlags(c)
	return c
}

func NewTestCmdContext(scClient *serverscom.Client) *base.CmdContext {
	manager := config.NewTestManager(&config.Config{
		DefaultContext: "test",
		Contexts: []config.Context{
			{
				Name: "test",
			},
		},
	})
	cli := client.NewWithClient(scClient)
	return base.NewCmdContext(manager, cli)
}
