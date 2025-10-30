output "db_arn" {
  value = aws_dynamodb_table.chat_history.arn
}

output "table_name" {
    value = aws_dynamodb_table.chat_history.name
}

output "user_id_key" {
    value = var.user_id
}

output "timestamp_key" {
    value = var.timestamp
}