config {
  module = true
  force = false
}

plugin "aws" {
  enabled = true
  deep_check = true
  version = "0.27.0"
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

rule "terraform_module_pinned_source" {
  enabled = true
}
