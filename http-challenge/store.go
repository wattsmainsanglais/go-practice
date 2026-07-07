package main

import (
	"fmt"
	"sync"
)

// Task is the resource this API manages — deliberately shaped like the
// "project task" concept in andamio-cli (a title + done/not-done state).
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// TaskStore is a thread-safe in-memory task store. Fully implemented —
// don't modify it. You've already built a mutex-protected shared structure
// in the worker pool exercise; this is the same idea, just guarding a map
// instead of a results slice. An HTTP server handles requests on multiple
// goroutines concurrently (net/http spins up one per request), so anything
// shared across handlers — like this store — needs the same protection a
// worker pool needs.
type TaskStore struct {
	mu     sync.Mutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

// Create adds a new task with the given title and returns it.
func (s *TaskStore) Create(title string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := &Task{ID: s.nextID, Title: title, Done: false}
	s.tasks[t.ID] = t
	s.nextID++
	return t
}

// List returns all tasks. Order is not guaranteed.
func (s *TaskStore) List() []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out
}

// ErrNotFound is returned by Get, MarkDone, and Delete when no task with
// the given ID exists.
var ErrNotFound = fmt.Errorf("task not found")

// Get returns the task with the given ID, or ErrNotFound.
func (s *TaskStore) Get(id int) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// MarkDone sets a task's Done field to true, or returns ErrNotFound.
func (s *TaskStore) MarkDone(id int) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	t.Done = true
	return t, nil
}

// Delete removes a task, or returns ErrNotFound.
func (s *TaskStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
