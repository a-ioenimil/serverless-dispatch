resource "aws_apigatewayv2_api" "http_api" {
  name          = "${var.project_name}-http-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"] # Lock this down in production!
    allow_methods = ["GET", "POST", "PUT", "DELETE"]
    allow_headers = ["Authorization", "Content-Type"]
    max_age       = 300
  }
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.http_api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito-authorizer"

  jwt_configuration {
    audience = [var.user_pool_client_id]
    issuer   = "https://cognito-idp.${var.region}.amazonaws.com/${var.user_pool_id}"
  }
}

# -----------------------------------------------------------------------------
# Route: Create Task (POST /tasks)
# -----------------------------------------------------------------------------
resource "aws_apigatewayv2_integration" "create_task" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = var.create_task_invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "create_task" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "POST /tasks"
  target             = "integrations/${aws_apigatewayv2_integration.create_task.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_lambda_permission" "create_task" {
  statement_id  = "AllowExecutionFromAPIGateway-CreateTask"
  action        = "lambda:InvokeFunction"
  function_name = var.create_task_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/POST/tasks"
}

# -----------------------------------------------------------------------------
# Route: Get Tasks (GET /tasks)
# -----------------------------------------------------------------------------
resource "aws_apigatewayv2_integration" "get_tasks" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = var.get_task_invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_tasks" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "GET /tasks"
  target             = "integrations/${aws_apigatewayv2_integration.get_tasks.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_lambda_permission" "get_tasks" {
  statement_id  = "AllowExecutionFromAPIGateway-GetTasks"
  action        = "lambda:InvokeFunction"
  function_name = var.get_task_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/GET/tasks"
}

# -----------------------------------------------------------------------------
# Route: Update Task (PUT /tasks/{taskId})
# -----------------------------------------------------------------------------
resource "aws_apigatewayv2_integration" "update_task" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = var.update_task_invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "update_task" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "PUT /tasks/{taskId}"
  target             = "integrations/${aws_apigatewayv2_integration.update_task.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_lambda_permission" "update_task" {
  statement_id  = "AllowExecutionFromAPIGateway-UpdateTask"
  action        = "lambda:InvokeFunction"
  function_name = var.update_task_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/PUT/tasks/*"
}
