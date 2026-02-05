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
