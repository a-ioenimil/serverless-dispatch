# 1. Database Module
module "database" {
  source = "./modules/database"

  project_name = var.project_name
  environment  = var.environment
}

# 2. Storage Module (Artifacts)
module "storage" {
  source = "./modules/storage"

  project_name = var.project_name
  environment  = var.environment
}

# 3. Compute Module (Lambdas)
# Compute now comes before Auth because Auth needs Lambda ARNs for triggers
module "compute" {
  source = "./modules/compute"

  project_name = var.project_name
  environment  = var.environment

  # Dependency Injection
  dynamodb_table_arn        = module.database.table_arn
  dynamodb_table_id         = module.database.table_id
  dynamodb_table_stream_arn = module.database.table_stream_arn
  artifact_bucket_id        = module.storage.bucket_id
  app_version               = var.app_version

  # Configuration
  allowed_email_domains = var.allowed_email_domains
  region                = var.region
  user_pool_id          = var.user_pool_id

  # Point to the Go Source Code relative to the module
  source_dir = abspath("${path.module}/../../functions/cmd")
}

# 3. Auth Module (Cognito)
module "auth" {
  source = "./modules/auth"

  project_name = var.project_name
  environment  = var.environment

  # Lambda Triggers from Compute
  pre_sign_up_lambda_arn       = module.compute.auth_pre_sign_up_arn
  post_confirmation_lambda_arn = module.compute.auth_post_confirmation_arn
}

# 4. API Gateway Module
module "api_gateway" {
  source = "./modules/api_gateway"

  project_name = var.project_name
  region       = var.region

  user_pool_id        = module.auth.user_pool_id
  user_pool_client_id = module.auth.user_pool_client_id

  create_task_function_name = module.compute.create_task_function_name
  create_task_invoke_arn    = module.compute.create_task_invoke_arn

  get_task_function_name = module.compute.get_task_function_name
  get_task_invoke_arn    = module.compute.get_task_invoke_arn

  update_task_function_name = module.compute.update_task_function_name
  update_task_invoke_arn    = module.compute.update_task_invoke_arn

  delete_task_function_name = module.compute.delete_task_function_name
  delete_task_invoke_arn    = module.compute.delete_task_invoke_arn

  list_users_function_name = module.compute.list_users_function_name
  list_users_invoke_arn    = module.compute.list_users_invoke_arn
}
