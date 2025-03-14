package base

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

// CombinePreRunE combines multiple pre-run functions into one
func CombinePreRunE(funcs ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range funcs {
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

// InitCmdContext inits cmd context and sets up necessary dependencies
func InitCmdContext(cmdContext *CmdContext) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		context, err := cmd.Flags().GetString("context")
		if err != nil {
			return err
		}

		m, err := setupConfigManager(configPath, context)
		if err != nil {
			return fmt.Errorf("failed to initialize config manager: %w", err)
		}
		if context != "" {
			if _, err := m.GetContext(context); err != nil {
				return fmt.Errorf("context %q not found in config", context)
			}
		}

		c := client.NewClient(
			m.GetToken(context),
			m.GetEndpoint(context),
		)
		version := cmd.Root().Version
		c.SetUserAgent(userAgent(version))

		cmdContext.manager = m
		cmdContext.client = c
		cmdContext.formatter = output.NewFormatter(cmd, m)

		return nil
	}
}

// CheckFormatterFlags checks flags related to formatter
func CheckFormatterFlags(cmdContext *CmdContext, entities map[string]entities.EntityInterface) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if entities == nil {
			return fmt.Errorf("entities is not initialized")
		}
		manager := cmdContext.GetManager()
		formatter := cmdContext.GetOrCreateFormatter(cmd)

		fieldList, err := manager.GetResolvedBoolValue(cmd, "field-list")
		if err != nil {
			return err
		}

		entity := findEntity(cmd, entities)
		if entity == nil {
			return fmt.Errorf("can't find entity")
		}
		if fieldList {
			formatter.ListEntityFields(entity.GetFields())
			os.Exit(0)
		}

		output := formatter.GetOutput()
		switch output {
		case "json", "yaml":
			return nil
		case "text":
		default:
			return fmt.Errorf("invalid output %q, allowed values: json, text, yaml", output)
		}

		tmpl := formatter.GetTemplateStr()
		if tmpl != "" {
			tmpl = strings.Trim(tmpl, " ")
			r := strings.NewReplacer(`\t`, "\t", `\n`, "\n")
			tmpl = r.Replace(tmpl)

			t, err := template.New("").Parse(tmpl)
			if err != nil {
				return err
			}
			formatter.SetTemplate(t)
			return nil
		}

		fields, err := manager.GetResolvedStringSliceValue(cmd, "field")
		if err != nil {
			return err
		}
		if len(fields) > 0 {
			if err := entity.Validate(fields); err != nil {
				return err
			}
		}

		return nil
	}
}

// CheckEmptyContexts returns error if no contexts found
func CheckEmptyContexts(cmdContext *CmdContext) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		manager := cmdContext.GetManager()

		if len(manager.GetContexts()) == 0 {
			return fmt.Errorf("no contexts found, log in first: 'srvctl login <context-name>'")
		}
		return nil
	}
}
