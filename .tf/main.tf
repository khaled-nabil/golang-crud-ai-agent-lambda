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
  source = "./modules/lambda"
  # fix cyclical loop to apply depends on
  # depends_on = [module.ai_agent_secrets]

  function_name = "ai-agent-lambda"
  app_name      = "ai-agent"
  environment   = "test"

  go_app_path  = "${path.module}/../cmd/main.go"
  binary_path  = "${path.module}/.build/bootstrap"
  archive_path = "${path.module}/.build/archive.zip"

  secrets_arn = module.ai_agent_secrets.secret_arn
}

module "ai_agent_secrets" {
  source = "./modules/secret-manager"

  function_name = module.ai-agent-lambda.ai_agent_localised_function_name
}