package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/common/logger"
	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/services"
	"github.com/aws/aws-lambda-go/events"
)

var (
	policySvc *services.AuthPolicyService
)

func init() {
	logger.InitLogger()
	// Initialize Service with Configuration (injected dependencies)
	allowedEnv := os.Getenv("ALLOWED_EMAIL_DOMAINS")
	policySvc = services.NewAuthPolicyService(allowedEnv)
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	email := event.Request.UserAttributes["email"]

	slog.Info("Validating new user signup", "email", email)

	if err := policySvc.ValidateSignup(email); err != nil {
		slog.Warn("Signup rejected", "reason", err, "email", email)
		return event, err
	}

	return event, nil
}
