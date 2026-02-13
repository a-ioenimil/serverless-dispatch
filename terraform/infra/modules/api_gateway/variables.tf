variable "project_name" {
  type = string
}

variable "region" {
  type = string
}

variable "user_pool_id" {
  type = string
}

variable "user_pool_client_id" {
  type = string
}

variable "create_task_function_name" {
  type = string
}

variable "create_task_invoke_arn" {
  type = string
}

variable "get_task_function_name" {
  type = string
}

variable "get_task_invoke_arn" {
  type = string
}

variable "update_task_function_name" {
  type = string
}

variable "update_task_invoke_arn" {
  type = string
}

variable "delete_task_function_name" {
  type = string
}

variable "delete_task_invoke_arn" {
  type = string
}

variable "list_users_function_name" {
  type = string
}

variable "list_users_invoke_arn" {
  type = string
}
