package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/infrastructure/dynamodb"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

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

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	slog.Info("Handling user list request", "request_id", request.RequestContext.RequestID)

	if request.RequestContext.Authorizer.JWT.Claims == nil {
		return response(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"}), nil
	}

	users, err := userService.ListUsers(ctx)
	if err != nil {
		slog.Error("Error listing users", "error", err)
		return response(http.StatusInternalServerError, map[string]string{"error": "Internal server error"}), nil
	}

	res := make([]userResponse, 0, len(users))
	for _, user := range users {
		res = append(res, userResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
			Status:   user.Status,
		})
	}

	return response(http.StatusOK, res), nil
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
