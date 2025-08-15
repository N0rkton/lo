package usecase_test

import (
	"context"
	"testing"
	"time"

	"lo/internal/entity"
	"lo/internal/usecase"

	"github.com/stretchr/testify/assert"
)

// Тестирование успешного создания новой задачи
func TestCreateTask_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tt := usecase.NewTracker(ctx)
	task := entity.Task{Name: "Test Task", Description: "This is a test"}
	createdID := tt.CreateTask(task)
	assert.NotEmpty(t, createdID)
}

// Проверка ошибок при попытке создать новую задачу дважды с одним именем
func TestCreateTask_DuplicateName_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tt := usecase.NewTracker(ctx)
	task := entity.Task{Name: "Duplicate Name", Description: "This should be unique"}
	firstID := tt.CreateTask(task)
	secondID := tt.CreateTask(task)

	assert.NotEqual(t, firstID, secondID)
}

// Получение существующей задачи по правильному ID
func TestGetTask_ExistingTask_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tt := usecase.NewTracker(ctx)
	task := entity.Task{Name: "Test Task", Description: "This is a test"}
	createdID := tt.CreateTask(task)
	retrievedTask, err := tt.GetTask(createdID)

	assert.NoError(t, err)
	assert.Equal(t, retrievedTask.Name, task.Name)
	assert.Equal(t, retrievedTask.Description, task.Description)
}

// Ошибка при получении несуществующей задачи
func TestGetTask_NonExistentTask_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tt := usecase.NewTracker(ctx)
	_, err := tt.GetTask("nonexistent-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
}

// Список всех созданных задач
func TestListTasks_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tt := usecase.NewTracker(ctx)
	task1 := entity.Task{Name: "Task One", Description: "First task"}
	task2 := entity.Task{Name: "Task Two", Description: "Second task"}

	tt.CreateTask(task1)
	tt.CreateTask(task2)

	listedTasks, err := tt.ListTasks()

	assert.NoError(t, err)
	assert.Len(t, listedTasks, 2)
}

// Отсутствие задач возвращает ошибку
func TestListTasks_NoTasks_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tt := usecase.NewTracker(ctx)
	_, err := tt.ListTasks()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no tasks found")
}
