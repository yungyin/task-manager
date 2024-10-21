package main

import (
	"log"
	"net/http"
	"task-manager/datastore"
	"task-manager/handlers"
)

func main() {
	taskMemStore := datastore.NewMemStore()
	tasksHandler := handlers.NewTasksHandler(taskMemStore)

	mux := http.NewServeMux()

	mux.Handle("/v1/tasks", tasksHandler)
	mux.Handle("/v1/tasks/", tasksHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
