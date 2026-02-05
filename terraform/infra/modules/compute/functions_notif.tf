module "async_notifier" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-async-notifier"
  handler       = "main"
  runtime       = "go1.x"

  create_package = true
  source_path = [
    {
      path     = "${var.source_dir}/async-notifier"
      commands = ["GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go", ":zip"]
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

  event_source_mappings = {
    dynamodb = {
      event_source_arn  = var.dynamodb_table_stream_arn
      starting_position = "LATEST"
      batch_size        = 1
    }
  }
}
