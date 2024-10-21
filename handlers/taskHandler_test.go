package handlers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
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
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var taskList []models.Task
	err := json.Unmarshal(recorder.Body.Bytes(), &taskList)
	assert.NoError(t, err)

	assert.Equal(t, mockTaskList, taskList)
}

func TestTasksHandler_ListTasks_InternalServerError(t *testing.T) {
	mockStore := new(MockStore)
	mockStore.On("List").Return([]models.Task{}, errors.New("store error"))

	handler := NewTasksHandler(mockStore)
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	recorder := httptest.NewRecorder()

	handler.ListTasks(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}
