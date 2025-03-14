package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Manager represents a configuration manager
type Manager struct {
	config     *Config
	configPath string
}

// NewManager creates a new Manager
func NewManager(configPath string) (*Manager, error) {
	if configPath == "" {
		var err error
		configPath, err = getConfigPath()
		if err != nil {
			return nil, err
		}
	}

	m := &Manager{
		configPath: configPath,
	}

	if err := m.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return m, nil
}

// NewManagerWithConfig returns a new Manager with specified config
func NewManagerWithConfig(config *Config) *Manager {
	if config == nil {
		return &Manager{
			configPath: "/dev/null",
			config:     &Config{},
		}
	}
	return &Manager{
		configPath: "/dev/null",
		config:     config,
	}
}

// GetContexts returns all contexts from the config
func (m *Manager) GetContexts() []Context {
	if m.config == nil {
		return []Context{}
	}
	return m.config.Contexts
}

// GetDefaultContextName returns the name of the default context
func (m *Manager) GetDefaultContextName() string {
	if m.config.DefaultContext != "" {
		return m.config.DefaultContext
	}
	if len(m.config.Contexts) > 0 {
		return m.config.Contexts[0].Name
	}
	return ""
}

// GetContext returns a context by name
func (m *Manager) GetContext(name string) (*Context, error) {
	for i := range m.config.Contexts {
		if m.config.Contexts[i].Name == name {
			return &m.config.Contexts[i], nil
		}
	}
	return nil, fmt.Errorf("context %q not found", name)
}

// IsDefaultContext checks if the given context is the default context
func (m *Manager) IsDefaultContext(name string) (bool, error) {
	defaultName := m.GetDefaultContextName()
	return defaultName == name, nil
}

// SetDefaultContext sets the default context
func (m *Manager) SetDefaultContext(name string) error {
	if _, err := m.GetContext(name); err != nil {
		return err
	}
	m.config.DefaultContext = name
	return nil
}

// DeleteContext deletes a context by name.
func (m *Manager) DeleteContext(name string) error {
	for i := range m.config.Contexts {
		if m.config.Contexts[i].Name == name {
			m.config.Contexts = append(m.config.Contexts[:i], m.config.Contexts[i+1:]...)

			if defaultName := m.GetDefaultContextName(); defaultName == name {
				m.config.DefaultContext = ""
			}
			return nil
		}
	}
	return fmt.Errorf("context %q not found", name)
}

// Load loads config.
// If config file or config dir doesn't exist it will be created.
func (m *Manager) Load() error {
	viper.SetConfigFile(m.configPath)

	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		m.config = &Config{
			Contexts: make([]Context, 0),
		}
		return m.Save()
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	m.config = &cfg
	return nil
}

