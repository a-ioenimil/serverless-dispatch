package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"log/slog"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"
	infra "github.com/a-ioenimil/serverless-dispatch/functions/internals/task/infrastructure/dynamodb"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/task/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	svc *services.TaskService
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

	client := dynamodb.NewFromConfig(cfg)
	repo := infra.NewDynamoDBTaskRepository(client, tableName)
	svc = services.NewTaskService(repo)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("Handling task list request", "request_id", request.RequestContext.RequestID)

	// Extract Auth Context
	// Note: Repeated logic - Consider refactoring into a middleware/utility in 'internals/common/auth'
	claims := request.RequestContext.Authorizer["claims"]
	claimsMap, ok := claims.(map[string]interface{})

	userID := "unknown"
	userRole := "MEMBER"

	if ok {
		if sub, ok := claimsMap["sub"].(string); ok {
			userID = sub
		}
		if groups, ok := claimsMap["cognito:groups"].(string); ok {
			if strings.Contains(groups, "Admins") {
				userRole = "ADMIN"
			}
		}
	} else {
		slog.Warn("No claims found in request context, defaulting to Member role")
	}

	tasks, err := svc.ListTasks(ctx, userRole, userID)
	if err != nil {
		// Logged in service
		return response(http.StatusInternalServerError, map[string]string{"error": "Internal server error"}), nil
	}

	// Map Domain to Response DTO
	startRes := make([]services.TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		startRes = append(startRes, services.TaskResponse{
			ID:         t.ID,
			Title:      t.Title,
			Status:     t.Status,
			Priority:   t.Priority,
			AssigneeID: t.AssigneeID,
			CreatedBy:  t.CreatedBy,
			CreatedAt:  t.CreatedAt.String(),
		})
	}

	return response(http.StatusOK, startRes), nil
}

func response(statusCode int, body interface{}) events.APIGatewayProxyResponse {
	b, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(b),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func main() {
	lambda.Start(handler)
}
