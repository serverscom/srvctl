package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
)

type CmdContext struct {
	manager *config.Manager
	client  *client.Client
}

func (c *CmdContext) SetManagerConfig(config *config.Config) {
	c.manager.SetConfig(config)
}

func (c *CmdContext) GetManager() *config.Manager {
	if c.manager != nil {
		return c.manager
	}
	return &config.Manager{}
}

func (c *CmdContext) GetClient() *client.Client {
	if c.client != nil {
		return c.client
	}
	return &client.Client{}
}

// NewCmdContext creates new cmd context with specified manager and client
func NewCmdContext(manager *config.Manager, client *client.Client) *CmdContext {
	return &CmdContext{
		manager: manager,
		client:  client,
	}
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
func ReadInputJSON(path string, in io.Reader, input interface{}) error {
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

// CheckEmptyContexts returns error if no contexts found
func CheckEmptyContexts(cmdContext *CmdContext) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		manager := cmdContext.GetManager()

		if len(manager.GetContexts()) == 0 {
			return fmt.Errorf("no contexts found")
		}
		return nil
	}
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
