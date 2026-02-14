package urlsanitizer

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

type URLSanitizer interface {
	Sanitize(raw string) (string, error)
}

type urlSanitizer struct{}

func NewUrlSanitizer() URLSanitizer {
	return &urlSanitizer{}
}

func (s *urlSanitizer) Sanitize(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" {
		return "", fmt.Errorf("url is empty")
	}

	parsed, err := parseWithSchemeFallback(trimmed, "https")

	if err != nil {
		return "", fmt.Errorf("url '%s' can't be parsed because: %s", trimmed, err)
	}

	if parsed.Scheme != "https" && parsed.Scheme != "http" { // reject non-http urls
		return "", fmt.Errorf("only http and https supported, but got '%s'", parsed)
	}
	if parsed.Port() != "" { // reject urls with port
		return "", fmt.Errorf("only urls without port are accepted, but got '%s'", parsed)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("url host is empty: %s", parsed)
	}

	parsed.User = nil // strip username and password

	if parsed.Path != "" {
		parsed.Path = path.Clean(parsed.Path) // clean path
	}

	parsed.RawQuery = parsed.Query().Encode()  // sort query params
	parsed.Host = strings.ToLower(parsed.Host) // normalize host case

	return parsed.String(), nil
}

func parseWithSchemeFallback(str string, protocol string) (*url.URL, error) {
	parsed, err := url.ParseRequestURI(str)

	if err != nil {
		if !strings.Contains(str, "://") {
			parsed, err = url.ParseRequestURI(protocol + "://" + str)
			if err != nil {
				return nil, fmt.Errorf("url '%s' can't be parsed because: %s", str, err)
			}
		} else {
			return nil, fmt.Errorf("url '%s' can't be parsed because: %s", str, err)
		}
	}

	return parsed, nil
}
