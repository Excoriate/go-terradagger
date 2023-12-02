variable "is_enabled" {
  type        = bool
  description = <<EOF
  Whether this module will be created or not. It is useful, for stack-composite
modules that conditionally includes resources provided by this module..
EOF
}

variable "aws_region" {
  type        = string
  description = "AWS region to deploy the resources"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources."
  default     = {}
}

/*
-------------------------------------
Custom input variables
-------------------------------------
*/
variable "module_config" {
  type = list(object({
    name = string
  }))
  description = <<EOF
Try to put a meaningful description here. Hopefully, referencing the
documentation of the module that is being instantiated.
EOF
  default     = null
}
