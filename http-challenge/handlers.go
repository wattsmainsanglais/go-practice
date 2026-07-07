package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Exercise: HTTP server — a tiny task-tracker API
//
// Concepts:
//   - net/http.ServeMux (the version built into Go 1.22+) supports method +
//     path-pattern routes directly: mux.HandleFunc("GET /tasks/{id}", ...).
//     No third-party router needed. `{id}` is a wildcard segment you read
//     back with r.PathValue("id").
//   - A handler has the signature `func(w http.ResponseWriter, r *http.Request)`.
//     You write the response by calling methods on w — and the ORDER
//     matters: any headers (w.Header().Set(...)) and w.WriteHeader(status)
//     must happen BEFORE you call w.Write / json.NewEncoder(w).Encode(...).
//     Once you've written body bytes, the status is locked in as whatever
//     it was (200 by default if you never called WriteHeader).
//   - encoding/json: `json.NewEncoder(w).Encode(v)` writes v as JSON
//     straight to the response body. `json.NewDecoder(r.Body).Decode(&v)`
//     reads a JSON request body into v. Both can fail — a malformed
//     request body should get a 400, not a crash.
//
// This is a bigger exercise than the earlier ones. Do the handlers in the
// order they appear — each one is a bit more involved than the last.
//
// The store (store.go) is fully implemented — use it, don't reimplement
// task storage here. routes.go wires these handlers to a *http.ServeMux —
// also given, don't modify it.

type taskHandlers struct {
	store *TaskStore
}

// createTaskRequest is the expected JSON body for POST /tasks.
type createTaskRequest struct {
	Title string `json:"title"`
}

// YOUR TASK 1: handleCreate
//
//	STEP 1: decode the request body into a createTaskRequest using
//	        json.NewDecoder(r.Body).Decode(&req). If it errors, write
//	        http.StatusBadRequest and return (nothing else to do).
//
//	STEP 2: if req.Title is empty, write http.StatusBadRequest and return.
//	        (A task with no title isn't valid input.)
//
//	STEP 3: call h.store.Create(req.Title) to get the new *Task.
//
//	STEP 4: set the response header `w.Header().Set("Content-Type", "application/json")`,
//	        then `w.WriteHeader(http.StatusCreated)` — in that order, and
//	        both BEFORE you encode the body.
//
//	STEP 5: encode the task as the JSON response body:
//	        `json.NewEncoder(w).Encode(task)`.
func (h *taskHandlers) handleCreate(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// YOUR TASK 2: handleList
//
//	STEP 1: call h.store.List() to get all tasks.
//
//	STEP 2: set the Content-Type header to "application/json" (no need to
//	        call WriteHeader — 200 is the default if you never set one).
//
//	STEP 3: encode the slice of tasks as the JSON response body.
//
// This one's short — a warm-up before the ID-based handlers below.
func (h *taskHandlers) handleList(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// parseID is provided — don't modify it. It reads the "{id}" path value
// from the request and parses it as an int, used by the three handlers
// below.
func parseID(r *http.Request) (int, error) {
	return parseIntStrict(r.PathValue("id"))
}

// YOUR TASK 3: handleGet
//
//	STEP 1: call parseID(r). If it errors, write http.StatusBadRequest and
//	        return.
//
//	STEP 2: call h.store.Get(id). If it returns ErrNotFound (check with
//	        `errors.Is(err, ErrNotFound)`), write http.StatusNotFound and
//	        return. If it's some other non-nil error, write
//	        http.StatusInternalServerError and return.
//
//	STEP 3: set Content-Type to "application/json" and encode the task.
func (h *taskHandlers) handleGet(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// YOUR TASK 4: handleMarkDone
//
// Same shape as handleGet, but calls h.store.MarkDone(id) instead of Get.
// Same error handling: ErrNotFound -> 404, anything else -> 500, success ->
// 200 with the updated task encoded as JSON.
func (h *taskHandlers) handleMarkDone(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// YOUR TASK 5: handleDelete
//
//	STEP 1: parse the ID same as before.
//
//	STEP 2: call h.store.Delete(id). ErrNotFound -> 404.
//
//	STEP 3: on success, write http.StatusNoContent (204) and return —
//	        no body to encode for a 204. Don't call json.Encode here; a 204
//	        response has an empty body by convention and encoding `nil`
//	        would write the four bytes "null", which isn't valid for this
//	        status code.
func (h *taskHandlers) handleDelete(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// writeError is a small helper you may find useful in the handlers above:
// writeError(w, http.StatusNotFound, "task not found")
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

var errBadID = errors.New("invalid id")
