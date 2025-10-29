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