module "dynamodb_table" {
  count  = var.is_enabled ? 1 : 0
  source = "terraform-aws-modules/dynamodb-table/aws"

  name     = var.table_name
  hash_key = "id"

  attributes = [
    {
      name = "id"
      type = "N"
    }
  ]

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
