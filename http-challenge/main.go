package main

import "log"

// Given — don't modify. Once the handlers in handlers.go are implemented,
// run this with `go run .` and try it against curl:
//
//	curl -X POST localhost:8080/tasks -d '{"title":"write go practice"}'
//	curl localhost:8080/tasks
//	curl localhost:8080/tasks/1
//	curl -X PATCH localhost:8080/tasks/1/done
//	curl -X DELETE localhost:8080/tasks/1
//
// Or run `go test ./...` — main_test.go exercises all five handlers
// via Fiber's app.Test(), no real listening socket needed.
func main() {
	h := &taskHandlers{store: NewTaskStore()}
	app := newApp(h)

	log.Println("listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
