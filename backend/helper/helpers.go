package helper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func IsAllowedQuery(query, expectedCommand string) bool {
	// Split the query to get the first word
	queryWords := strings.Fields(query)
	if len(queryWords) == 0 {
		return false
	}

	return strings.ToUpper(queryWords[0]) == expectedCommand
}

func ParseDomainRequest(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	referer := r.Referer()
	if referer == "" {
		return ""
	}

	// Parse the referer URL
	parsedURL, err := url.Parse(referer)
	if err != nil {
		return ""
	}

	// Get the scheme and host
	return fmt.Sprintf("%s%s", scheme, parsedURL.Host)
}