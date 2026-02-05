package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/task/domain"
	uuid "github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized: only admins can create tasks")
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, cmd CreateTaskRequest, userRole string, userID string) (*domain.Task, error) {
	// RBAC Check
	if userRole != "ADMIN" {
		return nil, ErrUnauthorized
	}

	taskID := uuid.New().String()

	task, err := domain.NewTask(
		taskID,
		cmd.Title,
		cmd.Description,
		userID,
		cmd.Priority,
	)
	if err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	// Assign immediately if requested
	if cmd.AssigneeID != nil {
		task.AssigneeID = cmd.AssigneeID
	}

	if err := s.repo.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("repository save failed: %w", err)
	}

	return task, nil
}
