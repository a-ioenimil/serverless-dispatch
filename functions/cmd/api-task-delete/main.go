package apitaskdelete
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
































































































}	lambda.Start(handler)func main() {}	}		},			"Access-Control-Allow-Origin": "*",			"Content-Type":                "application/json",		Headers: map[string]string{		Body:       string(b),		StatusCode: statusCode,	return events.APIGatewayV2HTTPResponse{	b, _ := json.Marshal(body)func response(statusCode int, body interface{}) events.APIGatewayV2HTTPResponse {}	return userID, userRole	slog.Warn("No claims found in request context")	}		return userID, userRole		}			}				userRole = "ADMIN"			if strings.Contains(groups, "Admins") {		if groups, ok := request.RequestContext.Authorizer.JWT.Claims["cognito:groups"]; ok {		}			userID = cognitoUsername		} else if cognitoUsername, ok := request.RequestContext.Authorizer.JWT.Claims["cognito:username"]; ok && strings.TrimSpace(cognitoUsername) != "" {			userID = username		} else if username, ok := request.RequestContext.Authorizer.JWT.Claims["preferred_username"]; ok && strings.TrimSpace(username) != "" {			userID = sub		if sub, ok := request.RequestContext.Authorizer.JWT.Claims["sub"]; ok && strings.TrimSpace(sub) != "" {	if request.RequestContext.Authorizer.JWT.Claims != nil {	userRole := "MEMBER"	userID := "unknown"func extractAuth(request events.APIGatewayV2HTTPRequest) (string, string) {}	return response(http.StatusNoContent, nil), nil	}		return response(http.StatusInternalServerError, map[string]string{"error": "Internal server error"}), nil		slog.Error("Error deleting task", "error", err)		}			return response(http.StatusNotFound, map[string]string{"error": err.Error()}), nil		if err == services.ErrTaskNotFound {		}			return response(http.StatusForbidden, map[string]string{"error": err.Error()}), nil		if err == services.ErrForbidden {	if err := svc.DeleteTask(ctx, taskID, userRole); err != nil {	_, userRole := extractAuth(request)	}		return response(http.StatusBadRequest, map[string]string{"error": "Missing task ID"}), nil	if !ok || taskID == "" {	taskID, ok := request.PathParameters["taskId"]	slog.Info("Handling task delete request", "request_id", request.RequestContext.RequestID)func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {}	svc = services.NewTaskService(repo)	repo := infra.NewDynamoDBTaskRepository(client, tableName)	client := dynamodb.NewFromConfig(cfg)	}		os.Exit(1)		log.Error("TABLE_NAME env var is required")	if tableName == "" {	tableName := os.Getenv("TABLE_NAME")	}		os.Exit(1)		log.Error("Unable to load SDK config", "error", err)	if err != nil {	cfg, err := config.LoadDefaultConfig(context.Background())	log := logger.InitLogger()func init() {)	svc *services.TaskServicevar ()	"github.com/aws/aws-sdk-go-v2/service/dynamodb"	"github.com/aws/aws-sdk-go-v2/config"	"github.com/aws/aws-lambda-go/lambda"	"github.com/aws/aws-lambda-go/events"	"github.com/a-ioenimil/serverless-dispatch/functions/internals/task/services"	infra "github.com/a-ioenimil/serverless-dispatch/functions/internals/task/infrastructure/dynamodb"	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"