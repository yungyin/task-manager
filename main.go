package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/v1/tasks", &tasksHandler{})
	mux.Handle("/v1/tasks/", &tasksHandler{})

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type tasksHandler struct {
}

func (h *tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
