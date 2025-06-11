package base

import "github.com/spf13/cobra"

func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("config", "", "config file path")
	cmd.PersistentFlags().String("context", "", "context name")
	cmd.PersistentFlags().String("proxy", "", "proxy url")
	cmd.PersistentFlags().Int("http-timeout", 30, "HTTP timeout ( seconds )")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringP("output", "o", "text", "output format (text/json/yaml)")
	// define help flag without shorthand before cobra adds it by default to avoid conflict with no-header flag shorthand
	cmd.PersistentFlags().Bool("help", false, "Print usage")
	cmd.PersistentFlags().BoolP("no-header", "h", false, "print output without headers")
}

func AddFormatFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayP("field", "f", []string{}, "output only these fields, can be specified multiple times")
	cmd.PersistentFlags().Bool("field-list", false, "list available fields")
	cmd.PersistentFlags().Bool("page-view", false, "use page view format")
	cmd.PersistentFlags().StringP("template", "t", "", "go template string to output in specified format")
}
