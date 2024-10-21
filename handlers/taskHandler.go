package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"task-manager/models"
)

var (
	TaskRegex = regexp.MustCompile(`^/v1/tasks$`)
)

type TasksHandler struct {
	store TasksStore
}

func NewTasksHandler(s TasksStore) *TasksHandler {
	return &TasksHandler{
		store: s,
	}
}

type TasksStore interface {
	List() ([]models.Task, error)
	Create(task models.Task) (models.Task, error)
	Update(taskId string, task models.Task) (models.Task, error)
	Delete(taskId string) error
}

func (h *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && TaskRegex.MatchString(r.URL.Path):
		h.ListTasks(w, r)
		return
	default:
		return
	}
}

func (h *TasksHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	taskList, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskList)
}
