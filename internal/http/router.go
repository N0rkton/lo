package http

import (
	"context"
	"encoding/json"
	"errors"
	"lo/internal/entity"
	"lo/internal/usecase"
	"log"
	"net/http"
)

type Router struct {
	taskTracker usecase.TaskTracker
}

func NewRouter(ctx context.Context, taskTracker usecase.TaskTracker, hostPort string, errorChan chan error) (func(context.Context) error, error) {
	router := &Router{
		taskTracker: taskTracker,
	}

	server := &http.Server{Addr: hostPort}

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			router.getTasks(w, req)
		case http.MethodPost:
			router.createTask(w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path != "/tasks/" {
			router.getTaskByID(w, r)
		} else {
			http.Error(w, "should use get method", http.StatusMethodNotAllowed)
		}
	})

	server.Handler = mux

	log.Printf("listening on %s\n", hostPort)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorChan <- err
			close(errorChan)
		}
	}()

	shutdownFunc := func(ctx context.Context) error {
		return server.Shutdown(ctx)
	}

	return shutdownFunc, nil
}

func (r *Router) createTask(w http.ResponseWriter, req *http.Request) {
	var task entity.Task

	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		http.Error(w, "can't decode task:"+err.Error(), http.StatusBadRequest)
		return
	}

	id := r.taskTracker.CreateTask(entity.Task{})
	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, "can't encode id: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) getTasks(w http.ResponseWriter, req *http.Request) {
	taskList, err := r.taskTracker.ListTasks()
	if err != nil {
		http.Error(w, "can't get task list: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		http.Error(w, "can't encode task: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) getTaskByID(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := r.taskTracker.GetTask(id)
	if err != nil {
		http.Error(w, "can't get task: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "can't encode task: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
