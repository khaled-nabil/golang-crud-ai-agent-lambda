resource "aws_dynamodb_table" "chat_history" {
  name         = var.table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "user_id"
  range_key    = "timestamp"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "S"
  }

  tags = {
    Group        = var.function_name
    Environment = var.environment
  }
}
