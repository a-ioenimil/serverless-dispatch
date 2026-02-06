locals {
  binary_name = "bootstrap"
  # Build key: CGO disabled, OS linux, Arch arm64. Output must be 'bootstrap' for provided.al2023
  build_command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go"
}

module "api_create_task" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-api-create-task"
  handler       = local.binary_name
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  create_package = true
  source_path = [
    {
      path     = "${var.source_dir}/api-task-create"
      commands = [local.build_command, ":zip"]
      patterns = ["*.go"]
    }
  ]

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
  source_path = [
    {
      path     = "${var.source_dir}/api-task-list"
      commands = [local.build_command, ":zip"]
      patterns = ["*.go"]
    }
  ]

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["dynamodb:GetItem", "dynamodb:Query"]
      Resource = var.dynamodb_table_arn
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
  source_path = [
    {
      path     = "${var.source_dir}/api-task-update"
      commands = [local.build_command, ":zip"]
      patterns = ["*.go"]
    }
  ]

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["dynamodb:UpdateItem", "dynamodb:GetItem"]
      Resource = var.dynamodb_table_arn
    }]
  })

  environment_variables = {
    TABLE_NAME = var.dynamodb_table_id
  }
}

