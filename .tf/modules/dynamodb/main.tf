resource "aws_dynamodb_table" "chat_history" {
  name         = var.table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = var.user_id
  range_key    = var.timestamp

  attribute {
    name = var.user_id
    type = "S"
  }

  attribute {
    name = var.timestamp
    type = "N"
  }

  tags = {
    Group       = var.function_name
    Environment = var.environment
  }
}
