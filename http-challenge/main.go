package main

import (
	"fmt"
	"log"
	"net/http"
)

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
// without needing a real running server (httptest.NewRecorder).
func main() {
	h := &taskHandlers{store: NewTaskStore()}
	mux := newMux(h)

	addr := ":8080"
	fmt.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
