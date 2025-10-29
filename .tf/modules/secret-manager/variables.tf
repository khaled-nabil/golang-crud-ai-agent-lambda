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
