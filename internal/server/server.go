package server

import (
	"log"
	"net/http"

	"github.com/gengjosh/windwall/internal/rules"
)

func Start() error {

	log.Println("Starting Server...")

	// HandleFunc to take on logic for incoming requests and outgoing responses to the client
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Windwall is running! \n"))
		log.Printf("Received request: Method=%s Path=%s Host=%s", r.Method, r.URL.Path, r.Host)

		// Create RequestInfo struct to pass to rules engine
		info := rules.RequestInfo{
			Host: r.Host,
			Path: r.URL.Path,
		}

		// Apply rules to incoming request
		result := rules.ApplyRules(info)
		if !result.Allow {
			log.Printf("Request blocked: %s", result.Reason)
			http.Error(w, "Forbidden: "+result.Reason, http.StatusForbidden)
			return
		}

		log.Println("Request allowed")
		w.Write([]byte("Request Allowed from:" + r.Host))

	})

	// Start server on port 8080
	return http.ListenAndServe("127.0.0.1:8080", nil)
}
