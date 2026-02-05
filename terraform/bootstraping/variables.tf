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

variable "project" {
  description = "Project name for tagging resources"
  type        = string
  default     = "serverless-dispatch"
}

variable "environment" {
  description = "Deployment environment (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket for Terraform state."
  type        = string
  default     = "serverless-dispatch-tfstate"
}

variable "s3_bucket_versioning_enabled" {
  description = "Enable versioning for the S3 bucket."
  type        = bool
  default     = true
}
