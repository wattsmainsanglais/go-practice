package main

import "net/http"

// newMux wires up the routes. Given — don't modify it. Note the Go 1.22+
// pattern syntax: "METHOD /path/{wildcard}". A request that matches no
// pattern gets net/http's default 404.
func newMux(h *taskHandlers) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", h.handleCreate)
	mux.HandleFunc("GET /tasks", h.handleList)
	mux.HandleFunc("GET /tasks/{id}", h.handleGet)
	mux.HandleFunc("PATCH /tasks/{id}/done", h.handleMarkDone)
	mux.HandleFunc("DELETE /tasks/{id}", h.handleDelete)
	return mux
}
