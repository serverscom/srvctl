package base

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
)

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

// LoadInputFromFile unmarshall json from path to input
func LoadInputFromFile(path string, input interface{}) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	err = json.Unmarshal(fileData, input)
	if err != nil {
		return fmt.Errorf("could not parse JSON: %w", err)
	}

	return nil
}

// CheckEmptyContexts returns error if no contexts found
func CheckEmptyContexts(cmd *cobra.Command, args []string) {
	manager, err := config.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	if len(manager.GetContexts()) == 0 {
		log.Fatal("no contexts found")

	}
}
