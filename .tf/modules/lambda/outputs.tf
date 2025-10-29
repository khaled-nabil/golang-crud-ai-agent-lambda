output "ai_agent_function_name" {
  description = "The name of the Lambda function."
  value       = aws_lambda_function.ai_lambda_function.function_name
}

output "ai_agent_function_arn" {
  description = "The ARN of the Lambda function."
  value       = aws_lambda_function.ai_lambda_function.arn
}

output "ai_agent_localised_function_name" {
  description = "The name of the Lambda function with localised environment."
  value       = "${var.function_name}-${var.environment}"
}