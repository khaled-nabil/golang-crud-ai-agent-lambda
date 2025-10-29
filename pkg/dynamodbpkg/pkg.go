package dynamodbpkg

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDB struct {
	Client    *dynamodb.Client
	TableName string
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
