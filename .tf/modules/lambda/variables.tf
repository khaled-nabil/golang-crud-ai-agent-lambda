variable "function_name" {
  description = "The name of the Lambda function."
  type        = string
}

variable "handler" {
  description = "The handler for the Lambda function."
  type        = string
  default     = "bootstrap"
}

variable "runtime" {
  description = "The runtime for the Lambda function."
  type        = string
  default     = "provided.al2023"
}

variable "go_app_path" {
  description = "The path to the Lambda function source code."
  type        = string
}

variable "binary_path" {
  description = "The output path for the zipped Lambda function."
  type        = string
}

variable "archive_path" {
  description = "The path to the Lambda function archive."
  type        = string
}

variable "app_name" {
  description = "The name of the application."
  type        = string
}

variable "environment" {
  description = "The environment for the application."
  type        = string
}

variable "memory_size" {
  description = "The amount of memory allocated to the Lambda function."
  type        = number
  default     = 128
}
