package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/infrastructure/dynamodb"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	userService           *services.UserService
	snsClient             *sns.Client
	notificationsTopicARN string
)

func init() {
	log := logger.InitLogger()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error("Unable to load SDK config", "error", err)
		os.Exit(1)
	}

	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Error("TABLE_NAME env var is required")
		os.Exit(1)
	}

	client := db.NewFromConfig(cfg)
	repo := dynamodb.NewDynamoDBUserRepository(client, tableName)
	userService = services.NewUserService(repo)

	snsClient = sns.NewFromConfig(cfg)
	notificationsTopicARN = os.Getenv("NOTIFICATIONS_TOPIC_ARN")
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	slog.Info("Processing PostConfirmation", "triggerSource", event.TriggerSource)

	if event.TriggerSource == "PostConfirmation_ConfirmSignUp" {
		email := event.Request.UserAttributes["email"]
		sub := event.Request.UserAttributes["sub"] // Cognito User ID
		username := event.Request.UserAttributes["preferred_username"]
		// Fallback for sub if not in attr? Usually sub is standard.
		if sub == "" {
			sub = event.UserName // Sometimes UserName is the sub
		}
		if username == "" {
			username = event.UserName
		}

		if email == "" || sub == "" {
			slog.Warn("Missing email or sub in event", "email", email, "sub", sub)
			return event, nil
		}

		err := userService.CreateUser(ctx, sub, email, username)
		if err != nil {
			slog.Error("Failed to create user", "error", err)
			return event, err
		}
		slog.Info("User created successfully", "id", sub, "username", username)

		if notificationsTopicARN == "" {
			slog.Warn("NOTIFICATIONS_TOPIC_ARN is empty, skipping notification subscription", "email", email)
			return event, nil
		}

		if err := subscribeUserToNotifications(ctx, email, sub, username); err != nil {
			slog.Error("Failed to subscribe user to notifications", "error", err, "email", email)
			return event, err
		}
	}

	return event, nil
}

func subscribeUserToNotifications(ctx context.Context, email, sub, username string) error {
	recipientFilterValues := make([]string, 0, 3)
	for _, value := range []string{email, sub, username} {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			recipientFilterValues = append(recipientFilterValues, trimmed)
		}
	}

	if len(recipientFilterValues) == 0 {
		return fmt.Errorf("at least one recipient filter value is required")
	}

	filterPolicy, err := json.Marshal(map[string][]string{
		"recipient": recipientFilterValues,
		"channel":   {"email"},
	})
	if err != nil {
		return err
	}

	_, err = snsClient.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn:              aws.String(notificationsTopicARN),
		Protocol:              aws.String("email"),
		Endpoint:              aws.String(email),
		Attributes:            map[string]string{"FilterPolicy": string(filterPolicy)},
		ReturnSubscriptionArn: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
