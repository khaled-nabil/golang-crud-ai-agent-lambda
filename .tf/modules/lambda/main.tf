data "aws_iam_policy_document" "assume_lambda_function_role" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ai_lambda_role" {
  name               = "${var.app_name}-${var.environment}-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_function_role.json
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.ai_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "dynamodb_access" {
  role       = aws_iam_role.ai_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}

resource "aws_iam_role_policy_attachment" "cloudwatch_logging" {
  role       = aws_iam_role.ai_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"
}

resource "null_resource" "build_binary" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${var.binary_path} ${var.go_app_path}"
  }
}

## ARCHIVE FILE
data "archive_file" "ai_lambda_archive" {
  depends_on = [null_resource.build_binary]

  type        = "zip"
  source_file = var.binary_path
  output_path = var.archive_path
}

resource "aws_lambda_function" "ai_lambda_function" {
  filename         = data.archive_file.ai_lambda_archive.output_path
  function_name    = var.function_name
  role             = aws_iam_role.ai_lambda_role.arn
  handler          = var.handler
  source_code_hash = data.archive_file.ai_lambda_archive.output_base64sha256
  runtime          = var.runtime
  memory_size      = var.memory_size

  tags = {
    Environment = var.environment
    Application = var.app_name
  }

  environment {
    variables = {
      SECRETS_ARN = var.secrets_arn
      GIN_MODE    = var.gin_mode
      DYNAMODB_TABLE_NAME = var.table_name
    }
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.ai_lambda_function.function_name}"
  retention_in_days = 7
}


## SECRET MANAGER
data "aws_iam_policy_document" "secrets_manager_read" {
  statement {
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [var.secrets_arn]
  }
}

resource "aws_iam_policy" "secrets_manager_read_policy" {
  name   = "${var.app_name}-${var.environment}-secrets-read-policy"
  policy = data.aws_iam_policy_document.secrets_manager_read.json
}

resource "aws_iam_role_policy_attachment" "secrets_manager_access" {
  role       = aws_iam_role.ai_lambda_role.name
  policy_arn = aws_iam_policy.secrets_manager_read_policy.arn
}
