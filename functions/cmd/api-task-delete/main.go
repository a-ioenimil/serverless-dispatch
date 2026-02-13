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

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	slog.Info("Handling task delete request", "request_id", request.RequestContext.RequestID)

	taskID, ok := request.PathParameters["taskId"]
	if !ok || taskID == "" {
		return response(http.StatusBadRequest, map[string]string{"error": "Missing task ID"}), nil
	}

	_, userRole := extractAuth(request)
	if err := svc.DeleteTask(ctx, taskID, userRole); err != nil {
		if err == services.ErrForbidden {
			return response(http.StatusForbidden, map[string]string{"error": err.Error()}), nil
		}
		if err == services.ErrTaskNotFound {
			return response(http.StatusNotFound, map[string]string{"error": err.Error()}), nil
		}
		slog.Error("Error deleting task", "error", err)
		return response(http.StatusInternalServerError, map[string]string{"error": "Internal server error"}), nil
	}

	return response(http.StatusNoContent, nil), nil
}

func extractAuth(request events.APIGatewayV2HTTPRequest) (string, string) {
	userID := "unknown"
	userRole := "MEMBER"

	if request.RequestContext.Authorizer.JWT.Claims != nil {
		if sub, ok := request.RequestContext.Authorizer.JWT.Claims["sub"]; ok && strings.TrimSpace(sub) != "" {
			userID = sub
		} else if username, ok := request.RequestContext.Authorizer.JWT.Claims["preferred_username"]; ok && strings.TrimSpace(username) != "" {
			userID = username
		} else if cognitoUsername, ok := request.RequestContext.Authorizer.JWT.Claims["cognito:username"]; ok && strings.TrimSpace(cognitoUsername) != "" {
			userID = cognitoUsername
		}
		if groups, ok := request.RequestContext.Authorizer.JWT.Claims["cognito:groups"]; ok {
			if strings.Contains(groups, "Admins") {
				userRole = "ADMIN"
			}
		}
		return userID, userRole
	}

	slog.Warn("No claims found in request context")
	return userID, userRole
}

func response(statusCode int, body interface{}) events.APIGatewayV2HTTPResponse {
	b, _ := json.Marshal(body)
	return events.APIGatewayV2HTTPResponse{
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
