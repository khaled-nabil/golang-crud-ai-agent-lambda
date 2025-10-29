resource "aws_secretsmanager_secret" "ai_agent_secrets" {
  name = "${var.function_name}-secrets"
}
