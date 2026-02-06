package ports

import (
	"context"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
)

// EmailSender defines the contract for sending emails
type EmailSender interface {
	Send(ctx context.Context, msg domain.NotificationMessage) error
}
