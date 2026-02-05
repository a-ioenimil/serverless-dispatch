resource "random_id" "bucket_suffix" {
  byte_length = 4

}

module "s3-bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "5.10.0"
  bucket  = "${var.s3_bucket_name}-${random_id.bucket_suffix.hex}"
  region  = var.region
  versioning = {
    enabled = var.s3_bucket_versioning_enabled
  }
  tags = {
    Name        = var.project
    Environment = var.environment
    ManagedBy   = var.managed_by
  }
}
