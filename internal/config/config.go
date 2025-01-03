package config

import (
	"os"
	"path/filepath"
)

type ConfigOptions = map[string]interface{}

type Context struct {
	Name     string        `yaml:"name"`
	Endpoint string        `yaml:"endpoint"`
	Token    string        `yaml:"token"`
	Config   ConfigOptions `yaml:"config"`
}

type Config struct {
	GlobalConfig   ConfigOptions `yaml:"globalConfig"`
	DefaultContext string        `yaml:"defaultContext"`
	Contexts       []Context     `yaml:"contexts"`
}

// getConfigPath returns config path from env in this priority:
// $SRVCTL_CONFIG_PATH
// $XDG_CONFIG_HOME
// $HOME
// Config will be created in $VAR/.config/srvctl/config.yaml.
// If none of these env vars are defined, config will be created
// in the current dir using path .config/srvctl/config.yaml
func getConfigPath() (string, error) {
	if configPath := os.Getenv("SRVCTL_CONFIG_PATH"); configPath != "" {
		return configPath, nil
	}

	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "srvctl", "config.yaml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".config", "srvctl", "config.yaml"), nil
	}

	return filepath.Join(home, ".config", "srvctl", "config.yaml"), nil
}
