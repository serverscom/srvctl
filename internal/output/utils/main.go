package utils

import (
	"strings"

	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Humanize(input string) string {
	snake := strcase.SnakeCase(input)
	words := strings.ReplaceAll(snake, "_", " ")
	return cases.Title(language.English).String(words)
}
