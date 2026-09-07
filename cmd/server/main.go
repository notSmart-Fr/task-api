package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, server!")
	})

	port := ":8000"
	fmt.Printf("Server starting on http://localhost%s...\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
