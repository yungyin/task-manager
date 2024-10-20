package datastore

import (
	"errors"
	"reflect"
	"task-manager/models"
	"testing"
)

func TestNewMemStore(t *testing.T) {
	memStore := NewMemStore()

	if memStore == nil {
		t.Fatalf("expected non-nil MemStore, but got nil")
	}
	if memStore.taskMap == nil {
		t.Fatalf("expected taskMap to be initialized, but got nil")
	}
	if len(memStore.taskMap) != 0 {
		t.Errorf("expected taskMap to be empty, but got %d elements", len(memStore.taskMap))
	}
}

func TestMemStore_List(t *testing.T) {
	memStore := MemStore{
		taskMap: map[string]models.Task{
			"1": {Id: "1", Name: "Task 1", Status: models.Incomplete},
			"2": {Id: "2", Name: "Task 2", Status: models.Complete},
		},
	}
	expected := []models.Task{}
	for _, task := range memStore.taskMap {
		expected = append(expected, task)
	}

	actual, err := memStore.List()

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, but got: %v", expected, actual)
	}
}

func TestMemStore_Create(t *testing.T) {
	memStore := MemStore{
		taskMap: make(map[string]models.Task),
	}
	inputTask := models.Task{Name: "name", Status: models.Incomplete}

	createdTask, err := memStore.Create(inputTask)

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if createdTask.Id == "" {
		t.Errorf("expected task to have an ID, but got an empty string")
	}
	if createdTask.Name != inputTask.Name {
		t.Errorf("expected task name: %v, but got: %v", inputTask.Name, createdTask.Name)
	}
	if createdTask.Status != inputTask.Status {
		t.Errorf("expected task status: %v, but got: %v", inputTask.Status, createdTask.Status)
	}

	validateStoredTask(t, memStore, createdTask.Id, inputTask)
}

func TestMemStore_Update_GivenNonExistedId_ReturnNotFoundError(t *testing.T) {
	memStore := MemStore{
		taskMap: make(map[string]models.Task),
	}
	inputId := "1"
	inputTask := models.Task{Name: "name", Status: models.Complete}

	_, err := memStore.Update(inputId, inputTask)

	if !errors.Is(err, NotFoundError) {
		t.Fatalf("expected to be NotFoundError, but got %v", err)
	}
}

func TestMemStore_Update_GivenExistedId_UpdateInStoreAndReturnUpdatedTask(t *testing.T) {
	memStore := MemStore{
		taskMap: map[string]models.Task{
			"1": {Id: "1", Name: "Task 1", Status: models.Incomplete},
			"2": {Id: "2", Name: "Task 2", Status: models.Complete},
		},
	}
	inputId := "1"
	inputTask := models.Task{Name: "name", Status: models.Complete}
	expectedTask := models.Task{Id: "1", Name: "name", Status: models.Complete}

	updatedTask, err := memStore.Update(inputId, inputTask)

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if !reflect.DeepEqual(expectedTask, updatedTask) {
		t.Errorf("expected: %v, but got: %v", expectedTask, updatedTask)
	}

	validateStoredTask(t, memStore, inputId, inputTask)
}

func TestMemStore_Delete_GivenNonExistedId_ReturnNotFoundError(t *testing.T) {
	memStore := MemStore{
		taskMap: make(map[string]models.Task),
	}
	inputId := "1"

	err := memStore.Delete(inputId)

	if !errors.Is(err, NotFoundError) {
		t.Fatalf("expected to be NotFoundError, but got %v", err)
	}
}

func TestMemStore_Delete_GivenExistedId_DeleteFromStore(t *testing.T) {
	memStore := MemStore{
		taskMap: map[string]models.Task{
			"1": {Id: "1", Name: "Task 1", Status: models.Incomplete},
			"2": {Id: "2", Name: "Task 2", Status: models.Complete},
		},
	}
	inputId := "1"

	err := memStore.Delete(inputId)

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	_, exists := memStore.taskMap[inputId]
	if exists {
		t.Errorf("expected task to be not found, but it was stored in memStore")
	}
}

func validateStoredTask(t *testing.T, memStore MemStore, inputId string, inputTask models.Task) {
	storedTask, exists := memStore.taskMap[inputId]
	if !exists {
		t.Errorf("expected task to be stored in memStore, but it was not found")
	}

	expectedTask := models.Task{Id: inputId, Name: inputTask.Name, Status: inputTask.Status}
	if !reflect.DeepEqual(expectedTask, storedTask) {
		t.Errorf("expected: %v, but got: %v", expectedTask, storedTask)
	}
}
