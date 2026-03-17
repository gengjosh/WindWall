package server

import (
	"io"
	"net"
	"net/http"
)

func handleHTTPSConnect(w http.ResponseWriter, r *http.Request) {
	// Connect to the target server
	targetConn, err := net.Dial("tcp", r.Host)

	if err != nil {
		http.Error(w, "Failed to connect to target: "+err.Error(), http.StatusBadGateway)
		return
	}

	defer targetConn.Close()

	// Hijack the client connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Failed to hijack connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer clientConn.Close()

	// Send 200 OK to client to indicate connection is established
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	if err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy data in both directions
	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {
	outReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create outbound request", http.StatusInternalServerError)
		return
	}

	outReq.Header = r.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(outReq)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
