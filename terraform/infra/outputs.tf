output "artifact_bucket_id_output" {
  description = "S3 artifact bucket store"
  value       = module.storage.bucket_id
}
output "api_endpoint" {
  description = "API endpoint url"
  value       = module.api_gateway.api_endpoint
}

output "cognito_user_pool_id" {
  description = "The ID of the User Pool"
  value       = module.auth.user_pool_id
}

output "cognito_user_pool_client_id" {
  description = "The ID of the User Pool Client"
  value       = module.auth.user_pool_client_id
}

output "aws_region" {
  description = "The AWS region"
  value       = var.region
}

output "amplify_app_id" {
  description = "Amplify app ID"
  value       = module.amplify.app_id
}

output "amplify_default_domain" {
  description = "Amplify default domain"
  value       = module.amplify.default_domain
}
