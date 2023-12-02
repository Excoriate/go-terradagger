<!-- BEGIN_TF_DOCS -->
# ☁️ Example module

## Description


Provide all the description that's required.

* 🚀 **Functionality** - Describe what the module does.

---

## Example

Examples of this module's usage are available in the [examples](./examples) folder.

```hcl
resource "random_id" "this" {
  byte_length = 8
}

resource "random_password" "this" {
  length  = 16
  special = false
}

resource "random_string" "this" {
  length  = 16
  special = false
}

resource "random_uuid" "this" {
}
```

---

## Module's documentation

(This documentation is auto-generated using [terraform-docs](https://terraform-docs.io))

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.5.1 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [random_id.this](https://registry.terraform.io/providers/hashicorp/random/3.5.1/docs/resources/id) | resource |
| [random_password.this](https://registry.terraform.io/providers/hashicorp/random/3.5.1/docs/resources/password) | resource |
| [random_string.this](https://registry.terraform.io/providers/hashicorp/random/3.5.1/docs/resources/string) | resource |
| [random_uuid.this](https://registry.terraform.io/providers/hashicorp/random/3.5.1/docs/resources/uuid) | resource |

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.6 |
| <a name="requirement_random"></a> [random](#requirement\_random) | 3.5.1 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | AWS region to deploy the resources | `string` | n/a | yes |
| <a name="input_is_enabled"></a> [is\_enabled](#input\_is\_enabled) | Whether this module will be created or not. It is useful, for stack-composite<br>modules that conditionally includes resources provided by this module.. | `bool` | n/a | yes |
| <a name="input_module_config"></a> [module\_config](#input\_module\_config) | Try to put a meaningful description here. Hopefully, referencing the<br>documentation of the module that is being instantiated. | <pre>list(object({<br>    name = string<br>  }))</pre> | `null` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | A map of tags to add to all resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_aws_region_for_deploy_this"></a> [aws\_region\_for\_deploy\_this](#output\_aws\_region\_for\_deploy\_this) | The AWS region where the module is deployed. |
| <a name="output_is_enabled"></a> [is\_enabled](#output\_is\_enabled) | Whether the module is enabled or not. |
| <a name="output_random_id"></a> [random\_id](#output\_random\_id) | The random id generated. |
| <a name="output_random_password"></a> [random\_password](#output\_random\_password) | The random pet generated. |
| <a name="output_random_string"></a> [random\_string](#output\_random\_string) | The random string generated. |
| <a name="output_random_uuid"></a> [random\_uuid](#output\_random\_uuid) | The random uuid generated. |
| <a name="output_tags_set"></a> [tags\_set](#output\_tags\_set) | The tags set for the module. |
<!-- END_TF_DOCS -->
