terraform {
  backend "s3" {
    bucket = "terraform-state"
    key    = "ai-agent/terraform.tfstate"
    region = "eu-central-1"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
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
  region = var.aws_region

  default_tags {
    tags = {
      app = "ai-agent-resources"
    }
  }
}

module "ai-agent-lambda" {
  source     = "./modules/lambda"
  depends_on = [module.ai_agent_secrets]

  function_name = var.function_name
  environment   = var.environment

  app_name     = "ai-agent"
  go_app_path  = "${path.module}/../cmd/main.go"
  binary_path  = "${path.module}/.build/bootstrap"
  archive_path = "${path.module}/.build/archive.zip"

  secrets_arn      = module.ai_agent_secrets.secret_arn
}

module "ai_agent_api_gateway" {
  source     = "./modules/api-gateway"
  depends_on = [module.ai-agent-lambda]

  api_name             = "ai-agent-api"
  stage_name           = "local"
  lambda_function_name = module.ai-agent-lambda.function_name
  lambda_invoke_arn    = module.ai-agent-lambda.invoke_arn
  aws_region           = var.aws_region
}

module "ai_agent_secrets" {
  source = "./modules/secret-manager"

  function_name = var.function_name
  environment   = var.environment

  secret_name = "${var.function_name}-${var.environment}-secret"
}

output "lambda_function_url" {
  description = "The URL of the Lambda function."
  value       = module.ai-agent-lambda.lambda_function_url
}
