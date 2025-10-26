# Makefile

.PHONY: build clean deploy invoke-health lint deps

# Define your AWS SAM parameters
STACK_NAME = your-stack-name
REGION = your-region
TEMPLATE = template.yaml

# Install dependencies
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint:
	golangci-lint run ./...

# Build the serverless application
build:
	sam build

# Clean up the build artifacts
clean:
	sam delete --stack-name $(STACK_NAME) --region $(REGION)

# Deploy the serverless application
deploy:
	sam deploy --stack-name $(STACK_NAME) --region $(REGION) --guided

# Invoke the Lambda function locally with the health check event
invoke-health:
	sam build
	sam local invoke -e .invokation/health.json --env-vars env.json

invoke-ai:
	sam build
	sam local invoke -e .invokation/ai.json --env-vars env.json