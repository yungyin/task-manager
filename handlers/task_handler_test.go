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

const (
	sampleTaskId        = "7494d1aa-21f1-4003-8504-70602e167839"
	pathWithoutId       = "/v1/tasks"
	pathWithValidId     = pathWithoutId + "/" + sampleTaskId
	pathWithMalformedId = "/v1/tasks/malformedId"
	taskName            = "Task 1"
)

var (
	task = models.Task{
		Name: taskName, Status: models.Incomplete,
	}
	invalidTask = models.Task{
		Name: taskName, Status: -1,
	}
	storeError = errors.New("store error")
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
	request := httptest.NewRequest(http.MethodGet, pathWithoutId, nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))

	var taskList []models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &taskList)
	assert.NoError(t, err)
	assert.Equal(t, mockTaskList, taskList)
}

func TestTasksHandler_ListTasks_GivenInvalidPath_NotFoundError(t *testing.T) {
	handler := NewTasksHandler(nil)
	request := httptest.NewRequest(http.MethodGet, pathWithMalformedId, nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_ListTasks_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("List").Return([]models.Task{}, storeError)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodGet, pathWithoutId, nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_CreateTask_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Create", mock.Anything).Return(task, nil)

	body, _ := json.Marshal(task)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPost, pathWithoutId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))

	var actualTask models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &actualTask)
	assert.NoError(t, err)
	assert.Equal(t, task, actualTask)
}

func TestTasksHandler_CreateTask_GivenInvalidPath_NotFoundError(t *testing.T) {
	handler := NewTasksHandler(nil)
	request := httptest.NewRequest(http.MethodGet, pathWithMalformedId, nil)
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_CreateTask_BadRequestError(t *testing.T) {
	mockStore := new(MockStore)

	body, _ := json.Marshal(invalidTask)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPost, pathWithoutId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_CreateTask_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Create", mock.Anything).Return(models.Task{}, storeError)

	body, _ := json.Marshal(task)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPost, pathWithoutId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.CreateTask(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_UpdateTask_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockTask := models.Task{
		Id: sampleTaskId, Name: taskName, Status: models.Complete,
	}
	mockStore.On("Update", mock.Anything, mock.Anything).Return(mockTask, nil)

	body, _ := json.Marshal(task)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, pathWithValidId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))

	var actualTask models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &actualTask)
	assert.NoError(t, err)
	assert.Equal(t, mockTask, actualTask)
}

func TestTasksHandler_UpdateTask_GivenInvalidStatus_BadRequestError(t *testing.T) {
	mockStore := new(MockStore)

	body, _ := json.Marshal(invalidTask)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, pathWithValidId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_UpdateTask_GivenMalformedTaskId_NotFoundError(t *testing.T) {
	mockStore := new(MockStore)

	body, _ := json.Marshal(task)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, pathWithMalformedId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_UpdateTask_GivenNonExistedTaskId_NotFoundError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Update", mock.Anything, mock.Anything).Return(models.Task{}, datastore.NotFoundError)

	body, _ := json.Marshal(task)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, pathWithValidId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_UpdateTask_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Update", mock.Anything, mock.Anything).Return(models.Task{}, storeError)

	body, _ := json.Marshal(task)
	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodPut, pathWithValidId, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.UpdateTask(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_DeleteTask_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Delete", mock.Anything).Return(nil)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, pathWithValidId, nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_DeleteTask_GivenMalformedTaskId_NotFoundError(t *testing.T) {
	mockStore := new(MockStore)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, pathWithMalformedId, nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_DeleteTask_GivenNonExistedTaskId_NotFoundError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Delete", mock.Anything).Return(datastore.NotFoundError)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, pathWithValidId, nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestTasksHandler_DeleteTask_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("Delete", mock.Anything).Return(storeError)

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodDelete, pathWithValidId, nil)
	recorder := httptest.NewRecorder()

	handler.DeleteTask(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}
