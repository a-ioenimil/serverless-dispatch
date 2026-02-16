package main

import (
	"context"
	"log/slog"
	"os"

	identitycognito "github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/infrastructure/cognito"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/infrastructure/sender"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/ports"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	identityprovider "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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
	notificationsTopicARN := os.Getenv("NOTIFICATIONS_TOPIC_ARN")

	if notificationsTopicARN != "" {
		slog.Info("Initializing SNS Sender", "topic_arn", notificationsTopicARN)
		snsClient := sns.NewFromConfig(cfg)
		emailSender = sender.NewSNSSender(snsClient, notificationsTopicARN)
	} else {
		slog.Warn("NOTIFICATIONS_TOPIC_ARN not set, defaulting to Logger Sender")
		emailSender = sender.NewLoggerEmailSender()
	}

	var recipientResolver ports.RecipientResolver
	userPoolID := os.Getenv("USER_POOL_ID")
	if userPoolID != "" {
		identityClient := identityprovider.NewFromConfig(cfg)
		recipientResolver = identitycognito.NewUserDirectory(identityClient, userPoolID)
	} else {
		slog.Warn("USER_POOL_ID not set, recipient resolution by username will be disabled")
	}

	notifierService = services.NewTaskNotifierService(emailSender, recipientResolver)
}

// Handler uses the Internal Service to process events
func handler(ctx context.Context, event events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	slog.Info("Processing Dynamodb Stream", "count", len(event.Records))

	failures := make([]events.DynamoDBBatchItemFailure, 0)

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
		if err := notifierService.ProcessTaskStream(ctx, record); err != nil {
			slog.Error("Failed to process task stream record", "event_id", record.EventID, "error", err)
			failures = append(failures, events.DynamoDBBatchItemFailure{ItemIdentifier: record.EventID})
		}
	}

	return events.DynamoDBEventResponse{BatchItemFailures: failures}, nil
}

func main() {
	lambda.Start(handler)
}
