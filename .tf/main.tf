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

  function_name = "ai-agent-lambda"
  app_name      = "ai-agent"
  environment   = "test"

  go_app_path  = "${path.module}/../cmd/main.go"
  binary_path  = "${path.module}/.build/bootstrap"
  archive_path = "${path.module}/.build/archive.zip"
}
