package base

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/spf13/cobra"
)

// CollectionFactory is a function type that creates a typed resource collection
// with configurable verbosity level
// type CollectionFactory[T any] func(verbose bool) serverscom.Collection[T]
type CollectionFactory[T any] func(verbose bool, args ...string) serverscom.Collection[T]

// NewListCmd base list command for different collections
func NewListCmd[T any](use string, entityName string, colFactory CollectionFactory[T], cmdContext *CmdContext, opts ...ListOptions[T]) *cobra.Command {
	aliases := []string{}
	if use == "list" {
		aliases = append(aliases, "ls")
	}
	cmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   "List " + entityName,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := SetupContext(cmd, manager)
			defer cancel()

			SetupProxy(cmd, manager)

			collection := colFactory(manager.GetVerbose(cmd), args...)
			for _, opt := range opts {
				opt.ApplyToCollection(collection)
			}

			items, err := fetchItems(ctx, collection, opts)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(items)
		},
	}

	for _, opt := range opts {
		opt.AddFlags(cmd)
	}

	return cmd
}
