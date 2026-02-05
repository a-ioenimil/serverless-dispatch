module "async_notifier" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-async-notifier"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  create_package = true
  source_path = [
    {
      path     = "${var.source_dir}/async-notifier"
      commands = ["GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go", ":zip"]
      patterns = ["*.go"]
    }
  ]

  # Permission to read from DynamoDB Stream and Send Emails
  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetRecords",
          "dynamodb:GetShardIterator",
          "dynamodb:DescribeStream",
          "dynamodb:ListStreams"
        ]
        Resource = var.dynamodb_table_stream_arn
      },
      {
        Effect   = "Allow"
        Action   = ["ses:SendEmail", "ses:SendRawEmail"]
        Resource = "*"
      }
    ]
  })

  environment_variables = {
    FROM_EMAIL = "notifications@amalitech.com" # Placeholder
  }
}

resource "aws_lambda_event_source_mapping" "dynamodb_stream" {
  event_source_arn               = var.dynamodb_table_stream_arn
  function_name                  = module.async_notifier.lambda_function_arn
  starting_position              = "LATEST"
  batch_size                     = 1
  bisect_batch_on_function_error = true
  maximum_retry_attempts         = 3
  function_response_types        = ["ReportBatchItemFailures"]

  destination_config {
    on_failure {
      destination_arn = aws_sqs_queue.notifier_dlq.arn
    }
  }
}

resource "aws_sqs_queue" "notifier_dlq" {
  name                      = "${var.project_name}-notifier-dlq"
  message_retention_seconds = 1209600 # 14 days
}
