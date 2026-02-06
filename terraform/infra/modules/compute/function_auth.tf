module "auth_pre_sign_up" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-auth-pre-sign-up"
  handler       = "main"
  runtime       = "go1.x"

  create_package = true
  source_path = [
    {
      path     = "${var.source_dir}/auth-pre-signup"
      commands = ["GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go", ":zip"]
      patterns = ["*.go"]
    }
  ]

  environment_variables = {
    ALLOWED_EMAIL_DOMAINS = join(",", var.allowed_email_domains)
  }
}

module "auth_post_confirmation" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.2.0"

  function_name = "${var.project_name}-auth-post-confirmation"
  handler       = "main"
  runtime       = "go1.x"

  create_package = true
  source_path = [
    {
      path     = "${var.source_dir}/auth-post-signup"
      commands = ["GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go", ":zip"]
      patterns = ["*.go"]
    }
  ]

  environment_variables = {
    TABLE_NAME = var.dynamodb_table_id
  }

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["dynamodb:PutItem"]
      Resource = var.dynamodb_table_arn
    }]
  })
}
