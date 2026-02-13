package dynamodb

import (
	"context"
	"fmt"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/task/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Metadata struct {
	PK     string      `dynamodbav:"PK"`
	SK     string      `dynamodbav:"SK"`
	GSI1PK string      `dynamodbav:"GSI1_PK,omitempty"`
	GSI1SK string      `dynamodbav:"GSI1_SK,omitempty"`
	Data   domain.Task `dynamodbav:",inline"` // Flattens task fields into top level
}

type DynamoDBTaskRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBTaskRepository(client *dynamodb.Client, tableName string) *DynamoDBTaskRepository {
	return &DynamoDBTaskRepository{
		client:    client,
		tableName: tableName,
	}
}

// Save creates or overwrites a task
func (r *DynamoDBTaskRepository) Save(ctx context.Context, task *domain.Task) error {
	pk := fmt.Sprintf("TASK#%s", task.ID)
	sk := "METADATA"

	item := Metadata{
		PK:   pk,
		SK:   sk,
		Data: *task,
	}

	// Add GSI attributes if assigned
	if task.AssigneeID != nil {
		item.GSI1PK = fmt.Sprintf("USER#%s", *task.AssigneeID)
		item.GSI1SK = fmt.Sprintf("TASK#%s", task.ID)
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put item to dynamodb: %w", err)
	}

	return nil
}

func (r *DynamoDBTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	pk := fmt.Sprintf("TASK#%s", id)
	sk := "METADATA"

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if out.Item == nil {
		return nil, nil // Not found
	}

	var metadata Metadata
	if err := attributevalue.UnmarshalMap(out.Item, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &metadata.Data, nil
}

func (r *DynamoDBTaskRepository) ListByAssignee(ctx context.Context, assigneeID string) ([]domain.Task, error) {
	// Query GSI1: PK = USER#{assigneeID}
	gsiPK := fmt.Sprintf("USER#%s", assigneeID)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1_PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: gsiPK},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	tasks := make([]domain.Task, 0, len(out.Items))
	for _, item := range out.Items {
		var metadata Metadata
		if err := attributevalue.UnmarshalMap(item, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		tasks = append(tasks, metadata.Data)
	}

	return tasks, nil
}

func (r *DynamoDBTaskRepository) ListAll(ctx context.Context) ([]domain.Task, error) {
	// SCAN Operation: Warning - expensive on large tables.
	// For production systems with high volume, use a GSI with low cardinality sharding or specific access patterns.
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		// Filter to only get Task items (PK starts with TASK#)
		FilterExpression: aws.String("begins_with(PK, :pk_prefix) AND SK = :sk_meta"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk_prefix": &types.AttributeValueMemberS{Value: "TASK#"},
			":sk_meta":   &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan tasks: %w", err)
	}

	tasks := make([]domain.Task, 0, len(out.Items))
	for _, item := range out.Items {
		var metadata Metadata
		if err := attributevalue.UnmarshalMap(item, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		tasks = append(tasks, metadata.Data)
	}

	return tasks, nil
}

func (r *DynamoDBTaskRepository) ListUnassigned(ctx context.Context) ([]domain.Task, error) {
	// SCAN Operation: Warning - expensive on large tables.
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		FilterExpression: aws.String(
			"begins_with(PK, :pk_prefix) AND SK = :sk_meta AND attribute_not_exists(assignee_id)",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk_prefix": &types.AttributeValueMemberS{Value: "TASK#"},
			":sk_meta":   &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan unassigned tasks: %w", err)
	}

	tasks := make([]domain.Task, 0, len(out.Items))
	for _, item := range out.Items {
		var metadata Metadata
		if err := attributevalue.UnmarshalMap(item, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		tasks = append(tasks, metadata.Data)
	}

	return tasks, nil
}

func (r *DynamoDBTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	// Re-save for now (since PutItem overwrites).
	// In production, we might want strict conditional updates.
	return r.Save(ctx, task)
}

func (r *DynamoDBTaskRepository) Delete(ctx context.Context, id string) error {
	pk := fmt.Sprintf("TASK#%s", id)
	sk := "METADATA"

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
