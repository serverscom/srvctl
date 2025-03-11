package main

import (
	"fmt"

	"github.com/serverscom/srvctl/cmd"
)

var (
	version string = "dev"
	commit  string = ""
)

func main() {
	if commit != "" && version == "dev" {
		version = fmt.Sprintf("%s-%s", version, commit)
	}
	rootCmd := cmd.NewRootCmd(version)

	_ = rootCmd.Execute()
}
