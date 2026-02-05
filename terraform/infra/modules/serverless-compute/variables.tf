variable "dynamodb_table_arn" {
    description = "The ARN of the DynamoDB table."
    type        = string
}
variable "dynamodb_table_id" {
    description = "The ID (name) of the DynamoDB table."
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