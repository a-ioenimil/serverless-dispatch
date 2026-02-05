package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	email := event.Request.UserAttributes["email"]

	// Security: Validate Email Domain
	if !isValidDomain(email) {
		// Returning an error here rejects the signup in Cognito
		return event, fmt.Errorf("sign-up restricted to authorized domains only")
	}

	// Auto-confirm the user if needed (Configuration decision)
	// event.Response.AutoConfirmUser = true

	return event, nil
}

func isValidDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	allowedEnv := os.Getenv("ALLOWED_EMAIL_DOMAINS")
	if allowedEnv == "" {
		// Fallback safe default or Deny All
		return false
	}

	allowedDomains := strings.Split(allowedEnv, ",")
	for _, d := range allowedDomains {
		if domain == strings.TrimSpace(d) {
			return true
		}
	}

	return false
}

func main() {
	lambda.Start(handler)
}
