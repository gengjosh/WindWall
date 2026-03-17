package rules

import (
	"net"
	"net/url"
	"strings"
)

type RequestInfo struct {
	Host         string
	Path         string
	Referer      string
	Origin       string
	Method       string
	IsHTTPS      bool
	IsThirdParty bool
}

type Result struct {
	Allow  bool
	Reason string
}

var blockedHostKeywords = []string{
	"badhost",
	"malicious",
	"blocked",
}

var blockedPathKeywords = []string{
	"/admin",
	"/secret",
	"/private",
}

// ApplyHTTPRules checks the incoming request against defined rules and returns a Result struct
func ApplyHTTPRules(info RequestInfo) Result {
	for _, keyword := range blockedHostKeywords {
		if strings.Contains(info.Host, keyword) {
			return Result{
				Allow:  false,
				Reason: "Blocked host keyword: " + keyword,
			}
		}
	}

	// Check for blocked path keywords
	for _, keyword := range blockedPathKeywords {
		if strings.Contains(info.Path, keyword) {
			return Result{
				Allow:  false,
				Reason: "Blocked path keyword: " + keyword,
			}
		}
	}

	return Result{Allow: true, Reason: "Request Allowed from: " + info.Host}
}

// ApplyHTTPSRules can have different logic than HTTP rules, but for now it checks the same host keywords
func ApplyHTTPSRules(info RequestInfo) Result {
	for _, keyword := range blockedHostKeywords {
		if strings.Contains(info.Host, keyword) {
			return Result{
				Allow:  false,
				Reason: "Blocked host keyword: " + keyword,
			}
		}
	}

	return Result{Allow: true, Reason: "HTTPS Request Allowed from: " + info.Host}
}

// Function to build RequestInfo struct from HTTP request data
func BuildRequestInfo(host string, path string, referer string, origin string, method string, ishttps bool, isthirdparty bool) RequestInfo {

	return RequestInfo{
		Host:         host,
		Path:         path,
		Referer:      referer,
		Origin:       origin,
		Method:       method,
		IsHTTPS:      ishttps,
		IsThirdParty: isthirdparty,
	}
}

// Function to determine if a request is third-party based on host, referer, and origin
func DetermineThirdParty(host, referer, origin string) bool {
	// Extract domains from host, referer, and origin
	hostDomain := extractDomain(host)

	if origin != "" {
		originDomain := extractDomain(origin)
		return hostDomain != originDomain
	}

	if referer != "" {
		refererDomain := extractDomain(referer)
		return hostDomain != refererDomain
	}

	return false
}

// Function to check if a request should be blocked based on various criteria
func ShouldBlock(info RequestInfo) (bool, string) {
	return false, ""
}

// Function to check if the host matches any blocked keywords
func hostMatchesBlockedKeyword(host string) (bool, string) {
	return false, ""
}

// Function to check if the path matches any blocked keywords
func pathMatchesBlockedKeyword(path string) (bool, string) {
	return false, ""
}

// Helper function to extract the hostname from a URL
func extractHostname(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return ""
	}

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		parsedURL, err := url.Parse(input)
		if err != nil {
			return ""
		}
		input = parsedURL.Host
	}

	host, _, err := net.SplitHostPort(input)
	if err == nil {
		return host
	}

	return input
}

// Helper function to extract the domain from a URL
func extractDomain(host string) string {
	host = extractHostname(host)
	if host == "" {
		return ""
	}

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return host
	}

	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}
