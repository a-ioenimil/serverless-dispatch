output "artifact_bucket_id_output" {
  value = module.storage.bucket_id
}
output "api_endpoint" {
  value = module.api_gateway.api_endpoint
}
