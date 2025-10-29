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
