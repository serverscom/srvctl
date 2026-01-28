package docgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cpuguy83/go-md2man/v2/md2man"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Generator generates documentation from cobra commands
type Generator struct {
	docsPath     string             // base path to docs/ folder with description/examples
	markdownPath string             // output path for generated markdown files
	manPath      string             // output path for generated man pages
	template     *template.Template // reusable template for doc generation
}

// ExtraContent holds additional markdown docs for a command
type ExtraContent struct {
	Description string
	Examples    string
}

const docTemplate = `# NAME

{{.Name}} - {{.Short}}

# SYNOPSIS

{{.UseLine}}

# DESCRIPTION

{{.Description}}

# OPTIONS

{{.Options}}

{{if .SubCommands}}
# SUB COMMANDS

{{.SubCommands}}
{{end}}

{{if .Examples}}
# EXAMPLES

{{.Examples}}
{{end}}
`

// NewGenerator creates a new documentation generator
func NewGenerator(docsPath, markdownPath, manPath string) *Generator {
	tmpl, err := template.New("doc").Parse(docTemplate)
	if err != nil {
		// This should never happen with a valid template
		panic(fmt.Sprintf("failed to parse documentation template: %v", err))
	}

	return &Generator{
		docsPath:     docsPath,
		markdownPath: markdownPath,
		manPath:      manPath,
		template:     tmpl,
	}
}

// Generate generates documentation for all commands
func (g *Generator) Generate(rootCmd *cobra.Command) error {
	if err := os.MkdirAll(g.markdownPath, 0755); err != nil {
		return fmt.Errorf("failed to create markdown directory: %w", err)
	}
	if err := os.MkdirAll(g.manPath, 0755); err != nil {
		return fmt.Errorf("failed to create man directory: %w", err)
	}

	if err := g.processCommand(rootCmd, ""); err != nil {
		return err
	}

	return nil
}

// processCommand processes a single command and its subcommands recursively
func (g *Generator) processCommand(cmd *cobra.Command, parentPath string) error {
	// Skip if command is hidden or deprecated
	if cmd.Hidden || cmd.Deprecated != "" {
		return nil
	}

	// Build command path (e.g., "hosts-ds-add" for "hosts ds add")
	commandPath := g.getCommandPath(cmd, parentPath)

	// Read extra content from docs/srvctl-<commandPath>/
	extra, err := g.readExtraContent(commandPath)
	if err != nil {
		return fmt.Errorf("failed to read extra content for %s: %w", commandPath, err)
	}

	// Generate markdown
	markdown, err := g.generateMarkdown(cmd, extra, commandPath)
	if err != nil {
		return fmt.Errorf("failed to generate markdown for %s: %w", commandPath, err)
	}

	// Write markdown file
	mdFilename := fmt.Sprintf("%s.md", commandPath)
	mdPath := filepath.Join(g.markdownPath, mdFilename)
	if err := os.WriteFile(mdPath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write markdown file %s: %w", mdPath, err)
	}

	// Convert to man page
	manFilename := fmt.Sprintf("%s.1", commandPath)
	manPath := filepath.Join(g.manPath, manFilename)
	if err := g.convertToMan(markdown, manPath); err != nil {
		return fmt.Errorf("failed to convert to man page %s: %w", manPath, err)
	}

	fmt.Printf("Generated: %s -> %s, %s\n", commandPath, mdPath, manPath)

	// Process subcommands recursively
	for _, subCmd := range cmd.Commands() {
		if err := g.processCommand(subCmd, commandPath); err != nil {
			return err
		}
	}

	return nil
}

// getCommandPath builds the command path (e.g., "srvctl-hosts-ds-add")
func (g *Generator) getCommandPath(cmd *cobra.Command, parentPath string) string {
	cmdName := cmd.Name()
	if parentPath == "" {
		return cmdName
	}
	return parentPath + "-" + cmdName
}

