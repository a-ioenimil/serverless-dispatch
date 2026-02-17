resource "aws_sns_topic" "notifications" {
  name = "${var.project_name}-notifications"
}


module "async_notifier" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-async-notifier"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  # Artifact Config: Local Build vs Pre-Built S3
  create_package = var.app_version == null ? true : false
  store_on_s3    = var.app_version == null ? true : false
  s3_bucket      = var.artifact_bucket_id

  # Local Build Source
  source_path = var.app_version == null ? [
    {
      path     = "${var.source_dir}/async-notifier"
      commands = ["GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go", ":zip"]
      patterns = [".*\\.go"]
    }
  ] : null

  # Pre-Built Source
  s3_existing_package = var.app_version != null ? {
    bucket = var.artifact_bucket_id
    key    = "builds/${var.app_version}/async-notifier.zip"
  } : null

  # Permission to read from DynamoDB Stream and publish notifications
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
        Action   = ["sns:Publish"]
        Resource = aws_sns_topic.notifications.arn
      },
      {
        Effect   = "Allow"
        Action   = ["sqs:SendMessage"]
        Resource = aws_sqs_queue.notifier_dlq.arn
      }
    ]
  })

  environment_variables = {
    NOTIFICATIONS_TOPIC_ARN = aws_sns_topic.notifications.arn
    USER_POOL_ID            = var.user_pool_id
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
