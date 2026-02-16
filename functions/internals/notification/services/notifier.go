package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/ports"
	"github.com/aws/aws-lambda-go/events"
)

type TaskNotifierService struct {
	sender   ports.EmailSender
	resolver ports.RecipientResolver
}

func NewTaskNotifierService(sender ports.EmailSender, resolver ports.RecipientResolver) *TaskNotifierService {
	return &TaskNotifierService{sender: sender, resolver: resolver}
}

// ProcessTaskStream handles DynamoDB stream events involving tasks
func (s *TaskNotifierService) ProcessTaskStream(ctx context.Context, record events.DynamoDBEventRecord) error {
	// 1. Determine Event Type & Context
	switch record.EventName {
	case "INSERT":
		return s.handleTaskCreated(ctx, record.Change.NewImage)
	case "MODIFY":
		return s.handleTaskUpdated(ctx, record.Change.OldImage, record.Change.NewImage)
	}

	return nil
}

func (s *TaskNotifierService) handleTaskCreated(ctx context.Context, newImage map[string]events.DynamoDBAttributeValue) error {
	assignee := extractString(newImage, "AssigneeID")
	// Safe fallback to lowercase
	if assignee == "" {
		assignee = extractString(newImage, "assignee_id")
	}

	if assignee == "" {
		return nil
	}

	title := extractString(newImage, "Title")
	if title == "" {
		title = extractString(newImage, "title")
	}

	recipient := s.resolveRecipient(ctx, assignee)
	if recipient == "" {
		return nil
	}

	msg := domain.NotificationMessage{
		Recipient: recipient,
		Subject:   "New Task Assigned",
		Body:      fmt.Sprintf("You have been assigned a new task: %s", title),
	}

	return s.sendNotification(ctx, msg)
}

func (s *TaskNotifierService) handleTaskUpdated(ctx context.Context, oldImage, newImage map[string]events.DynamoDBAttributeValue) error {
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
	oldAssignee := extractString(oldImage, "AssigneeID")
	if oldAssignee == "" {
		oldAssignee = extractString(oldImage, "assignee_id")
	}
	newAssignee := extractString(newImage, "AssigneeID")
	if newAssignee == "" {
		newAssignee = extractString(newImage, "assignee_id")
	}

	statusChanged := oldStatus != newStatus
	assigneeChanged := oldAssignee != newAssignee

	if !statusChanged && !assigneeChanged {
		return nil
	}

	adminEmails := s.adminRecipients(ctx)
	var firstErr error

	if statusChanged {
		for _, adminEmail := range adminEmails {
			err := s.sendNotification(ctx, domain.NotificationMessage{
				Recipient: adminEmail,
				Subject:   "Task Status Update",
				Body:      fmt.Sprintf("Task '%s' status changed from %s to %s", title, oldStatus, newStatus),
			})
			if firstErr == nil && err != nil {
				firstErr = err
			}
		}

		recipient := s.resolveRecipient(ctx, newAssignee)
		if recipient != "" {
			err := s.sendNotification(ctx, domain.NotificationMessage{
				Recipient: recipient,
				Subject:   "Task Status Updated",
				Body:      fmt.Sprintf("Your task '%s' is now %s", title, newStatus),
			})
			if firstErr == nil && err != nil {
				firstErr = err
			}
		}
	}

	if assigneeChanged && newAssignee != "" {
		recipient := s.resolveRecipient(ctx, newAssignee)
		if recipient != "" {
			err := s.sendNotification(ctx, domain.NotificationMessage{
				Recipient: recipient,
				Subject:   "Task Reassigned",
				Body:      fmt.Sprintf("You were assigned task '%s'", title),
			})
			if firstErr == nil && err != nil {
				firstErr = err
			}
		}

		for _, adminEmail := range adminEmails {
			err := s.sendNotification(ctx, domain.NotificationMessage{
				Recipient: adminEmail,
				Subject:   "Task Reassignment",
				Body:      fmt.Sprintf("Task '%s' reassigned from %s to %s", title, oldAssignee, newAssignee),
			})
			if firstErr == nil && err != nil {
				firstErr = err
			}
		}
	}

	if assigneeChanged && newAssignee == "" {
		for _, adminEmail := range adminEmails {
			err := s.sendNotification(ctx, domain.NotificationMessage{
				Recipient: adminEmail,
				Subject:   "Task Unassigned",
				Body:      fmt.Sprintf("Task '%s' was unassigned (previous assignee: %s)", title, oldAssignee),
			})
			if firstErr == nil && err != nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func (s *TaskNotifierService) sendNotification(ctx context.Context, msg domain.NotificationMessage) error {
	if err := s.sender.Send(ctx, msg); err != nil {
		slog.Error("Failed to send notification", "error", err, "recipient", msg.Recipient)
		return err
	}

	return nil
}

func (s *TaskNotifierService) resolveRecipient(ctx context.Context, usernameOrEmail string) string {
	if usernameOrEmail == "" {
		return ""
	}

	if strings.Contains(usernameOrEmail, "@") {
		return usernameOrEmail
	}

	if s.resolver == nil {
		return usernameOrEmail
	}

	email, err := s.resolver.ResolveUserEmailByUsername(ctx, usernameOrEmail)
	if err != nil {
		slog.Error("failed to resolve user email", "username", usernameOrEmail, "error", err)
		return usernameOrEmail
	}

	if strings.TrimSpace(email) == "" {
		return usernameOrEmail
	}

	return email
}

func (s *TaskNotifierService) adminRecipients(ctx context.Context) []string {
	if s.resolver == nil {
		return nil
	}

	emails, err := s.resolver.ListAdminEmails(ctx)
	if err != nil {
		slog.Error("failed to list admin emails", "error", err)
		return nil
	}

	return emails
}

// Helper to handle DynamoDB Map values robustly
func extractString(image map[string]events.DynamoDBAttributeValue, key string) string {
	val, ok := image[key]
	if !ok || val.DataType() != events.DataTypeString {
		return ""
	}
	return val.String()
}
