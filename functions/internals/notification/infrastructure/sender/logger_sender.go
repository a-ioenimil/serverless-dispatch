package sender

import (
	"context"
	"log/slog"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
)

// LoggerEmailSender implements EmailSender by logging to stdout (simulating SES)
type LoggerEmailSender struct{}

func NewLoggerEmailSender() *LoggerEmailSender {
	return &LoggerEmailSender{}
}

func (s *LoggerEmailSender) Send(ctx context.Context, msg domain.NotificationMessage) error {
	slog.Info("📧 SENDING EMAIL",
		"to", msg.Recipient,
		"subject", msg.Subject,
		"body", msg.Body,
	)
	return nil
}
