package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestServer spins up a fresh store + mux per test so tests don't
// interfere with each other's task IDs.
func newTestServer() *httptest.Server {
	h := &taskHandlers{store: NewTaskStore()}
	return httptest.NewServer(newMux(h))
}

func TestCreateAndList(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body, _ := json.Marshal(createTaskRequest{Title: "write tests"})
	resp, err := http.Post(srv.URL+"/tasks", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /tasks failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created task: %v", err)
	}
	if created.Title != "write tests" || created.Done {
		t.Fatalf("unexpected created task: %+v", created)
	}
	if created.ID == 0 {
		t.Fatalf("expected a non-zero ID, got %+v", created)
	}

	listResp, err := http.Get(srv.URL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks failed: %v", err)
	}
	defer listResp.Body.Close()

	var tasks []Task
	if err := json.NewDecoder(listResp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode task list: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task in list, got %d", len(tasks))
	}
}

func TestCreateEmptyTitleRejected(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body, _ := json.Marshal(createTaskRequest{Title: ""})
	resp, err := http.Post(srv.URL+"/tasks", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /tasks failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty title, got %d", resp.StatusCode)
	}
}

func TestGetTask(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body, _ := json.Marshal(createTaskRequest{Title: "find me"})
	createResp, _ := http.Post(srv.URL+"/tasks", "application/json", bytes.NewReader(body))
	var created Task
	json.NewDecoder(createResp.Body).Decode(&created)
	createResp.Body.Close()

	getResp, err := http.Get(fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID))
	if err != nil {
		t.Fatalf("GET /tasks/{id} failed: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	var got Task
	json.NewDecoder(getResp.Body).Decode(&got)
	if got.ID != created.ID || got.Title != "find me" {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestGetMissingTaskReturns404(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/tasks/999")
	if err != nil {
		t.Fatalf("GET /tasks/999 failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing task, got %d", resp.StatusCode)
	}
}

func TestMarkDone(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body, _ := json.Marshal(createTaskRequest{Title: "finish this"})
	createResp, _ := http.Post(srv.URL+"/tasks", "application/json", bytes.NewReader(body))
	var created Task
	json.NewDecoder(createResp.Body).Decode(&created)
	createResp.Body.Close()

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/tasks/%d/done", srv.URL, created.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /tasks/{id}/done failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updated Task
	json.NewDecoder(resp.Body).Decode(&updated)
	if !updated.Done {
		t.Fatalf("expected task to be marked done: %+v", updated)
	}
}

func TestDeleteTask(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body, _ := json.Marshal(createTaskRequest{Title: "delete me"})
	createResp, _ := http.Post(srv.URL+"/tasks", "application/json", bytes.NewReader(body))
	var created Task
	json.NewDecoder(createResp.Body).Decode(&created)
	createResp.Body.Close()

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /tasks/{id} failed: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	getResp, err := http.Get(fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID))
	if err != nil {
		t.Fatalf("GET after delete failed: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}
