package server

import (
	"log"
	"net/http"

	"github.com/gengjosh/windwall/internal/rules"
)

func Start() error {
	log.Println("Starting Server...")
	return http.ListenAndServe("127.0.0.1:8080", http.HandlerFunc(proxyHandler))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("TOP HANDLER HIT: Method=%q Host=%q URL=%q Path=%q",
		r.Method, r.Host, r.URL.String(), r.URL.Path)

	// Create RequestInfo struct to pass to rules engine
	info := rules.BuildRequestInfo(r.Host, r.URL.Path, r.Referer(), r.Header.Get("Origin"), r.Method, false, false)

	if r.Method == http.MethodConnect {
		// Handle HTTPS CONNECT requests
		log.Printf("CONNECT BRANCH HIT: Host=%s", r.Host)
		info.IsHTTPS = true

		result := rules.ApplyHTTPRules(info)

		if !result.Allow {
			log.Printf("HTTPS request blocked: %s", result.Reason)
			http.Error(w, "Forbidden: "+result.Reason, http.StatusForbidden)
			return
		}

		log.Println("Request allowed for HTTPS")
		w.Write([]byte("HTTPS Request Allowed from:" + r.Host))
		return
	}

	// Handle regular HTTP requests
	log.Printf("HTTP BRANCH HIT: Path=%s", r.URL.Path)

	result := rules.ApplyHTTPRules(info)

	if !result.Allow {
		log.Printf("HTTP request blocked: %s", result.Reason)
		http.Error(w, "Forbidden: "+result.Reason, http.StatusForbidden)
		return
	}

	log.Println("Request allowed for HTTP")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("HTTP request received"))
}
