package server

import (
	"log"
	"net/http"
)

func Start() error {

	log.Println("Starting Server...")

	// HandleFunc to take on logic for incoming requests and outgoing responses to the client
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Windwall is running!"))
		log.Printf("Received request: Method=%s Path=%s Host=%s", r.Method, r.URL.Path, r.Host)
	})

	// Start server on port 8080
	return http.ListenAndServe("127.0.0.1:8080", nil)
}
