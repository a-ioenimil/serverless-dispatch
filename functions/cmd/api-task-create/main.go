package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

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
		log.Error("unable to load SDK config", "error", err)
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
	slog.Info("Handling create task request", "request_id", request.RequestContext.RequestID)

	// 1. Parse Body
	var req services.CreateTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		slog.Warn("Invalid request body", "error", err)
		return response(http.StatusBadRequest, map[string]string{"error": "Invalid request body"}), nil
	}

	// 2. Extract Auth Context
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
		slog.Warn("No claims found in request context")
	}

	// 3. Call Service
	task, err := svc.CreateTask(ctx, req, userRole, userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			slog.Warn("Unauthorized attempt to create task", "user_id", userID)
			return response(http.StatusForbidden, map[string]string{"error": err.Error()}), nil
		}
		slog.Error("Error creating task", "error", err)
		return response(http.StatusInternalServerError, map[string]string{"error": "Internal server error"}), nil
	}

	// 4. Return Response
	resDTO := services.TaskResponse{
		ID:         task.ID,
		Title:      task.Title,
		Status:     task.Status,
		Priority:   task.Priority,
		AssigneeID: task.AssigneeID,
		CreatedBy:  task.CreatedBy,
		CreatedAt:  task.CreatedAt.String(),
	}

	slog.Info("Task created successfully", "task_id", task.ID)
	return response(http.StatusCreated, resDTO), nil
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
