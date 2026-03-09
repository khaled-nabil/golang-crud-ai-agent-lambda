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

- Basic conversational AI functionality (using Gemini).
- Maintains user conversation history in Postgres.
- (Under Progress) Retrieves user related context using Vector Similarity search, as user conversational history is stored as Embeddings.

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
- PostgreSQL (local via Docker)
- Terraform
- Localstack (for local environment deployment)

## Setup and Run

### Prerequisites

- Go 1.24.0 or later
- Docker
- Docker Compose (local)
- Localstack (local)
- Terraform
- Google Gemini API Key
- AWS Account with necessary permissions
- AWS CLI configured

### Considerations

Currently the Postgres DB is run through Docker, via docker compose, together with Localstack.
As Localstack free version does not support RDS.

When it comes to production, feel free to add your own AWS RDS Terraform configurations, or run Postgres remotely elsewhere.

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
4. Add the secret value to AWS Secrets Manager manually. Terraform creates the secret store; you must add the secret version with a JSON object containing these keys:

   | Key             | Description                                         |
   |-----------------|-----------------------------------------------------|
   | `GEMINI_API_KEY`| Your Google Gemini API key                          |
   | `MODEL_ID`      | Gemini model name (e.g. `gemini-1.5-flash`)         |
   | `DB_HOST`       | Database host (e.g. RDS endpoint or other Postgres) |
   | `DB_PORT`       | Database port (e.g. `5432`)                         |
   | `DB_NAME`       | Database name                                       |
   | `DB_USER`       | Database user                                       |
   | `DB_PASSWORD`   | Database password                                   |

   Example secret value (JSON):
   ```json
   {
     "GEMINI_API_KEY": "your-api-key",
     "MODEL_ID": "gemini-1.5-flash",
     "DB_HOST": "your-rds-endpoint",
     "DB_PORT": "5432",
     "DB_NAME": "aiagent",
     "DB_USER": "aiagent",
     "DB_PASSWORD": "prodpassword"
   }
   ```

5. Run database migrations (first time or when the DB schema changes).

   You can either use VSCode or run the migrate command directly.

   **Option A – VSCode launch config**

   Add the following configuration to your `.vscode/launch.json` (or adapt the existing one) and run the `Migrate` config once:

   ```jsonc
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "Run Service",
         "type": "go",
         "request": "launch",
         "mode": "auto",
         "program": "${workspaceFolder}/cmd/main.go",
         "env": {
           "SECRETS_ARN": "ai-agent-dev-secret-secrets",
           "GIN_MODE": "debug"
         }
       },
       {
         "name": "Migrate",
         "type": "go",
         "request": "launch",
         "mode": "auto",
         "program": "${workspaceFolder}/cmd/migrate/main.go",
         "env": {
           "SECRETS_ARN": "ai-agent-dev-secret-secrets",
           "GIN_MODE": "debug",
           "DB_MIGRATIONS_PATH": "${workspaceFolder}/migrations"
         }
       }
     ]
   }
   ```

   For production, you typically use a different `SECRETS_ARN` that points to your production secret in AWS Secrets Manager (for example `ai-agent-prod-secret-secrets`).

   **Option B – CLI**

   From the project root:

   ```bash
   SECRETS_ARN="ai-agent-prod-secret-secrets" \
   GIN_MODE=debug \
   DB_MIGRATIONS_PATH="$(pwd)/migrations" \
     go run ./cmd/migrate
   ```

Now the database is migrated, the application is deployed on AWS Lambda and ready to use.

### Running Locally with LocalStack

You can run the full stack locally using LocalStack to emulate AWS services. Postgres runs via Docker Compose (no RDS).

1. Copy `.env.example` to `.env` and set the Postgres credentials (used by docker-compose):
   ```bash
   cp .env.example .env
   # Edit .env if needed (defaults: POSTGRES_USER=aiagent, POSTGRES_PASSWORD=localpostgres, POSTGRES_DB=aiagent)
   ```

2. Start LocalStack and Postgres using Docker Compose:
   ```bash
   docker compose up -d
   ```
3. In a separate terminal, apply Terraform resources against LocalStack using `tflocal`:
   ```bash
   cd .tf
   tflocal init
   tflocal apply
   cd ..
   ```

   Note: you might run into an error if the S3 bucket does not exist. To resolve pass the configuration via CLI arguments alongside S3_HOSTNAME=localhost to fix internal DNS routing.
   ```bash
   S3_HOSTNAME=localhost tflocal init -backend-config="force_path_style=true" -reconfigure
   ```

