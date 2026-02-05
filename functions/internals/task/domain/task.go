package domain

import (
	"time"
)

// TaskStatus defines the state of a task
type TaskStatus string

const (
	TaskStatusOpen       TaskStatus = "OPEN"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"
)

// TaskPriority defines the importance of a task
type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
)

// Task represents a unit of work in the system
type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	AssigneeID  *string      `json:"assignee_id,omitempty"` // Nullable if unassigned
	CreatedBy   string       `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func NewTask(id, title, description, createdBy string, priority TaskPriority) (*Task, error) {
	if title == "" {
		return nil, ErrTaskTitleRequired
	}
	if len(description) > 2048 {
		return nil, ErrDescriptionTooLong
	}

	now := time.Now().UTC()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      TaskStatusOpen,
		Priority:    priority,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
