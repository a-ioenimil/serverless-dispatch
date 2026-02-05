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

