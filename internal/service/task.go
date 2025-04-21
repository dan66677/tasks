package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"tz/internal/domain"
	"tz/internal/repository"
	"tz/internal/worker"
)

type TaskService struct {
	repo   repository.TaskRepository
	worker *worker.TaskWorker
}

func NewTaskService(repo repository.TaskRepository, worker *worker.TaskWorker) *TaskService {
    return &TaskService{
        repo:   repo,
        worker: worker, 
    }
}


func (s *TaskService) CreateTask(ctx context.Context, req *domain.TaskRequest) (string, error) {
	task := &domain.Task{
		ID:        generateID(),
		Type:      req.Type,
		Status:    domain.StatusPending,
		Payload:   req.Payload,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, task); err != nil {
		return "", err
	}

	s.worker.EnqueueTask(task.ID)

	return task.ID, nil
}

func (s *TaskService) GetTaskStatus(ctx context.Context, taskID string) (domain.TaskStatus, error) {
	task, err := s.repo.Get(ctx, taskID)
	if err != nil {
		return "", err
	}

	return task.Status, nil
}

func (s *TaskService) GetTaskResult(ctx context.Context, taskID string) (*domain.TaskResult, error) {
	task, err := s.repo.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != domain.StatusCompleted {
		return nil, repository.ErrTaskNotCompleted
	}

	return &domain.TaskResult{
		ID:     task.ID,
		Result: task.Result,
	}, nil
}

func generateID() string {
	return uuid.New().String()
}

