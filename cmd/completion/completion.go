package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmd creates the completion command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completions [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for srvctl.

The shell code must be evaluated to provide interactive completion of srvctl
commands. This can be done by sourcing it from your shell configuration file.

Note for zsh users: zsh completions require zsh >= 5.2.

Bash:

  # Load completions into current shell
  source <(srvctl completions bash)

  # To load completions for each session, execute once:

  # Linux (requires bash-completion package):
  srvctl completions bash > /etc/bash_completion.d/srvctl

  # Linux (user-level):
  srvctl completions bash > ~/.bash_completion

  # macOS (requires bash-completion package from homebrew):
  # For Bash 3.2 (default macOS): brew install bash-completion
  # For Bash 4.1+: brew install bash-completion@2
  srvctl completions bash > $(brew --prefix)/etc/bash_completion.d/srvctl

Zsh:

  # Load completions into current shell
  source <(srvctl completions zsh)

  # To load completions for each session, add the following to ~/.zshrc:
  source <(srvctl completions zsh)

  # If shell completion is not already enabled in your environment,
  # you may need to add this to ~/.zshrc:
  autoload -U compinit; compinit

Fish:

  # Load completions into current shell
  srvctl completions fish | source

  # To load completions for each session, execute once:
  srvctl completions fish > ~/.config/fish/completions/srvctl.fish

PowerShell:

  # Load completions into current shell
  srvctl completions powershell | Out-String | Invoke-Expression

  # To load completions for each session, add to your PowerShell profile:

  # Option 1: Save to script and source from profile
  srvctl completions powershell > "$HOME\.srvctl\completion.ps1"
  Add-Content $PROFILE ". '$HOME\.srvctl\completion.ps1'"

  # Option 2: Lazy load with command check
  Add-Content $PROFILE "if (Get-Command srvctl -ErrorAction SilentlyContinue) {
    srvctl completions powershell | Out-String | Invoke-Expression
  }"
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootCmd := cmd.Root()
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletionV2(os.Stdout, true)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				return rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
	return cmd
}
