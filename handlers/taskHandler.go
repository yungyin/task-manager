package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"task-manager/datastore"
	"task-manager/models"
)

var (
	TaskRegex       = regexp.MustCompile(`^/v1/tasks$`)
	TaskRegexWithId = regexp.MustCompile(`^/v1/tasks/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)
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
	case r.Method == http.MethodGet:
		h.ListTasks(w, r)
		return
	case r.Method == http.MethodPost:
		h.CreateTask(w, r)
		return
	case r.Method == http.MethodPut:
		h.UpdateTask(w, r)
		return
	default:
		return
	}
}

func (h *TasksHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	if !TaskRegex.MatchString(r.URL.Path) {
		NotFoundErrorHandler(w, r)
		return
	}

	taskList, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskList)
}

func (h *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if !TaskRegex.MatchString(r.URL.Path) {
		NotFoundErrorHandler(w, r)
		return
	}

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil || task.Name == "" || (task.Status != models.Incomplete && task.Status != models.Complete) {
		BadRequestErrorHandler(w, r)
		return
	}

	createdTask, err := h.store.Create(task)
	if err != nil {
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if !TaskRegexWithId.MatchString(r.URL.Path) {
		NotFoundErrorHandler(w, r)
		return
	}
	matches := TaskRegexWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		NotFoundErrorHandler(w, r)
		return
	}

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil || task.Name == "" || (task.Status != models.Incomplete && task.Status != models.Complete) {
		BadRequestErrorHandler(w, r)
		return
	}

	updatedTask, err := h.store.Update(matches[1], task)
	if err != nil {
		if errors.Is(err, datastore.NotFoundError) {
			NotFoundErrorHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}
