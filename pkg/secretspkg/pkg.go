package secretspkg

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type (
	AppConfig struct {
		GeminiAPIKey string `json:"GEMINI_API_KEY"`
		ModelID      string `json:"MODEL_ID"`
		DBHost       string `json:"DB_HOST"`
		DBPort       string `json:"DB_PORT"`
		DBName       string `json:"DB_NAME"`
		DBUser       string `json:"DB_USER"`
		DBPassword   string `json:"DB_PASSWORD"`
	}
)

func New() (*AppConfig, error) {
	return loadSecrets()
}

func loadSecrets() (*AppConfig, error) {
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

	var appConfig AppConfig
	err = json.Unmarshal([]byte(*result.SecretString), &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
