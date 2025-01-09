package base

import "github.com/spf13/cobra"

func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("config", "", "config file path")
	cmd.PersistentFlags().String("context", "", "context name")
	cmd.PersistentFlags().String("proxy", "", "proxy url")
	cmd.PersistentFlags().Int("http-timeout", 30, "HTTP timeout ( seconds )")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringP("output", "o", "text", "output format (text/json/yaml)")
}
