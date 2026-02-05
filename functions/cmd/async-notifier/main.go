package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Simplified Notification Model
type NotificationPayload struct {
	To      string
	Subject string
	Body    string
}

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	slog.Info("Processing Dynamodb Stream", "count", len(event.Records))

	for _, record := range event.Records {
		// Only interest in Task Metadata changes (PK starts with TASK#)
		pk, ok := record.Change.Keys["PK"]
		if !ok || pk.DataType() != events.DataTypeString || len(pk.String()) < 5 || pk.String()[:5] != "TASK#" {
			continue
		}

		// Filter out non-metadata items if using Single Table Design aggressively
		sk, ok := record.Change.Keys["SK"]
		if !ok || sk.String() != "METADATA" {
			continue
		}

		slog.Info("Processing Task Event", "event_id", record.EventID, "event_name", record.EventName)

		switch record.EventName {
		case "INSERT":
			handleTaskCreated(record.Change.NewImage)
		case "MODIFY":
			handleTaskUpdated(record.Change.OldImage, record.Change.NewImage)
		}
	}
	return nil
}

func handleTaskCreated(newImage map[string]events.DynamoDBAttributeValue) {
	// Task Assigned on Creation
	assigneeID := extractString(newImage, "assignee_id") // JSON tag is assignee_id, ensure casing matches storage
	// Note: In DynamoDB storage it might be stored map keys depending on marshaller.
	// Our 'dynamodbav' tags in struct were like `json:"assignee_id"`.
	// Usually attributevalue marshaler uses struct field names or `dynamodbav` tags.
	// Let's assume standard 'assignee_id' if that's what was used.
	// Wait, standard struct uses Go Field Name if no tag?
	// The struct was: AssigneeID  *string `json:"assignee_id,omitempty"`
	// It did NOT have dynamodbav tag on that field in the previous step?
	// checking domain/task.go: `AssigneeID  *string      json:"assignee_id,omitempty"`
	// Default behavior of attributevalue marshal is to use struct field name "AssigneeID" unless `dynamodbav` tag exists.
	// NOTE: This is a common pitfall. The current repository implementation uses `item := Metadata{ Data: *task }` with `inline`.
	// The `Metadata` struct had `Data domain.Task dynamodbav:",inline"`.
	// So the fields will be at top level. Since `domain.Task` has NO `dynamodbav` tags on `AssigneeID`,
	// it will default to "AssigneeID".

	// Let's safe check both because strict coupling to field names is fragile.
	assignee := extractString(newImage, "AssigneeID")
	if assignee == "" {
		assignee = extractString(newImage, "assignee_id")
	}

	if assignee != "" {
		title := extractString(newImage, "Title")
		sendEmail(assignee, "New Task Assigned", fmt.Sprintf("You have been assigned a new task: %s", title))
	}
}

func handleTaskUpdated(oldImage, newImage map[string]events.DynamoDBAttributeValue) {
	title := extractString(newImage, "Title")

	// 1. Check Status Change
	oldStatus := extractString(oldImage, "Status")
	newStatus := extractString(newImage, "Status")

	if oldStatus != newStatus {
		// Notify Admins (mock group email)
		sendEmail("admins@amalitech.com", "Task Status Update",
			fmt.Sprintf("Task '%s' status changed from %s to %s", title, oldStatus, newStatus))

		// Notify Assignee
		assignee := extractString(newImage, "AssigneeID")
		if assignee != "" {
			sendEmail(assignee, "Task Status Updated",
				fmt.Sprintf("Your task '%s' is now %s", title, newStatus))
		}
	}
}

// Mock SES Sender
func sendEmail(to, subject, body string) {
	// In production, use "github.com/aws/aws-sdk-go-v2/service/ses"
	slog.Info("📧 SENDING EMAIL", "to", to, "subject", subject, "body", body)
}

func extractString(image map[string]events.DynamoDBAttributeValue, key string) string {
	val, ok := image[key]
	if !ok || val.DataType() != events.DataTypeString {
		return ""
	}
	return val.String()
}

func main() {
	lambda.Start(handler)
}
