package dynamodbpkg

import (
	"ai-agent/model/datamodels"
	"ai-agent/repositories/chatpersistance"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	Client    *dynamodb.Client
	TableName string
	IDKey     string
	SortKey   string
}

var (
	ctx = context.Background()
)

func New() (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &DynamoDB{
		Client:    dynamodb.NewFromConfig(cfg),
		TableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		IDKey:     os.Getenv("DYNAMODB_USER_ID_KEY"),
		SortKey:   os.Getenv("DYNAMODB_TIMESTAMP_KEY"),
	}, nil
}

func (d *DynamoDB) StoreItem(item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.TableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to store item: %w", err)
	}

	return nil
}

func (d *DynamoDB) RetrieveItems(id string, limit int32) ([]datamodels.HistoryContext, error) {
	out, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.TableName),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :uid", d.IDKey)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: id},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items: %w", err)
	}

	if len(out.Items) == 0 {
		return nil, nil
	}

	var items []chatpersistance.ChatMessage
	err = attributevalue.UnmarshalListOfMaps(out.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	var history []datamodels.HistoryContext
	for _, item := range items {
		history = append(history, *item.Conversation)
	}

	return history, nil
}
