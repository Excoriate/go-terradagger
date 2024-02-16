<!-- BEGIN_TF_DOCS -->
# ☁️ Example module

## Description


Provide all the description that's required.

* 🚀 **Functionality** - Describe what the module does.

  ---

  ```hcl
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

  ```
  ---

  ## Module's documentation
  (This documentation is auto-generated using [terraform-docs](https://terraform-docs.io))
  ## Providers

No providers.

  ## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_dynamodb_table"></a> [dynamodb\_table](#module\_dynamodb\_table) | terraform-aws-modules/dynamodb-table/aws | n/a |

  ## Resources

No resources.

  ## Requirements

No requirements.

  ## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_is_enabled"></a> [is\_enabled](#input\_is\_enabled) | Whether the table is enabled | `bool` | n/a | yes |
| <a name="input_table_name"></a> [table\_name](#input\_table\_name) | The name of the table | `string` | n/a | yes |

  ## Outputs

No outputs.
<!-- END_TF_DOCS -->
