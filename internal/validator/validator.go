package validator

import (
	"fmt"
	"net/url"
	"regexp"
)

var (
	contextNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
)

// ValidateContextName validates context name.
func ValidateContextName(name string) error {
	if !contextNameRegex.MatchString(name) {
		return fmt.Errorf("invalid context name %q: must contain only letters, numbers, underscores, periods, and hyphens, and must start with a letter or number", name)
	}
	return nil
}

// ValidateEndpoint validates endpoint.
func ValidateEndpoint(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid endpoint URL scheme, must be http or https")
	}

	return nil
}
