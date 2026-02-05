module "dynamodb_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "4.0.0"

  name      = "${var.project_name}-${var.environment}-table"
  hash_key  = "PK"
  range_key = "SK"

  attributes = [
    { name = "PK", type = "S" },
    { name = "SK", type = "S" },
    { name = "GSI1_PK", type = "S" },
    { name = "GSI1_SK", type = "S" }
  ]

  global_secondary_indexes = [
    {
      name            = "GSI1"
      hash_key        = "GSI1_PK"
      range_key       = "GSI1_SK"
      projection_type = "ALL"
      read_capacity   = 1
      write_capacity  = 1
    }
  ]

  billing_mode   = "PAY_PER_REQUEST"
  read_capacity  = 1
  write_capacity = 1
}
