package services

import (
	"context"
	"fmt"
	"log/slog"
)

// DeleteTask removes a task with RBAC enforcement
func (s *TaskService) DeleteTask(ctx context.Context, taskID string, userRole string) error {
	slog.Info("Deleting task", "task_id", taskID, "role", userRole)

	if userRole != "ADMIN" {
		return ErrForbidden
	}

	task, err := s.repo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to fetch task: %w", err)
	}
	if task == nil {
		return ErrTaskNotFound
	}

	if err := s.repo.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
