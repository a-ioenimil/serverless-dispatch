module "auth_pre_sign_up" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-auth-pre-sign-up"
  handler       = local.binary_name
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  # Artifact Config: Local Build vs Pre-Built S3
  create_package = var.app_version == null ? true : false
  store_on_s3    = var.app_version == null ? true : false
  s3_bucket      = var.artifact_bucket_id

  # Local Build Source
  source_path = var.app_version == null ? [
    {
      path     = "${var.source_dir}/auth-pre-signup"
      commands = [local.build_command, ":zip"]
      patterns = [".*\\.go"]
    }
  ] : null

  # Pre-Built Source
  s3_existing_package = var.app_version != null ? {
    bucket = var.artifact_bucket_id
    key    = "builds/${var.app_version}/auth-pre-signup.zip"
  } : null

  environment_variables = {
    ALLOWED_EMAIL_DOMAINS = join(",", var.allowed_email_domains)
  }
}

module "auth_post_confirmation" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-auth-post-confirmation"
  handler       = local.binary_name
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  # Artifact Config: Local Build vs Pre-Built S3
  create_package = var.app_version == null ? true : false
  store_on_s3    = var.app_version == null ? true : false
  s3_bucket      = var.artifact_bucket_id

  # Local Build Source
  source_path = var.app_version == null ? [
    {
      path     = "${var.source_dir}/auth-post-signup"
      commands = [local.build_command, ":zip"]
      patterns = [".*\\.go"]
    }
  ] : null

  # Pre-Built Source
  s3_existing_package = var.app_version != null ? {
    bucket = var.artifact_bucket_id
    key    = "builds/${var.app_version}/auth-post-signup.zip"
  } : null

  environment_variables = {
    TABLE_NAME              = var.dynamodb_table_id
    NOTIFICATIONS_TOPIC_ARN = aws_sns_topic.notifications.arn
  }

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["dynamodb:PutItem"]
        Resource = var.dynamodb_table_arn
      },
      {
        Effect = "Allow"
        Action = [
          "sns:Subscribe",
          "sns:SetSubscriptionAttributes",
          "sns:ListSubscriptionsByTopic"
        ]
        Resource = aws_sns_topic.notifications.arn
      }
    ]
  })
}
