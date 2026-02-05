package services

import "github.com/a-ioenimil/serverless-dispatch/functions/internals/task/domain"

type CreateTaskRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Priority    domain.TaskPriority `json:"priority"`
	AssigneeID  *string             `json:"assignee_id,omitempty"`
}

type TaskResponse struct {
	ID         string              `json:"id"`
	Title      string              `json:"title"`
	Status     domain.TaskStatus   `json:"status"`
	Priority   domain.TaskPriority `json:"priority"`
	AssigneeID *string             `json:"assignee_id,omitempty"`
	CreatedBy  string              `json:"created_by"`
	CreatedAt  string              `json:"created_at"`
}

type UpdateTaskRequest struct {
	ID         string               `json:"id"` // From Path Parameter
	Status     *domain.TaskStatus   `json:"status,omitempty"`
	AssigneeID *string              `json:"assignee_id,omitempty"`
	Priority   *domain.TaskPriority `json:"priority,omitempty"`
}
