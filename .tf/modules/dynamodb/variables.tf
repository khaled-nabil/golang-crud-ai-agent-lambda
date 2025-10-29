variable "function_name" {
  description = "The name of the Lambda function associated with the DynamoDB table."
  type        = string
}

variable "table_name" {
  description = "The name of the DynamoDB table."
  type        = string
}

variable "environment" {
  description = "The environment for the DynamoDB table."
  type        = string
}

variable "user_id" {
  description = "The user ID key for the DynamoDB table."
  type        = string
  default     = "user_id"
}

variable "timestamp" {
  description = "The timestamp key for the DynamoDB table."
  type        = string
  default     = "timestamp"
}
