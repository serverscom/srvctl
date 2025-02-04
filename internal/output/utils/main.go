package utils

import (
	"fmt"
	"strings"

	"github.com/jmespath/go-jmespath"
	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Humanize(input string) string {
	snake := strcase.SnakeCase(input)
	words := strings.ReplaceAll(snake, "_", " ")
	return cases.Title(language.English).String(words)
}

// GetFieldValue returns the value of a given field for an item.
func GetFieldValue(item interface{}, jsonPath string) (interface{}, error) {
	if jsonPath == "" || jsonPath == "." {
		return item, nil
	}

	result, err := jmespath.Search(jsonPath, item)
	if err != nil {
		return nil, fmt.Errorf("error getting value by path %s: %w", jsonPath, err)
	}
	return result, nil
}
