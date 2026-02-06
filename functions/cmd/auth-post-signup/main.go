package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/infrastructure/dynamodb"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	userService *services.UserService
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
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	slog.Info("Processing PostConfirmation", "triggerSource", event.TriggerSource)

	if event.TriggerSource == "PostConfirmation_ConfirmSignUp" {
		email := event.Request.UserAttributes["email"]
		sub := event.Request.UserAttributes["sub"] // Cognito User ID
		// Fallback for sub if not in attr? Usually sub is standard.
		if sub == "" {
			sub = event.UserName // Sometimes UserName is the sub
		}

		if email == "" || sub == "" {
			slog.Warn("Missing email or sub in event", "email", email, "sub", sub)
			return event, nil
		}

		err := userService.CreateUser(ctx, sub, email)
		if err != nil {
			slog.Error("Failed to create user", "error", err)
			return event, err
		}
		slog.Info("User created successfully", "id", sub)
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