4. Export AWS environment variables so the Go SDK and SAM talk to LocalStack:
   ```bash
   export AWS_ACCESS_KEY_ID=test
   export AWS_SECRET_ACCESS_KEY=test
   export AWS_REGION=eu-central-1
   export AWS_ENDPOINT_URL=http://localhost:4566
   ```
5. Add the secret value to AWS Secrets Manager (LocalStack). Use the same keys as above. For local runs, use `DB_HOST=host.docker.internal` so Lambda can reach the Postgres container on your host. Example via AWS CLI:
   ```bash
   aws --endpoint-url=http://localhost:4566 secretsmanager put-secret-value \
     --secret-id ai-agent-lambda-test-secret-secrets \
     --secret-string '{"GEMINI_API_KEY":"your-key","MODEL_ID":"gemini-1.5-flash","DB_HOST":"host.docker.internal","DB_PORT":"5432","DB_NAME":"aiagent","DB_USER":"aiagent","DB_PASSWORD":"localpostgres"}'
   ```

6. Run database migrations against the LocalStack Postgres (first time or on schema change).

   **Option A – VSCode launch config**

   Add a LocalStack-specific migrate configuration if you want to separate it:

   ```jsonc
   {
     "name": "Migrate (LocalStack)",
     "type": "go",
     "request": "launch",
     "mode": "auto",
     "program": "${workspaceFolder}/cmd/migrate/main.go",
     "env": {
       "SECRETS_ARN": "ai-agent-lambda-test-secret-secrets",
       "GIN_MODE": "debug",
       "AWS_ENDPOINT_URL": "http://localhost:4566",
       "DB_MIGRATIONS_PATH": "${workspaceFolder}/migrations"
     }
   }
   ```

   **Option B – CLI**

   ```bash
   SECRETS_ARN="ai-agent-lambda-test-secret-secrets" \
   GIN_MODE=debug \
   AWS_ENDPOINT_URL="http://localhost:4566" \
   DB_MIGRATIONS_PATH="$(pwd)/migrations" \
     go run ./cmd/migrate
   ```

7. Run the Lambda locally via SAM using `template.yaml` and one of the `invokation/*.json` events:
   ```bash
   sam local invoke AiAgentServerlessAPI \
     -t template.yaml \
     -e invokation/<event-file>.json
   ```

8. **Call the API via REST (API Gateway)**  
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

### Running Locally for developement

Following the steps above to run Localstack and Postgres locally, you may run the Go application also locally. To boostrap it with localstack and DB, you may provide the environment variables in the run command.

Since we're using [aws-lambda-web-adapter](https://github.com/awslabs/aws-lambda-web-adapter), the service is able to run also standalone not just as a lambda.

#### Running VSCode

In VSCode, you may defined a `launch.json` config file, which makes you able to run the application locally via the IDE.

Adapt the config below to match your set Environment variables.

```jsonc
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run Service",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/main.go",
      "env": {
        "SECRETS_ARN": "ai-agent-dev-secret-secrets",
        "GIN_MODE": "debug"
      }
    },
    {
      "name": "Migrate (Local Dev)",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/migrate/main.go",
      "env": {
        "SECRETS_ARN": "ai-agent-dev-secret-secrets",
        "GIN_MODE": "debug",
        "DB_MIGRATIONS_PATH": "${workspaceFolder}/migrations"
      }
    }
  ]
}
```

To run migrations for your local dev DB, just run the `Migrate (Local Dev)` configuration once from VSCode.

#### Running via CLI

You can also run both the migrate script and the service via CLI.

Run migrations:

```bash
SECRETS_ARN="ai-agent-dev-secret-secrets" \
GIN_MODE=debug \
DB_MIGRATIONS_PATH="$(pwd)/migrations" \
  go run ./cmd/migrate
```

Run the service:

```bash
SECRETS_ARN="ai-agent-dev-secret-secrets" \
GIN_MODE=debug \
  go run ./cmd/main.go
```

Once running, call the local HTTP endpoint (port depends on your Gin setup and adapter config), for example:

```bash
curl http://localhost:8080/api/v1/health
```
