package datastore

import (
	"github.com/stretchr/testify/assert"
	"task-manager/models"
	"testing"
)

const (
	nonExistedId = "0"
	existedId    = "1"
)

var memStore = MemStore{
	taskMap: map[string]models.Task{
		"1": {Id: "1", Name: "Task 1", Status: models.Incomplete},
		"2": {Id: "2", Name: "Task 2", Status: models.Complete},
	},
}
var inputTask = models.Task{Name: "name", Status: models.Complete}

func TestNewMemStore(t *testing.T) {
	store := NewMemStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.taskMap)
	assert.Equal(t, 0, len(store.taskMap))
}

func TestMemStore_List(t *testing.T) {
	expected := []models.Task{
		{Id: "1", Name: "Task 1", Status: models.Incomplete},
		{Id: "2", Name: "Task 2", Status: models.Complete},
	}

	actual, err := memStore.List()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMemStore_Create(t *testing.T) {
	createdTask, err := memStore.Create(inputTask)

	assert.NoError(t, err)
	assert.NotEmpty(t, createdTask.Id)
	assert.Equal(t, createdTask.Name, inputTask.Name)
	assert.Equal(t, createdTask.Status, inputTask.Status)

	validateStoredTask(t, memStore, createdTask.Id, inputTask)
}

func TestMemStore_Update_GivenNonExistedId_ReturnNotFoundError(t *testing.T) {
	_, err := memStore.Update(nonExistedId, inputTask)

	assert.Equal(t, err, NotFoundError)
}

func TestMemStore_Update_GivenExistedId_UpdateInStoreAndReturnUpdatedTask(t *testing.T) {
	updatedTask, err := memStore.Update(existedId, inputTask)

	assert.NoError(t, err)
	assert.Equal(t, existedId, updatedTask.Id)
	assert.Equal(t, updatedTask.Name, inputTask.Name)
	assert.Equal(t, updatedTask.Status, inputTask.Status)

	validateStoredTask(t, memStore, existedId, inputTask)
}

func TestMemStore_Delete_GivenNonExistedId_ReturnNotFoundError(t *testing.T) {
	err := memStore.Delete(nonExistedId)

	assert.Equal(t, err, NotFoundError)
}

func TestMemStore_Delete_GivenExistedId_DeleteFromStore(t *testing.T) {
	err := memStore.Delete(existedId)
	assert.NoError(t, err)

	_, exists := memStore.taskMap[existedId]
	assert.False(t, exists)
}

func validateStoredTask(t *testing.T, memStore MemStore, inputId string, inputTask models.Task) {
	storedTask, exists := memStore.taskMap[inputId]
	assert.True(t, exists)

	expectedTask := models.Task{Id: inputId, Name: inputTask.Name, Status: inputTask.Status}
	assert.Equal(t, expectedTask, storedTask)
}
