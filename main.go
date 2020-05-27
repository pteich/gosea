package main

import (
	"net/http"

	"github.com/pteich/gosea/status"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", status.Health)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	srv.ListenAndServe()
}
