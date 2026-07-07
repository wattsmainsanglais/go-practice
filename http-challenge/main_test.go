package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// newTestApp spins up a fresh store + app per test so tests don't
// interfere with each other's task IDs.
func newTestApp() *fiber.App {
	h := &taskHandlers{store: NewTaskStore()}
	return newApp(h)
}

// jsonRequest builds an httptest.Request with a JSON body and the right
// Content-Type header — Fiber's BodyParser needs that header set to know
// how to decode the body.
func jsonRequest(method, target string, body any) *http.Request {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, target, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestCreateAndList(t *testing.T) {
	app := newTestApp()

	req := jsonRequest(http.MethodPost, "/tasks", createTaskRequest{Title: "write tests"})
	resp, err := app.Test(req)
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

	listResp, err := app.Test(httptest.NewRequest(http.MethodGet, "/tasks", nil))
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
	app := newTestApp()

	req := jsonRequest(http.MethodPost, "/tasks", createTaskRequest{Title: ""})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /tasks failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty title, got %d", resp.StatusCode)
	}
}

func TestGetTask(t *testing.T) {
	app := newTestApp()

	createResp, _ := app.Test(jsonRequest(http.MethodPost, "/tasks", createTaskRequest{Title: "find me"}))
	var created Task
	json.NewDecoder(createResp.Body).Decode(&created)
	createResp.Body.Close()

	getResp, err := app.Test(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", created.ID), nil))
	if err != nil {
		t.Fatalf("GET /tasks/:id failed: %v", err)
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
	app := newTestApp()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/tasks/999", nil))
	if err != nil {
		t.Fatalf("GET /tasks/999 failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing task, got %d", resp.StatusCode)
	}
}

func TestMarkDone(t *testing.T) {
	app := newTestApp()

	createResp, _ := app.Test(jsonRequest(http.MethodPost, "/tasks", createTaskRequest{Title: "finish this"}))
	var created Task
	json.NewDecoder(createResp.Body).Decode(&created)
	createResp.Body.Close()

	resp, err := app.Test(httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/tasks/%d/done", created.ID), nil))
	if err != nil {
		t.Fatalf("PATCH /tasks/:id/done failed: %v", err)
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
	app := newTestApp()

	createResp, _ := app.Test(jsonRequest(http.MethodPost, "/tasks", createTaskRequest{Title: "delete me"}))
	var created Task
	json.NewDecoder(createResp.Body).Decode(&created)
	createResp.Body.Close()

	resp, err := app.Test(httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", created.ID), nil))
	if err != nil {
		t.Fatalf("DELETE /tasks/:id failed: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	getResp, err := app.Test(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", created.ID), nil))
	if err != nil {
		t.Fatalf("GET after delete failed: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}