// Save saves current config to file
func (m *Manager) Save() error {
	if m.config == nil {
		return fmt.Errorf("no config loaded")
	}

	// marshal config via yaml since viper doesn't support case sensitive keys
	data, err := yaml.Marshal(m.config)
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// SetContext adds new context to config or updates existing one
func (m *Manager) SetContext(ctx Context) error {
	for i := range m.config.Contexts {
		if m.config.Contexts[i].Name == ctx.Name {
			m.config.Contexts[i] = ctx
			return nil
		}
	}

	m.config.Contexts = append(m.config.Contexts, ctx)

	if len(m.config.Contexts) == 1 {
		m.config.DefaultContext = ctx.Name
	}

	return nil
}

// GetGlobalConfig returns global config options
func (m *Manager) GetGlobalConfig() ConfigOptions {
	if m.config == nil || m.config.GlobalConfig == nil {
		return make(ConfigOptions)
	}
	return m.config.GlobalConfig
}

// UpdateGlobalConfig updates global config
func (m *Manager) UpdateGlobalConfig(configOptions ConfigOptions) {
	if m.config.GlobalConfig == nil {
		m.config.GlobalConfig = make(ConfigOptions)
	}
	for k, v := range configOptions {
		if v == nil {
			delete(m.config.GlobalConfig, k)
		} else {
			m.config.GlobalConfig[k] = v
		}
	}
}

// UpdateContextConfig updates context config
func (m *Manager) UpdateContextConfig(contextName string, configOptions ConfigOptions) error {
	ctx, err := m.GetContext(contextName)
	if err != nil {
		return err
	}

	if ctx.Config == nil {
		ctx.Config = make(ConfigOptions)
	}

	for k, v := range configOptions {
		if v == nil {
			delete(ctx.Config, k)
		} else {
			ctx.Config[k] = v
		}
	}

	return nil
}

// GetToken returns token for specified context or for default context if empty context passed
func (m *Manager) GetToken(context string) string {
	ctx := m.config.DefaultContext
	if context != "" {
		ctx = context
	}

	if ctx == "" && len(m.config.Contexts) > 0 {
		return m.config.Contexts[0].Token
	}

	for _, c := range m.config.Contexts {
		if c.Name == ctx {
			return c.Token
		}
	}

	return ""
}

// GetEndpoint returns endpoint for specified context or for default context if empty context passed
func (m *Manager) GetEndpoint(context string) string {
	ctx := m.config.DefaultContext
	if context != "" {
		ctx = context
	}

	if ctx == "" && len(m.config.Contexts) > 0 {
		return m.config.Contexts[0].Endpoint
	}

	for _, c := range m.config.Contexts {
		if c.Name == ctx {
			return c.Endpoint
		}
	}

	return ""
}

// GetConfigValue returns config value for default context or global value
func (m *Manager) GetConfigValue(key string) any {
	ctx := m.config.DefaultContext

	if ctx == "" && len(m.config.Contexts) > 0 {
		if v, ok := m.config.Contexts[0].Config[key]; ok {
			return v
		}
		return nil
	}

	for _, c := range m.config.Contexts {
		if c.Name == ctx {
			if v, ok := c.Config[key]; ok {
				return v
			}
		}
	}
	if v, ok := m.config.GlobalConfig[key]; ok {
		return v
	}

	return nil
}

// GetResolvedStringValue returns resolved string value for a given config key.
// By resolved means value is taken from command line flag if provided or else from config file or default value.
func (m *Manager) GetResolvedStringValue(cmd *cobra.Command, flagName string) (string, error) {
	if cmd.Flags().Changed(flagName) {
		return cmd.Flags().GetString(flagName)
	}
	if configValue := m.GetConfigValue(flagName); configValue != nil {
		if v, ok := configValue.(string); ok {
			return v, nil
		}
		cmd.Printf("can't parse config value %q for %q as string, use default\n", configValue, flagName)
	}
	return cmd.Flags().GetString(flagName)
}

// GetResolvedIntValue returns resolved int value for a given config key.
// By resolved means value is taken from command line flag if provided or else from config file or default value.
func (m *Manager) GetResolvedIntValue(cmd *cobra.Command, flagName string) (int, error) {
	if cmd.Flags().Changed(flagName) {
		return cmd.Flags().GetInt(flagName)
	}
	if configValue := m.GetConfigValue(flagName); configValue != nil {
		if v, ok := configValue.(int); ok {
			return v, nil
		}
		cmd.Printf("can't parse config value %q for %q as int, use default\n", configValue, flagName)
	}
	return cmd.Flags().GetInt(flagName)
}

// GetResolvedBoolValue returns resolved bool value for a given config key.
// By resolved means value is taken from command line flag if provided or else from config file or default value.
func (m *Manager) GetResolvedBoolValue(cmd *cobra.Command, flagName string) (bool, error) {
	if cmd.Flags().Changed(flagName) {
		return cmd.Flags().GetBool(flagName)
	}
	if configValue := m.GetConfigValue(flagName); configValue != nil {
		if v, ok := configValue.(bool); ok {
			return v, nil
		}
		cmd.Printf("can't parse config value %q for %q as bool, use default\n", configValue, flagName)
	}
	return cmd.Flags().GetBool(flagName)
}

// GetResolvedStringSliceValue returns resolved slice of string value for a given config key.
// By resolved means value is taken from command line flag if provided or else from config file or default value.
func (m *Manager) GetResolvedStringSliceValue(cmd *cobra.Command, flagName string) ([]string, error) {
	if cmd.Flags().Changed(flagName) {
		return cmd.Flags().GetStringArray(flagName)
	}
	if configValue := m.GetConfigValue(flagName); configValue != nil {
		if v, ok := configValue.([]string); ok {
			return v, nil
		}
		cmd.Printf("can't parse config value %q for %q as slice of string, use default\n", configValue, flagName)
	}
	return cmd.Flags().GetStringArray(flagName)
}

// GetVerbose reads verbose flag from cmd or from config
func (m *Manager) GetVerbose(cmd *cobra.Command) bool {
	v, err := m.GetResolvedBoolValue(cmd, "verbose")
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// SetConfig sets config for manager
func (m *Manager) SetConfig(config *Config) {
	m.config = config
}
