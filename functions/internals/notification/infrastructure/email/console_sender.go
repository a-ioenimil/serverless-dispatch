package email

import (
	"context"
	"log/slog"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
)

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) Send(ctx context.Context, n domain.Notification) error {
	// Infrastructure implementation: This isolates the "How" (Console/SES/SMTP)
	slog.Info("📧 SENDING NOTIFICATION",
		"recipient", n.Recipient,
		"subject", n.Subject,
		"body", n.Body,
	)
	return nil
}
