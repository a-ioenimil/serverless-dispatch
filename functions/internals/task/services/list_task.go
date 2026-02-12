package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/task/domain"
)

// ListTasks retrieves a list of tasks based on the user's role.
// Admins see all tasks. members see only tasks assigned to them.
func (s *TaskService) ListTasks(ctx context.Context, userRole string, userID string) ([]domain.Task, error) {
	slog.Info("Listing tasks", "role", userRole, "user_id", userID)

	if userRole == "ADMIN" {
		tasks, err := s.repo.ListAll(ctx)
		if err != nil {
			slog.Error("Failed to list all tasks", "error", err)
			return nil, fmt.Errorf("failed to list all tasks: %w", err)
		}
		return tasks, nil
	}

	assignedTasks, err := s.repo.ListByAssignee(ctx, userID)
	if err != nil {
		slog.Error("Failed to list assigned tasks", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to list assigned tasks: %w", err)
	}

	unassignedTasks, err := s.repo.ListUnassigned(ctx)
	if err != nil {
		slog.Error("Failed to list unassigned tasks", "error", err)
		return nil, fmt.Errorf("failed to list unassigned tasks: %w", err)
	}

	combined := make(map[string]domain.Task, len(assignedTasks)+len(unassignedTasks))
	for _, task := range assignedTasks {
		combined[task.ID] = task
	}
	for _, task := range unassignedTasks {
		combined[task.ID] = task
	}

	result := make([]domain.Task, 0, len(combined))
	for _, task := range combined {
		result = append(result, task)
	}

	return result, nil
}
