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

  # Connect to API Gateway
  allowed_triggers = {
    APIGateway = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
    }
  }

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

  allowed_triggers = {
    APIGateway = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
    }
  }

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

  allowed_triggers = {
    APIGateway = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
    }
  }

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

# Define Routes
resource "aws_apigatewayv2_route" "create_task" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "POST /tasks"
  target             = "integrations/${aws_apigatewayv2_integration.create_task.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_integration" "create_task" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = module.api_create_task.lambda_function_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_tasks" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "GET /tasks"
  target             = "integrations/${aws_apigatewayv2_integration.get_task.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_integration" "get_task" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = module.api_get_task.lambda_function_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "update_task" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "PUT /tasks/{taskId}" # Path parameter
  target             = "integrations/${aws_apigatewayv2_integration.update_task.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_integration" "update_task" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = module.api_update_task.lambda_function_arn
  payload_format_version = "2.0"
}

