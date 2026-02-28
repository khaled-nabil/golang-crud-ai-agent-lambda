variable "api_name" {
  description = "Name of the REST API."
  type        = string
  default     = "ai-agent-api"
}

variable "stage_name" {
  description = "Name of the API Gateway stage."
  type        = string
  default     = "test"
}

variable "lambda_function_name" {
  description = "Name of the Lambda function to invoke."
  type        = string
}

variable "lambda_invoke_arn" {
  description = "Invoke ARN of the Lambda function."
  type        = string
}

variable "aws_region" {
  description = "AWS region (used for Lambda permission source_arn)."
  type        = string
}
