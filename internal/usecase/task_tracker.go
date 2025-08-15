package usecase

import (
	"context"
	"errors"
	"lo/internal/entity"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Tracker struct {
	tasks map[string]entity.Task
	logCh chan entity.Task
	m     sync.Mutex
}

type TaskTracker interface {
	CreateTask(entity.Task) string
	GetTask(string) (entity.Task, error)
	ListTasks() ([]entity.Task, error)
}

func NewTracker(ctx context.Context) TaskTracker {
	t := &Tracker{
		tasks: make(map[string]entity.Task),
		logCh: make(chan entity.Task),
	}

	go t.processTasks(ctx)

	return t
}

func (t *Tracker) CreateTask(task entity.Task) string {
	id := generateRandomID()

	t.m.Lock()

	newTask := entity.Task{
		ID:          id,
		Name:        task.Name,
		Description: task.Description,
		Status:      "new",
	}

	t.tasks[id] = newTask

	t.m.Unlock()

	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		t.logCh <- newTask
	}()

	return id
}

func (t *Tracker) GetTask(id string) (entity.Task, error) {
	t.m.Lock()
	defer t.m.Unlock()

	if task, ok := t.tasks[id]; ok {
		return task, nil
	}

	return entity.Task{}, errors.New("task not found")
}

func (t *Tracker) ListTasks() ([]entity.Task, error) {
	t.m.Lock()
	defer t.m.Unlock()

	if len(t.tasks) == 0 {
		return nil, errors.New("no tasks found")
	}

	taskList := make([]entity.Task, 0, len(t.tasks))

	for _, task := range t.tasks {
		taskList = append(taskList, task)
	}

	return taskList, nil
}

const IDlen = 5
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomID() string {
	b := make([]byte, IDlen)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

func (t *Tracker) processTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-t.logCh:
			log.Printf("New task created: %v\n", task)
			t.m.Lock()
			delete(t.tasks, task.ID)
			t.m.Unlock()
		}
	}
}
