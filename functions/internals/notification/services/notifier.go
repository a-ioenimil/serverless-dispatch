package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/ports"
	"github.com/aws/aws-lambda-go/events"
)

type TaskNotifierService struct {
	sender ports.EmailSender
}

func NewTaskNotifierService(sender ports.EmailSender) *TaskNotifierService {
	return &TaskNotifierService{sender: sender}
}

// ProcessTaskStream handles DynamoDB stream events involving tasks
func (s *TaskNotifierService) ProcessTaskStream(ctx context.Context, record events.DynamoDBStreamRecord) {
	// 1. Determine Event Type & Context
	switch record.EventName {
	case "INSERT":
		s.handleTaskCreated(ctx, record.NewImage)
	case "MODIFY":
		s.handleTaskUpdated(ctx, record.OldImage, record.NewImage)
	}
}

func (s *TaskNotifierService) handleTaskCreated(ctx context.Context, newImage map[string]events.DynamoDBAttributeValue) {
	assignee := extractString(newImage, "AssigneeID")
	// Safe fallback to lowercase
	if assignee == "" {
		assignee = extractString(newImage, "assignee_id")
	}

	if assignee != "" {
		title := extractString(newImage, "Title")
		if title == "" {
			title = extractString(newImage, "title")
		}

		msg := domain.NotificationMessage{
			Recipient: assignee,
			Subject:   "New Task Assigned",
			Body:      fmt.Sprintf("You have been assigned a new task: %s", title),
		}
		s.safeSend(ctx, msg)
	}
}

func (s *TaskNotifierService) handleTaskUpdated(ctx context.Context, oldImage, newImage map[string]events.DynamoDBAttributeValue) {
	title := extractString(newImage, "Title")
	if title == "" {
		title = extractString(newImage, "title")
	}

	oldStatus := extractString(oldImage, "Status")
	if oldStatus == "" {
		oldStatus = extractString(oldImage, "status")
	}

	newStatus := extractString(newImage, "Status")
	if newStatus == "" {
		newStatus = extractString(newImage, "status")
	}

	if oldStatus != newStatus {
		// Notify Admins
		s.safeSend(ctx, domain.NotificationMessage{
			Recipient: "admins@amalitech.com",
			Subject:   "Task Status Update",
			Body:      fmt.Sprintf("Task '%s' status changed from %s to %s", title, oldStatus, newStatus),
		})

		// Notify Assignee
		assignee := extractString(newImage, "AssigneeID")
		if assignee == "" {
			assignee = extractString(newImage, "assignee_id")
		}

		if assignee != "" {
			s.safeSend(ctx, domain.NotificationMessage{
				Recipient: assignee,
				Subject:   "Task Status Updated",
				Body:      fmt.Sprintf("Your task '%s' is now %s", title, newStatus),
			})
		}
	}
}

func (s *TaskNotifierService) safeSend(ctx context.Context, msg domain.NotificationMessage) {
	if err := s.sender.Send(ctx, msg); err != nil {
		slog.Error("Failed to send notification", "error", err, "recipient", msg.Recipient)
	}
}

// Helper to handle DynamoDB Map values robustly
func extractString(image map[string]events.DynamoDBAttributeValue, key string) string {
	val, ok := image[key]
	if !ok || val.DataType() != events.DataTypeString {
		return ""
	}
	return val.String()
}
