resource "aws_secretsmanager_secret" "ai_agent_secrets" {
  name = "${var.secret_name}-secrets"

  tags = {
    Name        = var.function_name
    Environment = var.environment
  }
}
