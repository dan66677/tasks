package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"tz/internal/domain"
	"tz/internal/repository"
)

type TaskWorker struct {
	repo     repository.TaskRepository
	taskChan chan string
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewTaskWorker(repo repository.TaskRepository) *TaskWorker {
	return &TaskWorker{
		repo:     repo,
		taskChan: make(chan string, 100),
		stopChan: make(chan struct{}),
	}
}

func (w *TaskWorker) Start() {
	for {
		select {
		case taskID := <-w.taskChan:
			w.wg.Add(1)
			go func(id string) {
				defer w.wg.Done()
				w.processTask(id)
			}(taskID)
		case <-w.stopChan:
			return
		}
	}
}

func (w *TaskWorker) Stop() {
	close(w.stopChan)
	w.wg.Wait()
	close(w.taskChan)
}

func (w *TaskWorker) processTask(taskID string) {
	ctx := context.Background()

	task, err := w.repo.Get(ctx, taskID)
	if err != nil {
		log.Printf("Error getting task %s: %v", taskID, err)
		return
	}

	task.Status = domain.StatusRunning
	task.UpdatedAt = time.Now()
	if err := w.repo.Update(ctx, task); err != nil {
		log.Printf("Error updating task %s: %v", taskID, err)
		return
	}

	result, taskErr := w.executeTask(task)

	if taskErr != nil {
		task.Status = domain.StatusFailed
		task.Error = taskErr.Error()
	} else {
		task.Status = domain.StatusCompleted
		task.Result = result
	}
	task.UpdatedAt = time.Now()
	if err := w.repo.Update(ctx, task); err != nil {
		log.Printf("Error updating task %s: %v", taskID, err)
	}
}

func (w *TaskWorker) executeTask(task *domain.Task) (interface{}, error) {
	time.Sleep(10 * time.Second) 

	switch task.Type {
	case "example_task":
		return processExampleTask(task.Payload)
	default:
		return nil, errors.New("unknown task type")
	}
}

func (w *TaskWorker) EnqueueTask(taskID string) error {
	select {
	case w.taskChan <- taskID:
		return nil
	default:
		return errors.New("task queue is full")
	}
}

var (
	globalWorker *TaskWorker
	once         sync.Once
)

func GetTaskWorker() *TaskWorker {
	once.Do(func() {
		globalWorker = NewTaskWorker(repository.NewTaskRepository())
		go globalWorker.Start()
	})
	return globalWorker
}

func processExampleTask(payload map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"success":    true,
		"processed":  payload,
		"timestamp":  time.Now().Unix(),
	}, nil
}