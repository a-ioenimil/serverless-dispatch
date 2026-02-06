resource "aws_cognito_user_pool" "main" {
  name = "${var.project_name}-user-pool"

  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  username_attributes = ["email"]

  # Prevent explicit user creation (Self-signup is allowed via client, but admin only creates admins)
  admin_create_user_config {
    allow_admin_create_user_only = false
  }

  lambda_config {
    pre_sign_up       = var.pre_sign_up_lambda_arn
    post_confirmation = var.post_confirmation_lambda_arn
  }
}

resource "aws_cognito_user_pool_client" "client" {
  name = "${var.project_name}-app-client"

  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret = false
  explicit_auth_flows = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]
}

resource "aws_cognito_user_group" "admins" {
  name         = "Admins"
  user_pool_id = aws_cognito_user_pool.main.id
  description  = "Administrator Group"
  precedence   = 1
}

resource "aws_cognito_user_group" "members" {
  name         = "Members"
  user_pool_id = aws_cognito_user_pool.main.id
  description  = "Standard Members Group"
  precedence   = 5
}

resource "aws_lambda_permission" "pre_signup" {
  statement_id  = "AllowExecutionFromCognitoPreSignup"
  action        = "lambda:InvokeFunction"
  function_name = var.pre_sign_up_lambda_arn
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.main.arn
}

resource "aws_lambda_permission" "post_confirmation" {
  statement_id  = "AllowExecutionFromCognitoPostConfirmation"
  action        = "lambda:InvokeFunction"
  function_name = var.post_confirmation_lambda_arn
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.main.arn
}

