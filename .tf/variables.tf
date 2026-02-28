variable "function_name" {
  description = "The name of the Lambda function."
  type        = string
  default     = "ai-agent-lambda"
}

variable "environment" {
  description = "The environment for the application."
  type        = string
  default     = "test"
}

variable "aws_region" {
  description = "AWS region to deploy to."
  type        = string
  default     = "eu-central-1"
}
