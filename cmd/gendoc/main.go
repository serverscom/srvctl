package main

import (
	"fmt"
	"log"
	"os"

	"github.com/serverscom/srvctl/cmd"
	"github.com/serverscom/srvctl/internal/docgen"
)

func main() {
	rootCmd := cmd.NewRootCmd("dev")

	baseDocsPath := "docs"
	markdownOutputPath := "docs/generated/markdown"
	manOutputPath := "docs/generated/man"

	generator := docgen.NewGenerator(
		baseDocsPath,
		markdownOutputPath,
		manOutputPath,
	)

	fmt.Println("Generating documentation...")
	if err := generator.Generate(rootCmd); err != nil {
		log.Fatalf("Failed to generate documentation: %v", err)
	}

	fmt.Println("\nDocumentation generated successfully!")
	fmt.Println("  Markdown: ", markdownOutputPath)
	fmt.Println("  Man pages: ", manOutputPath)

	os.Exit(0)
}
