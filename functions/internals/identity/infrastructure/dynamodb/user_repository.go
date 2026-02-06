package dynamodb

import (
	"context"
	"fmt"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/domain"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
