package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

// CmdContext represents the context for a command
type CmdContext struct {
	manager   *config.Manager
	client    *client.Client
	formatter *output.Formatter
}

// NewCmdContext creates new cmd context with specified manager and client
func NewCmdContext(manager *config.Manager, client *client.Client) *CmdContext {
	return &CmdContext{
		manager: manager,
		client:  client,
	}
}

// SetManagerConfig sets manager config
func (c *CmdContext) SetManagerConfig(config *config.Config) {
	c.manager.SetConfig(config)
}

// GetManager returns the manager from cmd context
func (c *CmdContext) GetManager() *config.Manager {
	if c.manager != nil {
		return c.manager
	}
	return &config.Manager{}
}

// GetOrCreateFormatter returns formatter for the command
func (c *CmdContext) GetOrCreateFormatter(cmd *cobra.Command) *output.Formatter {
	if c.formatter != nil {
		return c.formatter
	}
	c.formatter = output.NewFormatter(cmd, c.manager)
	return c.formatter
}

// GetClient returns the client from cmd context
func (c *CmdContext) GetClient() *client.Client {
	if c.client != nil {
		return c.client
	}
	return &client.Client{}
}

// SetupContext returns context with timeout based on 'http-timeout' from config or cli flag
func SetupContext(cmd *cobra.Command, manager *config.Manager) (context.Context, context.CancelFunc) {
	httpTimeout, err := manager.GetResolvedIntValue(cmd, "http-timeout")
	if err != nil {
		log.Fatal(err)
	}
	return context.WithTimeout(context.Background(), time.Duration(httpTimeout)*time.Second)
}

// SetupProxy setup proxy envs for client based on 'proxy' defined in config or cli flag
func SetupProxy(cmd *cobra.Command, manager *config.Manager) {
	proxy, err := manager.GetResolvedStringValue(cmd, "proxy")
	if err != nil {
		log.Fatal(err)
	}
	if proxy != "" {
		if strings.HasPrefix(proxy, "https") {
			os.Setenv("HTTPS_PROXY", proxy)
		} else {
			os.Setenv("HTTP_PROXY", proxy)
		}
	}
}

// ReadInputJSON reads input from file and unmarshals it into the given struct.
// If path is "-", it reads from stdin.
func ReadInputJSON(path string, in io.Reader, input any) error {
	var inputReader io.Reader

	if path != "" && path != "-" {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		inputReader = file
	} else {
		inputReader = in
	}

	decoder := json.NewDecoder(inputReader)
	if err := decoder.Decode(input); err != nil {
		return fmt.Errorf("could not parse JSON: %w", err)
	}

	return nil
}

// ParseLabels parses slice of labels and returns map
// expects that slice element would be: "foo=bar"
func ParseLabels(labels []string) (map[string]string, error) {
	labelsMap := make(map[string]string)

	for _, label := range labels {
		parts := strings.SplitN(label, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid label format: %s", label)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		labelsMap[key] = value
	}

	return labelsMap, nil
}

// userAgent returns user agent string
func userAgent(version string) string {
	return fmt.Sprintf("srvctl/%s (%s %s)", version, runtime.GOOS, runtime.GOARCH)
}

// findEntity search entity by cmd name in entities map
func findEntity(cmd *cobra.Command, entities map[string]entities.EntityInterface) entities.EntityInterface {
	name := cmd.Name()
	if entity, ok := entities[name]; ok {
		return entity
	}
	if cmd.Parent() != nil {
		return findEntity(cmd.Parent(), entities)
	}
	return nil
}

func UsageRun(cmd *cobra.Command, args []string) { _ = cmd.Usage() }

func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		helpStr := fmt.Sprintf("Run '%v --help' for usage.", cmd.CommandPath())
		return fmt.Errorf("unknown command %q for %q\n%s", args[0], cmd.CommandPath(), helpStr)
	}
	return nil
}
