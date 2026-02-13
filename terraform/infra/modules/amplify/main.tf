locals {
  default_env_vars = {
    AMPLIFY_DIFF_DEPLOY       = "false"
    AMPLIFY_MONOREPO_APP_ROOT = "frontend"
    VITE_API_BASE_URL         = var.api_base_url
    VITE_AWS_REGION           = var.aws_region
    VITE_COGNITO_CLIENT_ID    = var.user_pool_client_id
    VITE_COGNITO_USER_POOL_ID = var.user_pool_id
  }
}

resource "aws_amplify_app" "this" {
  name       = var.app_name
  repository = var.repository_url
  platform   = "WEB"

  enable_auto_branch_creation = false
  enable_basic_auth           = false
  enable_branch_auto_build    = false
  enable_branch_auto_deletion = false

  build_spec            = var.build_spec
  environment_variables = merge(local.default_env_vars, var.environment_variables)

  cache_config {
    type = "AMPLIFY_MANAGED_NO_COOKIES"
  }

  custom_rule {
    source = "/<*>"
    status = "404-200"
    target = "/index.html"
  }

  job_config {
    build_compute_type = "STANDARD_8GB"
  }

  tags = var.tags
}
