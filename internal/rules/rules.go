package rules

import "strings"

type RequestInfo struct {
	Host string
	Path string
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
