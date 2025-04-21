package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"tz/internal/domain"
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrTaskNotCompleted = errors.New("task not completed")
)

type TaskRepository interface {
	Save(ctx context.Context, task *domain.Task) error
	Get(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
}

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *InMemoryTaskRepository) Save(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepository) Get(ctx context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (r *InMemoryTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return nil
}

