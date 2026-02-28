# AI Agent Lambda

This project is a serverless AI agent built with Go and deployed on AWS Lambda.
It uses the Gin framework for routing and integrates with Google's Gemini AI model.
The application is designed to be a conversational AI, with plans to implement a sophisticated RAG (Retrieval-Augmented Generation) system.

## Technical Features

- Serverless architecture using AWS Lambda.
- RESTful API built with Gin.
- Integration with Google Gemini AI model for conversational capabilities.
- Dependency injection using Google Wire.

## Features

- Basic conversational AI functionality.
- Current RAG Systems implements a simple context window approach of 10 Conversations.
- Maintains single user conversation history in DynamoDB.

## Future Work

- **Create a smart RAG system for AI:**
  - In addition to the current context window approach, implement a more advanced RAG system that summarizes and retains important context over longer conversations.
    - When a user has more than 10 messages in history, it would summarize the user's chat into a new record in DynamoDB. This new summarization record would trim the context.
    - The system would include all summaries of the user, and the latest 10 messages.
    - When the user has more than X summaries (to be defined based on context window), we would use an in-memory vector search. All summaries are vectorized using an embedding model (to be decided on), and based on similarity search, we'd pull the relevant summaries that wouldn't exceed the context window.
- **Develop a frontend React application.**
- Implement user authentication and management using Cognito.
- Enhance error handling and logging mechanisms.
  - Handle AI hallucinations and API errors gracefully.
  - Implement error codes and logging for better monitoring and frontend integration.
- Implement knowledge base, using RAG System to pull relevant documents from S3 or other data sources.
  - Vectorize documents and store them in a vector database for efficient retrieval.
  - Integrate document retrieval into the conversational AI flow.
  - Add Prompts for the AI to use when relevant documents are found.

## Dependencies & Resources

- Go 1.24.0
- AWS Lambda Go
- AWS Secret Manager
- AWS DynamoDB
- AWS SAM CLI
- Terraform

## Setup and Run

### Prerequisites

- Go 1.24.0 or later
- AWS SAM CLI
- Docker
- Terraform
- Google Gemini API Key
- AWS Account with necessary permissions
- AWS CLI configured

### Install Dependencies

```bash
go mod tidy
```

### Configuration (AWS cloud)

1. Configure your AWS environment if not already done, preferably using AWS CLI V2 with SSO:
  ```bash
  aws configure sso
  ```
2. Initialise Terraform resources:
  ```bash
  make build
  ```
3. Plan and apply Terraform changes using your AWS profile:
  ```bash
  AWS_PROFILE="PROFILE_NAME" make plan
  AWS_PROFILE="PROFILE_NAME" make apply
  ```
4. Store your Google Gemini API key in AWS SecretsManager with the name `GEMINI_API_KEY`, and your Gemini model name with the name `MODEL_ID`.

Now the application is deployed on AWS Lambda and ready to use.

### Running Locally with LocalStack

You can run the full stack locally using LocalStack to emulate AWS services.

1. Start LocalStack using Docker Compose:
   ```bash
   docker compose up localstack
   ```
2. In a separate terminal, apply Terraform resources against LocalStack using `tflocal`:
   ```bash
   cd .tf
   tflocal init
   tflocal apply
   cd ..
   ```
3. Export AWS environment variables so the Go SDK and SAM talk to LocalStack:
   ```bash
   export AWS_ACCESS_KEY_ID=test
   export AWS_SECRET_ACCESS_KEY=test
   export AWS_REGION=eu-central-1
   export AWS_ENDPOINT_URL=http://localhost:4566
   ```
4. Run the Lambda locally via SAM using `template.yaml` and one of the `invokation/*.json` events:
   ```bash
   sam local invoke AiAgentServerlessAPI \
     -t template.yaml \
     -e invokation/<event-file>.json
   ```

5. **Call the API via REST (API Gateway)**  
   The stack exposes the Lambda through API Gateway. Get the base URL from Terraform (run from the project root). LocalStack uses the `/_aws/execute-api/<api_id>/<stage>` path:
   ```bash
   cd .tf && tflocal output -raw localstack_invoke_url && cd ..
   ```
   Then call your endpoints (append the path to the base URL):
   ```bash
   # Health check
   curl "$(cd .tf && tflocal output -raw localstack_invoke_url)/api/v1/health"

   # Send a message (replace <BASE_URL> with the output from above)
   curl -X POST "<BASE_URL>/api/v1/send" \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello, how are you?", "user_id": "754778bb-e286-4d26-8b32-77c939f5ee59"}'
   ```
