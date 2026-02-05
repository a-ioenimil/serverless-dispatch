package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	DomainAmalitech = "amalitech.com"
	DomainTraining  = "amalitechtraining.org"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	email := event.Request.UserAttributes["email"]

	// Security: Validate Email Domain
	if !isValidDomain(email) {
		// Returning an error here rejects the signup in Cognito
		return event, fmt.Errorf("sign-up restricted to @%s and @%s domains", DomainAmalitech, DomainTraining)
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
	return domain == DomainAmalitech || domain == DomainTraining
}

func main() {
	lambda.Start(handler)
}
