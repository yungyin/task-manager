package datastore

import (
	"errors"
	"github.com/google/uuid"
	"task-manager/models"
)

var (
	NotFoundError = errors.New("task not found")
)

type MemStore struct {
	taskMap map[string]models.Task
}

func NewMemStore() *MemStore {
	taskMap := make(map[string]models.Task)
	return &MemStore{
		taskMap,
	}
}

func (memStore MemStore) List() ([]models.Task, error) {
	taskList := []models.Task{}
	for _, task := range memStore.taskMap {
		taskList = append(taskList, task)
	}
	return taskList, nil
}

func (memStore MemStore) Create(task models.Task) (models.Task, error) {
	task.Id = uuid.New().String()
	memStore.taskMap[task.Id] = task
	return task, nil
}

func (memStore MemStore) Update(taskId string, task models.Task) (models.Task, error) {
	_, exists := memStore.taskMap[taskId]
	if !exists {
		return task, NotFoundError
	}

	task.Id = taskId
	memStore.taskMap[taskId] = task
	return task, nil
}

func (memStore MemStore) Delete(taskId string) error {
	_, exists := memStore.taskMap[taskId]
	if !exists {
		return NotFoundError
	}

	delete(memStore.taskMap, taskId)
	return nil
}
