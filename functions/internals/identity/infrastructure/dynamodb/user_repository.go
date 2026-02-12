package dynamodb

import (
	"context"
	"fmt"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBUserRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBUserRepository(client *dynamodb.Client, tableName string) *DynamoDBUserRepository {
	return &DynamoDBUserRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBUserRepository) Save(ctx context.Context, user domain.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Single Table Design modifications if needed (e.g. PK/SK)
	// Assuming simple table or specific User table from context,
	// but usually STD requires PK="USER#<ID>", SK="METADATA"
	// modifying item keys for STD support:
	item["PK"], _ = attributevalue.Marshal("USER#" + user.ID)
	item["SK"], _ = attributevalue.Marshal("METADATA")

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to put item to dynamodb: %w", err)
	}

	return nil
}

func (r *DynamoDBUserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	// Not implemented for Post-Signup, but needed for interface
	return nil, nil // Placeholder
}

func (r *DynamoDBUserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        &r.tableName,
		FilterExpression: aws.String("begins_with(PK, :pk_prefix) AND SK = :sk_meta"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk_prefix": &types.AttributeValueMemberS{Value: "USER#"},
			":sk_meta":   &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan users: %w", err)
	}

	users := make([]domain.User, 0, len(out.Items))
	for _, item := range out.Items {
		var user domain.User
		if err := attributevalue.UnmarshalMap(item, &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
