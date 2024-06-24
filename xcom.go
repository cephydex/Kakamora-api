package main

import (
	"errors"
	"fmt"
	"net/http"
	"svclookup/api"
)

func server(mux *http.ServeMux, port int) {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}

func handleRequests(mux *http.ServeMux) {
    // index page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("server: %s %s /\n", r.Method, r.RequestURI)
		fmt.Fprintf(w, `{"message": "Lookup service is live!!"}`)
	})

	mux.HandleFunc("POST /ping/{id}", api.CommentPost)
	// discovery page
	mux.HandleFunc("GET /check", api.JsonCheckPage)
	mux.HandleFunc("GET /discover", api.DiscoverPage)
	mux.HandleFunc("GET /discover2", api.DiscoverPageTt)
}
