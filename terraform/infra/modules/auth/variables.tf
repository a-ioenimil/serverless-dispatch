variable "project_name" {
  description = "The name of the project for tagging and naming resources."
  type        = string

}

variable "environment" {
  description = "The deployment environment (e.g., dev, staging, prod)."
  type        = string
}

variable "pre_sign_up_lambda_arn" {
  description = "ARN of the Lambda function for PreSignUp trigger"
  type        = string
}

variable "post_confirmation_lambda_arn" {
  description = "ARN of the Lambda function for PostConfirmation trigger"
  type        = string
}
