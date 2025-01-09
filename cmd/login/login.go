package login

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/validator"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewCmd(cmdContext *base.CmdContext, clientFactory client.ClientFactory) *cobra.Command {
	var (
		force      bool
		endpoint   string
		setDefault bool
	)

	cmd := &cobra.Command{
		Use:   "login <context-name>",
		Short: "Login to servers.com API",
		Long: `Login to servers.com API via token and save the credentials in a named context.
Example: srvctl login context-name`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			var contextName string
			if len(args) > 0 {
				contextName = args[0]
			} else {
				contextName = manager.GetDefaultContextName()
				if contextName == "" {
					return fmt.Errorf("no contexts found")
				}
			}

			if err := validator.ValidateContextName(contextName); err != nil {
				return err
			}

			existingCtx, _ := manager.GetContext(contextName)
			if existingCtx != nil && !force {
				return fmt.Errorf("context %q already exists. Use --force to override", contextName)
			}

			if err := validator.ValidateEndpoint(endpoint); err != nil {
				return err
			}

			token, err := readSecureInput(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read token: %w", err)
			}

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			apiClient := clientFactory.NewClient(token, endpoint)
			if err := apiClient.VerifyCredentials(ctx); err != nil {
				return fmt.Errorf("failed to verify credentials: %w", err)
			}

			cfgCtx := config.Context{
				Name:     contextName,
				Endpoint: endpoint,
				Token:    token,
				Config:   make(map[string]interface{}),
			}

			if err := manager.SetContext(cfgCtx); err != nil {
				return fmt.Errorf("failed to update context: %w", err)
			}

			if setDefault {
				if err := manager.SetDefaultContext(contextName); err != nil {
					return fmt.Errorf("failed to set default context: %w", err)
				}
			}

			if err := manager.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			cmd.Printf("Successfully logged in with context %q\n", contextName)
			if setDefault {
				cmd.Printf("Context %q set as default\n", contextName)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force override existing context")
	cmd.Flags().StringVar(&endpoint, "endpoint", "https://api.servers.com/v1", "API endpoint")
	cmd.Flags().BoolVarP(&setDefault, "default", "d", true, "Set as default context")

	return cmd
}

func readSecureInput(in io.Reader) (string, error) {
	if term.IsTerminal(int(syscall.Stdin)) {
		stdin := int(syscall.Stdin)
		oldState, err := term.GetState(stdin)
		if err != nil {
			return "", err
		}
		defer restoreTerminal(stdin, oldState)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		go func() {
			<-sigChan
			restoreTerminal(stdin, oldState)
			fmt.Println("\nInput interrupted.")
			os.Exit(1)
		}()

		fmt.Print("Enter API token: ")
		password, err := term.ReadPassword(stdin)
		fmt.Println()
		if err != nil {
			return "", err
		}
		return string(password), nil
	}

	// if not terminal read as usual
	reader := bufio.NewReader(in)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(token), nil
}

func restoreTerminal(stdin int, oldState *term.State) {
	if err := term.Restore(stdin, oldState); err != nil {
		log.Printf("failed to restore terminal: %v", err)
	}
}
