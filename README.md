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

### Configuration
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

### Running Locally
To run the application locally after creating your aws resources, use SAM CLI.
1. use `template.yaml` as your SAM template.
2. Use a template from `invokation/*.json` for event data.
