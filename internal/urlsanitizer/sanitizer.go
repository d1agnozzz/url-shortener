package urlsanitizer

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
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

	if parsed.Host == "" {
		return "", fmt.Errorf("host is empty: %s", parsed)
	}

	parsed.User = nil // strip username and password

	if parsed.Path != "" {
		parsed.Path = path.Clean(parsed.Path) // clean path
	}

	parsed.RawQuery = parsed.Query().Encode()  // sort query params
	parsed.Host = strings.ToLower(parsed.Host) // normalize host case

	// check that port is in 0-65535 range, if given
	if err := validatePort(parsed); err != nil {
		return "", err
	}

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

func validatePort(parsed *url.URL) error {
	port := parsed.Port()

	if port == "" {
		return nil
	}

	// check that port is in 0-65535 range
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return fmt.Errorf("port %s is invalid: must be in range 0-65535", port)
	}

	return nil
}
