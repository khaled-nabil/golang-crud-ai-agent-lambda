variable "function_name" {
  description = "The name of the Lambda function."
  type        = string
}

variable "environment" {
  description = "The environment for the secret."
  type        = string
}

variable "secret_name" {
  description = "The name of the secret in AWS Secrets Manager."
  type        = string
}

variable "secret_string" {
  description = "Initial secret value (JSON). Default placeholder for local dev; in AWS set via CLI/console or override."
  type        = string
  default     = "{\"GEMINI_API_KEY\":\"localstack-placeholder\",\"MODEL_ID\":\"gemini-1.5-flash\"}"
  sensitive   = true
}
