variable "app_name" {
  description = "Amplify app name"
  type        = string
}

variable "repository_url" {
  description = "Git repository URL connected to Amplify"
  type        = string
}

variable "api_base_url" {
  description = "Frontend API base URL"
  type        = string
}

variable "aws_region" {
  description = "AWS region for frontend environment variable"
  type        = string
}

variable "user_pool_id" {
  description = "Cognito user pool ID exposed to frontend"
  type        = string
}

variable "user_pool_client_id" {
  description = "Cognito user pool client ID exposed to frontend"
  type        = string
}

variable "environment_variables" {
  description = "Additional Amplify environment variables"
  type        = map(string)
  default     = {}
}

variable "tags" {
  description = "Tags applied to the Amplify app"
  type        = map(string)
  default     = {}
}

variable "build_spec" {
  description = "Amplify build spec"
  type        = string
  default     = <<-EOT
    version: 1
    applications:
      - frontend:
          phases:
            preBuild:
              commands:
                - npm ci --cache .npm --prefer-offline
            build:
              commands:
                - npm run build
          artifacts:
            baseDirectory: dist
            files:
              - '**/*'
          cache:
            paths:
              - .npm/**/*
        appRoot: frontend
  EOT
}
