package main

import "github.com/gofiber/fiber/v2"

// newApp wires up the routes. Given — don't modify it. Fiber's routing
// looks a lot like Express: app.Method(path, handler). ":id" is a param
// segment you read back with c.Params("id") inside the handler.
func newApp(h *taskHandlers) *fiber.App {
	app := fiber.New()
	app.Post("/tasks", h.handleCreate)
	app.Get("/tasks", h.handleList)
	app.Get("/tasks/:id", h.handleGet)
	app.Patch("/tasks/:id/done", h.handleMarkDone)
	app.Delete("/tasks/:id", h.handleDelete)
	return app
}
