variable "dynamodb_table_arn" {
  description = "The ARN of the DynamoDB table."
  type        = string
}
variable "dynamodb_table_id" {
  description = "The ID (name) of the DynamoDB table."
  type        = string
}

variable "dynamodb_table_stream_arn" {
  description = "The ARN of the DynamoDB table stream."
  type        = string
}

variable "source_dir" {
  description = "Path to the Go source code for Lambda functions."
  type        = string
}

variable "project_name" {
  description = "The name of the project for tagging and naming resources."
  type        = string

}

variable "allowed_email_domains" {
  description = "Comma-separated list of allowed email domains for user sign-up."
  type        = list(string)
}

variable "artifact_bucket_id" {
  description = "S3 Bucket ID for storing Lambda artifacts"
  type        = string
}

variable "environment" {
  description = "The deployment environment (e.g., dev, staging, prod)."
  type        = string
}


variable "region" {
  description = "AWS Region"
  type        = string
}

variable "app_version" {
  description = "The version of the app artifacts. If null, builds locally."
  type        = string
  default     = null
}
