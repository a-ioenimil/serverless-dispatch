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

variable "environment" {
  description = "The deployment environment (e.g., dev, staging, prod)."
  type        = string
}

variable "user_pool_id" {
  description = "Cognito User Pool ID"
  type        = string
}

variable "user_pool_client_id" {
  description = "Cognito User Pool Client ID"
  type        = string
}

variable "user_pool_arn" {
  description = "Cognito User Pool ARN"
  type        = string
}

variable "region" {
  description = "AWS Region"
  type        = string
}
variable "allowed_email_domains" {
  description = "Allowed email domains for user sign-up."
  type        = list(string)
}
