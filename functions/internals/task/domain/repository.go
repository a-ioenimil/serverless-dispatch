package domain

import "context"

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	ListByAssignee(ctx context.Context, assigneeID string) ([]Task, error)
	ListAll(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task *Task) error
}
