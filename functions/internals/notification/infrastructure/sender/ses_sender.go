package sender

import (
	"context"
	"fmt"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESSender struct {
	client      *ses.Client
	sourceEmail string
}

func NewSESSender(client *ses.Client, sourceEmail string) *SESSender {
	return &SESSender{
		client:      client,
		sourceEmail: sourceEmail,
	}
}

func (s *SESSender) Send(ctx context.Context, msg domain.NotificationMessage) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{msg.Recipient},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(msg.Body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(msg.Subject),
			},
		},
		Source: aws.String(s.sourceEmail),
	}

	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	return nil
}
