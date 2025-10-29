output "db_arn" {
  value = aws_dynamodb_table.chat_history.arn
}

output "table_name" {
    value = aws_dynamodb_table.chat_history.name
}