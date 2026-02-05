resource "aws_cognito_user_pool" "main" {
  name = "${var.project_name}-user-pool"
  auto_verified_attributes = ["email"]
  # ... config ...
}

resource "aws_cognito_user_group" "admins" {
  name         = "Admins"
  user_pool_id = aws_cognito_user_pool.main.id
}

