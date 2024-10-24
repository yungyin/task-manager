package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"task-manager/datastore"
	"task-manager/models"
)

const (
	ContentType     = "Content-Type"
	JsonContentType = "application/json"
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

func (handler *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		handler.ListTasks(w, r)
		return
	case r.Method == http.MethodPost:
		handler.CreateTask(w, r)
		return
	case r.Method == http.MethodPut:
		handler.UpdateTask(w, r)
		return
	case r.Method == http.MethodDelete:
		handler.DeleteTask(w, r)
	default:
		return
	}
}

func (handler *TasksHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	if !TaskRegex.MatchString(r.URL.Path) {
		NotFoundErrorHandler(w, r)
		return
	}

	taskList, err := handler.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskList)
}

func (handler *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	createdTask, err := handler.store.Create(task)
	if err != nil {
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (handler *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId := extractTaskId(r)
	if taskId == "" {
		NotFoundErrorHandler(w, r)
		return
	}

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil || task.Name == "" || (task.Status != models.Incomplete && task.Status != models.Complete) {
		BadRequestErrorHandler(w, r)
		return
	}

	updatedTask, err := handler.store.Update(taskId, task)
	if err != nil {
		if errors.Is(err, datastore.NotFoundError) {
			NotFoundErrorHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

func (handler *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := extractTaskId(r)
	if taskId == "" {
		NotFoundErrorHandler(w, r)
		return
	}

	err := handler.store.Delete(taskId)
	if err != nil {
		if errors.Is(err, datastore.NotFoundError) {
			NotFoundErrorHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r, err)
		return
	}

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusNoContent)
}

func extractTaskId(r *http.Request) string {
	if !TaskRegexWithId.MatchString(r.URL.Path) {
		return ""
	}
	matches := TaskRegexWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
