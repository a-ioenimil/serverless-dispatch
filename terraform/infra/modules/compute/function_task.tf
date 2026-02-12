module "api_create_task" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-api-create-task"
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
      path     = "${var.source_dir}/api-task-create"
      commands = [local.build_command, ":zip"]
      patterns = [".*\\.go"]
    }
  ] : null

  # Pre-Built Source
  s3_existing_package = var.app_version != null ? {
    bucket = var.artifact_bucket_id
    key    = "builds/${var.app_version}/api-task-create.zip"
  } : null

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["dynamodb:PutItem", "dynamodb:UpdateItem"]
      Resource = var.dynamodb_table_arn
    }]
  })

  environment_variables = {
    TABLE_NAME = var.dynamodb_table_id
  }
}

module "api_get_task" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-api-get-task"
  handler       = local.binary_name
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  create_package = true
  store_on_s3    = true
  s3_bucket      = var.artifact_bucket_id

  source_path = [
    {
      path     = "${var.source_dir}/api-task-list"
      commands = [local.build_command, ":zip"]
      patterns = [".*\\.go"]
    }
  ]

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = ["dynamodb:GetItem", "dynamodb:Query", "dynamodb:Scan"]
      Resource = [
        var.dynamodb_table_arn,
        "${var.dynamodb_table_arn}/index/GSI1"
      ]
    }]
  })

  environment_variables = {
    TABLE_NAME = var.dynamodb_table_id
  }
}


module "api_update_task" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-api-update-task"
  handler       = local.binary_name
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  create_package = true
  store_on_s3    = true
  s3_bucket      = var.artifact_bucket_id

  source_path = [
    {
      path     = "${var.source_dir}/api-task-update"
      commands = [local.build_command, ":zip"]
      patterns = [".*\\.go"]
    }
  ]

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:GetItem",
        "dynamodb:Query"
      ]
      Resource = [
        var.dynamodb_table_arn,
        "${var.dynamodb_table_arn}/index/*"
      ]
    }]
  })

  environment_variables = {
    TABLE_NAME = var.dynamodb_table_id
  }
}

module "api_list_users" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-api-list-users"
  handler       = local.binary_name
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  create_package = true
  store_on_s3    = true
  s3_bucket      = var.artifact_bucket_id

  source_path = [
    {
      path     = "${var.source_dir}/api-user-list"
      commands = [local.build_command, ":zip"]
      patterns = [".*\\.go"]
    }
  ]

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["dynamodb:Scan"]
      Resource = var.dynamodb_table_arn
    }]
  })

  environment_variables = {
    TABLE_NAME = var.dynamodb_table_id
  }
}

