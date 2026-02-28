output "rest_api_id" {
  description = "The ID of the REST API."
  value       = aws_api_gateway_rest_api.api.id
}

output "stage_name" {
  description = "The name of the deployed stage."
  value       = aws_api_gateway_stage.stage.stage_name
}

output "invoke_url" {
  description = "Base URL to invoke the API (AWS format). Append path e.g. /api/v1/health."
  value       = "${aws_api_gateway_stage.stage.invoke_url}/"
}

output "localstack_invoke_url" {
  description = "Base URL for LocalStack. Append path e.g. /api/v1/health."
  value       = "http://localhost:4566/_aws/execute-api/${aws_api_gateway_rest_api.api.id}/${aws_api_gateway_stage.stage.stage_name}"
}
