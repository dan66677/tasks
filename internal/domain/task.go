package domain

import (
	"time"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type TaskRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Status    TaskStatus             `json:"status"`
	Payload   map[string]interface{} `json:"payload"`
	Result    interface{}           `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type TaskResult struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result"`
}
