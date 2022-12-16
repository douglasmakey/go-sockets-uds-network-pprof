package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
)

var (
	socketPath = "/tmp/httpecho.sock"
	// Creating a new HTTP client that is configured to make HTTP requests over a Unix domain socket.
	httpClient = http.Client{
		Transport: &http.Transport{
			// Set the DialContext field to a function that creates
			// a new network connection to a Unix domain socket
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
)

func test(w http.ResponseWriter, req *http.Request) {
	resp, err := httpClient.Get("http://unix/")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(b)
}

func main() {
	http.HandleFunc("/test", test)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
