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
