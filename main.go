package main

import (
	"fmt"
	"log"

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

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
