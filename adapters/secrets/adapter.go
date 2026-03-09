package secrets

import (
	"context"
	"encoding/json"
	"os"

	"ai-agent/model/datamodels"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func New() (*datamodels.AppConfig, error) {
	return loadSecrets()
}

func loadSecrets() (*datamodels.AppConfig, error) {
	ctx := context.Background()

	secretsArn := os.Getenv("SECRETS_ARN")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	smClient := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretsArn,
	}

	result, err := smClient.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	var appConfig datamodels.AppConfig
	err = json.Unmarshal([]byte(*result.SecretString), &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
