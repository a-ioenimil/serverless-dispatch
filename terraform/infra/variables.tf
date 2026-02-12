variable "region" {
  description = "The AWS region to deploy resources in"
  type        = string
  default     = "eu-west-1"
}

variable "managed_by" {
  description = "Environment manager"
  type        = string
  default     = "Terraform"
}

variable "project_name" {
  description = "Project name for tagging resources"
  type        = string
  default     = "serverless-dispatch"
}

variable "environment" {
  description = "Deployment environment (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "allowed_email_domains" {
  description = "Comma-separated list of allowed email domains for user sign-up."
  type        = list(string)
  default     = ["amalitech.com", "amalitechtraining.org"]
}

variable "app_version" {
  description = "The version (e.g. git commit hash) of the application artifacts. If null, local build is used."
  type        = string
  default     = null
}

variable "user_pool_id" {
  description = "Cognito User Pool ID injected into backend lambdas that need directory lookups."
  type        = string
  default     = ""
}
