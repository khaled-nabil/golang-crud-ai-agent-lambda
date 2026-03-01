variable "aws_region" {
  description = "The AWS region to deploy the resources to."
  type        = string
  default     = "eu-central-1"
}

variable "function_name" {
  description = "The name of the Lambda function."
  type        = string
  default     = "ai-agent"
}

variable "environment" {
  description = "The environment (e.g., dev, prod)."
  type        = string
  default     = "dev"
}

