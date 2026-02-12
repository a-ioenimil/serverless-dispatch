package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/notification/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SNSSender struct {
	client   *sns.Client
	topicARN string
}

func NewSNSSender(client *sns.Client, topicARN string) *SNSSender {
	return &SNSSender{
		client:   client,
		topicARN: topicARN,
	}
}

type notificationPayload struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Channel   string `json:"channel"`
}

func (s *SNSSender) Send(ctx context.Context, msg domain.NotificationMessage) error {
	payload, err := json.Marshal(notificationPayload{
		Recipient: msg.Recipient,
		Subject:   msg.Subject,
		Body:      msg.Body,
		Channel:   "email",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal sns notification payload: %w", err)
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(s.topicARN),
		Message:  aws.String(string(payload)),
		Subject:  aws.String(msg.Subject),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"recipient": {
				DataType:    aws.String("String"),
				StringValue: aws.String(msg.Recipient),
			},
			"channel": {
				DataType:    aws.String("String"),
				StringValue: aws.String("email"),
			},
		},
	}

	_, err = s.client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish notification to sns: %w", err)
	}

	return nil
}
