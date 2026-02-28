resource "aws_secretsmanager_secret" "ai_agent_secrets" {
  name = "${var.secret_name}-secrets"

  tags = {
    Name        = var.function_name
    Environment = var.environment
  }
}

resource "aws_secretsmanager_secret_version" "ai_agent_secrets" {
  secret_id     = aws_secretsmanager_secret.ai_agent_secrets.id
  secret_string = var.secret_string
}
