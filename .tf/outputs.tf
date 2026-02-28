output "localstack_invoke_url" {
  description = "Base URL for LocalStack REST API. Append path e.g. /api/v1/health."
  value       = module.ai_agent_api_gateway.localstack_invoke_url
}

output "api_gateway_invoke_url" {
  description = "Base URL for API Gateway (AWS). Append path e.g. /api/v1/health."
  value       = module.ai_agent_api_gateway.invoke_url
}
