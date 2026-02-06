package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/infrastructure/sender"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/ports"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

var (
	notifierService *services.TaskNotifierService
)

func init() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("Unable to load SDK config", "error", err)
		os.Exit(1)
	}

	var emailSender ports.EmailSender
	fromEmail := os.Getenv("FROM_EMAIL")

	if fromEmail != "" {
		slog.Info("Initializing SES Sender", "from", fromEmail)
		sesClient := ses.NewFromConfig(cfg)
		emailSender = sender.NewSESSender(sesClient, fromEmail)
	} else {
		slog.Warn("FROM_EMAIL not set, defaulting to Logger Sender")
		emailSender = sender.NewLoggerEmailSender()
	}

	notifierService = services.NewTaskNotifierService(emailSender)
}

// Handler uses the Internal Service to process events
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

		slog.Info("Processing Task Event for Service", "event_id", record.EventID)
		notifierService.ProcessTaskStream(ctx, record)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
