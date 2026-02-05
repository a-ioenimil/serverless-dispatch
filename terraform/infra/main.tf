# 1. Database Module
module "database" {
  source = "./modules/database"

  project_name = var.project_name
  environment  = var.environment
}

# 2. Auth Module (Cognito)
module "auth" {
  source = "./modules/auth"

  project_name = var.project_name
  environment  = var.environment

  # Pass Lambda ARNs for Triggers (Circular dependency handled via variable injection if needed, 
  # but here we output the Pool ID for the Lambdas to use)
}

# 3. Compute Module (Lambdas)
module "compute" {
  source = "./modules/compute"

  project_name = var.project_name
  environment  = var.environment
  source_dir   = "${path.module}/../../functions/cmd"

  # Dependency Injection
  dynamodb_table_arn        = module.database.table_arn
  dynamodb_table_id         = module.database.table_id
  dynamodb_table_stream_arn = module.database.table_stream_arn
  user_pool_arn             = module.auth.user_pool_arn

  # Point to the Go Source Code relative to the module
  # We pass the absolute path to be safe
  source_dir = abspath("${path.module}/../functions/cmd")
}

# 4. Attach Triggers to Auth (Post-deployment wiring)
# Since Cognito needs Lambda ARNs, and Lambdas need Cognito permissions,
# we sometimes wire the specific "Trigger" attachment in the root or a dedicated wiring resource.
