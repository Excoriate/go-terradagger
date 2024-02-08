terraform {
  source = "modules/random"
}

generate "provider" {
  path = "provider.tf"
  if_exists = "overwrite_terragrunt"
    contents = <<EOF
provider "random" {
    version = "3.5.1"
}
EOF
}

inputs = {
  is_enabled = true
}
