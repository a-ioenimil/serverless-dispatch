output "create_task_invoke_arn" {
  value = module.api_create_task.lambda_function_invoke_arn
}
output "create_task_function_name" {
  value = module.api_create_task.lambda_function_name
}

output "get_task_invoke_arn" {
  value = module.api_get_task.lambda_function_invoke_arn
}

output "get_task_function_name" {
  value = module.api_get_task.lambda_function_name
}

output "update_task_invoke_arn" {
  value = module.api_update_task.lambda_function_invoke_arn
}

output "update_task_function_name" {
  value = module.api_update_task.lambda_function_name
}

output "auth_pre_sign_up_arn" {
  value = module.auth_pre_sign_up.lambda_function_arn
}

output "auth_post_confirmation_arn" {
  value = module.auth_post_confirmation.lambda_function_arn
}
