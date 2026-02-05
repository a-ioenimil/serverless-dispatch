package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/task/domain"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrForbidden    = errors.New("forbidden: insufficient permissions")
)

// UpdateTask handles task updates with RBAC enforcement
func (s *TaskService) UpdateTask(ctx context.Context, req UpdateTaskRequest, userRole string, userID string) (*domain.Task, error) {
	slog.Info("Updating task", "task_id", req.ID, "user_id", userID, "role", userRole)

	// 1. Fetch Existing Task
	task, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// 2. Validate Permissions & Apply Updates
	if userRole == "ADMIN" {
		// Admins can update everything
		if req.Status != nil {
			task.Status = *req.Status
		}
		if req.Priority != nil {
			task.Priority = *req.Priority
		}
		// Allow unassigning by passing explicit empty string?
		// For now simple replacement if field is present
		if req.AssigneeID != nil {
			task.AssigneeID = req.AssigneeID
		}
	} else {
		// MEMBER Restricted Logic
		// Members can only update tasks assigned to them
		if task.AssigneeID == nil || *task.AssigneeID != userID {
			slog.Warn("Member attempted to update unassigned or other's task", "task_id", req.ID, "user_id", userID)
			return nil, ErrForbidden
		}

		// Members can DOES NOT have permission to re-assign or change priority
		if req.AssigneeID != nil || req.Priority != nil {
			slog.Warn("Member attempted to change restricted fields", "task_id", req.ID)
			return nil, ErrForbidden
		}

		// Members can ONLY update status
		if req.Status != nil {
			task.Status = *req.Status
		}
	}

	// 3. Update Timestamp
	task.UpdatedAt = time.Now().UTC()

	// 4. Save
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task update: %w", err)
	}

	slog.Info("Task updated successfully", "task_id", task.ID, "new_status", task.Status)

	return task, nil
}
