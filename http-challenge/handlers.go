package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Exercise: HTTP server — a tiny task-tracker API (Fiber edition)
//
// Concepts:
//   - A Fiber handler has the signature `func(c *fiber.Ctx) error`. You
//     read the request and write the response through c, and the return
//     value is what Fiber's error-handling middleware acts on — on the
//     success path you can almost always just `return c.JSON(...)` etc.
//     directly, since those methods already return an error themselves.
//   - c.BodyParser(&v) decodes the request body into v — it looks at the
//     Content-Type header to pick JSON/form/etc, so requests need
//     `Content-Type: application/json` set for a JSON body to parse.
//   - c.Status(code) sets the status for the response that follows it in
//     the same chain — `c.Status(fiber.StatusCreated).JSON(task)`. If you
//     never call Status, a successful response defaults to 200.
//   - c.Params("id") reads a route param (the ":id" segment from
//     routes.go). It's always a string — you still need strconv to get an
//     int out of it, same as before.
//
// This is a bigger exercise than the earlier ones. Do the handlers in the
// order they appear.
//
// The store (store.go) is fully implemented — use it, don't reimplement
// task storage here. routes.go wires these handlers to a *fiber.App —
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
//	STEP 1: declare `var req createTaskRequest` and call
//	        c.BodyParser(&req). If it errors, return
//	        `c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})`
//	        — that single line both writes the response AND is the value
//	        this handler returns.
//
//	STEP 2: if req.Title is empty, return the same shape of response with
//	        message "title is required".
//
//	STEP 3: call h.store.Create(req.Title) to get the new *Task.
//
//	STEP 4: return `c.Status(fiber.StatusCreated).JSON(task)`.
func (h *taskHandlers) handleCreate(c *fiber.Ctx) error {
	var req createTaskRequest

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}

	result := h.store.Create(req.Title)

	return c.Status(fiber.StatusCreated).JSON(result)

}

// YOUR TASK 2: handleList
//
//	STEP 1: call h.store.List() to get all tasks.
//
//	STEP 2: return c.JSON(tasks) — no need to set a status, 200 is the
//	        default.
//
// This one's short — a warm-up before the ID-based handlers below.
func (h *taskHandlers) handleList(c *fiber.Ctx) error {
	result := h.store.List()

	return c.Status(200).JSON(result)
}

// parseID is provided — don't modify it. It reads the ":id" route param
// from the context and parses it as an int, used by the three handlers
// below.
func parseID(c *fiber.Ctx) (int, error) {
	return parseIntStrict(c.Params("id"))
}

// YOUR TASK 3: handleGet
//
//	STEP 1: call parseID(c). If it errors, return
//	        `c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})`.
//
//	STEP 2: call h.store.Get(id). If it returns ErrNotFound (check with
//	        `errors.Is(err, ErrNotFound)`), return a 404 in the same
//	        fiber.Map shape. If it's some other non-nil error, return a
//	        500 the same way.
//
//	STEP 3: return c.JSON(task).
func (h *taskHandlers) handleGet(c *fiber.Ctx) error {

	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_id"})
	}

	result, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(result)

}

// YOUR TASK 4: handleMarkDone
//
// Same shape as handleGet, but calls h.store.MarkDone(id) instead of Get.
// Same error handling: ErrNotFound -> 404, anything else -> 500, success ->
// c.JSON(task) with the updated task (200 by default).
func (h *taskHandlers) handleMarkDone(c *fiber.Ctx) error {

	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_id"})
	}

	result, err := h.store.MarkDone(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(result)
}

// YOUR TASK 5: handleDelete
//
//	STEP 1: parse the ID same as before.
//
//	STEP 2: call h.store.Delete(id). ErrNotFound -> 404.
//
//	STEP 3: on success, return `c.SendStatus(fiber.StatusNoContent)` — no
//	        body for a 204, so don't also call .JSON() on this one.
func (h *taskHandlers) handleDelete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_id"})
	}

	err = h.store.Delete(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

var errBadID = errors.New("invalid id")
