package main

import (
	"log"

	"github.com/serverscom/srvctl/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
