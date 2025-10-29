terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

provider "aws" {
  region = "eu-central-1"

  default_tags {
    tags = {
      app = "ai-agent-resources"
    }
  }
}

module "ai-agent-lambda" {
  source     = "./modules/lambda"
  depends_on = [module.ai_agent_secrets, module.ai_agent_dynamodb]

  function_name = var.function_name
  environment   = var.environment

  app_name     = "ai-agent"
  go_app_path  = "${path.module}/../cmd/main.go"
  binary_path  = "${path.module}/.build/bootstrap"
  archive_path = "${path.module}/.build/archive.zip"

  secrets_arn      = module.ai_agent_secrets.secret_arn
  db_table_name    = module.ai_agent_dynamodb.table_name
  db_user_id_key   = module.ai_agent_dynamodb.user_id_key
  db_timestamp_key = module.ai_agent_dynamodb.timestamp_key
}

module "ai_agent_secrets" {
  source = "./modules/secret-manager"

  function_name = var.function_name
  environment   = var.environment

  secret_name = "${var.function_name}-${var.environment}-secret"
}

module "ai_agent_dynamodb" {
  source = "./modules/dynamodb"

  function_name = var.function_name
  environment   = var.environment

  table_name = "${var.function_name}-${var.environment}-table"
}
