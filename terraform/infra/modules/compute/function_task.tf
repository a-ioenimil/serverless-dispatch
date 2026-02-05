module "api_create_task" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-api-create-task"
  handler       = "main"
  runtime       = "go1.x"

  create_package = true
  source_path = [
    {
      path = "${var.source_dir}/api-task-create"
      commands = ["GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go", ":zip"]
      patterns = ["*.go"]
    }
  ]

  attach_policy_json = true
  policy_json        = jsonencode({
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
  handler       = "main"
  runtime       = "go1.x"

  create_package = true
  source_path = [
    {
      path = "${var.source_dir}/api-task-get"
      commands = ["GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go", ":zip"]
      patterns = ["*.go"]
    }
  ]

  attach_policy_json = true
  policy_json        = jsonencode({
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

module "api_list_tasks" {
  source = 
}