// readExtraContent reads description.md and examples.md if they exist
func (g *Generator) readExtraContent(commandPath string) (*ExtraContent, error) {
	extra := &ExtraContent{}

	// Path to docs/srvctl-<commandPath>/
	docDir := filepath.Join(g.docsPath, commandPath)

	// Read description.md if exists
	descPath := filepath.Join(docDir, "description.md")
	if content, err := os.ReadFile(descPath); err == nil {
		extra.Description = string(content)
	}

	// Read examples.md if exists
	examplesPath := filepath.Join(docDir, "examples.md")
	if content, err := os.ReadFile(examplesPath); err == nil {
		extra.Examples = string(content)
	}

	return extra, nil
}

// generateMarkdown generates markdown documentation for a command
func (g *Generator) generateMarkdown(cmd *cobra.Command, extra *ExtraContent, commandPath string) (string, error) {
	// Build description section
	description := strings.TrimSpace(cmd.Long)
	if description == "" {
		description = cmd.Short
	}
	if extra.Description != "" {
		if description != "" {
			description += "\n\n"
		}
		description += strings.TrimSpace(extra.Description)
	}

	// Build options section
	options := g.buildOptionsSection(cmd)

	// Build subcommands section
	subCommands := g.buildSubCommandsSection(cmd)

	// Build examples section
	examples := strings.TrimSpace(extra.Examples)

	// Prepare data for template
	data := map[string]string{
		"Name":        commandPath,
		"Short":       cmd.Short,
		"UseLine":     cmd.UseLine(),
		"Description": description,
		"Options":     options,
		"SubCommands": subCommands,
		"Examples":    examples,
	}

	var buf bytes.Buffer
	if err := g.template.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// buildOptionsSection builds the OPTIONS section
func (g *Generator) buildOptionsSection(cmd *cobra.Command) string {
	var buf bytes.Buffer

	// Local flags (defined on this command)
	if cmd.LocalFlags().HasAvailableFlags() {
		g.formatFlags(&buf, cmd.LocalFlags())
	}

	// Inherited flags (from parent commands)
	if cmd.InheritedFlags().HasAvailableFlags() {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		g.formatFlags(&buf, cmd.InheritedFlags())
	}

	return strings.TrimSpace(buf.String())
}

// formatFlags formats flag set into buffer
func (g *Generator) formatFlags(buf *bytes.Buffer, flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}
		var flagTypeLong, flagTypeShort string
		if flag.Value.Type() != "bool" {
			typeName := flag.Value.Type()
			// Handle empty type or special cases
			if typeName == "" || typeName == "[]" {
				typeName = "value"
			}
			flagTypeLong = fmt.Sprintf("`=[<%s>]`", typeName)
			flagTypeShort = fmt.Sprintf("`[<%s>]`", typeName)
		}

		// Write flag name and type
		fmt.Fprintf(buf, "**--%s**%s", flag.Name, flagTypeLong)
		if flag.Shorthand != "" {
			fmt.Fprintf(buf, ", **-%s**%s", flag.Shorthand, flagTypeShort)
		}
		buf.WriteString("\n")

		if flag.Usage != "" {
			fmt.Fprintf(buf, "        %s", flag.Usage)
		}
		if flag.DefValue != "" && flag.DefValue != "false" && flag.DefValue != "[]" && flag.DefValue != "0" {
			fmt.Fprintf(buf, ". The default is `%s`.", flag.DefValue)
		}
		buf.WriteString("\n\n")
	})
}

// buildSubCommandsSection builds the SUB COMMANDS section
func (g *Generator) buildSubCommandsSection(cmd *cobra.Command) string {
	var buf bytes.Buffer

	for _, subCmd := range cmd.Commands() {
		if subCmd.Hidden || subCmd.Deprecated != "" {
			continue
		}
		fmt.Fprintf(&buf, "**%s**\n", subCmd.Name())
		if subCmd.Short != "" {
			fmt.Fprintf(&buf, "        %s\n\n", subCmd.Short)
		}
	}

	return strings.TrimSpace(buf.String())
}

// convertToMan converts markdown to man page format
func (g *Generator) convertToMan(markdown string, outputPath string) error {
	manContent := md2man.Render([]byte(markdown))

	if err := os.WriteFile(outputPath, manContent, 0644); err != nil {
		return fmt.Errorf("failed to write man page: %w", err)
	}

	return nil
}
