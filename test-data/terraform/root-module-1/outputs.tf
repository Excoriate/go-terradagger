output "is_enabled" {
  value       = local.is_enabled
  description = "Whether the module is enabled or not."
}

output "aws_region_for_deploy_this" {
  value       = local.aws_region_to_deploy
  description = "The AWS region where the module is deployed."
}

output "tags_set" {
  value       = var.tags
  description = "The tags set for the module."
}

/*
-------------------------------------
Custom outputs
-------------------------------------
*/
output "random_id" {
  value       = random_id.this
  description = "The random id generated."
}

output "random_password" {
  value       = random_password.this
  sensitive   = true
  description = "The random pet generated."
}

output "random_string" {
  value       = random_string.this
  description = "The random string generated."
}

output "random_uuid" {
  value       = random_uuid.this
  description = "The random uuid generated."
}
