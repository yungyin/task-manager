package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"task-manager/datastore"
	"task-manager/models"
	"testing"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) List() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockStore) Create(models.Task) (models.Task, error) {
	args := m.Called()
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockStore) Update(string, models.Task) (models.Task, error) {
	args := m.Called()
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockStore) Delete(string) error {
	args := m.Called()
	return args.Error(0)
}

func TestTasksHandler_ListTasks_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockTaskList := []models.Task{
		{Id: "1", Name: "Task 1", Status: models.Incomplete},
		{Id: "2", Name: "Task 2", Status: models.Complete},
	}
	mockStore.On("List").Return(mockTaskList, nil)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var taskList []models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &taskList)
	assert.NoError(t, err)

	assert.Equal(t, mockTaskList, taskList)
}

func TestTasksHandler_ListTasks_GivenInvalidPath_NotFoundError(t *testing.T) {
	handler := NewTasksHandler(nil)
	request := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_ListTasks_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("List").Return([]models.Task{}, errors.New("store error"))

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_CreateTask_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockTask := models.Task{
		Name: "Task 1", Status: models.Incomplete,
	}
	mockStore.On("Create", mock.Anything).Return(mockTask, nil)

	body, _ := json.Marshal(mockTask)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var task models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &task)
	assert.NoError(t, err)

	assert.Equal(t, mockTask, task)
}

func TestTasksHandler_CreateTask_GivenInvalidPath_NotFoundError(t *testing.T) {
	handler := NewTasksHandler(nil)
	request := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_CreateTask_BadRequestError(t *testing.T) {
	mockStore := new(MockStore)
	mockTask := models.Task{
		Name: "Task 1", Status: -1,
	}

	body, _ := json.Marshal(mockTask)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_CreateTask_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockTask := models.Task{
		Name: "Task 1", Status: models.Incomplete,
	}
	mockStore.On("Create", mock.Anything).Return(models.Task{}, errors.New("store error"))

	body, _ := json.Marshal(mockTask)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_UpdateTask_Success(t *testing.T) {
	inputId := "7494d1aa-21f1-4003-8504-70602e167839"
	inputTask := models.Task{
		Name: "Task 1", Status: models.Complete,
	}

	mockStore := new(MockStore)
	mockTask := models.Task{
		Id: inputId, Name: "Task 1", Status: models.Complete,
	}
	mockStore.On("Update", mock.Anything, mock.Anything).Return(mockTask, nil)

	mockPath := "/v1/tasks/" + inputId
	body, _ := json.Marshal(inputTask)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, mockPath, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var task models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &task)
	assert.NoError(t, err)

	assert.Equal(t, mockTask, task)
}

func TestTasksHandler_UpdateTask_GivenInvalidStatus_BadRequestError(t *testing.T) {
	task := models.Task{
		Name: "Task 1", Status: -1,
	}
	body, _ := json.Marshal(task)

	mockStore := new(MockStore)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, "/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_UpdateTask_GivenMalformedTaskId_NotFoundError(t *testing.T) {
	task := models.Task{
		Name: "Task 1", Status: models.Incomplete,
	}
	body, _ := json.Marshal(task)

	mockStore := new(MockStore)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, "/v1/tasks/malformedId", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_UpdateTask_GivenNonExistedTaskId_NotFoundError(t *testing.T) {
	task := models.Task{
		Name: "Task 1", Status: models.Complete,
	}
	body, _ := json.Marshal(task)

	mockStore := new(MockStore)
	mockStore.On("Update", mock.Anything, mock.Anything).Return(models.Task{}, datastore.NotFoundError)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, "/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_UpdateTask_InternalServerError(t *testing.T) {
	task := models.Task{
		Name: "Task 1", Status: models.Complete,
	}
	body, _ := json.Marshal(task)

	mockStore := new(MockStore)
	mockStore.On("Update", mock.Anything, mock.Anything).Return(models.Task{}, errors.New("store error"))

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, "/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_DeleteTask_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Delete", mock.Anything).Return(nil)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, "/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839", nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_DeleteTask_GivenMalformedTaskId_NotFoundError(t *testing.T) {
	mockStore := new(MockStore)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, "/v1/tasks/invalid", nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_DeleteTask_GivenNonExistedTaskId_NotFoundError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Delete", mock.Anything).Return(datastore.NotFoundError)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, "/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839", nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestTasksHandler_DeleteTask_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Delete", mock.Anything).Return(errors.New("store error"))

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, "/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839", nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}
