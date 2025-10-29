# AI Agent Lambda

This project is a serverless AI agent built with Go and deployed on AWS Lambda. It uses the Gin framework for routing and integrates with Google's Gemini AI model. The application is designed to be a conversational AI, with plans to implement a sophisticated RAG (Retrieval-Augmented Generation) system.

## Dependencies

- Go 1.24.0
- Gin
- AWS Lambda Go
- AWS SDK for Go V2
- Google Wire
- Google Gemini Go SDK

## Setup and Run

### Prerequisites

- Go 1.24.0 or later
- AWS SAM CLI
- Docker

### Configuration

To run with local SAM CLI, Create an `env.json` file in the root directory with the following content:

```json
{
  "AiAgentServerlessAPI": {
	"GEMINI_API_KEY": "YOUR_GEMINI_API_KEY",
    "MODEL_ID": "DESIRED_MODEL_ID"
  }
}
```

### Build

```bash
sam build
```

### Run locally

```bash
sam local start-api --env-vars env.json
```

### Deploy

```bash
sam deploy --guided
```

## Future Work

- **Create a smart RAG system for AI:**
  - It would use DynamoDB to store user's conversation by user's ID.
  - It would pull the latest 20 user messages.
  - When a user has more than 20 messages in history, it would summarize the user's chat into a new record in DynamoDB. This new summarization record would trim the context.
  - The system would include all summaries of the user, and the latest 20 messages.
  - When the user has more than X summaries (to be defined based on context window), we would use an in-memory vector search. All summaries are vectorized using an embedding model (to be decided on), and based on similarity search, we'd pull the relevant summaries that wouldn't exceed the context window.
- **Develop a frontend React application.**

